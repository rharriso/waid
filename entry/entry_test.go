package entry

import (
	"testing"
	"time"
)

func TestAll(t *testing.T) {
	t.Error("Nothing in All")
}

func TestEnded(t *testing.T) {
	e := Entry{}
	if e.Ended() {
		t.Errorf("Entry shouldn't be ended")
	}

	e.End = time.Now()
	if !e.Ended() {
		t.Error("Entry should be ended")
	}
}

func TestStarted(t *testing.T) {
	e := Entry{}
	if e.Started() {
		t.Error("Entry shouldn't be started")
	}

	e.Start = time.Now()
	if !e.Started() {
		t.Error("Entry should have started")
	}
}

func TestDuration(t *testing.T) {
	start := time.Date(1999, 11, 5, 12, 12, 12, 12, time.UTC)
	e := Entry{Start: start}

	diff := e.Duration().Nanoseconds() - time.Now().Sub(e.Start).Nanoseconds()
	if diff > 0 {
		t.Error("If not ended duration should be before now")
	}

	end := time.Date(1999, 11, 6, 12, 12, 12, 12, time.UTC)
	e.End = end
	if e.Duration() != end.Sub(e.Start) {
		t.Error("If ended duration should be before end time")
	}
}

func TestPreUpdate(t *testing.T) {
	t.Error("Nothing in preupdate")
}

func TestPreInsert(t *testing.T) {
	TestPreUpdate(t)
}

func TestPostGet(t *testing.T) {
	t.Error("Nothing in PostGet")
}
