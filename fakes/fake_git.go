// Code generated by counterfeiter. DO NOT EDIT.
package fakes

import (
	"sync"

	resource "github.com/aoldershaw/github-pr-resource"
)

type FakeGit struct {
	CheckoutStub        func(string, string, bool) error
	checkoutMutex       sync.RWMutex
	checkoutArgsForCall []struct {
		arg1 string
		arg2 string
		arg3 bool
	}
	checkoutReturns struct {
		result1 error
	}
	checkoutReturnsOnCall map[int]struct {
		result1 error
	}
	FetchStub        func(string, int, int, bool) error
	fetchMutex       sync.RWMutex
	fetchArgsForCall []struct {
		arg1 string
		arg2 int
		arg3 int
		arg4 bool
	}
	fetchReturns struct {
		result1 error
	}
	fetchReturnsOnCall map[int]struct {
		result1 error
	}
	GitCryptUnlockStub        func(string) error
	gitCryptUnlockMutex       sync.RWMutex
	gitCryptUnlockArgsForCall []struct {
		arg1 string
	}
	gitCryptUnlockReturns struct {
		result1 error
	}
	gitCryptUnlockReturnsOnCall map[int]struct {
		result1 error
	}
	InitStub        func(string) error
	initMutex       sync.RWMutex
	initArgsForCall []struct {
		arg1 string
	}
	initReturns struct {
		result1 error
	}
	initReturnsOnCall map[int]struct {
		result1 error
	}
	MergeStub        func(string, bool) error
	mergeMutex       sync.RWMutex
	mergeArgsForCall []struct {
		arg1 string
		arg2 bool
	}
	mergeReturns struct {
		result1 error
	}
	mergeReturnsOnCall map[int]struct {
		result1 error
	}
	PullStub        func(string, string, int, bool, bool) error
	pullMutex       sync.RWMutex
	pullArgsForCall []struct {
		arg1 string
		arg2 string
		arg3 int
		arg4 bool
		arg5 bool
	}
	pullReturns struct {
		result1 error
	}
	pullReturnsOnCall map[int]struct {
		result1 error
	}
	RebaseStub        func(string, string, bool) error
	rebaseMutex       sync.RWMutex
	rebaseArgsForCall []struct {
		arg1 string
		arg2 string
		arg3 bool
	}
	rebaseReturns struct {
		result1 error
	}
	rebaseReturnsOnCall map[int]struct {
		result1 error
	}
	RevParseStub        func(string) (string, error)
	revParseMutex       sync.RWMutex
	revParseArgsForCall []struct {
		arg1 string
	}
	revParseReturns struct {
		result1 string
		result2 error
	}
	revParseReturnsOnCall map[int]struct {
		result1 string
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeGit) Checkout(arg1 string, arg2 string, arg3 bool) error {
	fake.checkoutMutex.Lock()
	ret, specificReturn := fake.checkoutReturnsOnCall[len(fake.checkoutArgsForCall)]
	fake.checkoutArgsForCall = append(fake.checkoutArgsForCall, struct {
		arg1 string
		arg2 string
		arg3 bool
	}{arg1, arg2, arg3})
	fake.recordInvocation("Checkout", []interface{}{arg1, arg2, arg3})
	fake.checkoutMutex.Unlock()
	if fake.CheckoutStub != nil {
		return fake.CheckoutStub(arg1, arg2, arg3)
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.checkoutReturns
	return fakeReturns.result1
}

func (fake *FakeGit) CheckoutCallCount() int {
	fake.checkoutMutex.RLock()
	defer fake.checkoutMutex.RUnlock()
	return len(fake.checkoutArgsForCall)
}

func (fake *FakeGit) CheckoutCalls(stub func(string, string, bool) error) {
	fake.checkoutMutex.Lock()
	defer fake.checkoutMutex.Unlock()
	fake.CheckoutStub = stub
}

func (fake *FakeGit) CheckoutArgsForCall(i int) (string, string, bool) {
	fake.checkoutMutex.RLock()
	defer fake.checkoutMutex.RUnlock()
	argsForCall := fake.checkoutArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *FakeGit) CheckoutReturns(result1 error) {
	fake.checkoutMutex.Lock()
	defer fake.checkoutMutex.Unlock()
	fake.CheckoutStub = nil
	fake.checkoutReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeGit) CheckoutReturnsOnCall(i int, result1 error) {
	fake.checkoutMutex.Lock()
	defer fake.checkoutMutex.Unlock()
	fake.CheckoutStub = nil
	if fake.checkoutReturnsOnCall == nil {
		fake.checkoutReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.checkoutReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeGit) Fetch(arg1 string, arg2 int, arg3 int, arg4 bool) error {
	fake.fetchMutex.Lock()
	ret, specificReturn := fake.fetchReturnsOnCall[len(fake.fetchArgsForCall)]
	fake.fetchArgsForCall = append(fake.fetchArgsForCall, struct {
		arg1 string
		arg2 int
		arg3 int
		arg4 bool
	}{arg1, arg2, arg3, arg4})
	fake.recordInvocation("Fetch", []interface{}{arg1, arg2, arg3, arg4})
	fake.fetchMutex.Unlock()
	if fake.FetchStub != nil {
		return fake.FetchStub(arg1, arg2, arg3, arg4)
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.fetchReturns
	return fakeReturns.result1
}

func (fake *FakeGit) FetchCallCount() int {
	fake.fetchMutex.RLock()
	defer fake.fetchMutex.RUnlock()
	return len(fake.fetchArgsForCall)
}

func (fake *FakeGit) FetchCalls(stub func(string, int, int, bool) error) {
	fake.fetchMutex.Lock()
	defer fake.fetchMutex.Unlock()
	fake.FetchStub = stub
}

func (fake *FakeGit) FetchArgsForCall(i int) (string, int, int, bool) {
	fake.fetchMutex.RLock()
	defer fake.fetchMutex.RUnlock()
	argsForCall := fake.fetchArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3, argsForCall.arg4
}

func (fake *FakeGit) FetchReturns(result1 error) {
	fake.fetchMutex.Lock()
	defer fake.fetchMutex.Unlock()
	fake.FetchStub = nil
	fake.fetchReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeGit) FetchReturnsOnCall(i int, result1 error) {
	fake.fetchMutex.Lock()
	defer fake.fetchMutex.Unlock()
	fake.FetchStub = nil
	if fake.fetchReturnsOnCall == nil {
		fake.fetchReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.fetchReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeGit) GitCryptUnlock(arg1 string) error {
	fake.gitCryptUnlockMutex.Lock()
	ret, specificReturn := fake.gitCryptUnlockReturnsOnCall[len(fake.gitCryptUnlockArgsForCall)]
	fake.gitCryptUnlockArgsForCall = append(fake.gitCryptUnlockArgsForCall, struct {
		arg1 string
	}{arg1})
	fake.recordInvocation("GitCryptUnlock", []interface{}{arg1})
	fake.gitCryptUnlockMutex.Unlock()
	if fake.GitCryptUnlockStub != nil {
		return fake.GitCryptUnlockStub(arg1)
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.gitCryptUnlockReturns
	return fakeReturns.result1
}

func (fake *FakeGit) GitCryptUnlockCallCount() int {
	fake.gitCryptUnlockMutex.RLock()
	defer fake.gitCryptUnlockMutex.RUnlock()
	return len(fake.gitCryptUnlockArgsForCall)
}

func (fake *FakeGit) GitCryptUnlockCalls(stub func(string) error) {
	fake.gitCryptUnlockMutex.Lock()
	defer fake.gitCryptUnlockMutex.Unlock()
	fake.GitCryptUnlockStub = stub
}

func (fake *FakeGit) GitCryptUnlockArgsForCall(i int) string {
	fake.gitCryptUnlockMutex.RLock()
	defer fake.gitCryptUnlockMutex.RUnlock()
	argsForCall := fake.gitCryptUnlockArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeGit) GitCryptUnlockReturns(result1 error) {
	fake.gitCryptUnlockMutex.Lock()
	defer fake.gitCryptUnlockMutex.Unlock()
	fake.GitCryptUnlockStub = nil
	fake.gitCryptUnlockReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeGit) GitCryptUnlockReturnsOnCall(i int, result1 error) {
	fake.gitCryptUnlockMutex.Lock()
	defer fake.gitCryptUnlockMutex.Unlock()
	fake.GitCryptUnlockStub = nil
	if fake.gitCryptUnlockReturnsOnCall == nil {
		fake.gitCryptUnlockReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.gitCryptUnlockReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeGit) Init(arg1 string) error {
	fake.initMutex.Lock()
	ret, specificReturn := fake.initReturnsOnCall[len(fake.initArgsForCall)]
	fake.initArgsForCall = append(fake.initArgsForCall, struct {
		arg1 string
	}{arg1})
	fake.recordInvocation("Init", []interface{}{arg1})
	fake.initMutex.Unlock()
	if fake.InitStub != nil {
		return fake.InitStub(arg1)
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.initReturns
	return fakeReturns.result1
}

func (fake *FakeGit) InitCallCount() int {
	fake.initMutex.RLock()
	defer fake.initMutex.RUnlock()
	return len(fake.initArgsForCall)
}

func (fake *FakeGit) InitCalls(stub func(string) error) {
	fake.initMutex.Lock()
	defer fake.initMutex.Unlock()
	fake.InitStub = stub
}

func (fake *FakeGit) InitArgsForCall(i int) string {
	fake.initMutex.RLock()
	defer fake.initMutex.RUnlock()
	argsForCall := fake.initArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeGit) InitReturns(result1 error) {
	fake.initMutex.Lock()
	defer fake.initMutex.Unlock()
	fake.InitStub = nil
	fake.initReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeGit) InitReturnsOnCall(i int, result1 error) {
	fake.initMutex.Lock()
	defer fake.initMutex.Unlock()
	fake.InitStub = nil
	if fake.initReturnsOnCall == nil {
		fake.initReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.initReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeGit) Merge(arg1 string, arg2 bool) error {
	fake.mergeMutex.Lock()
	ret, specificReturn := fake.mergeReturnsOnCall[len(fake.mergeArgsForCall)]
	fake.mergeArgsForCall = append(fake.mergeArgsForCall, struct {
		arg1 string
		arg2 bool
	}{arg1, arg2})
	fake.recordInvocation("Merge", []interface{}{arg1, arg2})
	fake.mergeMutex.Unlock()
	if fake.MergeStub != nil {
		return fake.MergeStub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.mergeReturns
	return fakeReturns.result1
}

func (fake *FakeGit) MergeCallCount() int {
	fake.mergeMutex.RLock()
	defer fake.mergeMutex.RUnlock()
	return len(fake.mergeArgsForCall)
}

func (fake *FakeGit) MergeCalls(stub func(string, bool) error) {
	fake.mergeMutex.Lock()
	defer fake.mergeMutex.Unlock()
	fake.MergeStub = stub
}

func (fake *FakeGit) MergeArgsForCall(i int) (string, bool) {
	fake.mergeMutex.RLock()
	defer fake.mergeMutex.RUnlock()
	argsForCall := fake.mergeArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *FakeGit) MergeReturns(result1 error) {
	fake.mergeMutex.Lock()
	defer fake.mergeMutex.Unlock()
	fake.MergeStub = nil
	fake.mergeReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeGit) MergeReturnsOnCall(i int, result1 error) {
	fake.mergeMutex.Lock()
	defer fake.mergeMutex.Unlock()
	fake.MergeStub = nil
	if fake.mergeReturnsOnCall == nil {
		fake.mergeReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.mergeReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeGit) Pull(arg1 string, arg2 string, arg3 int, arg4 bool, arg5 bool) error {
	fake.pullMutex.Lock()
	ret, specificReturn := fake.pullReturnsOnCall[len(fake.pullArgsForCall)]
	fake.pullArgsForCall = append(fake.pullArgsForCall, struct {
		arg1 string
		arg2 string
		arg3 int
		arg4 bool
		arg5 bool
	}{arg1, arg2, arg3, arg4, arg5})
	fake.recordInvocation("Pull", []interface{}{arg1, arg2, arg3, arg4, arg5})
	fake.pullMutex.Unlock()
	if fake.PullStub != nil {
		return fake.PullStub(arg1, arg2, arg3, arg4, arg5)
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.pullReturns
	return fakeReturns.result1
}

func (fake *FakeGit) PullCallCount() int {
	fake.pullMutex.RLock()
	defer fake.pullMutex.RUnlock()
	return len(fake.pullArgsForCall)
}

func (fake *FakeGit) PullCalls(stub func(string, string, int, bool, bool) error) {
	fake.pullMutex.Lock()
	defer fake.pullMutex.Unlock()
	fake.PullStub = stub
}

func (fake *FakeGit) PullArgsForCall(i int) (string, string, int, bool, bool) {
	fake.pullMutex.RLock()
	defer fake.pullMutex.RUnlock()
	argsForCall := fake.pullArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3, argsForCall.arg4, argsForCall.arg5
}

func (fake *FakeGit) PullReturns(result1 error) {
	fake.pullMutex.Lock()
	defer fake.pullMutex.Unlock()
	fake.PullStub = nil
	fake.pullReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeGit) PullReturnsOnCall(i int, result1 error) {
	fake.pullMutex.Lock()
	defer fake.pullMutex.Unlock()
	fake.PullStub = nil
	if fake.pullReturnsOnCall == nil {
		fake.pullReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.pullReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeGit) Rebase(arg1 string, arg2 string, arg3 bool) error {
	fake.rebaseMutex.Lock()
	ret, specificReturn := fake.rebaseReturnsOnCall[len(fake.rebaseArgsForCall)]
	fake.rebaseArgsForCall = append(fake.rebaseArgsForCall, struct {
		arg1 string
		arg2 string
		arg3 bool
	}{arg1, arg2, arg3})
	fake.recordInvocation("Rebase", []interface{}{arg1, arg2, arg3})
	fake.rebaseMutex.Unlock()
	if fake.RebaseStub != nil {
		return fake.RebaseStub(arg1, arg2, arg3)
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.rebaseReturns
	return fakeReturns.result1
}

func (fake *FakeGit) RebaseCallCount() int {
	fake.rebaseMutex.RLock()
	defer fake.rebaseMutex.RUnlock()
	return len(fake.rebaseArgsForCall)
}

func (fake *FakeGit) RebaseCalls(stub func(string, string, bool) error) {
	fake.rebaseMutex.Lock()
	defer fake.rebaseMutex.Unlock()
	fake.RebaseStub = stub
}

func (fake *FakeGit) RebaseArgsForCall(i int) (string, string, bool) {
	fake.rebaseMutex.RLock()
	defer fake.rebaseMutex.RUnlock()
	argsForCall := fake.rebaseArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *FakeGit) RebaseReturns(result1 error) {
	fake.rebaseMutex.Lock()
	defer fake.rebaseMutex.Unlock()
	fake.RebaseStub = nil
	fake.rebaseReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeGit) RebaseReturnsOnCall(i int, result1 error) {
	fake.rebaseMutex.Lock()
	defer fake.rebaseMutex.Unlock()
	fake.RebaseStub = nil
	if fake.rebaseReturnsOnCall == nil {
		fake.rebaseReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.rebaseReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeGit) RevParse(arg1 string) (string, error) {
	fake.revParseMutex.Lock()
	ret, specificReturn := fake.revParseReturnsOnCall[len(fake.revParseArgsForCall)]
	fake.revParseArgsForCall = append(fake.revParseArgsForCall, struct {
		arg1 string
	}{arg1})
	fake.recordInvocation("RevParse", []interface{}{arg1})
	fake.revParseMutex.Unlock()
	if fake.RevParseStub != nil {
		return fake.RevParseStub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.revParseReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeGit) RevParseCallCount() int {
	fake.revParseMutex.RLock()
	defer fake.revParseMutex.RUnlock()
	return len(fake.revParseArgsForCall)
}

func (fake *FakeGit) RevParseCalls(stub func(string) (string, error)) {
	fake.revParseMutex.Lock()
	defer fake.revParseMutex.Unlock()
	fake.RevParseStub = stub
}

func (fake *FakeGit) RevParseArgsForCall(i int) string {
	fake.revParseMutex.RLock()
	defer fake.revParseMutex.RUnlock()
	argsForCall := fake.revParseArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeGit) RevParseReturns(result1 string, result2 error) {
	fake.revParseMutex.Lock()
	defer fake.revParseMutex.Unlock()
	fake.RevParseStub = nil
	fake.revParseReturns = struct {
		result1 string
		result2 error
	}{result1, result2}
}

func (fake *FakeGit) RevParseReturnsOnCall(i int, result1 string, result2 error) {
	fake.revParseMutex.Lock()
	defer fake.revParseMutex.Unlock()
	fake.RevParseStub = nil
	if fake.revParseReturnsOnCall == nil {
		fake.revParseReturnsOnCall = make(map[int]struct {
			result1 string
			result2 error
		})
	}
	fake.revParseReturnsOnCall[i] = struct {
		result1 string
		result2 error
	}{result1, result2}
}

func (fake *FakeGit) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.checkoutMutex.RLock()
	defer fake.checkoutMutex.RUnlock()
	fake.fetchMutex.RLock()
	defer fake.fetchMutex.RUnlock()
	fake.gitCryptUnlockMutex.RLock()
	defer fake.gitCryptUnlockMutex.RUnlock()
	fake.initMutex.RLock()
	defer fake.initMutex.RUnlock()
	fake.mergeMutex.RLock()
	defer fake.mergeMutex.RUnlock()
	fake.pullMutex.RLock()
	defer fake.pullMutex.RUnlock()
	fake.rebaseMutex.RLock()
	defer fake.rebaseMutex.RUnlock()
	fake.revParseMutex.RLock()
	defer fake.revParseMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeGit) recordInvocation(key string, args []interface{}) {
	fake.invocationsMutex.Lock()
	defer fake.invocationsMutex.Unlock()
	if fake.invocations == nil {
		fake.invocations = map[string][][]interface{}{}
	}
	if fake.invocations[key] == nil {
		fake.invocations[key] = [][]interface{}{}
	}
	fake.invocations[key] = append(fake.invocations[key], args)
}

var _ resource.Git = new(FakeGit)
