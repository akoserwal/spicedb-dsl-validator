/*
Copyright © 2023 Abhishek Koserwal
*/
package main

import (
	"flag"
	"github.com/golang/glog"
	"spicedb-dsl-validator/cmd"
)

func main() {
	// This is needed to make `glog` believe that the flags have already been parsed, otherwise
	// every log messages is prefixed by an error message stating the the flags haven't been
	// parsed.

	_ = flag.CommandLine.Parse([]string{})

	if err := flag.Set("logtostderr", "true"); err != nil {
		glog.Infof("Unable to set logtostderr to true")
	}

	cmd.Execute()
}
