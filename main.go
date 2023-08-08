package main

import (
	"os"
	//"runtime"
	"runtime/pprof"

	"github.com/andreistan26/top-string/cmd"
)


func main() {
    f, _ := os.Create("profile")    
    pprof.StartCPUProfile(f)
    defer pprof.StopCPUProfile()
    //runtime.SetCPUProfileRate(300)

    rootCmd := cmd.CreateMainCommand()
    rootCmd.AddCommand(cmd.CreateSenderCommand())
    rootCmd.AddCommand(cmd.CreateServerCommand())
    cmd, _ := cmd.CreateLocalCommand()
    rootCmd.AddCommand(cmd)

    rootCmd.Execute()
}
