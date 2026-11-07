package resources

import (
	"context"
	"ksef/internal/logging"
)

type destructor func(ctx context.Context, commandErr error) error

var destructors []destructor

// so basically the idea behind this is that we're creating instances of
// resources that the client uses in `resources` package. but we also want
// to have some convenient way of destructing them. and that's implemented
// here. so it's really simple - if you create a resource, just add the
// destructor here and remember to call `defer resources.Close(ctx)` - that's it.
// and before you mention that we could wire it up to PostRunE:
// https://github.com/spf13/cobra/issues/1893
func Close(ctx context.Context, commandErr error) {
	// iterate over destructors in reverse to reflect the order at which
	// resources were constructed
	for i := len(destructors) - 1; i >= 0; i-- {
		if err := destructors[i](ctx, commandErr); err != nil {
			logging.SeiLogger.Error("resource destructor failed", "err", err)
		}
	}
	destructors = nil
}
