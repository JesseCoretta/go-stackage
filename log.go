package stackage

import (
	"encoding/json"
	"io"
	"log"
	"os"
)

var (
	// io.Writer qualifiers for event dispatch
	stdout,
	stderr,
	devNull,

	// actual (default) logger instances, can be
	// supplanted by users using SetDefault[X]Logger
	sLogDefault,
	cLogDefault *log.Logger

	// actual (default) loglevel instances, can be
	// set by users using SetDefault[X]LogLevel
	sLogLevelDefault,
	cLogLevelDefault LogLevel
)

type logSystem struct {
	lvl logLevels
	log *log.Logger
}

/*
LogLevel is a uint16 type alias used to define compound logging
verbosity configuration values.

LogLevel consts zero (0) through four (4) are as follows:

  - NoLogLevels defines the lack of any logging level
  - LogLevel1 defines basic function and method call event logging
  - LogLevel2 defines the logging of events relating to configuration state changes
  - LogLevel3
  - LogLevel4
*/
type LogLevel uint16

const NoLogLevels LogLevel = 0 // silent

type logLevels uint16

var (
	logLevelNames map[LogLevel]string
	logLevelMap   map[string]LogLevel
)

const (
	LogLevel1      LogLevel = 1 << iota //     1 :: builtin: calls
	LogLevel2                           //     2 :: builtin: policy
	LogLevel3                           //     4 :: builtin: state
	LogLevel4                           //     8 :: builtin: debug
	LogLevel5                           //    16 :: builtin: errors
	LogLevel6                           //    32 :: builtin: trace
	UserLogLevel1                       //    64 :: user-defined
	UserLogLevel2                       //   128 :: user-defined
	UserLogLevel3                       //   256 :: user-defined
	UserLogLevel4                       //   512 :: user-defined
	UserLogLevel5                       //  1024 :: user-defined
	UserLogLevel6                       //  2048 :: user-defined
	UserLogLevel7                       //  4096 :: user-defined
	UserLogLevel8                       //  8192 :: user-defined
	UserLogLevel9                       // 16384 :: user-defined
	UserLogLevel10                      // 32768 :: user-defined

	AllLogLevels LogLevel = ^LogLevel(0) // 65535 :: log all of the above unconditionally (!!)
)

func (r logLevels) String() string {
	if int(r) == int(AllLogLevels) {
		return `ALL`
	} else if int(r) == 0 {
		return `NONE`
	}

	var levels []string
	for i := 0; i < 16; i++ {
		lvl := LogLevel(1 << i)
		if r.positive(lvl) {
			if name, found := logLevelNames[lvl]; found {
				levels = append(levels, name)
			}
		}
	}

	return join(levels, `,`)
}

/*
shift shall left-shift the bit value of the receiver to include
the addition of one (1) or more LogLevel specifiers (l) in variadic
fashion. Type instances that qualify are as follows:

  - string (case-insensitive loglevel "name", e.g.: "trace" or "none")
  - LogLevel (actual LogLevel constants, e.g.: NoLoglevels or LogLevel3)

If any of l's values are NoLogLevels, the receiver shall be set to
zero (0) and any remaining shifts shall be discarded. In this context,
"shift zero" translates to "log nothing".

Conversely, if any of l's values are LogLevel16, the receiver shall be
set to ^LogLevel(0) (uint16(65535)) and any remaining shifts shall be
discarded. Predictably, "shift 65535" translates to "log everything".
*/
func (r *logSystem) shift(l ...any) *logSystem {
	r.lvl.shift(l...)
	return r
}

/*
shift is a private method called by instances of the
logSystem struct type, et al.
*/
func (r *logLevels) shift(l ...any) *logLevels {
	for i := 0; i < len(l); i++ {
		var ll LogLevel // current iteration's resolved loglevel stored here
		var ok bool

		// Perform type-switch upon the
		// currently iterated 'any' (#i)
		switch tv := l[i].(type) {

		case string:
			// value could be a loglevel NAME. Try to
			// resolve it, and pay no regard to case
			// folding.
			ll, ok = logLevelMap[uc(tv)]

		case LogLevel:
			// value is a LogLevel instance. Just take
			// it at face value.
			ll = tv
			ok = true
		case int:
			ll = LogLevel(tv)
			ok = true
		}

		if logLevels(ll) == logLevels(0) {
			*r = logLevels(NoLogLevels)
			break
		} else if logLevels(ll) == ^logLevels(0) {
			*r = logLevels(AllLogLevels)
			break
		}

		// Loglevel is neither "all" nor "none",
		// meaning is a singular, discrete log
		// verbosity specifier; shift it into
		// current value, don't clobber.
		if ok {
			*r |= logLevels(ll)
		}
	}

	return r
}

