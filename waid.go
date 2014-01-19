package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/rharriso/waid/entry"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

var (
	// constants
	SERVER_URL = "http://localhost:3000"

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
	var e entry.Entry
	jsonRequest("GET", "/entries/latest", &e)

	// if this activity is still going
	if e.Started() && !e.Ended() {
		// do they want to close the old one?
		if confirm(fmt.Sprintf("End Activity (%s)", e.Msg)) {
			stop(false)
		} else {
			return
		}
	}
	//post the new one
	jsonRequest("POST", "/entries", &entry.Entry{Msg: *msg})
}

/*
	Find the active entry.
	Ask user to enter a message if none has been set yet.
	Set the end time to now, and save.
*/
func stop(fromCommand bool) {
	//if active entry, ask to end
	var e entry.Entry
	jsonRequest("GET", "/entries/latest", &e)

	// check for active entry
	if e.Started() && e.Ended() {
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
	path := fmt.Sprintf("/entries/%d", e.Id)
	jsonRequest("PUT", path, &e)
	fmt.Println("Activity Finished:", e.Msg, "|", e.TimeString())
}

/*
	add an entry directly
*/
func add() {
	e := entry.Entry{Msg: *msg}
	e.SetDuration(*dur)
	jsonRequest("POST", "/entries", &e)

	fmt.Println("Activity Added:", e.Msg, "|", e.TimeString())
}

/*
	show all the entries in the database
*/
func list() {
	// get all the entries
	var entries []entry.Entry
	jsonRequest("GET", "/entries", &entries)

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
	if confirm("Delete all the entries? ") {
		list()
		jsonRequest("DELETE", "/entries", nil)
		fmt.Println("Entries Deleted.")
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

/*
  Utility method for sending request of different types
*/
func jsonRequest(reqType string, path string, v interface{}) {
	var err error
	client := http.Client{}

	//prepare json data to send to the server
	jsonData, err := json.Marshal(v)
	doPanic(err)
	body := bytes.NewBuffer(jsonData)

	// create req uest
	req, err := http.NewRequest(reqType, fmt.Sprintf("%s%s", SERVER_URL, path), body)
	doPanic(err)
	resp, err := client.Do(req)
	doPanic(err)

	// error status code
	if resp.StatusCode >= 400 {
		panic(fmt.Sprintf("Request %d: %s:", resp.StatusCode, resp.Status))
	}

	// set returned "1data to interface, this will crash if the types aren't good
	entryData, err := ioutil.ReadAll(resp.Body)
	doPanic(err)

	json.Unmarshal(entryData, &v)
}
