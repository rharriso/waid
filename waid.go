package main

import (
	"bufio"
	"database/sql"
	"flag"
	"fmt"
	"github.com/coopernurse/gorp"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rharriso/waid/entry"
	"log"
	"os"
	"os/user"
	"time"
)

var (
	// list of commands
	cmdFlags = map[string]*flag.FlagSet{
		"start": flag.NewFlagSet("start", flag.ExitOnError),
		"stop":  flag.NewFlagSet("stop", flag.ExitOnError),
		"list":  flag.NewFlagSet("list", flag.ExitOnError),
		"clear": flag.NewFlagSet("clear", flag.ExitOnError),
	}

	// flag values
	msg   *string
	dbMap *gorp.DbMap
)

/*
	main ->
*/
func main() {
	// connect to db
	dbConnect()

	// get the command and flags
	flag.Parse()

	//if nothing passed do clear
	var cmd string
	if len(flag.Args()) == 0 {
		list()
		return
	} else {
		cmd = flag.Args()[0]
	}

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
	case "list":
		list()
	case "clear":
		clear()
	}
}

/*
	start
*/
func start() {
	// insert a new entry starting now
	e := entry.Entry{Msg: *msg}
	err := dbMap.Insert(&e)
	doPanic(err)
	fmt.Println("Starting activity: ", e.Msg)
}

/*
	stop the current entry.
*/
func stop() {
	// find most recent entry
	var entries []*entry.Entry
	_, err := dbMap.Select(&entries, "SELECT * FROM entries ORDER BY start_time DESC LIMIT 1")
	doPanic(err)

	// check for active entry
	if len(entries) == 0 || entries[0].Ended() {
		fmt.Println("No active entry")
		return
	}

	// update entry values
	e := entries[0]
	e.End = time.Now()

	// set msg to tagged value
	if *msg != "" {
		e.Msg = *msg
	}

	// if the message is empty, ask the ser for one
	if e.Msg == "" {
		fmt.Print("Enter a message for this entry: ")
		in := bufio.NewReader(os.Stdin)
		input, _, err := in.ReadLine()
		doPanic(err)
		e.Msg = string(input)
	}

	// update table entry
	dbMap.Update(e)
	fmt.Println("Activity Finished:", e.Msg, "|", e.TimeString())
}

/*
	show all the entries in the database
*/
func list() {
	entries, err := entry.All(dbMap)
	doPanic(err)

	fmt.Println("All Entries")
	fmt.Println("-------------------------------------")

	for _, e := range entries {
		fmt.Printf("-- %s\t%s\n", e.TimeString(), e.Msg)
	}

	fmt.Println("-------------------------------------")
}

/*
	empty the entries
*/
func clear() {
	entries, err := entry.All(dbMap)
	doPanic(err)

	for _, e := range entries {
		dbMap.Delete(e)
	}
}

/*
	dbConnect ->
		connect to databse and create tables maybe
*/
func dbConnect() {
	usr, err := user.Current()
	doPanic(err)
	db, err := sql.Open("sqlite3", usr.HomeDir+"/.waid.db")
	doPanic(err)

	dbMap = &gorp.DbMap{Db: db, Dialect: gorp.SqliteDialect{}}
	// add entries table
	dbMap.AddTableWithName(entry.Entry{}, "entries").SetKeys(true, "Id")
	dbMap.CreateTablesIfNotExists()
}

/*
	doPanic
*/
func doPanic(err error) {
	if err != nil {
		panic(err)
	}
}
