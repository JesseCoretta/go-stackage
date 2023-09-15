package main

/*
netstack - demonstrate a basic RESTful JSON interface
to an aliased Stack type instance with GET, DELETE and
POST capabilities.

This demo aims to extol as many of the (relevant)
features and capabilities of this package in a manner
most people can relate to.
*/

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"time"

	"github.com/JesseCoretta/go-stackage"
)

/*
APIRoot is the absolute root context for all API-related
resources. Changing this will result in a change of the
"/api" path element in all URLs throughout this example.
*/
const APIRoot = `/api`

/*
APIDefault defines the default API branch/context. This
is used in redirects in the event the user requests the
APIRoot URL without specifying a branch.
*/
const APIDefault = APIRoot + `/v1`

/*
MyStack is a type alias for the stackage.Stack type.

For simplicity, the only methods extended through this
type pertain to the intentions of this demo. We have NOT
extended ALL of stackage.Stack's many methods. Had this
occurred, the demo would be significantly larger for no
good reason.

However, we are able to ACCESS those methods by simply
type-casting MyStack->Stack.  For example, to access the
Len method in a MyStack instance:

  var mystack MyStack // assume this has been populated with something
  fmt.Printf("Length of %T: %d", mystack, stackage.Stack(mystack).Len())
  // Output: Length of main.MyStack: 2

Note that when we cast mystack within Stack, we're accessing
the underlying embedded pointer instance. If we need to make
changes to the instance, we need not worry about "writing any
changes back" to the instance of MyStack -- the pointer data
will reflect the changes, regardless of what type the pointer
instance is enveloped within.
*/
type MyStack stackage.Stack

/*
MyMuX is a single container in which we shall
store our HTTP listener.

This is largely for convenience, and is hardly
the only option. One could conceivably put the
contents in their own Auxiliary key/value pair
and do it that way.
*/
type MyMuX struct{
	*http.Server
	*http.ServeMux
}

/*
updateInfo contains a payload intended to alter
the served stack structure in some way, for
instance 'remove', 'insert', and others.

It would be delivered via a command such as

  curl -X POST <scheme://addr> -d '<this_type_marshaled_as_json_bytes>'
*/
type updateInfo struct {
       Action  string	      `json:"action"`		 // the stackage method name, such as 'insert' (case not important)
       Method  string	      `json:"-"`		 // typically POST or DELETE, not set by user explicitly, not visible in JSON
       Parameters map[string]any `json:"parameters,omitempty"` // input params required by various stackage.Stack methods
}

/*
slice is the abstraction of any particular
slice value within a stack.
*/
type slice struct{
        Type     string	`json:"type"`
	ID       string `json:"identifier,omitempty"`
	Kind     string `json:"kind,omitempty"`
	Category string `json:"category,omitempty"`
	Capacity int    `json:"capacity,omitempty"`
	Length   int    `json:"length,omitempty"`
        Index    int	`json:"index"`
	Nesting  bool	`json:"nesting,omitempty"`
	ReadOnly bool	`json:"readonly,omitempty"`
	FIFO     bool   `json:"fifo,omitempty"`
        Value    any	`json:"value"`
}

/*
our main function executes all of the top-level
processing actions, variable assembly, etc., and
launches a listener. This is a microcosm for your
hypothetical app, in essence.
*/
func main() {
	log.SetFlags(0)
        stackage.SetDefaultStackLogLevel(`ALL`)
        stackage.SetDefaultStackLogger(`stdout`)

        stackage.SetDefaultConditionLogLevel(`ALL`)
        stackage.SetDefaultConditionLogger(`stdout`)

	// create signal channel to gracefully
	// receive CTRL+C (^C) sequences. Note
	// this is an UNBUFFERED channel that
	// will block on the first signal it
	// receives.
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)

	// Load our pre-designed stack (MakeStack)
	// and initialize (but not execute) our
	// HTTP ServeMux/Server instances, with
	// support for listening on ALL local
	// addresses via TCP/8080.
	stack := MakeStack().SetListener(":8080")
	go func() {
		// Execute our listener and feed
		// it our signal channel.
		stack.Listen(quit)
	}()

	// BLOCK: do NOT close out ...
	<-quit

	// ... until you receive a signal, at which
	// point, shutdown gracefully (Shutdown) by
	// feeding it a background Context instance.
	ctx, cancel := context.WithTimeout(context.Background(), 60 * time.Second)
	defer cancel()
	if err := stack.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}
}

