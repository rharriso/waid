package entry

import (
	"bytes"
	"github.com/coopernurse/gorp"
	"math"
	"strconv"
	"time"
)

type Entry struct {
	Id        int64     `db:"id"`
	StartTime int64     `db:"start_time"`
	EndTime   int64     `db:"end_time"`
	Msg       string    `db:"message"`
	Start     time.Time `db:"-"` // db ignore
	End       time.Time `db:"-"` // db ignore
}

/*
	return array of all Entries
*/
func All(dbMap *gorp.DbMap) []*Entry {
	var entries []*Entry
	_, err := dbMap.Select(&entries, "SELECT * FROM entries")

	if err != nil {
		panic(err)
	}

	return entries
}

/*
	return the latest
*/
func Latest(dbMap *gorp.DbMap) *Entry {
	// find most recent entry
	var entries []*Entry
	_, err := dbMap.Select(&entries, "SELECT * FROM entries ORDER BY start_time DESC LIMIT 1")

	if err != nil {
		panic(err)
	}

	if len(entries) > 0 {
		return entries[0]
	}
	return nil

}

/*
	returns true if the entry has an End time
*/
func (e *Entry) Ended() bool {
	return e.End.Unix() >= 0
}

/*
	returns true if the entry has a Start time
*/
func (e *Entry) Started() bool {
	return e.Start.Unix() >= 0
}

/*
	Return formated string describing the entry
*/
func (e *Entry) TimeString() string {
	var buffer bytes.Buffer

	// get duration variables
	d := e.Duration()
	hours := int(math.Floor(d.Hours()))
	min := int(math.Floor(d.Minutes())) % 60
	sec := int(math.Floor(d.Seconds())) % 60

	buffer.WriteString(strconv.Itoa(hours))
	buffer.WriteString("h ")
	buffer.WriteString(strconv.Itoa(min))
	buffer.WriteString("m ")
	buffer.WriteString(strconv.Itoa(sec))
	buffer.WriteString("s")

	return buffer.String()
}

/*
	Return the time passed between the start time and the end time or now.
*/
func (e *Entry) Duration() time.Duration {
	var d time.Duration

	if !e.Ended() {
		d = time.Now().Sub(e.Start)
	} else {
		d = e.End.Sub(e.Start)
	}

	return d
}

/*
	Set the start and end times based on a duraction
*/
func (e *Entry) SetDuration(s string) {
	duration, err := time.ParseDuration(s)
	if err != nil {
		panic(err)
	}

	e.Start = time.Now()
	e.End = e.Start.Add(duration)
	e.setUnixTimes()
}

/*
	Gorp Hooks
*/
/*
	Set the startTime and endTime variables based on the time.Time members
*/
func (e *Entry) PreUpdate(s gorp.SqlExecutor) error {
	// default to now
	if !e.Started() {
		e.Start = time.Now()
	}
	e.setUnixTimes()

	return nil
}

/*
	see: PreUpdate
*/
func (e *Entry) PreInsert(s gorp.SqlExecutor) error {
	return e.PreUpdate(s)
}

/*
	Set the Start and End variables based on the database values
*/
func (e *Entry) PostGet(s gorp.SqlExecutor) error {
	e.Start = time.Unix(e.StartTime, 0)
	e.End = time.Unix(e.EndTime, 0)

	return nil
}

/*
	set unix time values based on go time
*/
func (e *Entry) setUnixTimes() {
	e.StartTime = e.Start.Unix()
	e.EndTime = e.End.Unix()
}
