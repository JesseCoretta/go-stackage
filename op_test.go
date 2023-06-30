package stackage

import (
	_ "fmt"
	"testing"
)

type customUserOperator uint8

const (
	equalityMatch           customUserOperator = customUserOperator(Eq)
	greaterThanOrEqualMatch customUserOperator = customUserOperator(Ge)
	lessThanOrEqualMatch    customUserOperator = customUserOperator(Le)
	extensibleRuleDN        customUserOperator = iota + customUserOperator(Ge)
	extensibleRuleDNOID
	extensibleRuleAttr
	extensibleRuleNoDN
	approximateMatch
)

var filterMatchOps map[customUserOperator]string = map[customUserOperator]string{
	equalityMatch:           Eq.String(),
	greaterThanOrEqualMatch: `>=`,
	lessThanOrEqualMatch:    `<=`,
	extensibleRuleDN:        `:dn:=`,
	extensibleRuleDNOID:     `:dn:`,
	extensibleRuleAttr:      `:=`,
	extensibleRuleNoDN:      `:`,
	approximateMatch:        `~=`,
}

func (r customUserOperator) String() string {
	if val, found := filterMatchOps[r]; found {
		return val
	}
	return ``
}

func (r customUserOperator) Context() string {
	return `filter`
}

func TestOperator_customBasic(t *testing.T) {
	dn := `uid=jesse,ou=People,dc=example,dc=com`
	c := Cond(`attr`, extensibleRuleDN, dn).NoPadding()
	want := `attr:dn:=uid=jesse,ou=People,dc=example,dc=com`
	got := c.String()
	if want != got {
		t.Errorf("%s failed: want '%s', got '%s'", t.Name(), want, got)
	}
}
