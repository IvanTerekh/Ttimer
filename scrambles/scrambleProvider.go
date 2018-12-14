// Package scrambles contains infrastructure for providing scrambles to user from tnoodle and other services.
package scrambles

import (
	"log"
	"os"
	"strings"
	"sync"
	"time"
)

// ScrambleProvider contains all scrambles. It provides functions to update and get scrambles.
type ScrambleProvider struct {
	scrambles map[string]string
	mx        sync.RWMutex
}

// Events contains all possible events.
var Events = []string{
	"2x2",
	"3x3",
	"4x4",
	"5x5",
	"6x6",
	"7x7",
	//"3x3fm",
	"3x3bld",
	"4x4bld",
	"5x5bld",
	"clock",
	"pyra",
	"mega",
	"sq1",
	"skewb",
}

// Get returns an actual scramble for given event.
func (sp *ScrambleProvider) Get(event string) string {
	event = trimSuffix(event)
	sp.mx.RLock()
	defer sp.mx.RUnlock()
	return sp.scrambles[event]
}

func trimSuffix(event string) string {
	event = strings.TrimSuffix(event, "wf")
	event = strings.TrimSuffix(event, "oh")
	return event
}

// NewScrambleProvider creates an instance of scrambling provider
func NewScrambleProvider() *ScrambleProvider {
	sp := ScrambleProvider{
		scrambles: make(map[string]string),
	}
	sp.startUpdating()
	return &sp
}

func (sp *ScrambleProvider) startUpdating() {
	n := int64(len(Events))
	dur, _ := time.ParseDuration(os.Getenv("SCRAMBLING_INTERVAL"))
	interval := time.Duration(dur.Nanoseconds() / n)
	for _, event := range Events {
		// Scrambles are generated with equal intervals to avoid scrambling server overload.
		time.Sleep(interval)
		go func(ev string) {
			for {
				sp.updateScramble(ev)
				time.Sleep(dur)
			}
		}(event)
	}
}

func (sp *ScrambleProvider) updateScramble(event string) {
	newScr, err := genScramble(event)
	if err != nil {
		handleError(err)
		return
	}
	sp.mx.Lock()
	defer sp.mx.Unlock()
	sp.scrambles[event] = newScr
}

func handleError(err error) {
	log.Println(err)
}
