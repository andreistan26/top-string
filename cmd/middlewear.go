package cmd

import (
	"fmt"
	"log"

	"github.com/andreistan26/top-string/internal/core"
	pqueue "github.com/andreistan26/top-string/internal/priority_queue"
	"github.com/andreistan26/top-string/internal/local"
	"github.com/spf13/cobra"
)

func ExecuteLocal(cmd *cobra.Command, args []string, opts *local.LocalOpts) error {
    opts.Paths = args

    for _, path := range opts.Paths {
        log.Print(path)       
    }

    hashes := core.GetHashStream(&opts.SenderOpts)

    var res []pqueue.FileHash
    go func() {
        res = core.CountStrings(hashes, &opts.ReceiverOpts)
    }()

    for idx, item := range res {
        fmt.Printf("%d:%s - [%d]\n", idx, item.Value, -item.Priority)
    }

    return nil
}
