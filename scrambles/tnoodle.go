package scrambles

import (
	"net/http"
	"io/ioutil"
	"encoding/json"
	"os"
	"log"
)

var tnoodleNames = map[string]string{
	"2x2":    "222",
	"3x3":    "333",
	"4x4":    "444",
	"5x5":    "555",
	"6x6":    "666",
	"7x7":    "777",
	"3x3bld": "333ni",
	"4x4bld": "444ni",
	"5x5bld": "555ni",
	"3x3fm":  "333fm",
	"clock":  "clock",
	"pyra":   "pyram",
	"mega":   "minx",
	"sq1":    "sq1",
	"skewb":  "skewb",
}

type ScrambleContainer struct {
	Scrambles []string `scrambles`
}

func genScramble(event string) (string, error) {
	event = tnoodleNames[event]

	resp, err := http.Get("http://" + os.Getenv("TNOODLE_IP") + ":2014/scramble/.json?=" + event + "*1")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	scrJson, err := ioutil.ReadAll(resp.Body)

	var scrContainer []ScrambleContainer
	err = json.Unmarshal(scrJson, &scrContainer)
	if err != nil {
		log.Println("Error while unmurshalling event " + event + " json: " + string(scrJson))
		return "", err
	}

	scramble := scrContainer[0].Scrambles[0]

	return scramble, nil
}
