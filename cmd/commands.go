package cmd

import (
	"fmt"
	"os"
    "net"

	"github.com/andreistan26/top-string/internal/core"
	"github.com/andreistan26/top-string/internal/remote"
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

func CreateLocalCommand() (*cobra.Command) {
    opts := core.SenderOpts{}
    localCommand := &cobra.Command {
        Use:    `local`,
        Short:  `start counting local strings`,
        Args:   fileArgsValidator,
        RunE:   func(cmd *cobra.Command, args []string) error {
            ExecuteLocal(cmd, args, opts)
            return nil
        },
    }

    localCommand.Flags().IntVar(&opts.QueryCount, "top", 5, "Number of strings returned")

    return localCommand
}

func CreateServerCommand() *cobra.Command {
    addressOpts := remote.ConnOpts{}
    serverCommand := &cobra.Command {
        Use:    `server`,
        Short:  `start a server for counting strings`,
        RunE:   func(cmd *cobra.Command, args []string) error {
            ExecuteServer(cmd, addressOpts)
            return nil
        },
    }

    serverCommand.Flags().IPVar(&addressOpts.Ip, "ip", net.IP("127.0.0.1"), "Server ip")
    serverCommand.Flags().IntVar(&addressOpts.Port, "port", 1069, "Server ip")

    return serverCommand
}

func CreateSenderCommand() *cobra.Command {
    senderOpts := core.SenderOpts{}
    addressOpts := remote.ConnOpts{}

    senderCommand := &cobra.Command {
        Use:    `sender`,
        Short:  `start sending `,
        Args:   fileArgsValidator, // TODO Add sender args validator
        RunE:   func(cmd *cobra.Command, args []string) error {
            ExecuteSender(cmd, args, senderOpts, addressOpts)
            return nil
        },
    }

    senderCommand.Flags().IntVar(&senderOpts.QueryCount, "top", 5, "Number of strings returned")
    senderCommand.Flags().IPVar(&addressOpts.Ip, "ip", net.IP("127.0.0.1"), "Server ip")
    senderCommand.Flags().IntVar(&addressOpts.Port, "port", 1069, "Server ip")
    

    return senderCommand
}

func fileArgsValidator(cmd *cobra.Command, args []string) error {
    fmt.Println(args)
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
