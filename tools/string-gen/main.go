package main

import (
	"github.com/andreistan26/top-string/tools/string-gen/pkg/generator"
	"github.com/spf13/cobra"
    "errors"
    "fmt"
)


func initCommand() (*cobra.Command, generator.StringGeneratorOpts)  {
    var opts generator.StringGeneratorOpts

    rootCommand := &cobra.Command {
        Use: `generate`,
        Short: `Generate random strings`,
        RunE: func(cmd *cobra.Command, args []string) error {
            if len(opts.SizeRange) != 2 {
                fmt.Print(opts.SizeRange)
                return errors.New("--size flag must have two parameters")
            }
            
            generator.GenerateStrings(opts)

            return nil
        },
    }

    rootCommand.Flags().StringVar(&opts.Path ,"path", "./rand-strings", "Directory relative path for output")
    rootCommand.Flags().UintVar(&opts.Count, "count", 10, "Number of strings to generate")
    rootCommand.Flags().Float32Var(&opts.RepeateChance, "chance", 0.2, "Number of strings to generate")
    rootCommand.Flags().UintSliceVar(&opts.SizeRange, "size", []uint{10, 100}, "Size range of a string")

    opts.Prefix = "data"

    rootCommand.MarkFlagRequired("size")
    rootCommand.MarkFlagRequired("count")
    rootCommand.MarkFlagRequired("path")
    rootCommand.MarkFlagRequired("chance")

    return rootCommand, opts
}

func main() {
    cmd, _ := initCommand()
    
    cmd.Execute()
}
