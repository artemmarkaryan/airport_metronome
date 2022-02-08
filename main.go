package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"
)

type C struct {
	sync.Mutex
	t              uint64
	secondDuration float64
}

func (c *C) start() {
	for {
		time.Sleep(time.Millisecond * time.Duration(c.secondDuration*1000))
		c.Mutex.Lock()
		c.t += 1
		c.Mutex.Unlock()
	}
}

func (c *C) getTime(w http.ResponseWriter, _ *http.Request) {
	var t uint64
	{
		c.Mutex.Lock()
		t = c.t
		c.Mutex.Unlock()
	}

	_, _ = fmt.Fprintf(w, "%d", t)
}

func (c *C) setSpeed(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		writeError(w, http.StatusBadRequest, "wrong method, use get")
		return
	}

	var q []string
	var ok bool
	if q, ok = req.URL.Query()["speed"]; !ok {
		writeError(w, http.StatusBadRequest, "speed parameter not found")
		return
	}

	if len(q) != 1 {
		writeError(w, http.StatusBadRequest, "speed parameter must be single")
		return
	}

	speed, err := strconv.ParseFloat(q[0], 64)
	if err != nil {
		writeError(w, http.StatusBadRequest, fmt.Sprintf("error parsing speed: %s", err.Error()))
		return
	}

	c.Lock()
	c.secondDuration = 1.0 / speed
	c.Unlock()

	w.WriteHeader(http.StatusOK)
	return
}

func main() {
	c := new(C)
	c.secondDuration = 1
	go c.start()

	http.HandleFunc("/time", c.getTime)
	http.HandleFunc("/setSpeed", c.setSpeed)

	http.ListenAndServe(":8090", nil)
}

func writeError(w http.ResponseWriter, status int, description string) {
	b, err := json.Marshal(
		struct {
			Msg string `json:"msg"`
		}{
			Msg: description,
		},
	)
	if err != nil {
		log.Printf("error converting error to json: %s", err.Error())
		return
	}

	w.WriteHeader(status)
	_, err = w.Write(b)
	if err != nil {
		log.Printf("error writing response: %s", err.Error())
		return
	}
}
