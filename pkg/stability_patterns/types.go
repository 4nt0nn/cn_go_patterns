package pkg_stp_types

import "context"

// Circuit type represents the signature of the function wrapped by either circuit breaker- or
// by the debounce_first/debounce_last pattern (example signature)
type Circuit func(context.Context) (string, error)

// Effector represents the function that interacts with the service using
// the retry pattern. (example signature)
type Effector func(context.Context) (string, error)

// SlowFunction represents the function that is slow to respond within
// the timeout pattern (example signature)
type SlowFunction func(string) (string, error)

// WithContext represents the wrapper of the SlowFunction which catches the
// case when a function lacks context as a paramter (example signature)
type WithContext func(context.Context, string) (string, error)
