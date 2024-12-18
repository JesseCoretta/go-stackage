package stackage

import (
	"bytes"
	"fmt"
	"log"
	"strings"
	"testing"
)

// any type would do, but lets do a struct
// because we don't use structs for operators
// within this package for built-ins.
type fakeOperator struct {
	Str string
	Ctx string
}

func (r fakeOperator) Context() string {
	return r.Ctx
}

func (r fakeOperator) String() string {
	return r.Str
}

func ExampleSetDefaultConditionLogLevel() {
	// define a custom loglevel cfg
	SetDefaultConditionLogLevel(
		LogLevel3 + // 4
			UserLogLevel1 + // 64
			UserLogLevel7, // 4096
	)
	custom := DefaultConditionLogLevel()

	// turn loglevel to none
	SetDefaultConditionLogLevel(NoLogLevels)
	off := DefaultConditionLogLevel()

	fmt.Printf("%d (custom), %d (off)", custom, off)
	// Output: 4164 (custom), 0 (off)

}

func ExampleDefaultConditionLogLevel() {
	fmt.Printf("%d", DefaultConditionLogLevel())
	// Output: 0
}

func ExampleCondition_IsEqual() {
	c1 := Cond(`Keyword`, Eq, 123)
	c2 := Cond(`Keyword`, Gt, 123)
	fmt.Printf("Conditions are equal: %t", c1.IsEqual(c2) == nil)
	// Output: Conditions are equal: false
}

func ExampleCondition_SetEqualityPolicy() {
	c1 := Cond(`KEYWORD`, Eq, `This is A value`)
	c2 := Cond(`Keyword`, Eq, `this is a VALUE`)

	// We'd like to compare the two Condition
	// instances -- but ignore case folding in
	// the keyword or value, and only compare
	// the string forms of the Eq comparison
	// operator (=), as Operator interface
	// qualifiers can manifest as any type...
	policy := func(a, b any) error {
		// Let's convert instances in case they are
		// type alias Conditions. Alternatively, if
		// we knew which type(s) to expect, we could
		// just live dangerously and wing it ...
		C1, ok1 := ConvertCondition(a)
		C2, ok2 := ConvertCondition(b)
		if !ok1 || !ok2 {
			return fmt.Errorf("Non-condition input failed equality assertion")
		}

		// Normalize keywords before equality
		// checks for maximum compatibility.
		kw1 := strings.ToLower(C1.Keyword())
		kw2 := strings.ToLower(C2.Keyword())
		if kw1 != kw2 {
			return fmt.Errorf("Keyword mismatch")
		}

		// Generally the operator is just a few
		// characters, e.g.: ">=" or "=", so we
		// will skip normalization. YMMV ...
		if C1.Operator().String() != C2.Operator().String() {
			return fmt.Errorf("Operator mismatch")
		}

		// Naturally if strings aren't the only
		// thing that might be encountered, you
		// may opt for more type coverage here.
		valA, okA := C1.Expression().(string)
		valB, okB := C2.Expression().(string)
		if !okA || !okB {
			return fmt.Errorf("Expression type mismatch")
		}

		// Compare values. Naturally, we could use
		// this opportunity to perform other kinds
		// of sanitation -- such as the removal of
		// any leading or trailing WHSP, et al ...
		if strings.ToLower(valA) != strings.ToLower(valB) {
			return fmt.Errorf("Expression content mismatch")
		}

		// Instances seem to be equal
		return nil
	}

	// Assign our function to whichever Condition
	// instance is used as the receiver (or both,
	// optionally) ...
	c1.SetEqualityPolicy(policy)

	fmt.Printf("Conditions are equal: %t", c1.IsEqual(c2) == nil)
	// Output: Conditions are equal: true
}

