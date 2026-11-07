package resources

import (
	"context"
	"sync"

	"ksef/internal/logging"
)

type destructor func(ctx context.Context, commandErr error) error

var (
	// destructorsMu guards the registry separately from `mu`: resources are often
	// constructed while `mu` is held (see clientLocked/invoiceDBLocked), and those
	// constructors register their destructors here.
	destructorsMu sync.Mutex
	destructors   []destructor
)

// so basically the idea behind this is that we're creating instances of
// resources that the client uses in `resources` package. but we also want
// to have some convenient way of destructing them. and that's implemented
// here. so it's really simple - if you create a resource, just add its
// destructor here and remember to call `defer resources.Close(ctx)` - that's it.
// and before you mention that we could wire it up to PostRunE:
// https://github.com/spf13/cobra/issues/1893
func Close(ctx context.Context, commandErr error) {
	destructorsMu.Lock()
	pending := destructors
	destructors = nil
	destructorsMu.Unlock()

	// iterate over destructors in reverse to reflect the order at which
	// resources were constructed
	for i := len(pending) - 1; i >= 0; i-- {
		if err := pending[i](ctx, commandErr); err != nil {
			logging.SeiLogger.Error("resource destructor failed", "err", err)
		}
	}
}

// OnClose registers a destructor for a resource. The constructor of a resource is
// responsible for registering it (i.e. `RateLimitsState` registers the storage's
// Close, `invoiceDBLocked` registers the DB's Close), so commands never have to
// remember what to clean up - they only `defer resources.Close(ctx, err)`.
//
// Registration is mutex-guarded because parallel downloads construct resources
// from several goroutines.
func OnClose(fn func(ctx context.Context, commandErr error) error) {
	destructorsMu.Lock()
	defer destructorsMu.Unlock()

	destructors = append(destructors, fn)
}
