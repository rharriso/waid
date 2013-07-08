package entry

import (
	"bytes"
	"math"
	"time"
)

type Entry struct {
	Id        int64     `db:"id"`
	startTime int64     `db:"start_time"`
	endTime   int64     `db:"end_time"`
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
func (e *Entry) String() string {
	// write to buffer
	var buffer bytes.Buffer
	buffer.WriteString(e.Msg)
	buffer.WriteString(": ")

	// get duration variables
	d := e.Duration()
	hours := int(math.Floor(d.Hours()))
	min := int(math.Floor(d.Minutes()))
	sec := int(math.Floor(d.Seconds()))

	// maybe write hours
	if hours > 0 {
		buffer.WriteString(string(hours))
		buffer.WriteString("h ")
	}

	// write minutes a seconds
	buffer.WriteString(string(min))
	buffer.WriteString("m ")
	buffer.WriteString(string(sec))
	buffer.WriteString("s ")

	return buffer.String()
}

/*
  Return the time passed between the start time and the end time or now.
*/
func (e *Entry) Duration() time.Duration {
	var d time.Duration

	if e.Ended() {
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
func (e *Entry) PreUpdate() {
	// default to now
	if !e.Started() {
		e.Start = time.Now()
	}
	e.startTime = e.Start.Unix()

	if !e.Ended() {
		e.endTime = e.End.Unix()
	}
}

/*
  see: PreUpdate
*/
func (e *Entry) PreInsert() {
	panic("FUCK")
	e.PreUpdate()
}

/*
  Set the Start and End variables based on the database values
*/
func (e *Entry) PostGet() {
	if e.Started() {
		e.Start = time.Unix(e.startTime, 0)
	}
	if e.Ended() {
		e.End = time.Unix(e.startTime, 0)
	}
}
