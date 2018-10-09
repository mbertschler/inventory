package guiapi

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/mbertschler/blocks/html"
)

// DefaultHandler is an empty guiapi handler ready for use.
var DefaultHandler = Handler{
	Functions: map[string]Callable{},
}

// Replace is a helper function that returns a Result that
// replaces the element chosen by the selector with the passed Block.
func Replace(selector string, block html.Block) (*Result, error) {
	out, err := html.RenderString(block)
	if err != nil {
		return nil, err
	}
	ret := &Result{
		HTML: []HTMLUpdate{
			{
				Operation: HTMLReplace,
				Selector:  selector,
				Content:   out,
			},
		},
	}
	return ret, nil
}

// Redirect lets the browser navigate to a given path
func Redirect(path string) (*Result, error) {
	args, err := json.Marshal(path)
	if err != nil {
		return nil, err
	}
	ret := &Result{
		JS: []JSCall{
			{
				Name:      "redirect",
				Arguments: args,
			},
		},
	}
	return ret, nil
}

// ============================================
// Logic
// ============================================

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, "guiapi request needs to use the POST method")
		return
	}
	var req Request
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&req)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintln(w, err.Error())
		return
	}
	resp := h.Handle(&req)
	enc := json.NewEncoder(w)
	err = enc.Encode(resp)
	if err != nil {
		log.Println("encoding error:", err)
	}
}

func (h Handler) Handle(req *Request) *Response {
	var resp Response
	for _, action := range req.Actions {
		var res = Result{
			ID:   action.ID,
			Name: action.Name,
		}
		fn, ok := h.Functions[action.Name]
		if !ok {
			res.Error = &Error{
				Code:    "undefinedFunction",
				Message: fmt.Sprint(action.Name, " is not defined"),
			}

		} else {
			r, err := fn(action.Args)
			if err != nil {
				res.Error = &Error{
					Code:    "error",
					Message: err.Error(),
				}
			}
			if r != nil {
				res.HTML = r.HTML
				res.JS = r.JS
			}
		}
		resp.Results = append(resp.Results, res)
	}
	return &resp
}

// ============================================
// Types
// ============================================

// Request is the sent body of a GUI API call
type Request struct {
	Actions []Action
}

type Action struct {
	ID   int    `json:",omitempty"` // ID can be used from the client to identify responses
	Name string // Name of the action that is called
	// Args as object, gets parsed by the called function
	Args json.RawMessage `json:",omitempty"`
}

type Handler struct {
	Functions map[string]Callable
}

type Callable func(args json.RawMessage) (*Result, error)

// Response is the returned body of a GUI API call
type Response struct {
	Results []Result
}

type Result struct {
	ID    int          `json:",omitempty"` // ID from the calling action is returned
	Name  string       // Name of the action that was called
	Error *Error       `json:",omitempty"`
	HTML  []HTMLUpdate `json:",omitempty"` // DOM updates to apply
	JS    []JSCall     `json:",omitempty"` // JS calls to execute
}

type Error struct {
	Code    string
	Message string
}

type HTMLOp int8

const (
	HTMLReplace HTMLOp = iota + 1
	HTMLDelete
	HTMLAppend
	HTMLPrepend
)

type HTMLUpdate struct {
	Operation HTMLOp // how to apply this update
	Selector  string // jQuery style selector: #id .class
	Content   string `json:",omitempty"` // inner HTML
	// Init calls are executed after the HTML is added
	Init []JSCall `json:",omitempty"`
	// Destroy calls are executed before the HTML is removed
	Destroy []JSCall `json:",omitempty"`
}

type JSCall struct {
	Name string // name of the function to call
	// Arguments as object, gets encoded by the called function
	Arguments json.RawMessage `json:",omitempty"`
}
