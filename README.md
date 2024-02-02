## Github PR resource

[graphql-api]: https://developer.github.com/v4
[original-resource]: https://github.com/telia-oss/github-pr-resource
[git-resource]: https://github.com/concourse/git-resource
[#example]: #example
[status-checks]: https://docs.github.com/en/github/collaborating-with-pull-requests/collaborating-on-repositories-with-code-quality-features/about-status-checks

A proof-of-concept Concourse resource type that can operate over both lists of
Github pull requests and individual pull requests. It can be used to perform
several tasks:
1. Track a list of pull requests to a GitHub repository (typically to set a
   pipeline for each pull request).
2. Track commits to a single pull request by its number (e.g. to run tests
   against new commits).
3. Update pull request [status checks][status-checks] to indicate
   success/failure of build steps.
4. Add/remove pull request comments.

Refer to [#example] for a full example.

The code is adapted from [the original][original-resource], but its usage
pattern differs slightly. In particular:

* [The original][original-resource] tracks all commits across all PRs that
  satisfy the criteria set in `source`.
  * When `source.number` is not specified, this resource tracks the list of PRs
    that satisfy the `source` criteria. A new resource version is emitted
    whenever the list of PRs changes (*not* when a new commit is made to a PR)
  * When `source.number` is specified, this resource tracks commits to a
    single PR

There are some benefits to this approach:

* Each PR has its own build/resource version history, and its own build status.
* Since [the original][original-resource] uses a single version history for all
  commits, you must use `version: every`, which doesn't play nicely with
  `passed` constraints (https://github.com/concourse/concourse/issues/736)

There are also some downsides:

* Webhooks can't currently be configured easily for tracking new commits to
  existing PRs, as you would need to configure a webhook for *each* PR pipeline
  * Note that you can fairly easily use webhooks when tracking the list of PRs,
    though
* Task caches can't be shared between PR pipelines
* If you have `n` open PRs, you will now have `n` independent resources, which
  means more containers running checks

## Source Configuration

As noted earlier, this resource can either track a list of PRs to a repository,
or track commits to a single PR. The different modes of operation have
different configuration options.

### List of PRs

By omitting the `source.number` parameter, this resource will list PRs in a
repository.

| Parameter                   | Required | Example                          | Description                                                                                                                                                                                                                                                                                 |
|-----------------------------|----------|----------------------------------|---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `repository`                | Yes      | `itsdalmo/test-repository`       | The repository to target.                                                                                                                                                                                                                                                                   |
| `access_token`              | Yes      |                                  | A Github Access Token with repository access (required for setting status on commits). N.B. If you want github-pr-resource to work with a private repository. Set `repo:full` permissions on the access token you create on GitHub. If it is a public repository, `repo:status` is enough.  |
| `hosting_endpoint`          | No       | `https://github.com`             | Endpoint under which repositories are hosted. Specifically, the resource will pull from `hosting_endpoint`/`repository`.                                                                                                                                                                    |
| `v3_endpoint`               | No       | `https://api.github.com`         | Endpoint to use for the V3 Github API (Restful).                                                                                                                                                                                                                                            |
| `v4_endpoint`               | No       | `https://api.github.com/graphql` | Endpoint to use for the V4 Github API (Graphql).                                                                                                                                                                                                                                            |
| `paths`                     | No       | `["terraform/*/*.tf"]`           | Only produce new versions if the PR includes changes to files that match one or more glob patterns or prefixes.                                                                                                                                                                             |
| `ignore_paths`              | No       | `[".ci/"]`                       | Inverse of the above. Pattern syntax is documented in [filepath.Match](https://golang.org/pkg/path/filepath/#Match), or a path prefix can be specified (e.g. `.ci/` will match everything in the `.ci` directory).                                                                          |
| `disable_ci_skip`           | No       | `true`                           | Disable ability to skip builds with `[ci skip]` and `[skip ci]` in the pull request title.                                                                                                                                                                                                  |
| `skip_ssl_verification`     | No       | `true`                           | Disable SSL/TLS certificate validation on API clients. Use with care!                                                                                                                                                                                                                       |
| `disable_forks`             | No       | `true`                           | Disable triggering of the resource if the pull request's fork repository is different to the configured repository.                                                                                                                                                                         |
| `ignore_drafts`             | No       | `false`                          | Disable triggering of the resource if the pull request is in Draft status.                                                                                                                                                                                                                  |
| `required_review_approvals` | No       | `2`                              | Disable triggering of the resource if the pull request does not have at least `X` approved review(s).                                                                                                                                                                                       |
| `pr_numbers_higher_than`    | No       | `1000`                           | Disable triggering of the resource if the pull request number is not higher than this number.                                                                                                                                                                                       |
| `base_branch`               | No       | `master`                         | Name of a branch. The pipeline will only trigger on pull requests against the specified branch.                                                                                                                                                                                             |
| `labels`                    | No       | `["bug", "enhancement"]`         | The labels on the PR. The pipeline will only trigger on pull requests having at least one of the specified labels.                                                                                                                                                                          |
| `states`                    | No       | `["OPEN", "MERGED"]`             | The PR states to select (`OPEN`, `MERGED` or `CLOSED`). The pipeline will only trigger on pull requests matching one of the specified states. Default is ["OPEN"].                                                                                                                          |

Notes:
- If any of `hosting_endpoint`, `v3_endpoint`, or `v4_endpoint` are set, all of them must be set.
- When using `required_review_approvals`, you may also want to enable GitHub's branch protection rules to [dismiss stale pull request approvals when new commits are pushed](https://help.github.com/en/articles/enabling-required-reviews-for-pull-requests).

### Single PR

By including the `source.number` parameter, this resource will track commits
for a single PR.

| Parameter                   | Required | Example                          | Description                                                                                                                                                                                                                                                                                 |
|-----------------------------|----------|----------------------------------|---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `repository`                | Yes      | `itsdalmo/test-repository`       | The repository to target.                                                                                                                                                                                                                                                                   |
| `number`                    | Yes      | `1234`                           | The PR number to track commits for.                                                                                                                                                                                                                                                         |
| `access_token`              | Yes      |                                  | A Github Access Token with repository access (required for setting status on commits). N.B. If you want github-pr-resource to work with a private repository. Set `repo:full` permissions on the access token you create on GitHub. If it is a public repository, `repo:status` is enough.  |
| `hosting_endpoint`          | No       | `https://github.com`             | Endpoint under which repositories are hosted. Specifically, the resource will pull from `hosting_endpoint`/`repository`.                                                                                                                                                                    |
| `v3_endpoint`               | No       | `https://api.github.com`         | Endpoint to use for the V3 Github API (Restful).                                                                                                                                                                                                                                            |
| `v4_endpoint`               | No       | `https://api.github.com/graphql` | Endpoint to use for the V4 Github API (Graphql).                                                                                                                                                                                                                                            |
| `paths`                     | No       | `["terraform/*/*.tf"]`           | Only produce new versions for commits that include changes to files that match one or more glob patterns or prefixes. Note: this differs from `source.paths` when listing PRs in that it applies on a commit-by-commit basis, whereas the former applies for the full PR.                   |
| `ignore_paths`              | No       | `[".ci/"]`                       | Inverse of the above. Pattern syntax is documented in [filepath.Match](https://golang.org/pkg/path/filepath/#Match), or a path prefix can be specified (e.g. `.ci/` will match everything in the `.ci` directory).                                                                          |
| `labels`                    | No       | `["bug", "enhancement"]`         | The labels on the PR. The pipeline will only trigger on pull requests having at least one of the specified labels.                                                                                                                                                                          |
| `disable_ci_skip`           | No       | `true`                           | Disable ability to skip builds with `[ci skip]` and `[skip ci]` in the commit message.                                                                                                                                                                                                      |
| `git_crypt_key`             | No       | `AEdJVENSWVBUS0VZAAAAA...`       | Base64 encoded git-crypt key. Setting this will unlock / decrypt the repository with git-crypt. To get the key simply execute `git-crypt export-key -- - | base64` in an encrypted repository.                                                                                             |
| `disable_git_lfs`           | No       | `true`                           | Disable Git LFS, skipping an attempt to convert pointers of files tracked into their corresponding objects when checked out into a working copy.                                                                                                                                           |
| `skip_ssl_verification`     | No       | `true`                           | Disable SSL/TLS certificate validation on API clients. Use with care!                                                                                                                                                                                                                       |

Notes:
- If any of `hosting_endpoint`, `v3_endpoint`, or `v4_endpoint` are set, all of them must be set.

## Behaviour

As noted earlier, this resource can either track a list of PRs to a repository,
or track commits to a single PR. The different modes of operation have
different behaviour.

### List of PRs

By omitting the `source.number` parameter, this resource will list PRs in a
repository.

#### `check`

Produces a version consisting of the list of all PRs that match the criteria
defined in the `source`, sorted by PR number, and the timestamp when that new
value was observed. A new version will only be emitted if the set of PRs has
changed (i.e. PRs were added/removed from the set).

#### `get`

Stores the list of PRs in the file `prs.json`, encoded as a list of JSON
objects . This file can then be loaded into the build's local var state via the
`load_var` step.

Refer to [#example] for a full example.

#### `put`

Unimplemented

### Single PR

By including the `source.number` parameter, this resource will act upon a
single PR.

#### `check`

Produces a version for each new commit made to a PR (under the criteria
specified in `source`).

#### `get`

| Parameter            | Required | Example  | Description                                                                        |
|----------------------|----------|----------|------------------------------------------------------------------------------------|
| `skip_download`      | No       | `true`   | Use with `get_params` in a `put` step to do nothing on the implicit get.           |
| `integration_tool`   | No       | `rebase` | The integration tool to use, `merge`, `rebase` or `checkout`. Defaults to `merge`. |
| `git_depth`          | No       | `1`      | Shallow clone the repository using the `--depth` Git option                        |
| `submodules`         | No       | `true  ` | Recursively clone git submodules. Defaults to false.                               |
| `list_changed_files` | No       | `true`   | Generate a list of changed files and save alongside metadata                       |
| `fetch_tags`         | No       | `true`   | Fetch tags from remote repository                                                  |

Clones the base (e.g. `master` branch) at the latest commit, and merges the pull request at the specified commit
into master. This ensures that we are both testing and setting status on the exact commit that was requested in
input. Because the base of the PR is not locked to a specific commit in versions emitted from `check`, a fresh
`get` will always use the latest commit in master and *report the SHA of said commit in the metadata*. Both the
requested version and the metadata emitted by `get` are available to your tasks as JSON:
- `.git/resource/version.json`
- `.git/resource/metadata.json`
- `.git/resource/changed_files` (if enabled by `list_changed_files`)

The information in `metadata.json` is also available as individual files in the `.git/resource` directory, e.g. the `base_sha`
is available as `.git/resource/base_sha`. For a complete list of available (individual) metadata files, please check the code
[here](https://github.com/aoldershaw/github-pr-resource/blob/master/pr/in.go#45).

When specifying `skip_download` the pull request volume mounted to subsequent tasks will be empty, which is a problem
when you set e.g. the pending status before running the actual tests. The workaround for this is to use an alias for
the `put` (see https://github.com/telia-oss/github-pr-resource/issues/32 for more details).
Example here:

```yaml
put: update-status # <-- Use an alias for the pull-request resource
resource: pull-request
params:
    path: pull-request
    status: pending
get_params: {skip_download: true}
```

git-crypt encrypted repositories will automatically be decrypted when the `git_crypt_key` is set in the source configuration.

Note that, should you retrigger a build in the hopes of testing the last commit to a PR against a newer version of
the base, Concourse will reuse the volume (i.e. not trigger a new `get`) if it still exists, which can produce
unexpected results (#5). As such, re-testing a PR against a newer version of the base is best done by *pushing an
empty commit to the PR*.

#### `put`

| Parameter                  | Required | Example                              | Description                                                                                                                                                   |
|----------------------------|----------|--------------------------------------|---------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `path`                     | Yes      | `pull-request`                       | The name given to the resource in a GET step.                                                                                                                 |
| `status`                   | No       | `success`                            | Set a status on a commit. One of `success`, `pending`, `failure` and `error`.                                                                                 |
| `base_context`             | No       | `concourse-ci`                       | Base context (prefix) used for the status context. Defaults to `concourse-ci`.                                                                                |
| `context`                  | No       | `unit-test`                          | A context to use for the status, which is prefixed by `base_context`. Defaults to `"status"`.                                                                 |
| `comment`                  | No       | `hello world!`                       | A comment to add to the pull request.                                                                                                                         |
| `comment_file`             | No       | `my-output/comment.txt`              | Path to file containing a comment to add to the pull request (e.g. output of `terraform plan`).                                                               |
| `target_url`               | No       | `$ATC_EXTERNAL_URL/builds/$BUILD_ID` | The target URL for the status, where users are sent when clicking details (defaults to the Concourse build page).                                             |
| `description`              | No       | `Concourse CI build failed`          | The description status on the specified pull request.                                                                                                         |
| `delete_previous_comments` | No       | `true`                               | Boolean. Previous comments made on the pull request by this resource will be deleted before making the new comment. Useful for removing outdated information. |
| `body_matcher_for_delete`  | No       | `Concourse CI`                       | Regex. When using `delete_previous_comments`, filter comments to delete by this regex. Useful when more than one pipeline/job comments on the same PR.        |
| `override_version_pr`      | No       | `Concourse CI`                       | String. The PR number of PR to update status and/or comment on. Must be set if `override_version_commit` is set. Useful when you want to update status/comment on a PR different than the one you have checked out. |
| `override_version_commit`  | No       | `Concourse CI`                       | String. The HEAD SHA of PR to update status and/or comment on. Must be set if `override_version_pr` is set. Useful when you want to update status/comment on a PR different than the one you have checked out. |

Note that `comment`, `comment_file`, `context,` and `target_url` will all expand environment variables, so in the examples above `$ATC_EXTERNAL_URL` will be replaced by the public URL of the Concourse ATCs.
See https://concourse-ci.org/implementing-resource-types.html#resource-metadata for more details about metadata that is available via environment variables.

## Example

Unlike the [original resource][original-resource], usage of `aoldershaw/github-pr-resource`
requires two pipeline templates.

1. There is one "parent" pipeline that track changes to the list of PRs and
   sets a group of "child" pipelines. It runs whenever the list of PRs changes
   (e.g. because a new PR is opened, or an existing PR is merged).
2. There is one "child" pipeline per active PR. Each "child" pipeline tracks
   new commits to the corresponding PR. These pipelines will be archived when
   they are removed from the list of PRs (e.g. when the PR is merged).

For this example, assume you have a resource named `ci`, a repo which contains
the following pipeline files:

`ci/pipelines/parent.yml`
```yaml
resource_types:
- name: pull-request
  type: registry-image
  source:
    repository: aoldershaw/github-pr-resource

resources:
- name: pull-requests
  type: pull-request
  source:
    repository: itsdalmo/test-repository
    access_token: ((github-access-token))

- name: ci
  type: git
  source:
    uri: https://github.com/concourse/ci

jobs:
- name: update-pr-pipelines
  plan:
  - get: ci
  - get: pull-requests
    trigger: true
  - load_var: pull_requests
    file: pull-requests/prs.json
  - across:
    - var: pr
      values: ((.:pull_requests))
    set_pipeline: prs
    file: ci/pipelines/child.yml
    instance_vars: {number: ((.:pr.number))}
```

`ci/pipelines/child.yml`
```yaml
resource_types:
- name: pull-request
  type: registry-image
  source:
    repository: aoldershaw/github-pr-resource

resources:
- name: pull-request
  type: pull-request
  source:
    repository: itsdalmo/test-repository
    access_token: ((github-access-token))
    number: ((number))

- name: test
  plan:
  - get: pull-request
    trigger: true
  - put: pull-request-status
    resource: pull-request
    params:
      path: pull-request
      status: pending
    get_params: {skip_download: true}
  - task: unit-test
    config:
      platform: linux
      image_resource:
        type: registry-image
        source: {repository: alpine/git, tag: latest}
      inputs:
        - name: pr
      run:
        path: /bin/sh
        args:
          - -xce
          - |
            cd pull-request
            git log --graph -n 10 --color --pretty=format:"%x1b[31m%h%x09%x1b[32m%d%x1b[0m%x20%s" > log.txt
            cat log.txt
  on_success:
    put: pull-request
    params:
      path: pr
      status: success
    get_params: {skip_download: true}
  on_failure:
    put: pull-request
    params:
      path: pr
      status: failure
    get_params: {skip_download: true}
```

## Costs

The Github API(s) have a rate limit of 5000 requests per hour (per user). For the V3 API this essentially
translates to 5000 requests, whereas for the V4 API (GraphQL)  the calculation is more involved:
https://developer.github.com/v4/guides/resource-limitations/#calculating-a-rate-limit-score-before-running-the-call

As noted earlier, this resource can either track a list of PRs to a repository,
or track commits to a single PR. The different modes of operation have
different costs (with respect to consuming the Github API rate limit).

### List of PRs

Ref the above, here are some examples of running `check` against large repositories and the cost of doing so:
- [concourse/concourse](https://github.com/concourse/concourse): 51 open pull requests at the time of testing. Cost 2.
- [torvalds/linux](https://github.com/torvalds/linux): 305 open pull requests. Cost 8.
- [kubernetes/kubernetes](https://github.com/kubernetes/kubernetes): 1072 open pull requests. Cost: 22.

For the other two operations the costing is a bit easier:
- `get`: 0 cost
- `put`: unimplemented

### Single PR

- `check`: 0 cost (checking is done through `git`, not through the Github API)
- `get`: Fixed cost of 1. Fetches the pull request at the given commit.
- `put`: Uses the V3 API and has a min cost of 1, +1 for each of `status`, `comment`, etc.