func ExampleCondition_SetUnmarshaler() {
	// convert the condition (the receiver) to
	// a map[string]string instance.
	var c Condition = Cond(`Keyword`, Eq, `this is the value`)
	tomap := func(_ ...any) (out []any, err error) {
		//err = errors.New("Set your error, if applicable")
		out = []any{map[string]string{
			`keyword`:  c.Keyword(),
			`operator`: c.Operator().String(),
			`value`:    c.Expression().(string)}}

		return out, err
	}

	// assign new func to Condition
	c.SetUnmarshaler(tomap)
	o, err := c.Unmarshal()
	if err != nil {
		fmt.Println(err)
		return
	}
	// Print output
	fmt.Println(o[0])
	// Output: map[keyword:Keyword operator:= value:this is the value]
}

func ExampleCondition_Auxiliary() {
	var c Condition
	c.Init()

	// one can put anything they wish into this map,
	// so we'll do a bytes.Buffer since it is simple
	// and commonplace.
	var buf *bytes.Buffer = &bytes.Buffer{}
	_, _ = buf.WriteString(`some .... data .....`)

	// Create our map (one could also use make
	// and populate it piecemeal as opposed to
	// in-line, as we do below).
	c.SetAuxiliary(map[string]any{
		`buffer`: buf,
	})

	//  Call our map and call its 'Get' method in one-shot
	if val, ok := c.Auxiliary().Get(`buffer`); ok {
		fmt.Printf("%s", val)
	}
	// Output: some .... data .....
}

func ExampleCondition_SetAuxiliary() {
	var c Condition
	c.Init()

	// alloc map
	aux := make(Auxiliary, 0)

	// populate map
	aux.Set(`somethingWeNeed`, struct {
		Type  string
		Value []string
	}{
		Type: `L`,
		Value: []string{
			`abc`,
			`def`,
		},
	})

	// assign map to condition rcvr
	c.SetAuxiliary(aux)

	// verify presence
	call := c.Auxiliary()
	fmt.Printf("%T found, length:%d", call, call.Len())
	// Output: stackage.Auxiliary found, length:1
}

/*
This example demonstrates the call and (failed) set of an
uninitialized Auxiliary instance. While no panic ensues,
the map instance is not writable.

The user must instead follow the procedures in the WithInit,
UserInit or ByTypeCast examples.
*/
func ExampleCondition_Auxiliary_noInit() {

	var c Condition = Cond(`keyword`, Eq, `value`)
	aux := c.Auxiliary()
	fmt.Printf("%T found, length:%d",
		aux, aux.Set(`testing`, `123`).Len())
	// Output: stackage.Auxiliary found, length:0
}

/*
This example continues the concepts within the NoInit example,
except in this case proper initialization occurs and a desirable
outcome is achieved.
*/
func ExampleCondition_Auxiliary_withInit() {

	var c Condition = Cond(`keyword`, Eq, `value`)
	c.SetAuxiliary() // auto-init
	aux := c.Auxiliary()
	fmt.Printf("%T found, length was:%d, is now:%d",
		aux,
		aux.Len(),                       // check initial (pre-set) length
		aux.Set(`testing`, `123`).Len()) // fluent Set/Len in one shot
	// Output: stackage.Auxiliary found, length was:0, is now:1
}

/*
This example demonstrates a scenario similar to that of the WithInit
example, except in this case the map instance is entirely created and
populated by the user in a traditional fashion.
*/
func ExampleCondition_Auxiliary_userInit() {

	var c Condition = Cond(`keyword`, Eq, `value`)
	aux := make(Auxiliary, 0)

	// user opts to just use standard map
	// key/val set procedure, and avoids
	// use of the convenience methods.
	// This is totally fine.
	aux[`value1`] = []int{1, 2, 3, 4, 5}
	aux[`value2`] = [2]any{float64(7.014), rune('#')}

	c.SetAuxiliary(aux)
	fmt.Printf("%T length:%d", aux, len(aux))
	// Output: stackage.Auxiliary length:2
}

