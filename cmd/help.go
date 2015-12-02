package cmd

import (
	"fmt"
)

type CmdHelp struct {
	Name    string
	Handler func()
}

var cmdhelps = []CmdHelp{
	{"dp", dpUsage},
	{"ep", epUsage},
	{"login", loginUsage},
	{"pub", pubUsage},
	{"pull", pullUsage},
	{"repo", repoUsage},
	{"subs", subsUsage},
}

func Help(login bool, args []string) (err error) {
	fmt.Println("Usage:\tdatahub COMMAND [arg...]\n\tdatahub COMMAND [ --help ]\n\tdatahub help [COMMAND]\n\nCommands:")
	if len(args) == 0 {
		for _, v := range Cmd {
			fmt.Printf("\t%-16s  %s\n", v.Name, v.Desc)
		}
		return nil
	} else {
		for _, v := range cmdhelps {
			if args[0] == v.Name {
				v.Handler()
				break
			}
		}
	}
	return
}
