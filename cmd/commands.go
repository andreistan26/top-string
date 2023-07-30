package cmd

import (
	"fmt"
	"os"

	"github.com/andreistan26/top-string/internal/local"
	"github.com/spf13/cobra"
)

type LocationType int
 
const (
    LOCAL_DATA = iota
    REMOTE_DATA
    LOCAL_COUNTER 
    REMOTE_COUNTER
    DISTRIBUTED_COUNTER
)

func CreateMainCommand() (*cobra.Command) {
    command := &cobra.Command{
        Use:    `top-string`,
        Short:  `top-string will count your strings`,
    } 

    return command
}

func CreateStartCommand() *cobra.Command {
    countCommand := &cobra.Command {
        Use:    `start`,
        Short:  `start an action server/send/local`,
    }

    return countCommand
}

func CreateLocalCommand() (*cobra.Command, *local.LocalOptions) {
    opts := &local.LocalOptions{}
    localCommand := &cobra.Command {
        Use:    `local`,
        Short:  `start counting local strings`,
        Args:   argsValidatorLocal,
        RunE:   func(cmd *cobra.Command, args []string) error {
            ExecuteLocal(cmd, args, opts)
            return nil
        },
    }

    localCommand.Flags().IntVar(&opts.QueryCount, "top", 5, "Number of strings returned")

    return localCommand, opts
}

func CreateServerCommand() *cobra.Command {
    serverCommand := &cobra.Command {
        Use:    `server`,
        Short:  `start a server for counting strings`,
        RunE:   func(cmd *cobra.Command, args []string) error {
            //TODO Add server starter
            return nil
        },
    }

    return serverCommand
}

func CreateSenderCommand() *cobra.Command {
    senderCommand := &cobra.Command {
        Use:    `sender`,
        Short:  `start sending `,
        Args:   nil, // TODO Add sender args validator
        RunE:   func(cmd *cobra.Command, args []string) error {
            //TODO Add sender starter
            return nil
        },
    }

    return senderCommand
}

func argsValidatorLocal(cmd *cobra.Command, args []string) error {
    fmt.Print(args)
    if err := cobra.MinimumNArgs(1)(cmd, args); err != nil {
        return err
    }

    for _, path := range args {
        if _, err := os.Stat(path); err != nil {
            return err
        }
    }



    return nil
}