/*
This example demonstrates building of the Auxiliary map in its generic
form (map[string]any) before being type cast to Auxiliary.
*/
func ExampleCondition_Auxiliary_byTypeCast() {

	var c Condition = Cond(`keyword`, Eq, `value`)
	proto := make(map[string]any, 0)
	proto[`value1`] = []int{1, 2, 3, 4, 5}
	proto[`value2`] = [2]any{float64(7.014), rune('#')}
	c.SetAuxiliary(Auxiliary(proto)) // cast proto and assign to stack
	aux := c.Auxiliary()             // call map to variable
	fmt.Printf("%T length:%d", aux, aux.Len())
	// Output: stackage.Auxiliary length:2
}

func TestCondition_SetAuxiliary(t *testing.T) {
	var c Condition = Cond(`keyword`, Eq, `value`)

	// alloc map
	aux := make(Auxiliary, 0)

	// populate map
	aux.Set(`somethingWeNeed`, struct {
		Type  string
		Value []string
	}{
		Type: `L`,
		Value: []string{
			`abc`,
			`def`,
		},
	})

	// assign map to stack rcvr
	c.SetAuxiliary(aux)

	// test that it was assigned properly
	var call Auxiliary = c.Auxiliary()
	if call == nil {
		t.Errorf("%s failed: %T nil", t.Name(), call)
		return
	}

	// make sure the contents are present
	want := 1
	if length := call.Len(); length != want {
		t.Errorf("%s failed: unexpected length; want %d, got %d",
			t.Name(), want, length)
		return
	}
}

func ExampleCond() {
	fmt.Printf("%s", Cond(`π`, Eq, float64(3.141592653589793)))
	// Output: π = 3.141592653589793
}

func ExampleCondition_CanNest() {
	c := Cond(`π`, Eq, float64(3.14159265358979323))
	c.NoNesting(true)
	fmt.Printf("Can nest: %t", c.CanNest())
	// Output: Can nest: false
}

func ExampleCondition_Category() {
	c := Cond(`π`, Eq, float64(3.14159265358979323))
	c.SetCategory(`mathematical constant`)
	fmt.Printf("Category: %s", c.Category())
	// Output: Category: mathematical constant
}

func ExampleCondition_SetCategory() {
	c := Cond(`π`, Eq, float64(3.14159265358979323))
	c.SetCategory(`mathematical constant`)
	fmt.Printf("Category: %s", c.Category())
	// Output: Category: mathematical constant
}

func ExampleCondition_ID() {
	c := Cond(`π`, Eq, float64(3.14159265358979323))
	c.SetID(`pi`)
	fmt.Printf("ID: %s", c.ID())
	// Output: ID: pi
}

func ExampleCondition_SetID() {
	c := Cond(`π`, Eq, float64(3.14159265358979323))
	c.SetID(`pi`)
	fmt.Printf("ID: %s", c.ID())
	// Output: ID: pi
}

func ExampleCondition_Encap_doubleQuote() {
	var c Condition
	c.Init()
	c.SetKeyword(`key`)
	c.SetOperator(Eq)
	c.SetExpression(`value`)
	c.Encap(`"`)
	fmt.Printf("%s", c)
	// Output: key = "value"
}

func ExampleCondition_Encap_tag() {
	var c Condition
	c.Init()
	c.SetKeyword(`key`)
	c.SetOperator(Eq)
	c.SetExpression(`value`)
	c.Encap([]string{`<`, `>`})
	fmt.Printf("%s", c)
	// Output: key = <value>
}

func ExampleCondition_IsEncap() {
	var c Condition = Cond(`key`, Eq, `value`)
	c.Encap([]string{`<`, `>`})
	fmt.Printf("%T expression is encapsulated: %t", c, c.IsEncap())
	// Output: stackage.Condition expression is encapsulated: true
}

func ExampleCondition_Err() {
	var c Condition = Cond(``, ComparisonOperator(7), `ThisIsBogus`)
	fmt.Printf("%v", c.Err())
	// Output: keyword value is zero
}

