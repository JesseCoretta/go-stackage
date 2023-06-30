package stackage

import (
	"fmt"
	"testing"
	_ "time"
)

var testParens []string = []string{`(`, `)`}

type customStack Stack // simulates a user-defined type that aliases a Stack

func (r customStack) String() string {
	return Stack(r).String()
}

func TestStack_And001(t *testing.T) {
	got := And().Paren().Fold().Push(
		`testing1`,
		`testing2`,
		`testing3`,
	)
	want := `(testing1 AND testing2 AND testing3)`
	if got.String() != want {
		t.Errorf("%s failed: want '%s', got '%s'", t.Name(), want, got)
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

	want := `(top_element_number_0 AND (sub_element_number_0 OR sub_element_number_1))`
	if got := A; got.String() != want {
		t.Errorf("%s failed: want '%s', got '%s'", t.Name(), want, got)
	}
}

func TestAnd_002(t *testing.T) {

	A := And().Paren().Push(
		`top_element_number_0`,

		// make each OR condition look like "<...>" (including quotes)
		Or().Paren().Encap(`"`, []string{`<`, `>`}).Push(
			`sub_element_number_0`,
			`sub_element_number_1`,
		),
	)

	want := `(top_element_number_0 AND ("<sub_element_number_0>" OR "<sub_element_number_1>"))`
	if got := A; got.String() != want {
		t.Errorf("%s failed: want '%s', got '%s'", t.Name(), want, got)
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
	}
}

/*
This example demonstrates ORed stack values using the double pipe (||) symbol
and custom value encapsulation.
*/
func ExampleOr_pipes() {
	or := Or().
		Paren().                                            // Add parenthesis
		Symbol(`||`).                                       // Use double pipes for OR
		Encap(` `, `"`).                                    // Encapsulate individual vals with double-quotes surrounded by spaces
		Push(`cn`, `sn`, `givenName`, `objectClass`, `uid`) // Push these values now

	fmt.Printf("%s\n", or)
	// Output: ( "cn" || "sn" || "givenName" || "objectClass" || "uid" )
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
	}

	slice, ok := A.Traverse(2, 0, 1)
	if !ok {
		t.Errorf("%s failed: got '%T' during traversal, want non-nil", t.Name(), slice)
	}

	strAssert, strOK := slice.(string)
	if !strOK {
		t.Errorf("%s failed: want '%T', got '%T' during assertion", t.Name(), ``, slice)
	}

	want = `unwanted_element_number_1`
	if want != strAssert {
		t.Errorf("%s failed: want '%T', got '%T'", t.Name(), want, strAssert)
	}
}

func TestList_001(t *testing.T) {

	A := List().JoinDelim(` `).Push(
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

func TestCustomStack001(t *testing.T) {
	A := List().JoinDelim(`,`).Push(
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
		//.SetMessageChan(ch)
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
	slice, _ := filter.Traverse(1, 1)  // Enter coordinates
	condAssert, _ := slice.(Condition) // The return is any, so assert to what we expect
	fmt.Printf("%s", condAssert)       // use its String method automagically
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
	}
}

func TestBasic_withCapacity(t *testing.T) {
	b := Basic(2)
	b.Push(
		float64(3.14159),
		float64(-9.378),
		float64(139.104),
	)

	if b.Len() != 2 || !b.CapReached() {
		t.Errorf("%s failed: maximum capacity not honored; want len:%d, got len:%d", t.Name(), 2, b.Len())
	}
}

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
	L := List().JoinDelim(`,`).Push(
		`item1`,
		`item2`,
	)
	want := `item1,item2`
	got := L.String()
	if want != got {
		t.Errorf("%s failed: want '%s', got '%s'", t.Name(), want, got)
	}
}

func ExampleStack_JoinDelim() {
	L := List().JoinDelim(`,`).Push(
		`item1`,
		`item2`,
	)
	fmt.Printf("%s", L)
	// Output: item1,item2
}
