package main

import (
	"github.com/axelarnetwork/utils/checks"
	"log"

	"github.com/spf13/cobra"
)

func main() {
	cmd := &cobra.Command{
		Use: "check",
	}
	cmd.AddCommand(checks.FieldDeclarations())

	err := cmd.Execute()
	if err != nil {
		log.Fatal(err)
	}
}
