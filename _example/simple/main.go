package main

import (
	"log"
	"math/rand"
	"net/http"
	"sync"
	"time"
	"github.com/ti/velox"
)

type Foo struct {
	//required velox state, adds sync state and a Push() method
	velox.State
	//optional mutex, prevents race conditions (foo.Push will make use of the sync.Locker interface)
	sync.Mutex
	A, B int
	C    Bar
}

type Bar struct {
	X, Y int
}

func main() {
	//state we wish to sync
	foo := &Foo{A: 21, B: 42}
	go func() {
		for {
			foo.Lock()
			foo.C.X = rand.Intn(99)
			foo.C.Y = rand.Intn(99)
			foo.Unlock()
			//push to all connections
			foo.Push()
			//do other stuff...
			time.Sleep(2500 * time.Millisecond)
		}
	}()
	//sync handlers
	http.Handle("/sync", velox.SyncHandler(foo))
	//index handler
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write(indexhtml)
	})
	//listen!
	log.Printf("Listening on 3000...")
	log.Fatal(http.ListenAndServe(":3000", nil))
}

var indexhtml = []byte(`
<div>Status: <b id="status">disconnected</b></div>
<pre id="example"></pre>
<script>
var evtSource = new EventSource('http://127.0.0.1:3000/sync');
evtSource.onmessage = function(e) {
   var v =  JSON.parse(e.data)
	document.querySelector("#example").innerHTML = JSON.stringify(v.body, null, 2);
};
evtSource.onopen = function() {
	document.querySelector("#status").innerHTML = "connected";
};
evtSource.onerror = function(e) {
	document.querySelector("#status").innerHTML = "disconnected";
};
</script>
`)