/*
This example demonstrates use of the Condition.SetErr method to assign
a custom error to the receiver. This is useful in cases where custom
processing is being conducted, and a custom error needs to be preserved
and accessible to the caller.
*/
func ExampleCondition_SetErr() {
	var c Condition = Cond(`mySupposedly`, Eq, `validCondition`)
	c.SetErr(fmt.Errorf("Hold on there buddy, this is garbage"))

	fmt.Printf("%v", c.Err())
	// Output: Hold on there buddy, this is garbage
}

/*
This example demonstrates use of the Condition.SetErr method to clear
an error instance from the receiver. This may be useful in the event
that a problem has potentially been solved and a re-run of testing is
needed.
*/
func ExampleCondition_SetErr_clear() {
	var c Condition = Cond(`mySupposedly`, Eq, `validCondition`)
	c.SetErr(nil)

	fmt.Printf("Nil: %t", c.Err() == nil)
	// Output: Nil: true
}

func ExampleCondition_Expression() {
	var c Condition
	c.Init()
	c.SetKeyword(`key`)
	c.SetOperator(Eq)
	c.SetExpression(`value`)
	fmt.Printf("%s", c.Expression())
	// Output: value
}

func ExampleCondition_Init() {
	var c Condition
	c.Init()
	fmt.Printf("%T is initialized: %t", c, c.IsInit())
	// Output: stackage.Condition is initialized: true
}

func ExampleCondition_IsInit() {
	var c Condition
	fmt.Printf("%T is initialized: %t", c, c.IsInit())
	// Output: stackage.Condition is initialized: false
}

func ExampleCondition_IsFIFO() {
	valueStack := And().SetFIFO(true).Push(`this`, `that`, `other`)
	c := Cond(`stack`, Eq, valueStack)
	fmt.Printf("%T is First-In/First-Out: %t", c, c.IsFIFO())
	// Output: stackage.Condition is First-In/First-Out: true
}

func ExampleCondition_IsNesting() {
	valueStack := And().Push(`this`, `that`, `other`)
	c := Cond(`stack`, Eq, valueStack)
	fmt.Printf("%T is nesting: %t", c, c.IsNesting())
	// Output: stackage.Condition is nesting: true
}

func ExampleCondition_IsPadded() {
	c := Cond(`keyword`, Eq, `value`)
	fmt.Printf("%T is padded: %t", c, c.IsPadded())
	// Output: stackage.Condition is padded: true
}

type customTestValue struct {
	Type  string
	Value string
}

func (r customTestValue) String() string {
	return sprintf("TYPE: %s; VALUE: %s", r.Type, r.Value)
}

type customTestKeyword uint8

const (
	testKeyword01 customTestKeyword = 1 << iota
	testKeyword02
	testKeyword03
	testKeyword04
	testKeyword05
	testKeyword06
)

func (r customTestKeyword) String() string {
	switch r {
	case testKeyword01:
		return `keyword_01`
	case testKeyword02:
		return `keyword_02`
	case testKeyword03:
		return `keyword_03`
	case testKeyword04:
		return `keyword_04`
	case testKeyword05:
		return `keyword_05`
	case testKeyword06:
		return `keyword_06`
	}

	return `<unknown_keyword>`
}

func TestCondition_001(t *testing.T) {
	c := Cond(`person`, Eq, `Jesse`)
	want := `person = Jesse`
	got := c.String()
	if want != got {
		t.Errorf("%s failed: want '%s', got '%s'", t.Name(), want, got)
	}
}

func TestCondition_customKeyword(t *testing.T) {
	var c Condition
	c.Init() // always init
	c.SetKeyword(testKeyword05)
	c.SetOperator(Gt)
	c.SetExpression(complex(2, 3))
	if c.Keyword() != testKeyword05.String() {
		t.Errorf("%s failed: want '%s', got '%s'", t.Name(), testKeyword05, c.Keyword())
	}

	want := `keyword_05 > (2+3i)`
	if got := c.String(); got != want {
		t.Errorf("%s failed: want '%s', got '%s'", t.Name(), want, got)
	}
}

