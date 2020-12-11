//+ build js,wasm

package main

import (
	"context"
	"fmt"
	"syscall/js"

	// "github.com/wabarc/wayback"
	"github.com/wabarc/wayback/config"
	"github.com/wabarc/wayback/logger"
	"github.com/wabarc/wayback/service/anonymity"
	"github.com/wabarc/wayback/service/telegram"
)

func main() {
	js.Global().Set("handle", js.FuncOf(handle))
	<-make(chan interface{})
}

func handle(_ js.Value, input []js.Value) interface{} {
	if !config.Opts.LogTime() {
		logger.DisableTime()
	}

	if config.Opts.HasDebugMode() {
		logger.EnableDebug()
	}

	logger.Info("%v", "Hello, WASM.")
	// TODO
	// callback := input[1]
	// callback.Invoke("Hello")
	srv := &service{daemon: []string{"telegram"}}
	if len(srv.daemon) > 0 {
		srv.serve(config.Opts)
	}

	return nil
}

type service struct {
	errCh  chan error
	daemon []string
}

func (srv *service) serve(opts *config.Options) {
	ctx := context.Background()
	ran := srv.run(ctx, opts)

	select {
	case err := <-ran.err():
		logger.Error(err.Error())
	case <-ctx.Done():
	}
}

func (srv *service) run(ctx context.Context, opts *config.Options) *service {
	telegram := telegram.New(opts)
	tor := anonymity.New(opts)

	srv.errCh = make(chan error, len(srv.daemon))
	for _, s := range srv.daemon {
		switch s {
		case "telegram":
			go func(errCh chan error) {
				errCh <- telegram.Serve(ctx)
			}(srv.errCh)
		case "web":
			go func(errCh chan error) {
				errCh <- tor.Serve(ctx)
			}(srv.errCh)
		default:
			fmt.Printf("Unrecognize %s in `--daemon`\n", s)
			srv.errCh <- ctx.Err()
		}
	}

	return srv
}

func (s *service) err() <-chan error { return s.errCh }
