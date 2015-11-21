package client

import (
	"fmt"
	"github.com/asiainfoLDP/datahub/cmd"
	"os"
	"strings"
)

func RunClient() {

	if len(os.Args) < 2 {
		ShowUsage()
		os.Exit(2)
	}

	command := os.Args[1]

	commandFound := false
	for _, v := range cmd.Cmd {
		if strings.EqualFold(v.Name, command) {
			commandFound = true
			if len(os.Args) > 2 && len(os.Args[2]) > 0 && os.Args[2][0] != '-' {
				subCmdFound := false
				for _, vv := range v.SubCmd {
					if strings.EqualFold(vv.Name, os.Args[2]) {
						command += os.Args[2]
						subCmdFound = true
						vv.Handler(v.NeedLogin, os.Args[3:])
					}
				}
				if !subCmdFound {
					v.Handler(v.NeedLogin, os.Args[2:])
				}
			} else {
				v.Handler(v.NeedLogin, os.Args[2:])
			}
		}
	}

	if command == "help" {
		cmd.Help(false, os.Args[2:])
		commandFound = true
	}
	if !commandFound {
		fmt.Println(command, "not found")
		ShowUsage()
	}
	/*

		if len(os.Args) > 2 && os.Args[2][0] != '-' {
			command += os.Args[2]
		}

		cmdFound := false
		for _, v := range cmd.Cmds {
			if strings.EqualFold(v.CmdName, command) {
				v.Handler(os.Args[1:], v)
				cmdFound = true
			}

		}
		if !cmdFound {
			fmt.Println(os.Args[1], "not found.")
			ShowUsage()
		}
	*/
	return

}

func ShowUsage() {
	for _, v := range cmd.Cmd {
		fmt.Printf("%-16s  %s\n", v.Name, v.Desc)

	}
}
