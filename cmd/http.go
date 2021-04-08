package cmd

import (
	httphandler "simpleBackend/handlers/http"
	"simpleBackend/log"
	httpserver "simpleBackend/servers/http"

	"simpleBackend/config/load"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

func httpservice() *cobra.Command {

	var configPath string
	srv := &cobra.Command{
		Use:   "http",
		Short: "http service",
		Run: func(c *cobra.Command, args []string) {
			c.Usage()
		},
	}

	run := &cobra.Command{
		Use:   "run",
		Short: "run http service",
		Run: func(c *cobra.Command, args []string) {
			log.Logger.Info("http server is runing...")

			conf, err := load.LoadFromFile(configPath)

			if err != nil {
				log.Logger.Error("failed to load config from file", zap.Error(err))
				return
			}

			log.Logger.Info("config", zap.Any("config", conf))

			h, err := httphandler.New("test")
			if err != nil {
				log.Logger.Error("failed to create http handler", zap.Error(err))
				return
			}

			hp, err := h.HTTPHandler()
			if err != nil {
				log.Logger.Error("failed to create http service handler", zap.Error(err))
				return
			}

			server, err := httpserver.New("127.0.0.1:3000", hp)
			if err != nil {
				log.Logger.Error("failed to create http server", zap.Error(err))
				return
			}

			if err := server.Run(); err != nil {
				log.Logger.Error("failed to run http server", zap.Error(err))
				return
			}
			log.Logger.Info("http server is existing...")
		},
	}

	run.Flags().StringVarP(&configPath, "config", "c", "config/example.yaml", "path of config file with yaml format")

	srv.AddCommand(
		run,
	)

	return srv
}
