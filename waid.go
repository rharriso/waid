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
	"time"
)

var (
	// list of commands
	cmdFlags = map[string]*flag.FlagSet{
		"start": flag.NewFlagSet("start", flag.ExitOnError),
		"stop":  flag.NewFlagSet("stop", flag.ExitOnError),
		"list":  flag.NewFlagSet("list", flag.ExitOnError),
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
	case "list":
		list()
	}
}

/*
	start
*/
func start() {
	// insert a new entry starting now
	e := &entry.Entry{Start: time.Now().Unix(), Msg: *msg}
	err := dbMap.Insert(e)
	doPanic(err)
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
	if len(entries) == 0 || entries[0].End != 0 {
		fmt.Println("No active entry")
		return
	}

	// update entry values
	e := *entries[0]
	e.End = time.Now().Unix()

	if *msg != "" {
		e.Msg = *msg
	}

	if e.Msg == "" {
		fmt.Print("Enter a message for this entry: ")
		in := bufio.NewReader(os.Stdin)
		input, _, err := in.ReadLine()
		doPanic(err)
		e.Msg = string(input)
	}

	// update table entry
	dbMap.Update(&e)
}

/*
	show all the entries in the database
*/
func list() {
	var entries []*entry.Entry
	_, err := dbMap.Select(&entries, "SELECT * FROM entries")
	doPanic(err)

	for _, e := range entries {
		fmt.Println(*e)
	}
}

/*
	dbConnect ->
		connect to databse and create tables maybe
*/
func dbConnect() {
	db, err := sql.Open("sqlite3", "./waid.db")
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