/*
unshift shall right-shift the bit value of the receiver to effect
the removal of one (1) or more logLevel instances (l) in variadic
fashion.

If any of l's values are NoLogLevels, the loop shall call continue, as
nothing can be done logically with that value here, though it is
not fatal, nor should it terminate processing.

If any of l's values are LogLevel16, the receiver shall be set to zero
(0) and any remaining shifts shall be discarded, as "unshift all"
in this context translates to "log nothing".
*/
func (r *logSystem) unshift(l ...any) *logSystem {
	r.lvl.unshift(l...)
	return r
}

func (r *logLevels) unshift(l ...any) *logLevels {
	for i := 0; i < len(l); i++ {
		var ll LogLevel
		var ok bool

		switch tv := l[i].(type) {
		case string:
			// value could be a loglevel NAME. Try to
			// resolve it, and pay no regard to case
			// folding.
			ll, ok = logLevelMap[uc(tv)]

		case LogLevel:
			// value is a LogLevel instance. Just take
			// it at face value.
			ll = tv
			ok = true
		case int:
			ll = LogLevel(tv)
			ok = true
		}

		if logLevels(ll) == logLevels(0) {
			continue
		} else if logLevels(ll) == ^logLevels(0) {
			break
		}

		if ok {
			*r = (*r &^ logLevels(ll))
		}
	}
	return r
}

/*
positive returns a Boolean value indicative of whether the receiver
contains the bit value for the specified logLevel. In context, this
means "logLevel <X>" is either active (true) or not (false).
*/
func (r logSystem) positive(l any) (posi bool) {
	if !r.isZero() {
		posi = r.lvl.positive(l)
	}
	return
}

func (r logLevels) positive(l any) (posi bool) {
	if r == logLevels(0) {
		return
	} else if r == ^logLevels(0) {
		posi = true
		return
	}

	var ll LogLevel
	var ok bool
	switch tv := l.(type) {
	case LogLevel:
		// value is a LogLevel instance. Just take
		// it at face value.
		ll = tv
		ok = true
	}

	if ok {
		posi = (r & logLevels(ll)) != 0
	}

	return
}

func (r *logSystem) isZero() (is bool) {
	if r != nil {
		is = r.log == nil && r.lvl == logLevels(NoLogLevels)
	}
	return
}

func (r logSystem) logger() (l *log.Logger) {
	if !r.isZero() {
		l = r.log
	}
	return
}

func (r *logSystem) setLogger(logger any) *logSystem {
	r.log = resolveLogger(logger)
	return r
}

func newLogSystem(logger any, l ...any) (lsys *logSystem) {
	lgr := devNull
	if logger != nil {
		lgr = resolveLogger(logger)
	}
	lsys = new(logSystem)
	lsys.log = lgr
	lsys.lvl = *new(logLevels).shift(l...)

	return
}

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

func DefaultConditionLogLevel() int {
	return int(cLogLevelDefault)
}

