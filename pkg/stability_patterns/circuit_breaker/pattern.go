package circuit_breaker

import (
	"context"
	"errors"
	"sync"
	"time"

	pkg_stp_types "github.com/4nt0nn/cn_go_patterns/pkg/stability_patterns"
)

func Breaker(circuit pkg_stp_types.Circuit, failureThreshold uint) pkg_stp_types.Circuit {
	var consecutiveFailures int = 0
	var lastAttempt = time.Now()
	var m sync.RWMutex

	return func(ctx context.Context) (string, error) {
		m.RLock() // establish a "read lock"

		d := consecutiveFailures - int(failureThreshold)

		if d >= 0 {
			shouldRetryAt := lastAttempt.Add(time.Second * 2 << d)
			if !time.Now().After(shouldRetryAt) {
				m.RUnlock()
				return "", errors.New("service unreachable")
			}
		}

		m.RUnlock() // release read lock

		response, err := circuit(ctx) // issue request

		m.Lock() // lock arround shared resources
		defer m.Unlock()

		lastAttempt = time.Now() // reccord time of attempt

		if err != nil { // circuit returned an error,
			consecutiveFailures++ // so we count the faiulure
			return response, err  // and return
		}

		consecutiveFailures = 0 // reset failures counter

		return response, nil
	}
}
