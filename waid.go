package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/rharriso/waid/entry"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

var (
	// list of commands
	cmdFlags = map[string]*flag.FlagSet{
		"help":  flag.NewFlagSet("help", flag.ExitOnError),
		"start": flag.NewFlagSet("start", flag.ExitOnError),
		"stop":  flag.NewFlagSet("stop", flag.ExitOnError),
		"add":   flag.NewFlagSet("add", flag.ExitOnError),
		"list":  flag.NewFlagSet("list", flag.ExitOnError),
		"clear": flag.NewFlagSet("clear", flag.ExitOnError),
	}

	// flag values
	msg *string
	dur *string
)

/*
	main ->
*/
func main() {
	// connect to db
	// dbConnect()

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
	dur = flags.String("t", "", "duration for activity")
	flags.Parse(flag.Args()[1:])

	// run the command with flags
	switch cmd {
	case "start":
		start()
	case "stop":
		stop(true)
	case "add":
		add()
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
	//if active entry, ask to end
	// e := entry.Latest(dbMap)
	// if e != nil && !e.Ended() {
	// 	// do they want to close the old one?
	// 	if confirm(fmt.Sprintf("End Activity (%s)", e.Msg)) {
	// 		stop(false)
	// 	} else {
	// 		return
	// 	}
	// }

	// // insert a new entry starting now
	// newEntry := entry.Entry{Msg: *msg}
	// err := dbMap.Insert(&newEntry)
	// doPanic(err)
	// fmt.Println("Starting activity: ", newEntry.Msg)
}

/*
	Find the active entry.
	Ask user to enter a message if none has been set yet.
	Set the end time to now, and save.
*/
func stop(fromCommand bool) {
	// e := entry.Latest(dbMap)

	// // check for active entry
	// if e == nil || e.Ended() {
	// 	fmt.Println("No active entry")
	// 	return
	// }

	// if fromCommand && *msg != "" {
	// 	e.Msg = *msg
	// }

	// if e.Msg == "" {
	// 	fmt.Print("Enter a message for this entry: ")
	// 	in := bufio.NewReader(os.Stdin)
	// 	input, _, err := in.ReadLine()
	// 	doPanic(err)
	// 	e.Msg = string(input)
	// }

	// e.End = time.Now()
	// dbMap.Update(e)
	// fmt.Println("Activity Finished:", e.Msg, "|", e.TimeString())
}

/*
	add an entry directly
*/
func add() {
	e := entry.Entry{Msg: *msg}
	e.SetDuration(*dur)

	//post json data to the server
	jsonData, err := json.Marshal(e)
	doPanic(err)
	resp, err := http.Post("http://localhost:3000/entries", "json", bytes.NewBuffer(jsonData))
	doPanic(err)

	entryBytes, err := ioutil.ReadAll(resp.Body)
	doPanic(err)
	json.Unmarshal(entryBytes, &e)

	fmt.Println("Activity Added:", e.Msg, "|", e.TimeString())
}

/*
	show all the entries in the database
*/
func list() {
	resp, err := http.Get("http://localhost:3000/entries")
	doPanic(err)
	entryData, err := ioutil.ReadAll(resp.Body)
	doPanic(err)

	var entries []entry.Entry
	json.Unmarshal(entryData, &entries)

	fmt.Println("\nAll Entries")
	fmt.Println("-------------------------------------")

	var total time.Duration

	for _, e := range entries {
		if e.Ended() {
			fmt.Printf("-- %d\t%s\t%s\n", e.Id, e.TimeString(), e.Msg)
		} else {
			fmt.Printf("-- \033[033m%d\t%s\t%s%s\n",
				e.Id,
				e.TimeString(),
				e.Msg,
				" <= active \033[0m")
		}

		total = total + e.Duration()
	}

	fmt.Println("-------------------------------------")
	fmt.Printf("Total - %dh %dm %ds\n\n",
		int(total.Hours()),
		int(total.Minutes())%60,
		int(total.Seconds())%60)
}

/*
	empty the entries
*/
func clear() {
	// if confirm("Delete all the entries? ") {
	// 	list()
	// 	entries := entry.All(dbMap)

	// 	for _, e := range entries {
	// 		dbMap.Delete(e)
	// 	}
	// 	fmt.Println("Entries Deleted.")
	// }
}

/*
	Help out the user
*/
func help() {
	fmt.Println("Usage: waid [command] [options]\n")

	fmt.Println("Commands:")
	fmt.Println("\tstart\t- start a new entry")
	fmt.Println("\tstop\t- complete current entry")
	fmt.Println("\tadd\t- add a complete entry")
	fmt.Println("\tlist\t- list all entry")
	fmt.Println("\tclear\t- clear list of entry")

	fmt.Println("Options:")
	fmt.Println("\t-m\t- add message to the current task on start or stop.")
	fmt.Println("\t-t\t- time (see go duration format), used for adding entries.")

	fmt.Println("")
}

/*
	doPanic
*/
func doPanic(err error) {
	if err != nil {
		panic(err)
	}
}

/*
	confirm
		ask the user if they really want to do that
*/
func confirm(msg string) bool {
	fmt.Print(fmt.Sprintf("%s (Y/n): ", msg))
	var answer string
	fmt.Scanf("%s", &answer)

	return strings.ToUpper(answer) == "Y"
}
