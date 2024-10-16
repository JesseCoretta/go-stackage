package stackage

import (
	"bytes"
	"fmt"
	// uncomment for TestStackagePerf runs
	//"log"
	//"net/http"
	//_ "net/http/pprof"
	"sort"
	"strconv"
	"strings"
	"testing"
	_ "time"
)

// Uncomment this test func (and the http+log imports above)
// to allow continuous pprof access via HTTP. Replace actions
// within the 'for loop' below with whatever expensive calls
// or activities you want to debug.
//
// You can run this test directly via the following command
// (which assumes your cwd is this package):
//
//   $ go test -run TestStackagePerf .
//
// While this test is running, in a separate terminal run
// the following command to acquire samples:
//
//   $ go tool pprof http://localhost:1234/debug/pprof/profile?seconds=60
//
// Alter this URL as needed (e.g.: use a different sampling
// duration, such as seconds=300).
//
// After N seconds (whatever duration you've used), you will
// be dropped into a pprof shell.  I like to enter 'svg' and
// obtain a profileNNN.svg file, but there are other options
// for use.
//
// Once you've been dropped into the pprof shell, this means
// the sample data has been acquired. At this point, you may
// kill the running test in the other terminal via CTRL+C,
// unless you plan on acquiring more samples.
//
// Don't forget to re-comment this function before running
// other tests else you'll hang at some point!  Definitely
// don't repackage/redistribute this package while this test
// function is UNcommented.
//func TestStackagePerf(t *testing.T) {
//        ch := make(chan bool)
//        go func() {
//		// edit this listener (e.g.: use TCP/8080) as
//		// seen fit.
//                log.Println(http.ListenAndServe(":1234", nil))
//        }()
//
//        for {
//        	custom := Cond(`outer`, Ne, customStack(And().Push(Cond(`keyword`, Eq, "somevalue"))))
//        	custom2 := Cond(`inner`, Eq, customStack(And().Push(Cond(`keyword`, Eq, "somevalue"))))
//		a1 := And().Push(
//                	`this4`,
//                	Not().Mutex().Push(
//                	        Or().Push(
//                	                custom2,
//                	                Cond(`dayofweek`, Ne, "Wednesday"),
//                	                Cond(`ssf`, Ge, "128"),
//                	                Cond(`greeting`, Ne, List().Push(List().Push(List().Push(``)))),
//                	        ),
//                	),
//		)
//                a2 := And().Push(
//                        `this4`,
//                        Not().Mutex().Push(
//                                Or().Push(
//                                        custom2,
//                                        Cond(`dayofweek`, Ne, "Wednesday"),
//                                        Cond(`ssf`, Gt, "128"),
//                                        Cond(`greeting`, Ne, List().Push(List().Push(List().Push(``)))),
//                                ),
//                        ),
//                )
//
//		//thisIsMyNightmare := And().Push(
//		_ = And().Push(
//		        `this1`,
//		        Or().Mutex().Push(
//		                custom,
//				a1,
//		                And().Push(
//		                        Or().Push(
//		                                Cond(`keyword2`, Lt, "someothervalue"),
//		                        ),
//		                ),
//		                Cond(`keyword`, Gt, Or().Push(`...`)),
//		        ),
//		        `this2`,
//		)
//
//		custom.IsEqual(custom2)
//		a1.IsEqual(a2)
//
//		//os.Stdout.Write([]byte(thisIsMyNightmare.String()))
//        }
//        <-ch
//}

/*
This example demonstrates basic support for stack sorting via the
[sort.Stable] method using basic string values.
*/
func ExampleStack_String_withSort() {
	var names Stack = List().SetDelimiter(' ')
	names.Push(`Frank`, `Anna`, `Xavier`, `Betty`, `aly`, `Jim`, `fargus`)
	sort.Stable(names)
	fmt.Println(names)
	// Output: Anna Betty Frank Jim Xavier aly fargus
}

/*
This example demonstrates the means of converting a custom [Stack]-alias
instance into a native [Stack] that contains the same values. Optionally,
users may shadow the first return value, and use the second (bool) value
to ascertain whether the instance was converted successfully.
*/
func ExampleConvertStack() {
	var native Stack = List().Push(`this`, `is`, `a`, `native`, `stack`)
	type YourStack Stack

	// A custom instance could have been
	// assembled in a variety of ways.
	//
	// For this example, we'll just use a
	// basic cast just for brevity.
	custom := YourStack(native)

	back, ok := ConvertStack(custom)
	if !ok {
		fmt.Printf("Failed to convert %T", custom)
		return
	}

	slice, found := back.Index(0) // any index would do
	if !found {
		fmt.Printf("Content not found")
		return
	}

	fmt.Printf("Value: %v", slice)
	// Output: Value: this
}

