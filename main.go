package main

import (
	"flag"
	"fmt"
	"log"
)

var (
	// list of commands
	cmdFlags = map[string]*flag.FlagSet{
		"start": flag.NewFlagSet("start", flag.ExitOnError),
		"stop":  flag.NewFlagSet("stop", flag.ExitOnError),
	}
)

func main() {
	flag.Parse()
	// get the command
	cmd := flag.Args()[0]
	flags, ok := cmdFlags[cmd]

	// check for valid command
	if !ok {
		log.Fatalf("no command %q for waid", cmd)
	}

	//get message from remaining flags
	msg := flags.String("m", "", "message for activity")
	flags.Parse(flag.Args()[1:])

	fmt.Println(cmd, "with message", *msg)
}
