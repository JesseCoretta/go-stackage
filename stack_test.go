package stackage

import (
	"bytes"
	"fmt"
	"log"
	"testing"
	_ "time"
)

func ExampleComparisonOperator_Context() {
	var cop ComparisonOperator = Ge
	fmt.Printf("Context: %s", cop.Context())
	// Output: Context: comparison
}

func ExampleComparisonOperator_String() {
	var cop ComparisonOperator = Ge
	fmt.Printf("Operator: %s", cop)
	// Output: Operator: >=
}

func ExampleStack_Auxiliary() {
	// make a stack ... any type would do
	l := List().Push(`this`, `that`, `other`)

	// one can put anything they wish into this map,
	// so we'll do a bytes.Buffer since it is simple
	// and commonplace.
	var buf *bytes.Buffer = &bytes.Buffer{}
	_, _ = buf.WriteString(`some .... data .....`)

	// Create our map (one could also use make
	// and populate it piecemeal as opposed to
	// in-line, as we do below).
	l.SetAuxiliary(map[string]any{
		`buffer`: buf,
	})

	//  Call our map and call its 'Get' method in one-shot
	if val, ok := l.Auxiliary().Get(`buffer`); ok {
		fmt.Printf("%s", val)
	}
	// Output: some .... data .....
}

func ExampleAuxiliary_Get() {
	var aux Auxiliary = make(Auxiliary, 0)
	aux.Set(`value`, 18)
	val, ok := aux.Get(`value`)
	if ok {
		fmt.Printf("%d", val)
	}
	// Output: 18
}

func ExampleAuxiliary_Set() {
	var aux Auxiliary = make(Auxiliary, 0)
	aux.Set(`value`, 18)
	aux.Set(`color`, `red`)
	aux.Set(`max`, 200)
	fmt.Printf("Len: %d", aux.Len())
	// Output: Len: 3
}

func ExampleAuxiliary_Unset() {
	var aux Auxiliary = make(Auxiliary, 0)
	aux.Set(`value`, 18)
	aux.Set(`color`, `red`)
	aux.Set(`max`, 200)
	aux.Unset(`max`)

	fmt.Printf("Len: %d", aux.Len())
	// Output: Len: 2
}

func ExampleAuxiliary_Len() {
	aux := Auxiliary{
		`value`: 18,
	}
	fmt.Printf("Len: %d", aux.Len())
	// Output: Len: 1
}

func ExampleSetDefaultStackLogLevel() {
	SetDefaultStackLogLevel(
		LogLevel1 + // 1
			LogLevel4 + // 8
			UserLogLevel2 + // 128
			UserLogLevel5, // 1024
	)
	custom := DefaultStackLogLevel()

	// turn loglevel to none
	SetDefaultStackLogLevel(NoLogLevels)
	off := DefaultStackLogLevel()

	fmt.Printf("%d (custom), %d (off)", custom, off)
	// Output: 1161 (custom), 0 (off)
}

func ExampleDefaultStackLogLevel() {
	fmt.Printf("%d", DefaultStackLogLevel())
	// Output: 0
}

/*
This example demonstrates setting a custom logger which writes
to a bytes.Buffer io.Writer qualifier. A loglevel of "all" is
invoked, and an event -- the creation of a Basic Stack -- shall
trigger log events that are funneled into our bytes.Buffer.

For the sake of idempotent test results, the length of the buffer
is checked only to ensure it is greater than 0.
*/
func ExampleSetDefaultStackLogger() {
	var buf *bytes.Buffer = &bytes.Buffer{}
	var customLogger *log.Logger = log.New(buf, ``, 0)
	SetDefaultStackLogger(customLogger)
	SetDefaultStackLogLevel(AllLogLevels) // highly verbose!

	// do something that triggers log events ...
	_ = Basic().Push(
		Cond(`π`, Eq, float64(3.14159265358979323)),
		Cond(`m`, Eq, `mc²`),
	)

	fmt.Printf("%T.Len>0 (bytes): %t", buf, buf.Len() > 0)
	// Output: *bytes.Buffer.Len>0 (bytes): true
}

var testParens []string = []string{`(`, `)`}

type customStack Stack // simulates a user-defined type that aliases a Stack

func (r customStack) String() string {
	return Stack(r).String()
}

