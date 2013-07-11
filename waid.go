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
	"math"
	"os"
	"os/user"
	"time"
)

var (
	// list of commands
	cmdFlags = map[string]*flag.FlagSet{
		"help":  flag.NewFlagSet("help", flag.ExitOnError),
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
		help()
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
	case "help":
		help()
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

	fmt.Println("\nAll Entries")
	fmt.Println("-------------------------------------")

	var totalHours, totalMinutes, totalSeconds int
	totalHours = 0

	for _, e := range entries {
		fmt.Printf("-- %s\t%s\n", e.TimeString(), e.Msg)
		totalHours += int(math.Floor(e.Duration().Hours()))
		totalMinutes += int(math.Floor(e.Duration().Minutes()))
		totalSeconds += int(math.Floor(e.Duration().Seconds()))
	}

	fmt.Println("-------------------------------------")
	fmt.Printf("Total - %dh %dm %ds\n\n", totalHours, totalMinutes, totalSeconds)
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
	Help out the user
*/
func help() {
	fmt.Println("Usage: waid [command] [options]\n")

	fmt.Println("Commands:")
	fmt.Println("\tstart\t- start a new task")
	fmt.Println("\tstop\t- complete current task")
	fmt.Println("\tlist\t- list all tasks")
	fmt.Println("\tclear\t- clear list of tasks")

	fmt.Println("Options:")
	fmt.Println("\t-m\t- add message to the current task on start or stop.")

	fmt.Println("")
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
