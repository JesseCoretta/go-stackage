package stackage

import (
	"fmt"
	"testing"
)

func TestCondition_001(t *testing.T) {
	c := Cond(`person`, Eq, `Jesse`)
	want := `person = Jesse`
	got := c.String()
	if want != got {
		t.Errorf("%s failed: want '%s', got '%s'", t.Name(), want, got)
	}
}

func TestCondition_002(t *testing.T) {
	st := List().Paren().Push(
		Cond(`person`, Eq, `Jesse`),
	)

	want := `( person = Jesse )`
	got := st.String()
	if want != got {
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
	var c Condition
	c.NoNesting(true)

	if c.CanNest() {
		t.Errorf("%s failed: want '%t', got '%t'", t.Name(), false, true)
	}

	c.SetKeyword(`myKeyword`)
	c.SetOperator(Eq)
	c.SetExpression(And().Push(`this`, `that`))

	// value should NOT have been assigned.
	if c.Expression() != nil {
		t.Errorf("%s failed: want '%t', got '%t' [%T]", t.Name(), false, true, c.Expression())
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
	c.SetKeyword(`myKeyword`)
	c.SetOperator(Eq)
	c.SetExpression(And().Push(`this`, `that`))

	if !c.IsNesting() {
		t.Errorf("%s failed: want '%t', got '%t'", t.Name(), true, false)
	}
}

func TestCondition_IsPadded(t *testing.T) {
	cond := Cond(`person`, Eq, `Jesse`).Paren().Encap(`"`).NoPadding()

	if got := cond.IsPadded(); got {
		t.Errorf("%s failed: want 'false', got '%t'", t.Name(), got)
	}
}

func TestCondition_IsEncap(t *testing.T) {
	cond := Cond(`person`, Eq, `Jesse`).Paren().Encap(`"`).NoPadding()

	if got := cond.IsEncap(); !got {
		t.Errorf("%s failed: want 'true', got '%t'", t.Name(), got)
	}
}

func ExampleCondition_basic() {
	c := Cond(`person`, Eq, `Jesse`).Paren().Encap(`"`).NoPadding()
	fmt.Printf("%s", c)
	// Output: (person="Jesse")
}

func ExampleCondition_stepByStep() {
	var c Condition
	c.Paren()
	c.SetKeyword(`myKeyword`)
	c.SetOperator(Eq)
	c.SetExpression(`value123`)

	fmt.Printf("%s", c)
	// Output: ( myKeyword = value123 )
}
