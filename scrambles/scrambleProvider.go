// Package scrambles contains infrastructure for providing scrambles to user from tnoodle and other services.
package scrambles

import (
	"log"
	"os"
	"strings"
	"time"
)

type scramble struct {
	event string
	moves string
}

// ScrambleProvider contains all scrambles. It provides functions to update and get scrambles.
type ScrambleProvider struct {
	scrambles map[string]string
	finish    chan int
	get       map[string]chan string
	request   chan string
	update    chan scramble
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
	sp.request <- event
	return <-sp.get[event]
}

func trimSuffix(event string) string {
	event = strings.TrimSuffix(event, "wf")
	event = strings.TrimSuffix(event, "oh")
	return event
}

func (sp *ScrambleProvider) manage() {
	select {
	case event := <-sp.request:
		sp.get[event] <- sp.scrambles[event]
	case scramble := <-sp.update:
		sp.scrambles[scramble.event] = scramble.moves
	case <-sp.finish:
		return
	}
}

// NewScrambleProvider creates an instance of scrambling provider
func NewScrambleProvider(finish chan int) *ScrambleProvider {
	sp := ScrambleProvider{
		scrambles: make(map[string]string),
		finish:    finish,
		get:       createGetChans(),
		request:   make(chan string),
		update:    make(chan scramble),
	}
	sp.startUpdating()
	go func() {
		for {
			sp.manage()
		}
	}()
	return &sp
}

func (sp *ScrambleProvider) startUpdating() {
	n := int64(len(Events))
	dur, _ := time.ParseDuration(os.Getenv("SCRAMBLING_INTERVAL"))
	interval := time.Duration(dur.Nanoseconds() / n)
	for _, event := range Events {
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
	sp.update <- scramble{event: event, moves: newScr}
}

func createGetChans() map[string]chan string {
	get := make(map[string]chan string)
	for _, event := range Events {
		get[event] = make(chan string)
	}
	return get
}

func handleError(err error) {
	log.Println(err)
}
