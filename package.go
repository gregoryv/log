/*
   Package log provides loggers for writing messages.
   Simple usage is using the Log method as a first class citizen

      Log := NewSyncLog(os.Stdout).Log
      Log("some", "nice", "message")

   Logging errors can also be simplified with the filtered variation.

      Log := NewSyncLog(os.Stdout).FilterEmpty().Log
      Log("") // will not be logged

      // Log errors only if there are any
      Log(nil) // nothing, it's nil
      Log(io.EOF)
*/
package fox

import (
	"fmt"
	"io"
	"sync"
)

func NewSyncLog(w io.Writer) *SyncLog {
	return &SyncLog{w: w}
}

type Logger interface {
	Log(...interface{})
}

type SyncLog struct {
	mu sync.Mutex
	w  io.Writer
}

// Log synchronizes calls to the underlying writer and makes sure
// each message ends with one new line
func (l *SyncLog) Log(v ...interface{}) {
	out := fmt.Sprint(v...)
	l.mu.Lock()
	fmt.Fprint(l.w, out)
	if len(out) == 0 || out[len(out)-1] != newline {
		l.w.Write([]byte{newline})
	}
	l.mu.Unlock()
}

var newline byte = '\n'

// FilterEmpty returns a wrapper filtering out empty and nil values
func (l *SyncLog) FilterEmpty() *FilterEmpty {
	return &FilterEmpty{l}
}

func (l *SyncLog) SetOutput(w io.Writer) { l.w = w }

type FilterEmpty struct {
	sync *SyncLog
}

// Log calls the underlying logger only if v is non empty
func (l *FilterEmpty) Log(v ...interface{}) {
	switch len(v) {
	case 0:
		return
	case 1:
		if v[0] == nil {
			return
		}
	}
	out := fmt.Sprint(v...)
	if out == "" {
		return
	}
	l.sync.Log(out)
}
