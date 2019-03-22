package cmds

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/alekns/tinyrstats/internal/config"
	"github.com/alekns/tinyrstats/internal/helpers/tasks"
	"github.com/alekns/tinyrstats/internal/helpers/tracer"
	"github.com/alekns/tinyrstats/internal/monitor/app"
	"github.com/alekns/tinyrstats/internal/presentation/external"
	cmdhelper "github.com/alekns/tinyrstats/pkg/helpers/cmd"
	monlog "github.com/alekns/tinyrstats/pkg/logger"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	stdopentracing "github.com/opentracing/opentracing-go"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func serveCommand() *cobra.Command {
	const (
		kresourcefromfile = "preload-from-file"
		kdefaultprotocol  = "default-protocol"
	)

	cmd := &cobra.Command{
		Use:   "serve",
		Short: "Serve resources",

		Run: func(cmd *cobra.Command, args []string) {
			//
			// Process signals
			//
			rootContext, cancelByContext := context.WithCancel(context.Background())
			defer cancelByContext()

			appSignals := make(chan os.Signal, 0)
			signal.Notify(appSignals, syscall.SIGINT, syscall.SIGTERM)
			go func() {
				// sig :=
				<-appSignals
				// log.Info("stop process. Catch signal", sig)
				cancelByContext()
			}()

			//
			// Prepare settings and Configure logger
			//
			settings := config.GetSettings(
				cmdhelper.ConfigureViper("tinyrstats", "config", cmd, viper.GetViper()))

			var logger log.Logger = log.NewLogfmtLogger(os.Stdout)
			logger = monlog.SetLevelLogger(logger, settings.Logger.ConsoleLevel)

			//
			// Init tracer (if configured through env)
			//
			tracer, close, err := tracer.NewJaeger(logger, "tinyrstats")
			if err != nil {
				panic(fmt.Sprintf("tracer failed to initialize: %s", err.Error()))
			}
			defer close.Close()
			stdopentracing.SetGlobalTracer(tracer)

			//
			// Create application registry
			//
			registry := app.NewRegistryImpl(settings, logger)
			if err := registry.Start(rootContext); err != nil {
				panic(err.Error())
			}

			//
			// Preload resources from the file (optional)
			//
			csvFilePath := cmd.Flag(kresourcefromfile).Value.String()
			if len(csvFilePath) > 0 {
				level.Info(logger).Log("msg", "preload resources from csv file", "file", csvFilePath)
				if _, err := os.Stat(csvFilePath); os.IsNotExist(err) {
					panic(err.Error())
				}

				defaultProtocol := cmd.Flag(kdefaultprotocol).Value.String()
				if defaultProtocol != "http" && defaultProtocol != "https" {
					panic(kdefaultprotocol + " has invalid protocol value, only http or https are allowed")
				}

				if scheduleTasks, err := tasks.ReadTasksFromCsvFile(defaultProtocol, csvFilePath); err != nil {
					panic(err.Error())
				} else {
					tasks.ScheduleTasksSlice(rootContext, registry.TaskApp(), scheduleTasks)
				}
			}

			//
			// Start external presentation
			//
			level.Info(logger).Log("msg", "start and serve")
			external.BootstrapAndServe(rootContext, settings, cmd, logger, registry)

			registry.Stop(rootContext)
		},
	}

	flags := cmd.PersistentFlags()
	flags.String(kresourcefromfile, "", "")
	flags.String(kdefaultprotocol, "http", "")

	return cmd
}