/*
MakeStack is simply an assembly function for data
to serve. Note this data is quite illogical and
serves only to provide content for the sake of
content. Ideally a user would fill this with some
kind of meaningful data.

Note that we CAST every stackage.Stack instance
as MyStack. Strictly speaking, its not required. One
can easily just *use* the Stack type as opposed to
creating a custom derivative type. In this demo,
however, we are using a custom type so that we can
extend various demo-related methods for simplicity
(such as Shutdown, Listen, SetListener, et al).

We also CAST Stack instances so that we can access
its myriad methods without having to manually write
wrapper to extend them piecemeal.
*/
func MakeStack() MyStack {
	var basicValues []any = []any{
		'L', 'x', '#', '🤣', rune(59), '临', '亶', 'ヨ', '⊕',
	}

	// top level of structure (r) is
	// returned as MyStack instance.
	r := MyStack(stackage.List().SetID(`netstack`).Push(
		MyStack(stackage.Basic().SetID(`_random`).Push(basicValues...)),
		MyStack(stackage.And().SetID(`_addr`).Push(
			`blue`,
			`red`,
			`green`,
		)),
		MyStack(stackage.And().SetID(`_random`).Push(
			MyStack(stackage.Or().SetID(`_random`).Push(
				stackage.Cond(`name`,stackage.Eq,`Mike`),
				stackage.Cond(`name`,stackage.Eq,`Courtney`),
			)),
			MyStack(stackage.Not().SetID(`negations`).Push(
				MyStack(stackage.Or(4).SetID(`or_negation`).Push(
					stackage.Cond(`terminated`,stackage.Eq,true),
					stackage.Cond(`disgruntled`,stackage.Eq,true),
				)),
			)),
		)),
	))

	return r
}

/*
Listen is a custom method extended through instances of the MyStack
type.

Listen shall launch the underlying *http.Server instance's listener
method (ListenAndServe) and serve the API as defined. An error is
returned when the process terminates or is terminated. The input
channel instance (stop) shall allow graceful intercept of CTRL+C
(^C) input signals, allowing the Server thread to *properly* close
down, as opposed to being clobbered into oblivion.
*/
func (r MyStack) Listen(stop chan os.Signal) (err error) {
	srv := r.Server()
	log.Printf("%s: listening on %s ... \n",
		stackage.Stack(r).ID(), srv.Addr)

	return r.Server().ListenAndServe()
}

/*
Shutdown gracefully terminates the underlying *http.ServeMux and
*http.Server listener system, returning any error resulting from
the attempt.

The context.Context input argument (ctx) is provided by the caller
of this method, which is the main() function in this demo.
*/
func (r MyStack) Shutdown(ctx context.Context) (err error) {
	log.Printf("\nReceived termination signal")

        srv := r.Server()
	if err = srv.Shutdown(ctx); err != nil {
		return
	}

	log.Printf("%s: stopping listener on %s...\n",
		stackage.Stack(r).ID(), srv.Addr)
	return
}

/*
SetListener assembles and configures the listener in
anticipation of serving content. This method does not
actually launch the listener.

Note that this method only impacts the receiver instance
directly. It will not create similar listener instances
in any nested stack or stack aliases found within the
receiver instance. Only the "top-level" or "outermost"
stack instance requires this setup procedure.
*/
func (r MyStack) SetListener(addr ...string) MyStack {
        var address string = ":8080"
        if len(addr) > 0 {
                address = addr[0]
        }

	// Assembly of our multiplexer
	// and http server
	mux := http.NewServeMux()
	srv := &http.Server{Addr: address, Handler: mux}
	ls := MyMuX{srv,mux}

	// Pass the receiver (MyStack) into the Configure
	// method. This creates routes and handlers used
	// for traversal and navigation within the stack
	// via the JSON RESTful interface.
	ls.Configure(r)

	// SetAuxiliary will initialize the underlying map
	// instance (Auxiliary) for administrative use. In
	// particular, it is needed to store the listener
	// system assembled above ...
	stackage.Stack(r).		// cast MyStack->Stack
		SetAuxiliary().		// init map (only needed once)
		Auxiliary().		// call map
		Set(`listener`, ls)	// set listener instance (ls) as map value identified by key `listener`

	return r
}