func TestComparisonOperator_String(t *testing.T) {
	for _, cop := range []ComparisonOperator{
		Eq, Ne, Lt, Gt, Le, Ge,
	} {
		if str := cop.String(); str == badOp {
			t.Errorf("%s failed: got '%s' (unexpected)", t.Name(), str)
			return
		}
	}
}

func TestCondition_customType(t *testing.T) {
	cv := customTestValue{Type: `Color`, Value: `Red`}

	var c Condition
	c.Init() // always init
	c.SetKeyword(testKeyword05)
	c.SetOperator(Ne)
	c.SetExpression(cv)

	want := `keyword_05 != TYPE: Color; VALUE: Red`
	if got := c.String(); got != want {
		t.Errorf("%s failed: want '%s', got '%s'", t.Name(), want, got)
	}
}

func TestCondition_002(t *testing.T) {
	st := List().Paren().Push(
		Cond(`person`, Eq, `Jesse`),
	)

	want := `( person = Jesse )`
	if got := st.String(); want != got {
		t.Errorf("%s failed: want '%s', got '%s'", t.Name(), want, got)
	}
}

func TestCondition_IsParen(t *testing.T) {
	cond := List().Paren().Push(
		Cond(`person`, Eq, `Jesse`),
	)

	if got := cond.IsParen(); !got {
		t.Errorf("%s failed: want 'true', got '%t'", t.Name(), got)
	}
}

func TestCondition_NoNesting(t *testing.T) {
	var c Condition = Cond(`myKeyword`, Eq, `temporary_value`)
	c.NoNesting(true)

	if c.CanNest() {
		t.Errorf("%s failed: want '%t', got '%t'", t.Name(), false, true)
		return
	}

	c.SetExpression(And().Push(`this`, `that`))

	// value should NOT have been assigned.
	if _, asserted := c.Expression().(string); !asserted {
		t.Errorf("%s failed: want '%t', got '%t' [%T]", t.Name(), false, asserted, c.Expression())
	}
}

func TestCondition_CanNest(t *testing.T) {
	var c Condition
	c.NoNesting(true)
	if c.CanNest() {
		t.Errorf("%s failed: want '%t', got '%t'", t.Name(), false, true)
	}
}

func TestCondition_IsNesting(t *testing.T) {
	var c Condition
	c.Init() // always init
	c.SetKeyword(`myKeyword`)
	c.SetOperator(Eq)
	c.SetExpression(And().Push(`this`, `that`))

	if !c.IsNesting() {
		t.Errorf("%s failed: want '%t', got '%t'", t.Name(), true, false)
	}
}

//func TestCondition_IsPadded(t *testing.T) {
//	cond := Cond(`person`, Eq, `Jesse`).Paren().Encap(`"`).NoPadding(false)
//
//	if got := cond.IsPadded(); got {
//		t.Errorf("%s failed: want 'false', got '%t'", t.Name(), got)
//	}
//}

func TestCondition_IsEncap(t *testing.T) {
	cond := Cond(`person`, Eq, `Jesse`).Paren().Encap(`"`).NoPadding()

	if got := cond.IsEncap(); !got {
		t.Errorf("%s failed: want 'true', got '%t'", t.Name(), got)
	}
}

/*
This example demonstrates the act of a one-shot creation of an instance
of Condition using the Cond package-level function. All values must be
non-zero/non-nil. This is useful in cases where users have all of the
information they need and simply want to fashion a Condition quickly.

Note: line-breaks added for readability; no impact on functionality,
so long as dot sequence "connectors" (.) are honored where required.
*/
func ExampleCondition_oneShot() {
	c := Cond(`person`, Eq, `Jesse`).
		Paren().
		Encap(`"`).
		NoPadding()

	fmt.Printf("%s", c)
	// Output: (person="Jesse")
}

