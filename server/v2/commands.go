package serverv2

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"cosmossdk.io/log"
)

func Commands(logger log.Logger, homePath string, components ...ServerComponent) (CLIConfig, error) {
	if len(components) == 0 {
		// TODO figure if we should define default components
		// and if so it should be done here to avoid uncessary dependencies
		return CLIConfig{}, errors.New("no modules provided")
	}

	server := NewServer(logger, components...)
	flags := server.StartFlags()

	if _, err := os.Stat(homePath); os.IsNotExist(err) {
		_ = server.WriteConfig(homePath)
	}

	startCmd := &cobra.Command{
		Use:                "start",
		Short:              "Run the application",
		DisableFlagParsing: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			// get viper from context
			flags := *pflag.FlagSet{} // magic flag parsing
			// bind flags to viper

			srvConfig := Config{StartBlock: true}
			ctx := cmd.Context()
			ctx = context.WithValue(ctx, ServerContextKey, srvConfig)
			ctx, cancelFn := context.WithCancel(ctx)
			go func() {
				sigCh := make(chan os.Signal, 1)
				signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
				sig := <-sigCh
				cancelFn()
				cmd.Printf("caught %s signal\n", sig.String())

				if err := server.Stop(ctx); err != nil {
					cmd.PrintErrln("failed to stop servers:", err)
				}
			}()

			if err := server.Start(ctx); err != nil {
				return fmt.Errorf("failed to start servers: %w", err)
			}

			return nil
		},
	}

	cmds := server.CLICommands()
	cmds.Commands = append(cmds.Commands, startCmd)

	return cmds, nil
}

func AddCommands(rootCmd *cobra.Command, logger log.Logger, homePath string, components ...ServerComponent) error {
	cmds, err := Commands(logger, homePath, components...)
	if err != nil {
		return err
	}

	rootCmd.AddCommand(cmds.Commands...)
	return nil
}
