package stackage

import (
	"encoding/json"
	"io"
	"log"
	"os"
)

var (
	stdout,
	stderr,
	devNull,
	sLogDefault,
	cLogDefault *log.Logger
)

type logLevel uint16
type logLevels uint16

const (
	logCalls  logLevel = 1 << iota // log func/meth calls and their returns
	logPolErr                      // log errors thrown by policies, cap, r/o
	logNCfg                        // log state changes to underlying nodeCfg
	logInt                         // log internals (very verbose!)

	logNone logLevel = 0 // log nothing
)

/*
SetDefaultConditionLogger is a package-level function that will define
which logging facility new instances of Condition or equivalent type
alias shall be assigned during initialization procedures.

Logging is available but is set to discard all events by default. Note
that enabling this will have no effect on instances already created.

An active logging subsystem within any given Condition shall supercede
this default package logger.

The following types/values are permitted:

  - string: `none`, `off`, `null`, `discard` will turn logging off
  - string: `stdout` will set basic STDOUT logging
  - string: `stderr` will set basic STDERR logging
  - int: 0 will turn logging off
  - int: 1 will set basic STDOUT logging
  - int: 2 will set basic STDERR logging
  - *log.Logger: user-defined *log.Logger instance will be set

Case is not significant in the string matching process.

Logging may also be set for individual Condition instances using the
SetLogger method. Similar semantics apply.
*/
func SetDefaultConditionLogger(logger any) {
	cLogDefault = resolveLogger(logger)
}

/*
SetDefaultStackLogger is a package-level function that will define
which logging facility new instances of Stack or equivalent type
alias shall be assigned during initialization procedures.

Logging is available but is set to discard all events by default.
Note that enabling this will have no effect on instances already
created.

An active logging subsystem within any given Stack shall supercede
this default package logger.

The following types/values are permitted:

  - string: `none`, `off`, `null`, `discard` will turn logging off
  - string: `stdout` will set basic STDOUT logging
  - string: `stderr` will set basic STDERR logging
  - int: 0 will turn logging off
  - int: 1 will set basic STDOUT logging
  - int: 2 will set basic STDERR logging
  - *log.Logger: user-defined *log.Logger instance will be set

Case is not significant in the string matching process.

Logging may also be set for individual Stack instances using the
SetLogger method. Similar semantics apply.
*/
func SetDefaultStackLogger(logger any) {
	sLogDefault = resolveLogger(logger)
}

func resolveLogger(logger any) (l *log.Logger) {
	switch tv := logger.(type) {
	case *log.Logger:
		l = tv
	case int:
		l = intResolveLogger(tv)
	case string:
		l = stringResolveLogger(tv)
	}

	// We need something to fallback to,
	// regardless of the user's logging
	// intentions; impose devNull if we
	// find ourselves with a *log.Logger
	// instance that is nil.
	if l == nil || logger == nil {
		l = devNull
	}

	return l
}

func intResolveLogger(logger int) (l *log.Logger) {
	switch logger {
	case 0:
		l = devNull
	case 1:
		l = stdout
	case 2:
		l = stderr
	}

	return
}

func stringResolveLogger(logger string) (l *log.Logger) {
	switch lc(logger) {
	case `none`, `off`, `null`, `discard`:
		l = devNull
	case `stderr`:
		l = stderr
	case `stdout`:
		l = stdout
	}

	return
}

func logDiscard(logger *log.Logger) bool {
	return logger.Writer() == io.Discard
}

/*
Message is an optional type for use when a user-supplied Message channel has
been initialized and provided to one (1) or more Stack or Condition instances.

Instances of this type shall contain diagnostic, error and debug information
pertaining to current operations of the given Stack or Condition instance.
*/
type Message struct {
	ID   string             `json:"id"`
	Msg  string             `json:"txt"`
	Tag  string             `json:"tag"`
	Type string             `json:"type"`
	Addr string             `json:"addr,omitempty"`
	Time string             `json:"time"` // YYYYMMDDhhmmss.nnnnnnnnn
	Len  int                `json:"len"`
	Cap  int                `json:"max_len"`
	Data map[string]string  `json:"data,omitempty"`
	PPol PresentationPolicy `json:"-"`
}

func (r *Message) setText(txt any) (ok bool) {
	switch tv := txt.(type) {
	case error:
		if tv == nil {
			break
		}

		r.Msg = tv.Error()
	case string:
		if len(tv) == 0 {
			break
		}

		r.Msg = tv
	}

	if ok = len(r.Msg) > 0; !ok {
		r.Msg = sprintf("Unidentified or zero debug payload (%T)", txt)
	}

	return
}

/*
String is a stringer method that returns the string representation of
the receiver instance. By default, this method returns JSON content.

Instances of this type are used for the logging subsystem only, and do
not serve any purpose elsewhere. Use of the logging system as a whole
is entirely optional.

Users may author their own stringer method by way of the PresentationPolicy
closure type and override the string representation procedure for instances
of this type (thus implementing any syntax or format they wish, i.e.: XML,
YAML, et al).
*/
func (r Message) String() string {
	if !r.Valid() {
		return ``
	} else if r.PPol != nil {
		return r.PPol()
	}

	b, err := json.Marshal(&r)
	if err != nil {
		return ``
	}

	var replacements [][]string = [][]string{
		{`\\u`, `\u`},
		{`<nil>`, `nil`},
		{` <= `, ` LE `},
		{` >= `, ` GE `},
		{` < `, ` LT `},
		{` > `, ` GT `},
		{`&&`, `SYMBOLIC_AND`},
		{`&`, `AMPERSAND`},
		{`||`, `SYMBOLIC_OR`},
	}

	var data string = string(b)
	for _, repl := range replacements {
		if str, err := uq(rplc(qt(data), repl[0], repl[1])); err == nil {
			data = string(json.RawMessage(str))
		}
	}

	return data
}

/*
Valid returns a Boolean value indicative of whether the receiver
is perceived to be valid.
*/
func (r Message) Valid() bool {
	return (r.Type != `UNKNOWN` &&
		len(r.Time) > 0 &&
		len(r.Msg) > 0 &&
		len(r.Tag) > 0)
}

func getLogID(elem string) (id string) {
	id = `[no_id]`
	if _id := elem; len(_id) > 0 {
		id = elem
	}
	return
}

func init() {
	stderr = log.New(os.Stderr, ``, 0)
	stdout = log.New(os.Stdout, ``, 0)
	devNull = log.New(io.Discard, ``, 0)
	sLogDefault = devNull
	cLogDefault = devNull
}
