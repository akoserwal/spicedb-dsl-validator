/*
Copyright © 2023 Abhishek Koserwal
*/
package main

import (
	"fmt"
	"github.com/akoserwal/spicedb-dsl-validator/cmd"
)

var (
	version string
)

func main() {
	cmd.Execute()
	fmt.Printf("\nversion=%s\n", version)
}
