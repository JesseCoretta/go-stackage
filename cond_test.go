package stackage

import (
	"bytes"
	"fmt"
	"log"
	"testing"
)

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

/*
This example demonstrates setting a custom logger which writes
to a bytes.Buffer io.Writer qualifier. A loglevel of "all" is
invoked, and an event -- the creation of a Condition -- shall
trigger log events that are funneled into our bytes.Buffer.

For the sake of idempotent test results, the length of the buffer
is checked only to ensure it is greater than 0.
*/
func ExampleSetDefaultConditionLogger() {
	var buf *bytes.Buffer = &bytes.Buffer{}
	var customLogger *log.Logger = log.New(buf, ``, 0)
	SetDefaultConditionLogger(customLogger)
	SetDefaultConditionLogLevel(AllLogLevels) // highly verbose!

	// do something that triggers log events ...
	_ = Cond(`π`, Eq, float64(3.14159265358979323))

	fmt.Printf("%T.Len>0 (bytes): %t", buf, buf.Len() > 0)
	// Output: *bytes.Buffer.Len>0 (bytes): true
}

func ExampleCond() {
	fmt.Printf("%s", Cond(`π`, Eq, float64(3.14159265358979323)))
	// Output: π = 3.14
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
	// Output: stackage.Condition keyword value is zero
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
