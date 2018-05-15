// Simple web service that counts views for type+id.
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

type app struct {
	views          map[string]map[int]int
	viewsLock      sync.RWMutex
	addr, strTypes string
	types          []string
}

func main() {
	a := app{
		views: make(map[string]map[int]int),
	}

	flag.StringVar(&a.addr, "addr", ":8090", "address listen to")
	flag.StringVar(&a.strTypes, "types", "item,news", "allowed types")
	flag.Parse()

	// register handlers and start listening
	a.registerHandlers()
	err := http.ListenAndServe(a.addr, nil)
	if err != nil {
		log.Fatal(err)
	}
}

// registerHandlers register /<type>/<id> handlers for track views and /pop for get data.
func (a *app) registerHandlers() {
	a.types = strings.Split(a.strTypes, ",")
	for _, t := range a.types {
		a.views[t] = make(map[int]int)
		http.HandleFunc("/"+t+"/", a.handleView(t))
	}

	http.HandleFunc("/pop", a.handlePop())
}

// handleView handles views for /type/id.
func (a *app) handleView(t string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		d := strings.SplitN(r.URL.Path, "/", 3)
		if len(d) == 2 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// get id, type in d[1]
		id, err := strconv.Atoi(d[2])
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		// store data
		a.viewsLock.Lock()
		defer a.viewsLock.Unlock()
		v := a.views[d[1]][id] + 1
		a.views[d[1]][id] = v

		// returns new count if data=1 is set
		if r.URL.Query().Get("data") == "1" {
			fmt.Fprint(w, v)
		}
	}
}

// handlePop returns all data in JSON format.
func (a *app) handlePop() http.HandlerFunc {
	type response struct {
		Type      string
		ID, Views int
	}

	return func(w http.ResponseWriter, r *http.Request) {
		var views []response

		// collect data
		a.viewsLock.Lock()
		for t, ids := range a.views {
			for id, v := range ids {
				views = append(views, response{t, id, v})
			}

			// delete data if keep=1 is not set.gi
			if r.URL.Query().Get("keep") != "1" {
				a.views[t] = make(map[int]int)
			}
		}
		a.viewsLock.Unlock()

		// return data
		if x, err := json.Marshal(views); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			w.Header().Set("Content-Type", "application/json")
			w.Write(x)
		}
	}
}