func ExampleStack_SetAuxiliary() {

	// Always alloc stack somehow, in this
	// case just use List because its simple
	// and (unlike Basic) it is feature-rich.
	var list Stack = List()

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
	list.SetAuxiliary(aux)

	// verify presence
	call := list.Auxiliary()
	fmt.Printf("%T found, length:%d", call, call.Len())
	// Output: stackage.Auxiliary found, length:1
}

func ExampleStack_Kind() {
	var myStack Stack = And()
	fmt.Printf("Kind: '%s'", myStack.Kind())
	// Output: Kind: 'AND'
}

/*
This example demonstrates the call and (failed) set of an
uninitialized Auxiliary instance. While no panic ensues,
the map instance is not writable.

The user must instead follow the procedures in the WithInit,
UserInit or ByTypeCast examples.
*/
func ExampleStack_Auxiliary_noInit() {

	var list Stack = List()
	aux := list.Auxiliary()
	fmt.Printf("%T found, length:%d",
		aux, aux.Set(`testing`, `123`).Len())
	// Output: stackage.Auxiliary found, length:0
}

/*
This example continues the concepts within the NoInit example,
except in this case proper initialization occurs and a desirable
outcome is achieved.
*/
func ExampleStack_Auxiliary_withInit() {

	var list Stack = List().SetAuxiliary() // no args triggers auto-init
	aux := list.Auxiliary()
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
func ExampleStack_Auxiliary_userInit() {

	var list Stack = List()
	aux := make(Auxiliary, 0)

	// user opts to just use standard map
	// key/val set procedure, and avoids
	// use of the convenience methods.
	// This is totally fine.
	aux[`value1`] = []int{1, 2, 3, 4, 5}
	aux[`value2`] = [2]any{float64(7.014), rune('#')}

	list.SetAuxiliary(aux)
	fmt.Printf("%T length:%d", aux, len(aux))
	// Output: stackage.Auxiliary length:2
}

/*
This example demonstrates building of the Auxiliary map in its generic
form (map[string]any) before being type cast to Auxiliary.
*/
func ExampleStack_Auxiliary_byTypeCast() {

	var list Stack = List()
	proto := make(map[string]any, 0)
	proto[`value1`] = []int{1, 2, 3, 4, 5}
	proto[`value2`] = [2]any{float64(7.014), rune('#')}
	list.SetAuxiliary(Auxiliary(proto)) // cast proto and assign to stack
	aux := list.Auxiliary()             // call map to variable
	fmt.Printf("%T length:%d", aux, aux.Len())
	// Output: stackage.Auxiliary length:2
}

func TestStack_SetAuxiliary(t *testing.T) {
	var list Stack = List()

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
	list.SetAuxiliary(aux)

	// test that it was assigned properly
	var call Auxiliary = list.Auxiliary()
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

func TestStack_noInit(t *testing.T) {
	var x Stack
	x.Push(`well?`) // first off, this should not panic

	// next, make sure it really
	// did not work ...
	if length := x.Len(); length != 0 {
		t.Errorf("%s failed [noInit]: want '%d' elements, got '%d'", t.Name(), 0, length)
		return
	}

	// lastly, we should definitely see an indication
	// of this issue during validity checks ...
	if err := x.Valid(); err == nil {
		t.Errorf("%s failed [noInit]: want 'error', got '%T'", t.Name(), nil)
		return
	}
}

func TestStack_And001(t *testing.T) {
	got := And().Paren().Fold().Push(
		`testing1`,
		`testing2`,
		`testing3`,
	)

	want := `( testing1 and testing2 and testing3 )`
	if got.String() != want {
		t.Errorf("%s failed: want '%s', got '%s'", t.Name(), want, got)
		return
	}
}

func TestStack_Insert(t *testing.T) {
	got := And().Paren().Fold().Push(
		`testing1`,
		`testing2`,
		`testing3`,
	)

	got.Insert(`testing0`, 0)

	want := `( testing0 and testing1 and testing2 and testing3 )`
	if got.String() != want {
		t.Errorf("%s.1 failed: want '%s', got '%s'", t.Name(), want, got)
		return
	}

	got.Insert(`testing2.5`, 3)
	want = `( testing0 and testing1 and testing2 and testing2.5 and testing3 )`
	if got.String() != want {
		t.Errorf("%s.2 failed: want '%s', got '%s'", t.Name(), want, got)
		return
	}

	got.Insert(`testing4`, 15)
	want = `( testing0 and testing1 and testing2 and testing2.5 and testing3 and testing4 )`
	if got.String() != want {
		t.Errorf("%s.3 failed: want '%s', got '%s'", t.Name(), want, got)
		return
	}
}

func TestStack_Replace(t *testing.T) {
	s := List().SetDelimiter(rune(44)).Push(
		`testing1`,
		`testing2`,
		`testing3`,
	)

	want := `testing0,testing2,testing3`
	s.Replace(`testing0`, 0)
	got := s.String()
	if want != got {
		t.Errorf("%s.1 failed: want '%s', got '%s'", t.Name(), want, got)
		return
	}

	if ok := s.Replace(`testingX`, 6); ok {
		t.Errorf("%s.2 failed: want '%t', got '%t'", t.Name(), false, ok)
		return
	}

}

func TestStack_IsParen(t *testing.T) {
	stk := And().Paren().Fold().Push(
		`testing1`,
		`testing2`,
		`testing3`,
	)

	if got := stk.IsParen(); !got {
		t.Errorf("%s failed: want 'true', got '%t'", t.Name(), got)
		return
	}
}

func TestStack_Transfer(t *testing.T) {
	stk := And().Push(
		`testing1`,
		`testing2`,
		`testing3`,
	)

	or := Or()
	if ok := stk.Transfer(or); or.Len() != 3 || !ok {
		t.Errorf("%s failed [post-transfer len comparisons]: want len:%d, got len:%d", t.Name(), stk.Len(), or.Len())
		return
	}

	stk.Reset() // reset source, removing all slices
	if or.Len() != 3 && stk.Len() != 0 {
		t.Errorf("%s failed [post-transfer src reset]: want slen:%d; dlen:%d, got slen:%d; dlen:%d",
			t.Name(), 0, 3, stk.Len(), or.Len())
		return
	}
}

func TestStack_IsPadded(t *testing.T) {
	stk := And().Paren().Fold().Push(
		`testing1`,
		`testing2`,
		`testing3`,
	)

	if got := stk.IsPadded(); !got {
		t.Errorf("%s failed: want 'true', got '%t'", t.Name(), got)
		return
	}
}

func TestStackAnd_001(t *testing.T) {

	A := And().Paren().Push(
		`top_element_number_0`,
		Or().Paren().Push(
			`sub_element_number_0`,
			`sub_element_number_1`,
		),
	)

	want := `( top_element_number_0 AND ( sub_element_number_0 OR sub_element_number_1 ) )`
	if got := A; got.String() != want {
		t.Errorf("%s failed: want '%s', got '%s'", t.Name(), want, got)
		return
	}
}

func TestStack_IsNesting(t *testing.T) {

	A := And().Paren().Push(
		`top_element_number_0`,
		Or().Paren().Push(
			`sub_element_number_0`,
			`sub_element_number_1`,
		),
	)

	want := true
	got := A.IsNesting()
	if want != got {
		t.Errorf("%s failed [isNesting]: want '%t', got '%t'", t.Name(), want, got)
		return
	}
}

func TestAnd_002(t *testing.T) {

	// make each OR condition look like "<...>" (including quotes)
	O := Or().Paren().Encap(`"`, []string{`<`, `>`}).Push(
		`sub_element_number_0`,
		`sub_element_number_1`,
	)

	A := And().Paren().Push(
		`top_element_number_0`,
		O,
	)

	want := `( top_element_number_0 AND ( "<sub_element_number_0>" OR "<sub_element_number_1>" ) )`
	if got := A; got.String() != want {
		t.Errorf("%s failed: want '%s', got '%s'", t.Name(), want, got)
		return
	}

	if !O.IsEncap() {
		t.Errorf("%s failed [IsEncap]: want '%t', got '%t'", t.Name(), true, false)
		return
	}
}

func TestAnd_003(t *testing.T) {

	A := And().Symbol('&').Paren().LeadOnce().Encap(testParens).NoPadding().Push(
		`top_element_number_0`,

		Or().Symbol('|').Paren().LeadOnce().Encap(testParens).NoPadding().Push(
			`sub_element_number_0`,
			`sub_element_number_1`,
		),
	)

	want := `(&(top_element_number_0)(|(sub_element_number_0)(sub_element_number_1)))`
	if got := A; got.String() != want {
		t.Errorf("%s failed: want '%s', got '%s'", t.Name(), want, got)
		return
	}
}

func TestAnd_004(t *testing.T) {
	// make a helper function that incorporates
	// the same options into each of the three
	// stacks we're going to make. it just looks
	// cleaner.
	enabler := func(r Stack) Stack {
		return r.Paren().LeadOnce().Encap(testParens).NoPadding()
	}

	A := enabler(And().Symbol('&')).Push(
		`top_element_number_0`,

		enabler(Or().Symbol('|')).Push(
			`sub_element_number_0`,
			`sub_element_number_1`,
		),

		enabler(Not().Symbol('!')).Push(
			`unwanted_element_number_0`,
			`unwanted_element_number_1`,
		),
	)

	want := `(&(top_element_number_0)(|(sub_element_number_0)(sub_element_number_1))(!(unwanted_element_number_0)(unwanted_element_number_1)))`
	if got := A; got.String() != want {
		t.Errorf("%s failed: want '%s', got '%s'", t.Name(), want, got)
		return
	}
}

/*
This example demonstrates ORed stack values using the double pipe (||) symbol
and custom value encapsulation.
*/
func ExampleOr_symbolicOr() {
	or := Or().
		Paren().                                            // Add parenthesis
		Symbol(`||`).                                       // Use double pipes for OR
		Encap(` `, `"`).                                    // Encapsulate individual vals with double-quotes surrounded by spaces
		Push(`cn`, `sn`, `givenName`, `objectClass`, `uid`) // Push these values now

	fmt.Printf("%s\n", or)
	// Output: ( "cn" || "sn" || "givenName" || "objectClass" || "uid" )
}

/*
This example demonstrates ANDed stack values using the double ampersand (&&) symbol.
*/
func ExampleAnd_symbolicAnd() {
	and := And().
		Paren().                                       // Add parenthesis
		Symbol(`&&`).                                  // Use double pipes for OR
		Push(`condition1`, `condition2`, `condition3`) // Push these values now

	fmt.Printf("%s\n", and)
	// Output: ( condition1 && condition2 && condition3 )
}

func TestAnd_005_nestedWithTraversal(t *testing.T) {

	A := And().NoPadding().Push(
		`top_element_number_0`,

		Or().NoPadding().Paren().Push(
			`sub_element_number_0`,
			`sub_element_number_1`,
		),

		Not().NoPadding().Paren().Push(
			Or().NoPadding().Push(
				`unwanted_element_number_0`,
				`unwanted_element_number_1`, // traversal target
			),
		),
	)

	want := `top_element_number_0 AND (sub_element_number_0 OR sub_element_number_1) AND NOT (unwanted_element_number_0 OR unwanted_element_number_1)`
	if got := A; got.String() != want {
		t.Errorf("%s failed: want '%s', got '%s'", t.Name(), want, got)
		return
	}

	slice, ok := A.Traverse(2, 0, 1)
	if !ok {
		t.Errorf("%s failed: got '%T' during traversal, want non-nil", t.Name(), slice)
		return
	}

	strAssert, strOK := slice.(string)
	if !strOK {
		t.Errorf("%s failed: want '%T', got '%T' during assertion", t.Name(), ``, slice)
		return
	}

	want = `unwanted_element_number_1`
	if want != strAssert {
		t.Errorf("%s failed: want '%s', got '%s'", t.Name(), want, strAssert)
		return
	}
}

func TestList_001(t *testing.T) {

	A := List().SetDelimiter(rune(32)).Push(
		`top_element_number_0`,
		`top_element_number_1`,
		`top_element_number_2`,
		`top_element_number_3`,
	)

	want := `top_element_number_0 top_element_number_1 top_element_number_2 top_element_number_3`
	if got := A; got.String() != want {
		t.Errorf("%s failed: want '%s', got '%s'", t.Name(), want, got)
	}
}

func TestList_001_withNoDelim(t *testing.T) {

	A := List().Push(
		`(top_element_number_0)`,
		`(top_element_number_1)`,
		`(top_element_number_2)`,
		`(top_element_number_3)`,
	)

	want := `(top_element_number_0)(top_element_number_1)(top_element_number_2)(top_element_number_3)`
	if got := A; got.String() != want {
		t.Errorf("%s failed: want '%s', got '%s'", t.Name(), want, got)
	}
}

func TestCustomStack001(t *testing.T) {
	A := List().SetDelimiter(`,`).Push(
		`top_element_number_0`,
		`top_element_number_1`,
		`top_element_number_2`,
		`top_element_number_3`,
	)
	T := customStack(A)
	want := `top_element_number_0,top_element_number_1,top_element_number_2,top_element_number_3`
	if got := T; got.String() != want {
		t.Errorf("%s failed: want '%s', got '%s'", t.Name(), want, got)
	}
}

func TestCustomStack002_ldapFilter(t *testing.T) {
	maker := func(r Stack) Stack {
		return r.Paren().LeadOnce().NoPadding()
	}

	Ands := maker(And().Symbol('&')).Encap(testParens)
	Ors := maker(Or().Symbol('|')).Encap(testParens)
	Nots := maker(Not().Symbol('!')).Encap(testParens)

	// Begin filter at AND
	filter := Ands.Push(
		`objectClass=employee`, // AND condition #1

		// Begin OR (which is AND condition #2)
		Ors.Push(
			`objectClass=engineeringLead`, // OR condition #1
			`objectClass=shareholder`,     // OR condition #2
		),

		// Begin NOT (which is AND condition #3)
		Nots.Push(
			`drink=beer`, // NOT condition #1
			`c=RU`,       // NOT condition #2
		),
	)

	T := customStack(filter)
	want := `(&(objectClass=employee)(|(objectClass=engineeringLead)(objectClass=shareholder))(!(drink=beer)(c=RU)))`
	if got := T; got.String() != want {
		t.Errorf("%s failed: want '%s', got '%s'", t.Name(), want, got)
	}
}

func TestCustomStack003_nested(t *testing.T) {
	maker := func(r Stack) Stack {
		return r.Paren().LeadOnce().NoPadding()
	}

	Ands := maker(And().Symbol('&')).Encap(testParens).SetID(`filtery`)
	Ors := maker(Or().Symbol('|')).Encap(testParens).Push(
		`objectClass=engineeringLead`, // OR condition #1
		`objectClass=shareholder`,     // OR condition #2
	)

	// Begin filter at AND
	filter := Ands.Push(
		`objectClass=employee`, // AND condition #1

		// Begin OR (which is AND condition #2)
		customStack(Ors),
	)

	T := customStack(filter)
	want := `(&(objectClass=employee)(|(objectClass=engineeringLead)(objectClass=shareholder)))`
	if got := T; got.String() != want {
		t.Errorf("%s failed: want '%s', got '%s'", t.Name(), want, got)
	}
}

/*
This example demonstrates traversing a Stack instance containing Condition
slice instances. We use a path sequence of []int{1, 1} to target slice #1
on the first level, and value #1 of the second level.
*/
func ExampleStack_Traverse() {

	// An optional Stack "maker" for configuring
	// a specific kind of Stack. Just looks neater
	// than doing fluent execs over and over ...
	sMaker := func(r Stack) Stack {
		// encapsulate in parens, use symbol rune
		// for logical operators only once, and
		// use no padding between kw, op and value
		return r.Paren().LeadOnce().NoPadding()
	}

	// An optional Condition "maker", same logic
	// as above ...
	cMaker := func(r Condition) Condition {
		// encapsulate in parens, no padding between
		// kw, op and value ...
		return r.Paren().NoPadding()
	}

	// Let's make a faux LDAP Search Filter ...
	//
	// This will be our top level, which is an AND
	Ands := sMaker(And().Symbol('&'))

	// This will be our second level Stack, which
	// is an OR containing a couple of Condition
	// instances in Equality form.
	Ors := sMaker(Or().Symbol('|')).Push(
		cMaker(Cond(`objectClass`, Eq, `engineeringLead`)), // OR condition #1
		cMaker(Cond(`objectClass`, Eq, `shareholder`)),     // OR condition #2 // **our traversal target**
	)

	// Begin filter at AND, and push our
	// desired elements into the Stack.
	filter := Ands.Push(
		cMaker(Cond(`objectClass`, Eq, `employee`)), // AND condition #1
		Ors, // Begin OR (which is AND's condition #2)
	)

	// Bool returns shadowed only for brevity.
	// Generally you should not do that ...
	slice, _ := filter.Traverse(1, 1)   // Enter coordinates
	condAssert, ok := slice.(Condition) // The return is any, so assert to what we expect
	if !ok {
		fmt.Printf("Type Assertion failed: %T is not expected value\n", slice)
		return
	}

	fmt.Printf("%s", condAssert) // use its String method automagically
	// Output: (objectClass=shareholder)
}

/*
This example demonstrates the LeadOnce feature, which limits a given Stack's
logical operator usage to once-per, and only at the beginning (i.e.: the
form of an LDAP Search Filter's operators when conditions are nested).
*/
func ExampleStack_LeadOnce() {
	maker := func(r Stack) Stack {
		return r.Paren().LeadOnce().NoPadding()
	}

	Ands := maker(And().Symbol('&'))
	Ors := maker(Or().Symbol('|')).Push(
		Cond(`objectClass`, Eq, `engineeringLead`).NoPadding().Paren(), // OR condition #1
		Cond(`objectClass`, Eq, `shareholder`).NoPadding().Paren(),     // OR condition #1
	)

	// Begin filter at AND
	filter := Ands.Push(
		Cond(`objectClass`, Eq, `employee`).NoPadding().Paren(), // AND condition #1
		Ors, // Begin OR (which is AND condition #2)
	)

	fmt.Printf("%s", filter)
	// Output: (&(objectClass=employee)(|(objectClass=engineeringLead)(objectClass=shareholder)))
}

func TestBasic(t *testing.T) {
	b := Basic()
	b.Push(
		float64(3.14159),
		float64(-9.378),
	)

	want := float64(3.14159)
	got, ok := b.Index(0)
	if !ok {
		t.Errorf("%s failed: assertion error; want %T, got %T", t.Name(), want, got)
		return
	}
	if want != got.(float64) {
		t.Errorf("%s failed: want '%f', got: '%f'", t.Name(), want, got.(float64))
		return
	}
}

func TestBasic_withCapacity(t *testing.T) {
	cp := 2
	b := Basic(cp)

	// make sure the instance correctly
	// reflects the capacity imposed.
	if got := b.Cap(); got != cp {
		t.Errorf("%s failed: unexpected capacity value found; want %d, got %d", t.Name(), cp, got)
		return
	}

	// try to push one too many slices
	b.Push(
		float64(3.14159),
		float64(-9.378),
		float64(139.104),
	)

	// make sure only two (2) made it in
	if b.Len() != cp || !b.CapReached() {
		t.Errorf("%s failed: maximum capacity (%d) not honored; want len:%d, got len:%d", t.Name(), b.Cap(), cp, b.Len())
		return
	}
}

func TestBasic_availableCapacity(t *testing.T) {
	// allocate w/ max 2
	cp := 2
	b := Basic(cp)
	b.Push(
		float64(3.14159),
		float64(-9.378),
		float64(139.104),
	)

	if b.Avail() != 0 {
		t.Errorf("%s failed: unexpected available slice count; want len:%d, got len:%d", t.Name(), cp, b.Avail())
		return
	}

	// reallocate w/ max 5
	cp = 5
	if b = Basic(cp); b.Avail() != cp {
		t.Errorf("%s failed: unexpected available slice count; want len:%d, got len:%d", t.Name(), cp, b.Avail())
		return
	}
}

func TestReset(t *testing.T) {
	b := Basic()
	b.Push(
		float64(3.14159),
		float64(-9.378),
		float64(139.104),
	)
	if b.Len() != 3 {
		t.Errorf("%s failed: want '%d', got: '%d' [%s]", t.Name(), 3, b.Len(), b)
		return
	}

	b.Reset()

	if b.Len() != 0 {
		sl, _ := b.Index(0)
		t.Errorf("%s failed: want '%d', got: '%d' [%#v]", t.Name(), 0, b.Len(), sl)
		return
	}
}

/*
This example demonstrates the creation of a basic stack
and a call to its first (0th) index. Type assertion is
performed to reveal a float64 instance.
*/
func ExampleBasic() {
	b := Basic()
	b.Push(
		float64(3.14159),
		float64(-9.378),
	)
	idx, _ := b.Index(0)       // call index
	assert, _ := idx.(float64) // type assert to float64
	fmt.Printf("%.02f", assert)
	// Output: 3.14
}

/*
This example demonstrates the creation of a basic stack
with (and without) read-only controls enabled.
*/
func ExampleBasic_setAsReadOnly() {
	b := Basic()
	b.Push(
		float64(3.14159),
		float64(-9.378),
	)
	b.ReadOnly() // set readonly
	b.Remove(1)  // this ought to fail ...
	//b.Pop()	 // alternative to b.Remove(1) in this case
	first := b.Len() // record len

	b.ReadOnly()      // unset readonly
	b.Remove(1)       // retry removal
	second := b.Len() // record len again

	fmt.Printf("first try: %d vs. second try: %d", first, second)
	// Output: first try: 2 vs. second try: 1
}

func ExampleStack_Pop_lIFO() {
	b := Basic()
	b.Push(
		float64(3.14159),
		float32(-9.378),
		-1,
		`banana`,
	)

	popped, _ := b.Pop()
	fmt.Printf("%T, length now: %d", popped, b.Len())
	// Output: string, length now: 3
}

func ExampleStack_Pop_fIFO() {
	b := Basic()
	b.SetFIFO(true)
	b.Push(
		float64(3.14159),
		float32(-9.378),
		-1,
		`banana`,
	)

	popped, _ := b.Pop()
	fmt.Printf("%T, length now: %d", popped, b.Len())
	// Output: float64, length now: 3
}

/*
This example demonstrates the creation of a basic stack
and an enforced capacity constraint.
*/
func ExampleBasic_withCapacity() {
	b := Basic(2)
	b.Push(
		float64(3.14159),
		float64(-9.378),
		float64(139.104),
	)
	fmt.Printf("%d", b.Len())
	// Output: 2
}

func TestList_Join(t *testing.T) {
	L := List().SetDelimiter(`,`).Push(
		`item1`,
		`item2`,
	)
	want := `item1,item2`
	got := L.String()
	if want != got {
		t.Errorf("%s failed: want '%s', got '%s'", t.Name(), want, got)
		return
	}
}

/*
This example demonstrates the creation of a list stack
using comma delimitation.
*/
func ExampleStack_SetDelimiter() {
	// note: one could also use a rune
	// e.g: ',' or rune(44) for comma.
	L := List().SetDelimiter(`,`).Push(
		`item1`,
		`item2`,
	)
	fmt.Printf("%s", L)
	// Output: item1,item2
}

/*
This example demonstrates the creation of a list stack
using comma delimitation and the retrieval of the same
delimiter value.
*/
func ExampleStack_Delimiter() {
	// note: one could also use a rune
	// e.g: ',' or rune(44) for comma.
	L := List().SetDelimiter(`,`).Push(
		`item1`,
		`item2`,
	)
	fmt.Printf("%s", L.Delimiter())
	// Output: ,
}

func TestInterface(t *testing.T) {
	var elem Interface
	elem = Cond(`greeting`, Eq, `Hello`)
	want := `greeting = Hello`
	got := elem.String()
	if want != got {
		t.Errorf("%s failed: want '%s', got '%s'", t.Name(), want, got)
		return
	}
}

func TestStack_Reveal_experimental001(t *testing.T) {
	thisIsMyNightmare := And().Push(
		`this1`,
		Or().Mutex().Push(
			And().Push(Cond(`keyword`, Eq, "somevalue")),
			And().Push(
				`this4`,
				Not().Mutex().Push(
					Or().Push(
						Cond(`dayofweek`, Ne, "Wednesday"),
						Cond(`ssf`, Ge, "128"),
					),
				),
			),
			And().Push(
				Or().Push(
					Cond(`keyword2`, Lt, "someothervalue"),
				),
			),
		),
		`this2`,
	)

	type row struct {
		Index int
		Path  [][]int
		Want  string
		Value any
	}

	table := []row{
		{0, [][]int{{1, 2, 0, 0}, {1, 2}}, ``, nil},
		{1, [][]int{{1, 1, 1, 0, 0}, {1, 1, 1, 0, 0}}, ``, nil},
	}

	// Scan the target values we'll be using for
	// comparison after processing completes.
	for idx, tst := range table {
		slice, _ := thisIsMyNightmare.Traverse(tst.Path[0]...)
		var c Condition
		var ok bool
		if c, ok = slice.(Condition); !ok {
			t.Errorf("%s failed [idx:%d;pre-mod:path:%v]: unexpected assertion result:\nwant: %T\ngot:  %s",
				t.Name(), idx, tst.Path[0], c, slice.(string))
			return
		}
		table[idx].Want = c.String()
		tst.Value = c
	}

	// save for comparison later
	var want string = thisIsMyNightmare.String()

	// do reveal recursion
	thisIsMyNightmare.Reveal()

	// make sure the complete string is identical
	// before and after.
	if got := thisIsMyNightmare.String(); want != got {
		t.Errorf("%s failed [main strcmp]:\nwant '%s'\ngot  '%s',",
			t.Name(), want, got)
		return
	}

	// Scan the updated values at the defined paths
	// and compare to original target values.
	for idx, tst := range table {
		slice, _ := thisIsMyNightmare.Traverse(tst.Path[1]...)
		val, ok := slice.(Condition)
		if !ok {
			t.Errorf("%s failed [idx:%d;post-mod]: unexpected assertion result:\nwant: %T\ngot:  %T",
				t.Name(), idx, val, slice)
			return
		}

		if gval := val.String(); tst.Want != gval {
			t.Errorf("%s failed [idx:%d;strcmp]:\nwant '%s'\ngot  '%s',",
				t.Name(), idx, tst.Want, gval)
		}
	}
}

func TestDefrag_experimental_001(t *testing.T) {
	// this list contains an assortment of
	// values mixed in with nils and a couple
	// hierarchies tossed in, too.
	var l Stack = List().SetLogLevel(LogLevel(45)).Push(
		`this`,
		nil,
		`that`,
		nil,
		`those`,
		nil,
		List().Push(
			nil, nil, nil, `other`, `level`, `of`, nil, nil, `stuff`,
			And().Push(
				`yet`, nil, nil, nil, nil, `more`, `JUNK`,
			),
		),
		nil,
		nil,
		nil,
		nil,
		3.14159,
		`moar`,
		`moar`,
		`moar`,
		`moar`,
		`moar`,
		`moar`,
		`moar`,
		`moar`,
		nil,
		`moar`,
		`moar`,
		nil,
		nil,
		nil,
		nil,
	).Paren()

	offset := 13         // number of nil occurrences
	beforeLen := l.Len() // record preop len

	// verify no errors resulted from the attempt
	// to defragment our stack
	if err := l.Defrag().Err(); err != nil {
		t.Errorf("%s failed: %v", t.Name(), err)
		return
	}

	// verify starting length minus offset is equal
	// to the result length, meaning <offset> slices
	// were truncated.
	offsetLen := beforeLen - offset
	if afterLen := l.Len(); offsetLen != afterLen {
		t.Errorf("%s failed: unexpected length; want '%d', got '%d'",
			t.Name(), offsetLen, afterLen)
		return
	}
}

func TestStack_codecov(t *testing.T) {
	var s Stack
	// panic checks
	s.debug(``)
	s.debug(nil)
	s.error(``)
	s.error(nil)
	s.trace(``)
	s.trace(nil)
	s.state(``)
	s.state(nil)
	s.calls(``)
	s.calls(nil)
	s.debug(nil)
	s.Len()
	s.IsZero()
	s.IsInit()
	s.Valid()

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
	ll.unshift()
	ll.unshift(0)
	ll.unshift(65535)

	var lsys *logSystem = newLogSystem(sLogDefault)
	if lsys.isZero() {
		t.Errorf("%s failed: nil %T",
			t.Name(), lsys.logger())
	}

	SetDefaultStackLogLevel(`none`)
	SetDefaultStackLogLevel(0)
	SetDefaultStackLogLevel(nil)
	SetDefaultStackLogLevel('a')
}

func init() {
	// just used for testing, etc.
	//SetDefaultStackLogger(`stdout`)
	//SetDefaultConditionLogger(`stdout`)
}
