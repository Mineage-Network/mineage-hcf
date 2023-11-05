package eotw

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// logFile will temporary (hackishly) store the timestamp of the EOTW.
const logFile = "eotw.log"

// eotw represents the private global object of the EOTW time.
var eotw time.Time

// init ...
func init() {
	if _, err := os.ReadFile(logFile); err != nil {
		_ = os.WriteFile(logFile, []byte(fmt.Sprint(time.Now().Unix())), 0777)
		return
	}

	d, _ := os.ReadFile(logFile)
	i, _ := strconv.ParseInt(string(d), 10, 64)
	eotw = time.Now().Add(time.Duration(i))
}

// Start starts the EOTW timer.
func Start() {
	// TODO: customizable time...
	eotw = time.Now().Add(1 * time.Hour)
}

// End ends the EOTW timer.
func End() {
	eotw = time.Time{}
}

// Save will save the EOTW time to the log file.
func Save() {
	s := fmt.Sprint(int64(time.Until(eotw)))
	_ = os.WriteFile(logFile, []byte(s), 0777)
}

// Running returns the EOTW time and if it is running.
func Running() (time.Time, bool) {
	return eotw, !eotw.IsZero()
}