func ExampleCondition_String() {
	c := Cond(`person`, Eq, `Jesse`).
		Paren().
		Encap(`"`).
		NoPadding()

	fmt.Printf("%s", c)
	// Output: (person="Jesse")
}

func ExampleCondition_Valid() {
	c := Cond(`person`, Eq, `Jesse`).
		Paren().
		Encap(`"`).
		NoPadding()

	fmt.Printf("Valid: %t", c.Valid() == nil)
	// Output: Valid: true
}

/*
This example demonstrates the act of procedural assembly of an instance
of Condition. The variable c is defined and not initialized. Its values
are not set until "later" in the users code, and are set independently
of one another.

Note that with this technique a manual call of the Init method is always
required as the first action following variable declaration. When using
the "one-shot" technique by way of the Cond package-level function, the
process of initialization is handled automatically for the user.
*/
func ExampleCondition_procedural() {
	var c Condition
	c.Init() // always init
	c.Paren()
	c.SetKeyword(`myKeyword`)
	c.SetOperator(Eq)
	c.SetExpression(`value123`)

	fmt.Printf("%s", c)
	// Output: ( myKeyword = value123 )
}

/*
This example demonstrates use of the Addr method to obtain the string
representation of the memory address for the underlying embedded receiver
instance.

NOTE: Since the address will usually be something different each runtime,
we can't reliably test this in literal fashion, so we do a simple prefix
check as an alternative.
*/
func ExampleCondition_Addr() {
	var c Condition
	c.Init()

	fmt.Printf("Address prefix valid: %t", c.Addr()[:2] == `0x`)
	// Address prefix valid: true
}

func ExampleCondition_SetExpression() {
	var c Condition
	c.Init()
	c.SetKeyword(`keyword`)
	c.SetOperator(Ge)
	c.SetExpression(1.456)
	fmt.Printf("Expr type: %T", c.Expression())
	// Output: Expr type: float64
}

func ExampleCondition_Operator() {
	var c Condition
	c.Init()
	c.SetKeyword(`keyword`)
	c.SetOperator(Ge)
	c.SetExpression(1.456)
	fmt.Printf("Operator: %s", c.Operator())
	// Output: Operator: >=
}

func ExampleCondition_SetOperator() {
	var c Condition
	c.Init()
	c.SetKeyword(`keyword`)
	c.SetOperator(Ge)
	c.SetExpression(1.456)
	fmt.Printf("Operator: %s", c.Operator())
	// Output: Operator: >=
}

func ExampleCondition_Paren() {
	var c Condition
	c.Init()
	c.Paren() // toggle to true from false default
	c.SetKeyword(`keyword`)
	c.SetOperator(Ge)
	c.SetExpression(1.456)
	fmt.Printf("%s", c)
	// Output: ( keyword >= 1.456 )
}

func ExampleCondition_IsParen() {
	var c Condition
	c.Init()
	c.Paren() // toggle to true from false default
	c.SetKeyword(`keyword`)
	c.SetOperator(Ge)
	c.SetExpression(1.456)
	fmt.Printf("Is parenthetical: %t", c.IsParen())
	// Output: Is parenthetical: true
}

func ExampleCondition_IsZero() {
	var c Condition
	fmt.Printf("Zero: %t", c.IsZero())
	// Output: Zero: true
}

func ExampleCondition_SetKeyword() {
	var c Condition
	c.Init()
	c.SetKeyword(`my_keyword`)
	fmt.Printf("Keyword: %s", c.Keyword())
	// Output: Keyword: my_keyword
}

func ExampleCondition_Keyword() {
	var c Condition
	c.Init()
	c.SetKeyword(`my_keyword`)
	fmt.Printf("Keyword: %s", c.Keyword())
	// Output: Keyword: my_keyword
}

