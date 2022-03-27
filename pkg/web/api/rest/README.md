### When you write software, there are two <br>kinds of problems that you run into...

- Problems that stretch your fundamental knowledge <br>
  of how things work and as a result of solving them <br>
  you become one step closer to unlocking the secrets <br>
  to immortality and transcending beyond mere human <br>
  limitations.
- Exceedingly stupid typos that static analysis tools <br>
  can't be taught how to catch and thus dooms humans <br>
  to feel like they wasted so much time on something <br>
  so trivial.
- Off-by-one errors.

```go
// APIServer can be used to create API's
type APIServer struct {
base string
resc map[string]Resource
m *http.ServeMux
}

// NewAPIServer creates and returns a new server instance
func NewAPIServer(base string) *APIServer {
return &APIServer{
base: base,
resc: make(map[string]Resource),
m: http.NewServeMux(),
}
}

// RegisterResource registers a `Resource` with the server
func (srv *APIServer) RegisterResource(name string, re Resource) {
srv.m.Handle(re.Path(), re)
srv.resc[name] = re
}

// ServeHTTP is the APIServer's default handler
func (srv *APIServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
// check for options request
if r.Method == http.MethodOptions {
handleOptions(w, r)
return
}
rh, found := lookupResourceHandler(r)
if !found {
handleNotFound(w, r)
return
}
switch apiRequestType(r) {
case apiReturnAll:
rh.returnAll(w, r)
case apiReturnOne:
rh.returnOne(w, r)
case apiInsertOne:
rh.insertOne(w, r)
case apiUpdateOne:
rh.updateOne(w, r)
case apiDeleteOne:
rh.deleteOne(w, r)
default:
handleBadRequest(w, r)
}
}

// ServeHTTP is the APIServer's default handler
func (srv *APIServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
// lookup resource handler
rh, found := lookupResourceHandler(r)
if !found {
handleNotFound(w, r)
return
}
// call default handler of the resource handler
rh.ServeHTTP(w, r)
return
}

```

