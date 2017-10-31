# velox

[![GoDoc](https://godoc.org/github.com/jpillora/velox?status.svg)](https://godoc.org/github.com/jpillora/velox)

Real-time JS object synchronisation over SSE and WebSockets in Go and JavaScript (Node.js and browser)

### Features

* Simple API
* Synchronise any JSON marshallable struct in Go
* Supports [Server-Sent Events (EventSource)](https://en.wikipedia.org/wiki/Server-sent_events) and [WebSockets](https://en.wikipedia.org/wiki/WebSocket)
* SSE [client-side poly-fill](https://github.com/remy/polyfills/blob/master/EventSource.js) to fallback to long-polling in older browsers (IE8+).
* Implement delta queries (return all results, then incrementally return changes)

### Quick Usage

Server (Go)

``` go
//syncable struct
type Foo struct {
	velox.State
	A, B int
}
foo := &Foo{}
//serve velox sync endpoint for foo
http.Handle("/sync", velox.SyncHandler(foo))
//make changes
foo.A = 42
foo.B = 21
//push to client
foo.Push()
```

Client (Node and Browser)

``` js
// load script /velox.js
var evtSource = new EventSource('http://127.0.0.1:3000/sync');
evtSource.onmessage = function(e) {
   var v =  JSON.parse(e.data)
};
```

### API

Server API (Go)

[![GoDoc](https://godoc.org/github.com/jpillora/velox?status.svg)](https://godoc.org/github.com/jpillora/velox)

Server API (Node)

* `velox.handle(object)` *function* returns `v` - Creates a new route handler for use with express
* `velox.state(object)` *function* returns `state` - Creates or restores a velox state from a given object
* `state.handle(req, res)` *function* returns `Promise` - Handle the provided `express` request/response. Resolves on connection close. Rejects on any error.

### Example

See this [simple `example/`](example/) and view it live here: https://velox.jpillora.com



*Here is a screenshot from this example page, showing the messages arriving as either a full replacement of the object or just a delta. The server will send which ever is smaller.*

### Notes

* JS object properties beginning with `$` will be ignored to play nice with Angular.
* JS object with an `$apply` function will automatically be called on each update to play nice with Angular.
* `velox.SyncHandler` is just a small wrapper around `velox.Sync`:

  ```go
  func SyncHandler(gostruct interface{}) http.Handler {
  	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
  		if conn, err := Sync(gostruct, w, r); err == nil {
  			conn.Wait()
  		}
  	})
  }
  ```

### Known issues

* Object synchronization is currently one way (server to client) only.
* Object diff has not been optimized. It is a simple property-by-property comparison.

### TODO

* WebRTC support
* Plain [`http`](https://nodejs.org/api/http.html#http_http_createserver_requestlistener) server support in Node

#### MIT License

Copyright © 2017 Jaime Pillora &lt;dev@jpillora.com&gt;

Permission is hereby granted, free of charge, to any person obtaining
a copy of this software and associated documentation files (the
'Software'), to deal in the Software without restriction, including
without limitation the rights to use, copy, modify, merge, publish,
distribute, sublicense, and/or sell copies of the Software, and to
permit persons to whom the Software is furnished to do so, subject to
the following conditions:

The above copyright notice and this permission notice shall be
included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED 'AS IS', WITHOUT WARRANTY OF ANY KIND,
EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY
CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT,
TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE
SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
