package main

import (
	"github.com/andreistan26/top-string/cmd"
)


func main() {

    rootCmd := cmd.CreateMainCommand()
    rootCmd.AddCommand(cmd.CreateSenderCommand())
    rootCmd.AddCommand(cmd.CreateServerCommand())
    rootCmd.AddCommand(cmd.CreateLocalCommand())

    rootCmd.Execute()
}
