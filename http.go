package hufx

import (
	"context"
	"errors"
	"net"
	"net/http"

	"github.com/rs/zerolog"
	"go.uber.org/fx"
)

func HTTPServerStartStopHook(srv *http.Server, log *zerolog.Logger) fx.Hook {
	return fx.Hook{
		OnStart: func(ctx context.Context) error {
			ln, err := net.Listen("tcp", srv.Addr)
			if err != nil {
				log.Error().Err(err).Str("listen_addr", srv.Addr).Msg("HTTP Server start: listen error!")
				return err
			}
			go func() {
				if err = srv.Serve(ln); err != nil && !errors.Is(err, http.ErrServerClosed) {
					log.Error().Err(err).Msg("HTTP Server start: serve error!")
				}
			}()
			return nil
		},

		OnStop: func(ctx context.Context) (err error) {
			if err = srv.Shutdown(ctx); err != nil {
				log.Error().Err(err).Msg("HTTP Server stop: error!")
			}
			return
		},
	}
}
