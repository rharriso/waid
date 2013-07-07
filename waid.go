package main

import (
	"database/sql"
	"flag"
	"fmt"
	"github.com/coopernurse/gorp"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rharriso/waid/entry"
	"log"
	"time"
)

var (
	// list of commands
	cmdFlags = map[string]*flag.FlagSet{
		"start": flag.NewFlagSet("start", flag.ExitOnError),
		"stop":  flag.NewFlagSet("stop", flag.ExitOnError),
	}
	// flag values
	msg   *string
	dbMap *gorp.DbMap
)

func main() {
	// connect to db
	dbConnect()

	// get the command and flags
	flag.Parse()
	cmd := flag.Args()[0]
	flags, ok := cmdFlags[cmd]

	// check for valid command
	if !ok {
		log.Fatalf("no command %q for waid", cmd)
	}

	//get message from remaining flags
	msg = flags.String("m", "", "message for activity")
	flags.Parse(flag.Args()[1:])

	// run the command with flags
	switch cmd {
	case "start":
		start()
	case "stop":
		stop()
	}
}

func start() {
	entry := entry.Entry{Start: time.Now().Unix(), Msg: *msg}
	fmt.Println("start", entry)
}

func stop() {
	fmt.Println("stop", *msg)
}

func dbConnect() {
	db, err := sql.Open("sqlite3", "waid.db")
	if err != nil {
		panic(err)
	}

	dbmap = &gorp.DbMap{suffix: "waid"}
	// add entries tabless
	entries := dbmap.AddTableWithName(entry.Entry{}, "entries").SetKeys(true, "Id")

	fmt.Println(db.Driver())
}
