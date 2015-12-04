package cmd

import (
	"fmt"
	"github.com/asiainfoLDP/datahub/ds"
)

func Version(needLogin bool, args []string) (err error) {
	fmt.Println("datahub", ds.DATAHUB_VERSION)
	return nil
}

func verUsage() {
	fmt.Println("Show datahub version")
}