/*
Handler is a convenience method that accesses the underlying admin
map (Auxiliary) and calls the listener key/value pair. In the case
of this method, the process is used to return the instance of the
*http.ServeMux, which was assembled within the SetListener method.
*/
func (r MyStack) Handler() *http.ServeMux {
	val, found := stackage.Stack(r).Auxiliary().Get(`listener`)
	if !found {
		return nil
	}
	return val.(MyMuX).ServeMux
}

/*
Server is a convenience method that accesses the underlying admin
map (Auxiliary) and calls the listener key/value pair. In the case
of this method, the process is used to return the instance of the
*http.Server, which was assembled within the SetListener method.
*/
func (r MyStack) Server() *http.Server {
        val, found := stackage.Stack(r).Auxiliary().Get(`listener`)
        if !found {
                return nil
        }
        return val.(MyMuX).Server
}

/*
Configure assigns redirects and handlers as needed within the receiver.

Users would likely want to significantly customize this, adding
their own paths, etc., as needed to suit their data structure.
*/
func (r *MyMuX) Configure(stack MyStack) MyMuX {

	// Begin setting up our fundamental handlers
        r.ServeMux.HandleFunc("/", rootHandler)					// absolute root of listener "/"
        r.ServeMux.Handle(APIRoot,    http.RedirectHandler(APIDefault,302))	// "/api"
        r.ServeMux.Handle(APIDefault, http.HandlerFunc(v1Handler))		// "/api/v1"

	// This is an unimplemented handler, always
	// returns HTTP 501: Not Implemented.
	apiV2 := APIRoot + `/v2`
        r.ServeMux.Handle(apiV2, http.HandlerFunc(v2Handler))			// "/api/v2" (dead-end)

	// "/slice" by itself brings the user to
	// the main landing page for the slice,
	// which contains a list of available
	// slice indices.
        slicesPath := APIDefault + "/slice"
        r.ServeMux.HandleFunc(slicesPath, stack.slicesHandler)

	// "/slice/<N>" (e.g.: "/slice/1") will
	// bring the user to the page for the
	// specific slice(s) they have called.
	// This accepts traversal paths as well
	// as single indices. Some examples:
	//
	//   - /slice/1		"just return slice index 1"
	//   - /slice/1/0/2/1	"hierarchically traverse slice indices 1->0->2 and call slice index 1"
	//
	// StripPrefix is used so the handler
	// function can more easily split the
	// path into discrete indices as shown
	// in the above example.
        r.ServeMux.Handle(slicesPath + `/`, http.StripPrefix(
		slicesPath + `/`,
		http.HandlerFunc(stack.sliceHandler),
	))

	return *r
}

/*
ServeHTTP is used by net/http to serve content via the mux. It need
not contain any actual code, only the signature is needed.
*/
func (r MyMuX) ServeHTTP(http.ResponseWriter, *http.Request) {}