/*
This example demonstrates the means for assigning a custom marshaler
function or method to a [Stack] instance.
*/
func ExampleStack_SetMarshaler() {
	var r Stack = List()

	// Lets write a marshaler that converts
	// string numbers to ints, each of which
	// are added to Stack r.
	r.SetMarshaler(func(in ...any) (err error) {
		if len(in) == 0 {
			err = fmt.Errorf("No input content")
			return
		}

		for i := 0; i < len(in); i++ {
			if assert, ok := in[i].(string); ok {
				var num int
				if num, err = strconv.Atoi(assert); err != nil {
					break
				} else {
					r.Push(num)
				}
			} else {
				err = fmt.Errorf("Cannot assert %T as a string", in[i])
				break
			}
		}

		return
	})

	if err := r.Marshal(`1`, `2`, `3`, `4`); err != nil {
		fmt.Println(err)
		return
	}

	slice, _ := r.Index(1)

	fmt.Printf("Slice is %T", slice)
	// Output: Slice is int
}

/*
This example demonstrates the means for assigning a custom unmarshaler
function or method to a [Stack] instance.
*/
func ExampleStack_SetUnmarshaler() {
	var r Stack = List().Push(`1`, `2`, `3`, `4`)
	// Write an unmarshaler that extracts just the
	// values, depositing them within the []any
	// output instance.
	r.SetUnmarshaler(func(_ ...any) (out []any, err error) {
		for i := 0; i < r.Len(); i++ {
			slice, _ := r.Index(i)
			assert, ok := slice.(string)
			// We only want string types
			if ok {
				out = append(out, assert)
			} else {
				err = fmt.Errorf("Unsupported slice type: %T", slice)
				break
			}
		}

		return
	})

	out, err := r.Unmarshal()
	if err != nil {
		fmt.Println(err)
		return
	}

	lout := len(out)
	fmt.Printf("Output has %d strings: %t", lout, lout == 4)
	// Output: Output has 4 strings: true
}

/*
This example demonstrates a simple "swapping" of indexed values by way
of the [Stack.Swap] method.
*/
func ExampleStack_Swap() {
	var names Stack = List().SetDelimiter(' ')
	names.Push(`Frank`, `Anna`, `Xavier`, `Betty`, `aly`, `Jim`, `fargus`)
	names.Swap(0, 5)
	slice, _ := names.Index(0)
	fmt.Println(slice)
	// Output: Jim
}

func ExampleStack_Free() {
	var names Stack = List() // initialize
	names.Push(`Frank`, `Anna`, `Xavier`, `Betty`, `aly`, `Jim`, `fargus`)
	names.Free()
	fmt.Printf("%T is zero: %t", names, names.IsZero())
	// Output: stackage.Stack is zero: true
}

/*
This example demonstrates a simple order-based comparison using two
string values by way of the [Stack.Less] method.  In this scenario,
a value of false is returned because "Frank" does not alphabetically
precede "Anna".

Note that the semantics of [Stack.Less] apply, in that certain conditions
with regards to the values must be satified, else a custom Less closure
needs to be devised by the end-user.  See the [Stack.SetLessFunc] method.
*/
func ExampleStack_Less() {
	var names Stack = List().SetDelimiter(' ')
	names.Push(`Frank`, `Anna`, `Xavier`, `Betty`, `aly`, `Jim`, `fargus`)
	fmt.Println(names.Less(0, 1)) // is Frank before Anna?
	// Output: false
}

/*
This example demonstrates custom stack sorting by way of a user-defined
closure set within the stack via the [Stack.SetLessFunc] value, thereby
allowing support with the [sort.Stable] method, et al.
*/
func ExampleStack_SetLessFunc() {
	var names Stack = List().SetDelimiter(' ')
	names.Push(`Frank`, `Anna`, `Xavier`, `Betty`, `aly`, `Jim`, `fargus`)

	// Submit our custom closure
	names.SetLessFunc(func(i, j int) bool {
		// Note we are just presuming the values are
		// strings merely for simplicity.  Index and
		// type checks should always be done in real
		// life scenarios, else you risk panics, etc.
		slice1, _ := names.Index(i)
		slice2, _ := names.Index(j)
		s1 := strings.ToLower(slice1.(string))
		s2 := strings.ToLower(slice2.(string))
		switch strings.Compare(s1, s2) {
		case -1:
			return true
		}
		return false
	})

	sort.Stable(names)
	fmt.Println(names)
	// Output: aly Anna Betty fargus Frank Jim Xavier
}

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

