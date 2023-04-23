ARG base_image=alpine
ARG builder_image=golang

FROM ${builder_image} as builder
WORKDIR /go/src/github.com/aoldershaw/github-pr-resource
RUN curl -sL https://taskfile.dev/install.sh | sh
ADD . /go/src/github.com/aoldershaw/github-pr-resource
RUN ./bin/task build

FROM ${base_image} as resource

RUN apk add --update --no-cache \
    git \
    git-lfs \
    openssh \
    git-crypt

COPY scripts/askpass.sh /usr/local/bin/askpass.sh

COPY --from=builder /go/src/github.com/aoldershaw/github-pr-resource/build /opt/resource

FROM resource
