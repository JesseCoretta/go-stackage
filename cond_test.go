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

func ExampleCondition_basic() {
	c := Cond(`person`, Eq, `Jesse`).Paren().Encap(`"`).NoPadding()
	fmt.Printf("%s", c)
	// Output: (person="Jesse")
}
