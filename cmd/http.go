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

	// root of srv
	srv := &cobra.Command{
		Use:   "http",
		Short: "http service",
		Run: func(c *cobra.Command, args []string) {
			c.Usage()
		},
	}

	// run
	run := &cobra.Command{
		Use:   "run",
		Short: "run http service",
		Run: func(c *cobra.Command, args []string) {

			// config
			conf, err := load.LoadFromFile(configPath)
			if err != nil {
				log.Logger.Error("failed to load config from file", zap.Error(err))
				return
			}

			// http handler
			h, err := httphandler.New(conf.HTTPSrv.Mode,
				httphandler.WithNasaAPIKey(conf.ThirdParty.Nasa.Key),
				httphandler.WithMainDatabase(conf.DB.Main.Type, conf.DB.Main.DSN, conf.DB.Main.MaxOpenConns, conf.DB.Main.MaxIdleConns, conf.DB.Main.MaxLifeTime),
			)
			if err != nil {
				log.Logger.Error("failed to create http handler", zap.Error(err))
				return
			}

			hp, err := h.HTTPHandler()
			if err != nil {
				log.Logger.Error("failed to create http service handler", zap.Error(err))
				return
			}

			// server
			server, err := httpserver.New(conf.HTTPSrv.Addr, hp)
			if err != nil {
				log.Logger.Error("failed to create http server", zap.Error(err))
				return
			}
			log.Logger.Info("server is staring...", zap.String("at", conf.HTTPSrv.Addr), zap.String("mode", conf.HTTPSrv.Mode))
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