/*
SetDefaultConditionLogLevel sets the instance of LogLevel (lvl)
as a LITERAL value, as the verbosity indicator. When set with
appropriate level identifiers, this will increase or decrease
log verbosity accordingly. This value shall be used for logging
verbosity (or lack thereof) for any newly created (and qualified)
instances.

Note that the input value(s) are NOT shifted. Users are expected
to either sum the values and cast the product as a LogLevel, OR
settle for one of the predefined LogLevel constants.

The default is NoLogLevels, which implies a loglevel of zero (0).
*/
func SetDefaultConditionLogLevel(lvl any) {
	var level LogLevel
	switch tv := lvl.(type) {
	case string:
		ll, found := logLevelMap[uc(tv)]
		if found {
			level = ll
		}
	case int:
		if 0 <= tv && tv <= int(^uint16(0)) {
			level = LogLevel(tv)
		}
	case LogLevel:
		level = tv
	default:
		level = NoLogLevels
	}

	cLogLevelDefault = level
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

func DefaultStackLogLevel() int {
	return int(sLogLevelDefault)
}

/*
SetDefaultStackLogLevel sets the instance of LogLevel (lvl)
as a LITERAL value, as the verbosity indicator. When set with
appropriate level identifiers, this will increase or decrease
log verbosity accordingly. This value shall be used for logging
verbosity (or lack thereof) for any newly created (and qualified)
instances.

Note that the input value(s) are NOT shifted. Users are expected
to either sum the values and cast the product as a LogLevel, OR
settle for one of the predefined LogLevel constants.

The default is NoLogLevels, which implies a loglevel of zero
(0).
*/
func SetDefaultStackLogLevel(lvl any) {
	var level LogLevel
	switch tv := lvl.(type) {
	case string:
		ll, found := logLevelMap[uc(tv)]
		if found {
			level = ll
		}
	case int:
		if 0 <= tv && tv <= int(^uint16(0)) {
			level = LogLevel(tv)
		}
	case LogLevel:
		level = tv
	default:
		level = NoLogLevels
	}

	sLogLevelDefault = level
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
	r.Msg = sprintf("Unidentified or zero debug payload (%T)", txt)

	switch tv := txt.(type) {
	case error:
		if ok = tv != nil; ok {
			r.Msg = tv.Error()
		}
	case string:
		if ok = len(tv) != 0; ok {
			r.Msg = tv
		}
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
func (r Message) String() (data string) {
	if r.PPol != nil {
		data = r.PPol()
		return
	}

	if r.Valid() {
		if b, err := json.Marshal(&r); err == nil {
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

			data = string(b)
			for _, repl := range replacements {
				if str, err := uq(rplc(qt(data), repl[0], repl[1])); err == nil {
					data = string(json.RawMessage(str))
				}
			}
		}
	}

	return
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
	id = `no_identifier`
	if _id := elem; len(_id) > 0 {
		id = elem
	}
	return
}

func init() {
	stderr = log.New(os.Stderr, ``, 0)
	stdout = log.New(os.Stdout, ``, 0)
	devNull = log.New(io.Discard, ``, 0)

	// log events dispatched into oblivion
	sLogDefault = devNull
	cLogDefault = devNull

	// log events are ignored regardless of severity/category
	sLogLevelDefault = NoLogLevels
	cLogLevelDefault = NoLogLevels

	// instance->str
	logLevelNames = map[LogLevel]string{
		NoLogLevels:    `NONE`,
		LogLevel1:      `CALLS`,
		LogLevel2:      `POLICY`,
		LogLevel3:      `STATE`,
		LogLevel4:      `DEBUG`,
		LogLevel5:      `ERROR`,
		LogLevel6:      `TRACE`,
		UserLogLevel1:  `USER1`,
		UserLogLevel2:  `USER2`,
		UserLogLevel3:  `USER3`,
		UserLogLevel4:  `USER4`,
		UserLogLevel5:  `USER5`,
		UserLogLevel6:  `USER6`,
		UserLogLevel7:  `USER7`,
		UserLogLevel8:  `USER8`,
		UserLogLevel9:  `USER9`,
		UserLogLevel10: `USER10`,
		AllLogLevels:   `ALL`,
	}

	// str->instance
	logLevelMap = map[string]LogLevel{
		`NONE`:   NoLogLevels,
		`CALLS`:  LogLevel1,
		`POLICY`: LogLevel2,
		`STATE`:  LogLevel3,
		`DEBUG`:  LogLevel4,
		`ERROR`:  LogLevel5,
		`TRACE`:  LogLevel6,
		`USER1`:  UserLogLevel1,
		`USER2`:  UserLogLevel2,
		`USER3`:  UserLogLevel3,
		`USER4`:  UserLogLevel4,
		`USER5`:  UserLogLevel5,
		`USER6`:  UserLogLevel6,
		`USER7`:  UserLogLevel7,
		`USER8`:  UserLogLevel8,
		`USER9`:  UserLogLevel9,
		`USER10`: UserLogLevel10,
		`ALL`:    AllLogLevels,
	}
}