/*
slicesHandler handles requests pertaining to the root slice index resource,
in which a specific index is not requested. Instead, a manifest of available
slice contexts is presented, allowing the visitor to select (choose) one or
more slice indices to call.

This handler allows GET, DELETE and POST methods.  GET shall render the manifest
object, while POST shall allow the manipulation of multiple indices from the top-level
(i.e.: remove, insert, et al). DELETE shall remove an indicated slice altogether,
using either pop or remove, each of which operate with subtle differences).
*/
func (r MyStack) slicesHandler(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Set("Content-Type", "application/json")

	// Make sure the client isn't requesting some
	// inappropriate HTTP method. Only support the
	// ones listed here, and bail on anything else.
	if !methodAllowed(req.Method, []string{
		http.MethodPost,
		http.MethodGet,
		http.MethodDelete,
	}) {
		http.Error(resp, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// We defined this type here because frankly it need
	// not be globally usable.
	type mainDirectory struct {
		ID       string `json:"identifier,omitempty"`
		Kind     string `json:"kind,omitempty"`
		Category string `json:"category,omitempty"`
		Length   int    `json:"len"`
		Capacity int    `json:"cap,omitempty"`
		ReadOnly bool   `json:"readonly,omitempty"`
		Nesting  bool   `json:"nesting,omitempty"`
		Slices []int    `json:"slices"`
	}

	// If method is POST or DELETE, just run the contents
	// of this section and bail out. Go no further!
	if req.Method == http.MethodDelete || req.Method == http.MethodPost {
		if err := r.deleteOrPost(resp, req); err != nil {
			// tell the ADMIN by console msg (optional)
			log.Println(err)
		}
		return
	}

	// type cast our custom stack instance
	// to access a few methods ...
	stk := stackage.Stack(r)

	// craft the landing page's contents
	// to guide the user from there ...
	directory := mainDirectory{
		ID:       stk.ID(),
		Kind:     stk.Kind(),
		Category: stk.Category(),
		Capacity: stk.Cap(),
		Length:   stk.Len(),
		ReadOnly: stk.IsReadOnly(),
                Nesting:  stk.IsNesting(),
	}

	if directory.Length > 0 {
		// don't alloc unless we have a reason
		// to do so ...
		directory.Slices = make([]int, 0)

		// iterate slices, and obtain information
		// from each to store within the directory.
		for i := 0; i < stk.Len(); i++{
			if _, ok := stk.Index(i); ok {
				directory.Slices = append(directory.Slices, i)
			}
		}
	}

	// marshal directory type instance into JSON bytes.
	j, err := json.MarshalIndent(directory,"", "\t")
	if err != nil {
		http.Error(resp, "Internal server error", 500)	// client-facing err - no sensitive details
		log.Println(err)				// admin (console) err
		return
	}

	// write successfully marshaled JSON bytes
	// to the http.ResponseWriter, thus sending
	// them to the client.
	fmt.Fprintf(resp, string(j))
	return
}

/*
deleteOrPost is a compartmentalized method executed by slicesHandler. It
is separated here just to keep any single method from being too large.

This method handles generalized delete or post operations relating to
data within the receiver.
*/
func (r MyStack) deleteOrPost(resp http.ResponseWriter, req *http.Request) (err error) {
	var info updateInfo
	if info, err = marshalUpdateInfo(req); err != nil {
	        http.Error(resp, "Invalid POST bytes", http.StatusBadRequest) // tell the USER
	        return
	}
        action := strings.ToLower(info.Action)

	metherr := fmt.Errorf("Inappropriate resource (%s) for method %s, or resource not implemented", action, req.Method)

        if req.Method == http.MethodDelete {
		// method is HTTP DELETE

		switch action {
		case `remove`,`pop`,`reset`:
			// reset requires a body parameter
			// within the updateInfo instance
			// for 'confirm', which is a bool.
			// This is required to trash the
			// entire stack. reset returns no
			// data, and is a 204 (no content).
			//
			// remove returns a 204 (no content)
			// with no body. remove requires a
			// body parameter within updateInfo
			// for 'idx', which is an integer
			// and targets a specific slice for
			// deletion.
			//
			// pop returns a 200 with the popped
			// element expressed as JSON. pop
			// requires no request body.
                        err = r.popResetOrRemoveRequest(resp, info)
		default:
			err = metherr				// admin-facing (console) error
                        http.Error(resp, err.Error(), 400)	// client-facing error (same error is fine here)
		}

        } else if req.Method == http.MethodPost {
		// method is HTTP POST
		//
                // Now we can choose the appropriate modifier
                // function to invoke for the request. Many of
                // these actions are direct links to similarly
                // named stackage.Stack methods.
                switch action {
                case `insert`:
                        err = r.insertRequest(resp,info)
                case `replace`:
                        err = r.replaceRequest(resp,info)
                case `push`:
                        err = r.pushRequest(resp,info)

                default:
			err = metherr				// admin-facing (console) error
                        http.Error(resp, err.Error(), 400)	// client-facing error (same error is fine here)
                }
        }

	return
}

/*
sliceHandler shall handle requests pertaining to the calling of a specific slice
index, either direct (e.g.: /0) or via traversal path (e.g.: /0/1/2/1/1/0/...).
*/
func (r MyStack) sliceHandler(resp http.ResponseWriter, req *http.Request) {
	resp.Header().Set("Content-Type", "application/json")

	// Make sure the client isn't requesting some
	// inappropriate HTTP method. Only support the
	// ones listed here, and bail on anything else.
	if !methodAllowed(req.Method, []string{
		http.MethodPost,
		http.MethodGet,
		http.MethodDelete,
	}) {
		http.Error(resp, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	path := make([]int,0)
	raw := strings.TrimSuffix(req.URL.Path,`/`)
	indices := strings.Split(raw, `/`)
	for i := 0 ; i < len(indices); i++ {
		idx, err := strconv.Atoi(indices[i])
		if err != nil {
			http.Error(resp, "Unavailable", 503)
			return
		}
		path = append(path,idx)
	}

	var sl any
	var ok bool
	var idx int

	// perform a path length switch,
	// handling the result in one (1)
	// of three (3) different ways.
	switch len(path) {
	case 0:
		// no result: get lost
		http.Error(resp, "Not Found", 400)
		return
	case 1:
		// one (1) path element indicates we
		// can just call the Index method and
		// be done with it.
		if sl, ok = stackage.Stack(r).Index(path[0]); !ok || sl == nil {
			http.Error(resp, "Not Found", 400)
			return
		}

		if req.Method != http.MethodGet {
			if err := r.deleteOrPost(resp, req); err != nil {
				http.Error(resp, "Internal Server Error", 500)
			}
			return
			// GO NO FURTHER IF NOT 'GET' METHOD
		}
		idx = path[0]
	default:
		// >1 path elements indicates Traversal needed ...
		if sl, ok = stackage.Stack(r).Traverse(path...); !ok || sl == nil {
			http.Error(resp, "Not Found", 400)
			return
		}

		if req.Method != http.MethodGet {

			// Make sure the slice we received back is
			// something we can traverse.
                        var stack MyStack
                        switch tv := sl.(type) {
                        case stackage.Stack:
                                stack = MyStack(tv)
                        case MyStack:
                                stack = tv
                        default:
                                http.Error(resp, "Not traversable, cannot remove from target", 400)
                                return
                        }

			// execute deleteOrPost method using the newly asserted
			// MyStack instance (stack), and NOT the receiver (r),
			// else we'd trash the wrong data!
			if err := stack.deleteOrPost(resp, req); err != nil {
				http.Error(resp, "Internal Server Error", 500)
			}
			return
			// GO NO FURTHER IF NOT 'GET' METHOD
		}

		idx = path[len(path)-1]
	}

	// marshal sliceInfo instance into JSON bytes.
        j, err := json.MarshalIndent(sliceInfo(sl,idx), "", "\t")
        if err != nil {
                http.Error(resp, "Internal server error", 500)	// user-facing error
		fmt.Println(err)				// admin-facing (console) err
                return
        }

	// write successfully marshaled JSON bytes
	// to the http.ResponseWriter, thus sending
	// them to the client.
	fmt.Fprintf(resp, string(j))
	return
}

/*
v1Handler returns the root /api/<default_version> resource, which in this
this case is a JSON payload of usage information and structural schemata.
*/
func v1Handler(resp http.ResponseWriter, req *http.Request) {
	// Only GET allowed for this resource.
	if req.Method != http.MethodGet {
                http.Error(resp, "Method not allowed", 405)
		return
	}

	// If URL does not begin with /api, get bent.
	if !strings.HasPrefix(req.URL.Path, "/api") {
		// SUGGESTIONS/IDEAS/ALTERNATIVES:
		// - more helpful info written to resp?
		// - redirect to help page?
		http.NotFound(resp, req)
		return
	}

	// becomes: /api/v1
	apiRoot := APIDefault

	help := struct{
		APIRoot string `json:"api_root"`
		Methods []string `json:"methods_available"`
		Resources map[string]any `json:"resources"`
	}{
		APIRoot: apiRoot,
		Methods: []string{http.MethodGet},
		Resources: map[string]any{
			"slice": struct{
				Path string `json:"path"`
				Description string `json:"description"`
			}{
				Path: apiRoot + "/slice",
				Description: "Main index of all slice indices present",
			},
		},
	}

        j, err := json.MarshalIndent(help, "", "\t")
        if err != nil {
                http.Error(resp, "Internal server error", 500)
                return
        }

	fmt.Fprintf(resp, string(j))
}

/*
v2Handler is a bogus parallel API branch that is present merely for users
to customize into something usable, or to act as a dead-end alternative to
the v1Handler, thereby allowing the testing of an APIDefault fallback.

In its current state, this will return a 'Not Implemented' (501) when used
properly.
*/
func v2Handler(resp http.ResponseWriter, req *http.Request) {
        // Only GET allowed for this resource.
        if req.Method != http.MethodGet {
                http.Error(resp, "Method not allowed", 405)
                return
        }

        // If URL does not begin with /api, get bent.
        if !strings.HasPrefix(req.URL.Path, "/api") {
                // SUGGESTIONS/IDEAS/ALTERNATIVES:
                // - more helpful info written to resp?
                // - redirect to help page?
                http.NotFound(resp, req)
                return
        }

        http.Error(resp, "Resource not implemented", 501)
	return
}

/*
rootHandler is the handlerfunc instance for any request which
lands upon the root (/) context.
*/
func rootHandler(resp http.ResponseWriter, req *http.Request) {
        resp.Header().Set("Content-Type", "text/html")

        if req.URL.Path != "/" {
                http.NotFound(resp, req)
                return
        }

        // Here we write basic HTML to the responsewriter,
        // to (cleanly) welcome clients who are visiting
        // the resource with a GUI browser as opposed to
        // curl/wget...
        fmt.Fprintf(resp,`
                <html>
                <head>
                        <h2 style="font-family:monospace";>Welcome to the Netstack API Demo</h2>
                </head>
                <body>
                        <p style="font-family:monospace";>coming soon! maybe! See the <a href="/api">API</a> section for access to stack data</p>
                </body>
                </html>
        `)
}

/*
valueHandler will take appropriate steps to ensure a value
is string represented properly and is JSON friendly ...
*/
func stackValueHandler(x any) (slices []slice) {
	switch tv := x.(type) {
	case stackage.Stack:
		slices = make([]slice, tv.Len())
		for i := 0; i < tv.Len(); i++ {
			sl, _ := tv.Index(i)
			switch uv := sl.(type) {
			case stackage.Stack:
				// cast stackage.Stack to main.MyStack
				slices[i] = sliceInfo(MyStack(uv),i)
			default:
				slices[i] = sliceInfo(uv,i)
			}
		}

	case MyStack:
		return stackValueHandler(stackage.Stack(tv))
	}

	return
}

/*
sliceInfo attempts to marshal x -- as an explicit value 
and expressed as "meta data" -- into an instance of slice.
*/
func sliceInfo(x any, idx int) (s slice) {
        s.Index = idx

        switch tv := x.(type) {
	case MyStack:
		// cast + self-execute for simplicity
		s = sliceInfo(stackage.Stack(tv), idx)
		s.Type = fmt.Sprintf("%T", tv)

        case stackage.Stack:
                s.Value = stackValueHandler(tv)
		s.Type = fmt.Sprintf("%T", tv)
		s.Kind = tv.Kind()
		s.Length = tv.Len()
		s.Capacity = tv.Cap()
                s.Nesting =  tv.IsNesting()
                s.ReadOnly = tv.IsReadOnly()
		s.FIFO = tv.IsFIFO()
		s.ID = tv.ID()
		s.Category = tv.Category()

	case rune:
		s.Type = fmt.Sprintf("%T", tv)
		s.Value = fmt.Sprintf("%c", tv)

	case int, int8, int16, int64,
		uint, uint8, uint16, uint32, uint64,
		float32, float64:
		s.Value = tv
		s.Type = fmt.Sprintf("%T", tv)

	case string:
		s.Type = fmt.Sprintf("%T", tv)
		s.Length = len(tv)
		s.Value = tv

	case []any:
		s.Type = fmt.Sprintf("%T", tv)
		s.Length = len(tv)
		s.Capacity = cap(tv)
		s.Value = tv

        default:
		s.Type = fmt.Sprintf("%T", tv)
                s.Value = fmt.Sprintf("%v", tv)

        }

	return
}

/*
methodAllowed scans the available methods to determine whether meth
is among them. Case is not significant in the matching process.
*/
func methodAllowed(meth string, methods []string) bool {
	for i := 0; i < len(methods); i++ {
		if strings.EqualFold(meth,methods[i]) {
			return true
		}
	}
	return false
}

/*
marshalUpdateInfo will attempt to marshal the request body
bytes into an instance of updateInfo, which can be used by
eligible handler funcs to effect changes upon the stack
and/or its contents.
*/
func marshalUpdateInfo(req *http.Request) (info updateInfo, err error) {
	// read resp. body (don't need
	// to close manually though).
	var body []byte
	if body, err = io.ReadAll(req.Body); err != nil {
	        return
	}

	// define unmarshal destination
	// and unmarshal content into it.
	// Make a note of the method as
	// well ...
	info.Method = http.MethodPost
	err = json.Unmarshal(body, &info)
	return
}

/*
insertRequest wraps stackage.Stack's Insert method.

See pkg.go.dev/github.com/JesseCoretta/go-stackage#Stack.Insert
for signature and other details.
*/
func (r MyStack) insertRequest(resp http.ResponseWriter, info updateInfo) (err error) {
	raw, found := info.Parameters[`left`]
	if !found {
		err = fmt.Errorf("Missing 'left' key parameter")
		return
	}

	assert, ok := raw.(float64)
	if !ok {
		err = fmt.Errorf("Bad request; 'left' body parameter NaN")
		http.Error(resp, err.Error(), 400)
		return
	}
	left := int(assert)

	value, found := info.Parameters[`value`]
	if !found {
		err = fmt.Errorf("Missing 'value' key parameter")
		http.Error(resp, err.Error(), 400)
		return
	}

	switch tv := value.(type) {
	case []any:
                stackage.Stack(r).Insert(revEngStack(tv),left)
	default:
                stackage.Stack(r).Insert(tv,left)
	}

	return
}

/*
replaceRequest wraps stackage.Stack's Replace method.

See pkg.go.dev/github.com/JesseCoretta/go-stackage#Stack.Replace
for signature and other details.
*/
func (r MyStack) replaceRequest(resp http.ResponseWriter, info updateInfo) (err error) {
        idx, found := info.Parameters[`idx`]
        if !found {
                err = fmt.Errorf("Missing 'idx' key parameter")
                return
        }

	assert, ok := idx.(float64)
	if !ok {
		http.Error(resp, "Bad request; 'idx' body parameter NaN", 400)
		err = fmt.Errorf("'idx' parameter must be an integer")
		return
	}
        index := int(assert)

        value, found := info.Parameters[`value`]
        if !found {
                err = fmt.Errorf("Missing 'value' key parameter")
                return
        }

        switch tv := value.(type) {
        case []any:
                stackage.Stack(r).Replace(revEngStack(tv),index)

	default:
                stackage.Stack(r).Replace(tv,index)
        }

        return
}

/*
pushRequest wraps stackage.Stack's Push method.

See pkg.go.dev/github.com/JesseCoretta/go-stackage#Stack.Push
for signature and other details.
*/
func (r MyStack) pushRequest(resp http.ResponseWriter, info updateInfo) (err error) {
        value, found := info.Parameters[`value`]
        if !found {
                err = fmt.Errorf("Missing 'value' key parameter")
                return
        }

        switch tv := value.(type) {
        case []any:
		if len(tv) > 0 {
			stackage.Stack(r).Push(revEngStack(tv))
			resp.WriteHeader(http.StatusNoContent)	// 204
			return
		}

        default:
		if tv != nil {
			stackage.Stack(r).Push(tv)
			resp.WriteHeader(http.StatusNoContent)	// 204
			return
		}
        }

	http.Error(resp, "Bad request; unknown push type, or empty payload", 400)

        return
}

/*
revEngStack attempts to reverse-engineer a stack based on the state and
composition of an instance of []any (slices of interface) that was sent
by the user via a POST operation containing a JSON payload.

Patterns are as follows.

  - Index 0 with a string value of `or` shall execute stackage.Or
  - Index 0 with a string value of `and` shall execute stackage.And
  - Index 0 with a string value of `not` shall execute stackage.Not
  - Index 0 with a string value of `list` shall execute stackage.List

If no desired pattern is discerned, a fallback stack type initialized by
stackage.Basic shall be used.

Once a stack instance is initialized, the contents of x (slices #1 and up)
are pushed into it using the Push method.
*/
func revEngStack(x []any) (stack any) {
        word, _ := x[0].(string) // fallback to basic if assertion fails

        switch strings.ToLower(word) {
        case `or`, `and`, `not`, `list`:
                stack = stackFromBooleanOperator(word)
        default:
                stack = stackFromBooleanOperator(`basic`)
        }

	for i := 1; i < len(x); i++ {
		stackage.Stack(stack.(MyStack)).Push(x[i])
	}

	return
}

/*
popResetOrRemoveRequest handles any stackage.Reset/Pop/Remove request
made by the user via JSON instruction payload.
*/
func (r MyStack) popResetOrRemoveRequest(resp http.ResponseWriter, info updateInfo) (err error) {
	switch strings.ToLower(info.Action) {
	case `reset`:
                target, found := info.Parameters[`confirm`]
                if !found {
                        http.Error(resp, "Bad request; missing 'confirm' body parameter", 400)
                        err = fmt.Errorf("Missing 'confirm' key parameter; this is required for 'reset'")
                        break
                }
		assert, ok := target.(bool)
		if !ok {
			http.Error(resp, "Bad request; 'confirm' body parameter not a boolean", 400)
			err = fmt.Errorf("'confirm' body parameter must be a bool")
			break
		}

		if assert {
			stackage.Stack(r).Reset()
		}
		resp.WriteHeader(http.StatusNoContent)	// 204

	case `pop`:
		var idx int = stackage.Stack(r).Len()-1
		if stackage.Stack(r).IsFIFO() {
			idx = 0
		}

		popped, ok := stackage.Stack(r).Pop()
		if !ok {
			err = fmt.Errorf("Failed to pop element from %T", r)
			break
		}

		var j []byte
		if j, err = json.MarshalIndent(sliceInfo(popped,idx), "", "\t"); err != nil {
			http.Error(resp, "Internal server error", 500)
			break
		}

		fmt.Fprintln(resp, string(j))

	case `remove`:
	        target, found := info.Parameters[`idx`]
	        if !found {
			http.Error(resp, "Bad request; missing 'idx' body parameter", 400)
	                err = fmt.Errorf("Missing 'idx' key parameter; this is required for 'remove'")
	                break
	        }

		assert, ok := target.(float64)
		if !ok {
			http.Error(resp, "Bad request; 'idx' body parameter NaN", 400)
			err = fmt.Errorf("'idx' parameter must be an integer")
			break
		}

		log.Printf("---> %s", stackage.Stack(r).ID())
		stackage.Stack(r).Remove(int(assert))
		resp.WriteHeader(http.StatusNoContent)	// 204
	}

        return
}

/*
stackFromBooleanOperator returns a ready-to-use instance of
stackage.Stack -- initialized by the appropriate function --
cast as an instance of MyStack.

The fallback stack type used for invalid requests is Basic.
*/
func stackFromBooleanOperator(word string) (stack MyStack) {
	switch strings.ToLower(word) {
	case `or`:
	        stack = MyStack(stackage.Or())
	case `and`:
	        stack = MyStack(stackage.And())
	case `not`:
	        stack = MyStack(stackage.Not())
	case `list`:
	        stack = MyStack(stackage.List())
	default:
	        stack = MyStack(stackage.Basic())
	}

	return
}

