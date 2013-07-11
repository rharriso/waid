package utils

import (
	"time"
)

/*
  Return formated string describing the entry
*/
func (d time.Duration) TimeString() string {
    var buffer bytes.Buffer

    // get duration variables
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