package main

import (
	"fmt"
	"os"
	"github.com/jheck90/terraform-copy/cmd" // Import your command package
)

func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
