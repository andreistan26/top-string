package cmd

import (
	"log"

	"github.com/andreistan26/top-string/internal/local"
	"github.com/spf13/cobra"
)

func ExecuteLocal(cmd *cobra.Command, args []string, opts *local.LocalOptions) error {
    opts.Paths = args

    for _, path := range opts.Paths {
        log.Print(path)       
    }

    local.ReceiveHashes(opts)
    return nil
}
