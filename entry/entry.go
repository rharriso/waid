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
	e.StartTime = e.Start.Unix()

	if !e.Ended() {
		e.EndTime = e.End.Unix()
	}

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
