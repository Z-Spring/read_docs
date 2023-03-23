package cmd

import (
	"bufio"
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"os"
)

var (
	name    string
	RootCmd = &cobra.Command{
		Use:   "sp",
		Short: "cat your file contents",
		Run:   CatFile,
	}
)

func CatFile(cmd *cobra.Command, args []string) {
	file, err := os.Open(args[1])
	if err != nil {
		panic(err)
	}
	scanner := bufio.NewScanner(file)
	line := 0

	for scanner.Scan() {
		fmt.Printf("line: %d | Contents: %q \n", line, scanner.Bytes())
		line++
	}
	err = scanner.Err()
	if err != nil {
		log.Fatalf("scanner err: %v", err)
	}
}
