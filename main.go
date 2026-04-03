package main

import (
	"os"

	"github.com/btbytes/bleat/cmd"
)

func main() {
	args := os.Args[1:]

	if len(args) == 0 {
		cmd.RunPorts(false, "")
		return
	}

	switch args[0] {
	case "ps":
		all := false
		for _, a := range args[1:] {
			if a == "--all" || a == "-a" {
				all = true
			}
		}
		cmd.RunPs(all)
	case "clean":
		cmd.RunClean()
	case "watch":
		cmd.RunWatch()
	case "help", "--help", "-h":
		cmd.RunHelp()
	case "--all", "-a":
		cmd.RunPorts(true, "")
	default:
		if len(args) == 1 {
			all := false
			for _, a := range args {
				if a == "--all" || a == "-a" {
					all = true
				}
			}
			if all {
				cmd.RunPorts(true, "")
			} else {
				cmd.RunPorts(false, args[0])
			}
		} else {
			cmd.RunHelp()
		}
	}
}