func ExampleStack_Addr() {
	var c Stack = List().Push(`this`, `and`, `that`)
	fmt.Printf("Address ID has '0x' prefix: %t", c.Addr()[:2] == `0x`)
	// Output: Address ID has '0x' prefix: true
}

func ExampleStack_Reverse() {
	var c Stack = List().Push(0, 1, 2, 3, 4)
	fmt.Printf("%s", c.Reverse())
	// Output: 4 3 2 1 0
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

var testParens []string = []string{`(`, `)`}

type customStack Stack // simulates a user-defined type that aliases a Stack

func (r customStack) String() string {
	return Stack(r).String()
}

type customStruct struct {
	Type  string
	Value any
}

func (r customStruct) String() string {
	return `struct_value`
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

	_, _ = got.Traverse(0, 2, 6, 19)
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
	s := List().SetDelimiter(rune(44)).NoPadding(true).Push(
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

func TestStack_Traverse(t *testing.T) {

	stk := Basic().Push(
		1,
		Basic().Push(2),
		Cond(`keyword`, Eq, Basic().Push(3)),
	)

	slice, _ := stk.Traverse(0)
	if _, ok := slice.(int); !ok {
		t.Errorf("%s failed: unexpected type %T", t.Name(), slice)
		return
	}

	slice, _ = stk.Traverse(1)
	if _, ok := slice.(Stack); !ok {
		t.Errorf("%s failed: unexpected type %T at depth 2", t.Name(), slice)
		return
	}

	slice, _ = stk.Traverse(1, 0)
	if _, ok := slice.(int); !ok {
		t.Errorf("%s failed: unexpected type %T at depth 2", t.Name(), slice)
		return
	}

	slice, _ = stk.Traverse(2, 0)
	if _, ok := slice.(int); !ok {
		t.Errorf("%s failed: unexpected type %T at depth 3", t.Name(), slice)
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

func TestStack_IsEqual(t *testing.T) {
	channel := make(chan error, 2)

	var str1 *string = new(string)
	*str1 = `this is crazy`

	var str2 *string = new(string)
	*str2 = `this is crazy`

	var iface1 any = struct{}{}
	var iface2 any = struct{}{}

	A := And().Paren().Push(
		Cond(`Test`, Eq, List().Push(
			`sub_nested_number_0`,
			`sub_nested_number_1`,
		)),
		[]string{`1`, `2`, `3`, `4`},
		iface1,
		map[int]any{
			0: map[any]any{
				uint8(0):  `hehe`,
				uint16(1): `omg`,
				true:      false,
				`runes`:   []rune{'H', 'E', 'L', 'P', ' ', 'M', 'E'},
			},
			1: nil,
			2: []uint{1, 2, 3, 4},
		},
		Or().Paren().Push(
			And().Push(`deep_string`),
			`sub_element_number_0`,
			`sub_element_number_1`,
		),
		//TestStack_IsEqual,
		&str1, // **string
		channel,
		struct{}{},
		&struct{}{},
		rune(76),
		nil,
	)

	B := And().Paren().Push(
		Cond(`Test`, Eq, List().Push(
			`sub_nested_number_0`,
			`sub_nested_number_1`,
		)),
		[4]string{`1`, `2`, `3`, `4`},
		iface2,
		map[int]any{
			0: map[any]any{
				uint8(0):  `hehe`,
				uint16(1): `omg`,
				true:      false,
				`runes`:   []rune{'H', 'E', 'L', 'P', ' ', 'M', 'E'},
			},
			1: nil,
			2: []uint{1, 2, 3, 4},
		},
		Or().Paren().Push(
			And().Push(`deep_string`),
			`sub_element_number_0`,
			`sub_element_number_1`,
		),
		//TestStack_Unmarshal_default,
		&str2, // **string
		channel,
		struct{}{},
		&struct{}{},
		rune(76),
		nil,
	)

	if err := A.IsEqual(B); err != nil {
		t.Errorf("%s failed: %v", t.Name(), err)
		return
	}
}

func TestStack_Unmarshal_default(t *testing.T) {

	// Make a stack with senseless garbage
	A := And().Paren().Push(
		// Condition with Stack expr
		Cond(`Test`, Eq, List().Push(
			`sub_nested_number_0`,
			`sub_nested_number_1`,
		)),
		// just a string
		`top_element_number_0`,
		// Stack
		Or().Paren().Push(
			And().Push(`deep_string`),
			`sub_element_number_0`,
			`sub_element_number_1`,
		),
		// A random rune
		rune(76),
		// A byte -- 0x2B (ASCII #43, "+")
		byte(43),
	)

	slices, err := A.Unmarshal()
	if err != nil {
		t.Errorf("%s failed: %v", t.Name(), err)
		return
	}

	want := 6
	if got := len(slices); got != want {
		t.Errorf("%s failed: want %d, got %d", t.Name(), want, got)
		return
	}

	var n Stack
	if err = n.Marshal(slices); err != nil {
		t.Errorf("%s failed: %v", t.Name(), err)
		return
	}

	//printf("Top-Stack: %s(%d)\n", n.Kind(), n.Len())
	//n.debugStack(0)
	//printf("End-of-Stack\n")
}

/*
func (r Stack) debugStack(lvl int) {
	for i := 0; i < r.Len(); i++ {
		slice, _ := r.Index(i)
		switch tv := slice.(type) {
		case Condition:
			expr := tv.Expression()
			if tv.IsNesting() {
				if sub, ok := stackTypeAliasConverter(expr); ok {
					printf("%d:%d - Condition-Nested-Stack::[KW]%s [OP]%s [%s(%d)]\n",
						i, lvl, tv.Keyword(), tv.Operator(), sub.Kind(), sub.Len())
					sub.debugStack(lvl+1)
					printf("%d:%d - End-of-Condition-Nested-Stack\n", i, lvl)
				}
			} else {
				printf("%d:%d - Condition::[%s]\n", lvl, i, tv)
	                        if isKnownPrimitive(tv) {
	                                printf("%d:%d - Primitive-Condition-Expr: %s\n",
						i, lvl,	primitiveStringer(tv))
	                        } else if meth := getStringer(expr); meth != nil {
	                                printf("%d:%d - Stringer-Condition-Expr: %s\n", i, lvl, meth())
	                        }
			}
		case Stack:
			printf("%d:%d - Nested-Stack(%s)[%d]::\n",
				i, lvl, tv.Kind(), tv.Len())
			tv.debugStack(lvl+1)
			printf("%d:%d - End-of-Nested-Stack\n", i, lvl)
		default:
			if isKnownPrimitive(tv) {
				printf("-:%d - Primitive-Stack-Value: (%T) %s\n",
					i, tv, primitiveStringer(tv))
			} else if meth := getStringer(tv); meth != nil {
				printf("-:%d - Stringer-Stack-Value: %s\n", i, meth())
			} else {
				printf("-:%d - NoStringer-Stack-Value: %T\n", i, tv)
			}
		}
	}
}
*/

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
This example demonstrates the creation of a List stack, a stack type suitable for
general use.
*/
func ExampleList() {
	l := List().Push(
		1.234,
		1+3i,
		`hello mr thompson`,
	)

	l.SetDelimiter('?')

	// alternatives ...
	//l.JoinDelim(`,`) 	//strings ok too!
	//l.JoinDelim(`delim`)

	fmt.Printf("%s", l)
	// Output: 1.234 ? (1+3i) ? hello mr thompson
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

	A := List().NoPadding(true).Push(
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
	A := List().SetDelimiter(`,`).NoPadding().Push(
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

func TestFactorNegIndex(t *testing.T) {
	b := Basic().NegativeIndices(true)
	b.Push(
		float64(3.14159),
		float64(-9.378),
		float64(109.9103),
	)

	// +1 offset for hidden cfg slice
	for i, v := range map[int]int{
		-1:  2,
		-2:  1,
		-3:  0,
		-4:  2,
		-5:  1,
		-12: 0,
		-24: 0,
	} {
		I := factorNegIndex(i, b.Len()) - 1
		if I != v && I-(I*2) <= b.Len() {
			t.Errorf("%s failed [key:%d]: want '%d', got '%d'", t.Name(), i, v, I)
			return
		}
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

/*
This example demonstrates the creation of a list stack
using comma delimitation.
*/
func ExampleStack_SetDelimiter() {
	// note: one could also use a rune
	// e.g: ',' or rune(44) for comma.
	L := List().SetDelimiter(`+`).Push(
		`item1`,
		`item2`,
	)
	fmt.Printf("%s", L)
	// Output: item1 + item2
}

func ExampleStack_Transfer() {
	var source Stack = Basic().Push(1, 2, 3, 4)
	var dest Stack = Basic()
	source.Transfer(dest)
	slice, _ := dest.Index(2)
	fmt.Printf("%d", slice.(int))
	// Output: 3
}

func ExampleStack_Insert() {
	var stk Stack = Basic().Push(1, 2, 3, 5)
	add := 4
	idx := 3
	stk.Insert(add, idx)
	slice, _ := stk.Index(idx)
	fmt.Printf("%d", slice.(int))
	// Output: 4
}

func ExampleStack_ForwardIndices() {
	var stk Stack = Basic().Push(1, 2, 3, 4)
	stk.ForwardIndices(true)
	slice, _ := stk.Index(1000)
	fmt.Printf("%d", slice.(int))
	// Output: 4
}

func ExampleStack_NegativeIndices() {
	var stk Stack = Basic().Push(1, 2, 3, 4)
	stk.NegativeIndices(true)
	slice, _ := stk.Index(-1)
	fmt.Printf("%d", slice.(int))
	// Output: 4
}

func ExampleStack_SetID_random() {
	var stk Stack = Basic().Push(1, 2, 3, 4)

	// can't predict what ID will
	// be, so we'll check length
	// which should always be 24.
	stk.SetID(`_random`)
	fmt.Printf("Random ID len: %d", len(stk.ID()))
	// Output: Random ID len: 24
}

func ExampleStack_SetID_pointerAddress() {
	var stk Stack = Basic().Push(1, 2, 3, 4)

	// can't predict what ID will be,
	// so we'll check the prefix to
	// be certain it begins with '0x'.
	stk.SetID(`_addr`)
	fmt.Printf("Address ID has '0x' prefix: %t", stk.ID()[:2] == `0x`)
	// Output: Address ID has '0x' prefix: true
}

func ExampleStack_Category() {
	var stk Stack = Basic().Push(1, 2, 3, 4)
	stk.SetCategory(`basic_stuff`)
	fmt.Printf("Category: %s", stk.Category())
	// Output: Category: basic_stuff
}

func ExampleStack_SetCategory() {
	var stk Stack = Basic().Push(1, 2, 3, 4)
	stk.SetCategory(`basic_stuff`)
	fmt.Printf("Category: %s", stk.Category())
	// Output: Category: basic_stuff
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
	custom := Cond(`outer`, Ne, customStack(And().Push(Cond(`keyword`, Eq, "somevalue"))))

	thisIsMyNightmare := And().Push(
		`this1`,
		Or().Mutex().Push(
			custom,
			And().Push(
				`this4`,
				Not().Mutex().Push(
					Or().Push(
						Cond(`dayofweek`, Ne, "Wednesday"),
						Cond(`ssf`, Ge, "128"),
						Cond(`greeting`, Ne, List().Push(List().Push(List().Push(``)))),
					),
				),
			),
			And().Push(
				Or().Push(
					Cond(`keyword2`, Lt, "someothervalue"),
				),
			),
			Cond(`keyword`, Gt, Or().Push(`...`)),
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
		//{2, [][]int{{1, 1, 1, 0, 2, 0, 0, 0}, {1, 1, 1, 0, 2}}, ``, nil},
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
		Cond(`keyword`, Eq, Basic().Push(1, 2, nil, 3, nil, nil, 4)),
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
	if err := l.Defrag(-1).Err(); err != nil {
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

func TestStack_withCap(t *testing.T) {
	src := Basic().Push(
		`element0`,
		`element1`,
		`element2`,
		`element3`,
		`element4`,
		`element5`,
		`element6`,
		`element7`,
		`element8`,
		`element9`,
	)

	dst := Basic(3).Push(
		`element0a`,
		`element0b`,
		`element0c`,
	)

	src.Transfer(dst)

	if dst.Len() != 3 {
		t.Errorf("%s failed; transfer process bypassed capacity constraints; want len:%d, got len:%d",
			t.Name(), 3, dst.Len())
		return
	}

	dst.Insert(struct{}{}, 2)
	if dst.Len() != 3 {
		t.Errorf("%s failed; insert process bypassed capacity constraints; want len:%d, got len:%d",
			t.Name(), 3, dst.Len())
		return
	}
}

func TestStack_codecov(t *testing.T) {
	var s Stack
	// panic checks
	s.Len()
	s.IsZero()
	s.IsInit()
	s.Valid()
	s.Pop()

	SetDefaultStackLogger(nil)
	_ = lonce.String()

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

	s.CanMutex()
	s.Avail()
	s.string()
	s.IsEncap()
	s.SetAuxiliary(nil)
	s.SetLogger(nil)
	s.NoNesting()
	s.NoNesting(true)
	s.NoNesting(false)
	s.Logger()
	s.IsEqual(nil)

	s = List()
	s.SetUnmarshaler()
	s.SetMarshaler()
	s.Marshal()
	s.Unmarshal()
	s.IsEqual(struct{}{})
	s.IsEqual(s)
	s.Pop()
	s.Push(-1, -2, -3, -4, -5)
	s.Front()
	s.Back()
	s.NegativeIndices(true)

	for i := 0; i < s.Len(); i++ {
		offset := i - (i + 1)

		slice, _ := s.Index(i)
		want, _ := slice.(int)

		slice, _ = s.Index(offset - want)
		got, _ := slice.(int)

		if want != got {
			t.Errorf("%s failed [loop:%d]: want %d, got %d",
				t.Name(), i, want, got)
			return
		}
	}

	s.config()
	s.IsEncap()
	s.ReadOnly()
	s.SetReadOnly(true)
	s.SetReadOnly(false)
	s.IsReadOnly()
	s.Avail()
	s.traverse()
	s.string()
	s.Paren()
	s.Paren(true)
	s.Paren(false)
	s.NoNesting()
	s.NoNesting(true)
	s.CanNest()
	s.NoNesting(false)
	s.CanMutex()

	s.SetReadOnly()
	s.SetReadOnly(true)
	s.Free()
	s.SetReadOnly(false)

	s.Push(customStruct{`keyword`, `vaLUE`})
	_ = s.String()

	s.SetAuxiliary(nil)
	s.SetLogger(nil)
	s.Logger()

	SetDefaultStackLogLevel(`none`)
	SetDefaultStackLogLevel(0)
	SetDefaultStackLogLevel(nil)
	SetDefaultStackLogLevel('a')
	s.SetLogger(`stderr`)
	s.SetLogger(2)
	s.SetLogger(`stdout`)
	s.SetLogger(1)
	s.SetLogger(sLogDefault)

	s.Less(3, 6)
	s.Less(6, 3)
	s.Less(1, 2)
	s.Swap(101, 48)
	s.Swap(-1, 480)
	s.Push(
		customStack(List().Push(1, `&`)),
		customStack(List().Push(`_n`, `?!?`, 5)),
		``)
	sort.Stable(s)
	s.Less(0, 1)
	s.Less(0, 2)
	s.SetLessFunc(nil)
	s.SetLessFunc()

	s.SetLogLevel(8, 16)
	s.UnsetLogLevel(8, 16)

	s.SetLogLevel(LogLevel4, LogLevel5)
	s.LogLevels()
	s.UnsetLogLevel(LogLevel4, LogLevel5)
	s.Logger()

	s.SetErr(errorf(`this is a serious error`))         // txt err
	s.SetErr(errorf(errorf(`this is a serious error`))) // err err
	s.Err()
	s.SetErr(nil)

	k := List(8).Push(`1`, `2`, `3`, `4`).SetFIFO(true)
	k.Front()
	k.Back()

	l := List().Push(`1`, `2`, `3`, `4`)
	l.IsEqual(k)
	l.SetEqualityPolicy(func(x, y any) error {
		return nil
	})
	l.IsEqual(k)

	m := And().Push(`1`, `2`, `3`, `4`)
	m.SetEqualityPolicy()
	m.IsEqual(l)
	m.Marshal()

	var ak Stack
	ak.Marshal([]any{`CONDITION`, `Keyword`, Eq, `value`})

	ak = List()
	ak.Marshal([]any{`CONDITION`, `Keyword`, Eq, `value`})
	ak.Marshal([]any{`other`, `stack`})
	ak.Marshal([]any{`CONDITION`, `Keyword`, Eq}) // missing value
	ak.Marshal()
	ak.Marshal([]any{5, `bogus`})
	marshalDefault([]any{})

	ak = stackByWord(`NOT`)
	ak = stackByWord(`blarg`)

	s.SetLogger(`off`)
	s.SetLogger(0)
}

func init() {
	// just used for testing, etc.
	//SetDefaultStackLogger(`stdout`)
	//SetDefaultConditionLogger(`stdout`)
}