func ExampleCondition_LogLevels() {
	var buf *bytes.Buffer = &bytes.Buffer{}
	var customLogger *log.Logger = log.New(buf, ``, 0)
	var c Condition
	c.Init()
	c.SetLogger(customLogger)
	c.SetLogLevel(LogLevel1, LogLevel3)
	fmt.Printf("Loglevels: %s", c.LogLevels())
	// Output: Loglevels: CALLS,STATE
}

func ExampleCondition_NoNesting() {
	var c Condition
	c.Init()
	c.Paren() // toggle to true from false default
	c.SetKeyword(`keyword`)
	c.SetOperator(Ge)
	c.NoNesting(true)

	c.SetExpression(
		And().Push(`this`, `won't`, `work`),
	)

	fmt.Printf("%v", c.Expression())
	// Output: <nil>
}

func ExampleCondition_NoPadding() {
	var c Condition
	c.Init()
	c.NoPadding(true)
	c.SetKeyword(`keyword`)
	c.SetOperator(Ge)
	c.SetExpression(`expression`)

	fmt.Printf("%s", c)
	// Output: keyword>=expression
}

func ExampleCondition_Len() {
	var c Condition
	c.Init()
	c.Paren() // toggle to true from false default
	c.SetKeyword(`keyword`)
	c.SetOperator(Ge)
	noLen := c.Len()

	c.SetExpression(`just_a_string`)
	strLen := c.Len()

	// Overwrite the above string
	// with a stack containing
	// three (3) strings.
	S := And().Push(`this`, `won't`, `work`)
	c.SetExpression(S)
	stkLen := c.Len()

	fmt.Printf("length with nothing: %d\nlength with string: %d\nlength with %T: %d", noLen, strLen, S, stkLen)
	// Output:
	// length with nothing: 0
	// length with string: 1
	// length with stackage.Stack: 3
}

func ExampleCondition_Logger() {
	var buf *bytes.Buffer = &bytes.Buffer{}
	var customLogger *log.Logger = log.New(buf, ``, 0)
	var c Condition
	c.Init()
	c.SetLogger(customLogger)
	fmt.Printf("%T", c.Logger())
	// Output: *log.Logger
}

func ExampleCondition_SetLogLevel() {
	var buf *bytes.Buffer = &bytes.Buffer{}
	var customLogger *log.Logger = log.New(buf, ``, 0)
	var c Condition
	c.Init()
	c.SetLogger(customLogger)
	c.SetLogLevel(LogLevel1, LogLevel3) // calls (1) + state(4)
	fmt.Printf("LogLevels active: %s", c.LogLevels())
	// Output: LogLevels active: CALLS,STATE
}

func ExampleCondition_UnsetLogLevel() {
	var buf *bytes.Buffer = &bytes.Buffer{}
	var customLogger *log.Logger = log.New(buf, ``, 0)
	var c Condition
	c.Init()
	c.SetLogger(customLogger)
	c.SetLogLevel(LogLevel1, LogLevel3) // calls (1) + state(4)
	c.UnsetLogLevel(LogLevel1)          // -1
	fmt.Printf("LogLevels active: %s", c.LogLevels())
	// Output: LogLevels active: STATE
}

func ExampleCondition_SetID_random() {
	var c Condition = Cond(`keyword`, Ne, `bogus`)

	// can't predict what ID will
	// be, so we'll check length
	// which should always be 24.
	c.SetID(`_random`)
	fmt.Printf("Random ID len: %d", len(c.ID()))
	// Output: Random ID len: 24
}

func ExampleCondition_SetID_pointerAddress() {
	var c Condition = Cond(`keyword`, Ne, `bogus`)

	// can't predict what ID will be,
	// so we'll check the prefix to
	// be certain it begins with '0x'.
	c.SetID(`_addr`)
	fmt.Printf("Address ID has '0x' prefix: %t", c.ID()[:2] == `0x`)
	// Output: Address ID has '0x' prefix: true
}

