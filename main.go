package main

import (
	"fmt"
	"os"

	"cloud-utils/commands"
)

func main() {
	if err := commands.RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
