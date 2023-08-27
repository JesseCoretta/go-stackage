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

func TestStack_noInit(t *testing.T) {
	var x Stack
	x.Push(`well?`) // first off, this should not panic

	// next, make sure it really
	// did not work ...
	if length := x.Len(); length != 0 {
		t.Errorf("%s failed [noInit]: want '%d' elements, got '%d'", t.Name(), 0, length)
	}

	// lastly, we should definitely see an indication
	// of this issue during validity checks ...
	if err := x.Valid(); err == nil {
		t.Errorf("%s failed [noInit]: want 'error', got '%T'", t.Name(), nil)
	}
}

func TestStack_And001(t *testing.T) {
	got := And().Paren().Fold().Push(
		`testing1`,
		`testing2`,
		`testing3`,
	)

	want := `( testing1 AND testing2 AND testing3 )`
	if got.String() != want {
		t.Errorf("%s failed: want '%s', got '%s'", t.Name(), want, got)
	}
}

func TestStack_Insert(t *testing.T) {
	got := And().Paren().Fold().Push(
		`testing1`,
		`testing2`,
		`testing3`,
	)

	got.Insert(`testing0`, 0)

	want := `( testing0 AND testing1 AND testing2 AND testing3 )`
	if got.String() != want {
		t.Errorf("%s.1 failed: want '%s', got '%s'", t.Name(), want, got)
	}

	got.Insert(`testing2.5`, 3)
	want = `( testing0 AND testing1 AND testing2 AND testing2.5 AND testing3 )`
	if got.String() != want {
		t.Errorf("%s.2 failed: want '%s', got '%s'", t.Name(), want, got)
	}

	got.Insert(`testing4`, 15)
	want = `( testing0 AND testing1 AND testing2 AND testing2.5 AND testing3 AND testing4 )`
	if got.String() != want {
		t.Errorf("%s.3 failed: want '%s', got '%s'", t.Name(), want, got)
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
	}

	if ok := s.Replace(`testingX`, 6); ok {
		t.Errorf("%s.2 failed: want '%t', got '%t'", t.Name(), false, ok)
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
	}

	stk.Reset() // reset source, removing all slices
	if or.Len() != 3 && stk.Len() != 0 {
		t.Errorf("%s failed [post-transfer src reset]: want slen:%d; dlen:%d, got slen:%d; dlen:%d",
			t.Name(), 0, 3, stk.Len(), or.Len())
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
	}

	if !O.IsEncap() {
		t.Errorf("%s failed [IsEncap]: want '%t', got '%t'", t.Name(), true, false)
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
	cMaker := func(r *Condition) *Condition {
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
	condAssert, _ := slice.(*Condition) // The return is any, so assert to what we expect
	fmt.Printf("%s", condAssert)        // use its String method automagically
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

func TestBasic_availableCapacity(t *testing.T) {
	b := Basic(2)
	b.Push(
		float64(3.14159),
		float64(-9.378),
		float64(139.104),
	)

	if b.Avail() != 0 {
		t.Errorf("%s failed: unexpected available slice count; want len:%d, got len:%d", t.Name(), 2, b.Avail())
	}

	b = Basic(5)
	if b.Avail() != 5 {
		t.Errorf("%s failed: unexpected available slice count; want len:%d, got len:%d", t.Name(), 5, b.Avail())
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
	}

	b.Reset()

	if b.Len() != 0 {
		sl, _ := b.Index(0)
		t.Errorf("%s failed: want '%d', got: '%d' [%#v]", t.Name(), 0, b.Len(), sl)
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
