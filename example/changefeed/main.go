package main

import (
	"log"
	"math/rand"
	"net/http"

	"sync"
	"time"

	"github.com/ti/velox"
)

type Results struct {
	//required velox state, adds sync state and a Push() method
	velox.State
	//optional mutex, prevents race conditions (foo.Push will make use of the sync.Locker interface)
	sync.Mutex
	//realtime database results
	X, Y, Z int
}

func main() {
	//sync handlers
	router := http.NewServeMux()
	router.HandleFunc("/sync", func(w http.ResponseWriter, r *http.Request) {
		results := &Results{}
		//hijack request
		conn, err := velox.Sync(results, w, r)
		if err != nil {
			log.Printf("[velox] sync handler error: %s", err)
			return
		}
		connected := true
		//connected, now query for results
		go func() {
			//load all results
			results.X = rand.Intn(99)
			results.Y = 2
			results.Z = 3
			results.Push()
			//then poll db for delta
			//OR push delta from db if has support (https://rethinkdb.com/docs/changefeeds/)
			for connected {
				results.Y = rand.Intn(99)
				time.Sleep(1 * time.Second)
				results.Push()
			}
		}()
		//wait here
		log.Printf("[%s] connected", results.ID())
		conn.Wait()
		log.Printf("[%s] disconnected", results.ID())
		//disconnected
		connected = false
	})
	//index handler
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write(indexhtml)
	})

	//listen!
	log.Printf("Listening on 3000...")
	log.Fatal(http.ListenAndServe(":3000", router))
}

var indexhtml = []byte(`
var evtSource = new EventSource('http://127.0.0.1:3000/sync');
    evtSource.onmessage = function(e) {
       var v =  JSON.parse(e.data)
        if (v.id !== undefined) {
            document.querySelector("#vid").textContent = v.id;
        }
        document.querySelector("#example").innerHTML = JSON.stringify(v.body, null, 2);
    };
    evtSource.onopen = function() {
        document.querySelector("#status").innerHTML = "connected";
    };
    evtSource.onerror = function(e) {
        document.querySelector("#status").innerHTML = "disconnected";
    };
`)
