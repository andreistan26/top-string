package cmd

import (
	"fmt"
	"log"

	"github.com/andreistan26/top-string/internal/core"
	pqueue "github.com/andreistan26/top-string/internal/priority_queue"
	"github.com/andreistan26/top-string/internal/remote"
	"github.com/spf13/cobra"
)

func ExecuteLocal(cmd *cobra.Command, args []string, opts core.SenderOpts) error {
    opts.Paths = args

    for _, path := range opts.Paths {
        log.Print(path)       
    }

    hashes := core.GetHashStream(opts)

    res := []pqueue.FileHash{}
    go func() {
        res = core.CountStrings(hashes, opts.QueryCount)
    }()

    for idx, item := range res {
        fmt.Printf("%d:%s - [%d]\n", idx, item.Value, -item.Priority)
    }

    return nil
}

func ExecuteSender(cmd *cobra.Command, args []string, senderOpts core.SenderOpts, connOpts remote.ConnOpts) error {
    senderOpts.Paths = args
    for _, path := range senderOpts.Paths {
        log.Println(path)       
    }

    remote.SendFiles(senderOpts, connOpts)

    return nil
}

func ExecuteServer(cmd *cobra.Command, connOpts remote.ConnOpts) error {
    server, err := remote.StartServer(&connOpts)
    if err != nil {
        panic(err)
    }

    server.Run()

    return nil
}