/*
This example offers a naïve demonstration of a user-authored
type alias for the Condition type.

In a practical scenario, a user would likely want to write
custom methods to perform new tasks and/or wrap any existing
Condition methods desired in order to avoid the frequent need
to cast a custom type to the native stackage.Condition type.

As this is an example, we'll perform a simple cast merely to
demonstrate the type was cast successfully.
*/
func ExampleCondition_typeAlias() {
	type customCondition Condition
	var c Condition = Cond(`keyword`, Ne, `Value`)
	fmt.Printf("%s (%T)",
		Condition(customCondition(c)),
		customCondition(c),
	)
	// Output: keyword != Value (stackage.customCondition)
}

func TestCondition_codecov(t *testing.T) {
	var c Condition
	// panic checks
	c.Len()
	c.IsZero()
	c.IsInit()
	c.Expression()
	c.Operator()
	c.Keyword()
	c.Valid()
	_ = cond.String()

	SetDefaultConditionLogger(nil)

	var ll logLevels
	ll.shift(`trace`)
	ll.shift(nil)
	ll.shift('a')
	ll.shift()
	ll.shift(0)
	ll.shift(65535)

	ll.positive(`trace`)
	ll.positive(nil)
	ll.positive('a')
	ll.positive(0)
	ll.positive(65535)

	ll.unshift(`trace`)
	ll.unshift(nil)
	ll.unshift('a')
	ll.unshift(0)
	ll.unshift(65535)

	c.SetAuxiliary(nil)

	var lsys *logSystem = newLogSystem(cLogDefault)
	if lsys.isZero() {
		t.Errorf("%s failed: nil %T",
			t.Name(), lsys.logger())
	}

	c.Init()
	c.SetEqualityPolicy()
	c.SetUnmarshaler()
	c.Unmarshal()
	c.Paren()
	c.Paren(true)
	c.Paren(false)
	c.SetKeyword(`valid_keyword`)
	c.Valid()
	c.SetOperator(Ne)
	c.SetAuxiliary(nil)
	c.Valid()
	c.SetExpression(rune(1438))
	_ = c.String()
	c.SetExpression(true)
	_ = c.String()
	c.SetExpression(float32(3.6663))
	_ = c.String()
	c.SetExpression(1438)
	_ = c.String()
	c.SetExpression(uint8(10))
	_ = c.String()
	c.SetExpression(`string`)
	_ = c.String()
	c.SetOperator(ComparisonOperator(77))
	c.Encap()
	c.Encap(`<<`, `>>`)
	c.Encap(``, ``)
	c.Encap(`<<`, `<<`)
	c.Encap(`"`)
	c.Valid()
	_ = c.condition.cfg.log.lvl.String()

	c.SetLogLevel(LogLevel5)
	c.SetLogLevel(AllLogLevels)
	_ = c.condition.cfg.log.lvl.String()
	c.UnsetLogLevel(`all`)

	SetDefaultConditionLogLevel(`none`)
	SetDefaultConditionLogLevel(0)
	SetDefaultConditionLogLevel(nil)
	SetDefaultConditionLogLevel('a')

	type customCondition Condition
	var cx Condition = Cond(`keyword`, Ne, `Value`)
	if Cx, ok := conditionTypeAliasConverter(customCondition(cx)); !ok {
		t.Errorf("%s failed: %T->%T conversion failure", t.Name(), cx, Cx)
	}

	var fo Operator = fakeOperator{Str: `>`, Ctx: `gtContext`}
	var cz Condition = Cond(`koyword`, fo, `Value`)
	cx.IsEqual(cz)

	cz.Free()
	fo = fakeOperator{Str: `!=`, Ctx: `gtContext`}
	cz = Cond(`keyword`, fo, `Value`)
	cx.IsEqual(cz)
	cx.Free()
	cz.SetReadOnly()
	cz.Free()
	cz.IsReadOnly()
	cz.SetReadOnly()
	cz.Free()

	subc := []any{`CONDITION`, `Keywerdd`, Gt, 5}
	extractConditionValues([]any{`CONDITION`, `Keyword`, Eq, subc})
}
