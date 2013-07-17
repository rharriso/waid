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
	"strings"
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
		stop(true)
	case "list":
		list()
	case "clear":
		clear()
	case "help":
		help()
	}

}

/*
	Checks for active entry, asks to end it.
	Creates new activity
*/
func start() {
	e := entry.Latest(dbMap)

	//if active entry, ask to end
	if e != nil && !e.Ended() {
		fmt.Printf("End Activity (%s) (Y/n):", e.Msg)
		var answer string
		fmt.Scanf("%s", &answer)

		// if stopping then close old answer
		if strings.ToUpper(answer) == "Y" {
			stop(false)
		} else {
			return
		}
	}

	// insert a new entry starting now
	newEntry := entry.Entry{Msg: *msg}
	err := dbMap.Insert(&newEntry)
	doPanic(err)
	fmt.Println("Starting activity: ", newEntry.Msg)
}

/*
	Find the active entry.
	Ask user to enter a message if none has been set yet.
	Set the end time to now, and save.
*/
func stop(fromCommand bool) {
	e := entry.Latest(dbMap)

	// check for active entry
	if e == nil || e.Ended() {
		fmt.Println("No active entry")
		return
	}

	if fromCommand && *msg != "" {
		e.Msg = *msg
	}

	if e.Msg == "" {
		fmt.Print("Enter a message for this entry: ")
		in := bufio.NewReader(os.Stdin)
		input, _, err := in.ReadLine()
		doPanic(err)
		e.Msg = string(input)
	}

	e.End = time.Now()
	dbMap.Update(e)
	fmt.Println("Activity Finished:", e.Msg, "|", e.TimeString())
}

/*
	show all the entries in the database
*/
func list() {
	entries := entry.All(dbMap)

	fmt.Println("\nAll Entries")
	fmt.Println("-------------------------------------")

	var totalHours, totalMinutes, totalSeconds float64
	totalHours = 0

	for _, e := range entries {
		fmt.Printf("-- %s\t%s\n", e.TimeString(), e.Msg)
		totalHours += e.Duration().Hours()
		totalMinutes += e.Duration().Minutes()
		totalSeconds += e.Duration().Seconds()
	}

	fmt.Println("-------------------------------------")
	fmt.Printf("Total - %dh %dm %ds\n\n",
		int(totalHours),
		int(totalMinutes)%60,
		int(totalSeconds)%60)
}

/*
	empty the entries
*/
func clear() {
	entries := entry.All(dbMap)

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
	fmt.Println("\tstart\t- start a new entry")
	fmt.Println("\tstop\t- complete current entry")
	fmt.Println("\tlist\t- list all entry")
	fmt.Println("\tclear\t- clear list of entry")

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
