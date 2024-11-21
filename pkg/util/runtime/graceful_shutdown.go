package runtime

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

func SetupGracefulShutdown(ctx context.Context, f func(sig os.Signal)) {
	go func() {
		c := make(chan os.Signal, 1)
		defer close(c)

		signal.Notify(c, syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM)

		sig := <-c
		f(sig)
	}()
}
