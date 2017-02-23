// This file was generated by counterfeiter
package providerfakes

import (
	"net/http"
	"sync"

	"code.cloudfoundry.org/lager"
	"github.com/concourse/atc/auth/provider"
)

type FakeVerifier struct {
	VerifyStub        func(lager.Logger, *http.Client) (bool, error)
	verifyMutex       sync.RWMutex
	verifyArgsForCall []struct {
		arg1 lager.Logger
		arg2 *http.Client
	}
	verifyReturns struct {
		result1 bool
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeVerifier) Verify(arg1 lager.Logger, arg2 *http.Client) (bool, error) {
	fake.verifyMutex.Lock()
	fake.verifyArgsForCall = append(fake.verifyArgsForCall, struct {
		arg1 lager.Logger
		arg2 *http.Client
	}{arg1, arg2})
	fake.recordInvocation("Verify", []interface{}{arg1, arg2})
	fake.verifyMutex.Unlock()
	if fake.VerifyStub != nil {
		return fake.VerifyStub(arg1, arg2)
	}
	return fake.verifyReturns.result1, fake.verifyReturns.result2
}

func (fake *FakeVerifier) VerifyCallCount() int {
	fake.verifyMutex.RLock()
	defer fake.verifyMutex.RUnlock()
	return len(fake.verifyArgsForCall)
}

func (fake *FakeVerifier) VerifyArgsForCall(i int) (lager.Logger, *http.Client) {
	fake.verifyMutex.RLock()
	defer fake.verifyMutex.RUnlock()
	return fake.verifyArgsForCall[i].arg1, fake.verifyArgsForCall[i].arg2
}

func (fake *FakeVerifier) VerifyReturns(result1 bool, result2 error) {
	fake.VerifyStub = nil
	fake.verifyReturns = struct {
		result1 bool
		result2 error
	}{result1, result2}
}

func (fake *FakeVerifier) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.verifyMutex.RLock()
	defer fake.verifyMutex.RUnlock()
	return fake.invocations
}

func (fake *FakeVerifier) recordInvocation(key string, args []interface{}) {
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

var _ provider.Verifier = new(FakeVerifier)
