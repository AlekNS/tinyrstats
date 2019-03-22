package cmds

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	// "github.com/alekns/tinymonitor/backend/internal/monitor"
	"github.com/spf13/cobra"
)

func serveCommand() *cobra.Command {
	const kresourcefromfile = "preload-from-file"

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

			// Preload resources from the file
			csvFilePath := cmd.Flag(kresourcefromfile).Value.String()
			if len(csvFilePath) > 0 {
				if _, err := os.Stat(csvFilePath); os.IsNotExist(err) {
					panic(err.Error())
				}
				// if err := helpers.ReadTasksFromCsvFile(csvFilePath); err != nil {
				// 	panic(err.Error())
				// }
			}

			// app.BootstrapAndServe(rootContext, cmd)

			<-rootContext.Done()
		},
	}

	flags := cmd.PersistentFlags()
	flags.String(kresourcefromfile, "", "")

	return cmd
}
