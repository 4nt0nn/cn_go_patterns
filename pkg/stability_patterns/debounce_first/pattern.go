package debounce_first

import (
	"context"
	"sync"
	"time"

	pkg_stp_types "github.com/4nt0nn/cn_go_patterns/pkg/stability_patterns"
)

func DebounceFirst(circuit pkg_stp_types.Circuit, d time.Duration) pkg_stp_types.Circuit {
	var threshold time.Time
	var result string
	var err error
	var m sync.Mutex

	return func(ctx context.Context) (string, error) {
		m.Lock()

		defer func() {
			threshold = time.Now().Add(d)
			m.Unlock()
		}()

		if time.Now().Before(threshold) {
			return result, err
		}

		result, err = circuit(ctx)

		return result, err
	}
}
