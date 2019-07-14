package assert

import (
	"fmt"
	"reflect"
	"runtime"
	"strings"

	"github.com/davecgh/go-spew/spew"
	"github.com/pkg/errors"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func (t *T) ErrorEquals(actual error, expected interface{}) {
	helper(t).Helper()
	hadError := actual != nil
	switch val := expected.(type) {
	case nil:
		if hadError {
			t.Fatalf("%s: unexpected error %+v", t.Name(), actual)
		}
	case bool:
		if val {
			if !hadError {
				t.Fatalf("%s: expected an error, but call succeeded", t.Name())
			}
		} else if hadError {
			t.Fatalf("%s: unexpected error %+v", t.Name(), actual)
		}
	case error:
		cause := errors.Cause(actual)
		if cause.Error() != val.Error() {
			t.Fatalf("%s: expected error %s but got %+v", t.Name(), val.Error(), actual)
		}
	case string:
		if val != "" {
			if !hadError {
				t.Fatalf("%s: expected an error containing %s, but call succeeded", t.Name(), val)
			} else if !strings.Contains(actual.Error(), val) {
				t.Fatalf("%s: expected an error containing %s, but got %+v", t.Name(), val, actual)
			}
		} else if hadError {
			t.Fatalf("%s: unexpected error %+v", t.Name(), actual)
		}
	default:
		t.Fatalf("ErrorEquals: expected must be an error or a string or nil")
	}
}

func (t *T) Nil(actual interface{}) {
	helper(t).Helper()
	if actual != nil {
		t.Fatalf("%s: should be nil, not %v", t.Name(), actual)
	}
}

func (t *T) NotNil(actual interface{}) {
	helper(t).Helper()
	if actual == nil {
		t.Fatalf("%s: should not be nil", t.Name())
	}
}

func (t *T) Len(actual interface{}, length int) {
	helper(t).Helper()
	t.Equal(length, reflect.ValueOf(actual).Len())
}

func (t *T) NotEqual(expected, actual interface{}, opts ...cmp.Option) {
	helper(t).Helper()
	// Add special flags for comparing structs from this service
	opts = append(opts,
		cmpopts.EquateEmpty(),
	)
	if cmp.Equal(actual, expected, opts...) {
		t.Fatalf("%s: actual equals expected:\n%s", t.Name(), spew.Sdump(expected))
	}
}

func (t *T) Equal(expected, actual interface{}, opts ...cmp.Option) {
	helper(t).Helper()
	// Add special flags for comparing structs from this service
	opts = append(opts,
		cmpopts.EquateEmpty(),
	)
	if !cmp.Equal(actual, expected, opts...) {
		diff := fmt.Sprintf("expected %s\ngot %s\n", spew.Sdump(expected), spew.Sdump(actual))
		if len(diff) > 200 {
			diff = cmp.Diff(actual, expected, opts...)
		}
		t.Fatalf("%s: differs: (-got +want)\n%s", t.Name(), diff)
	}
}

func (t *T) Contains(haystack, needle string) {
	helper(t).Helper()
	if !strings.Contains(haystack, needle) {
		t.Fatalf("%q does not contain %q", shortStr(haystack), shortStr(needle))
	}
}

func (t *T) Panics(f func(), msgContains string) {
	helper(t).Helper()

	if msg, ok := checkPanics(f); ok {
		if !strings.Contains(msg, msgContains) {
			t.Fatalf("panic message %q does not contain %q", shortStr(msg), shortStr(msgContains))
		}
	} else {
		t.Fatalf("expected function %s to panic", functionName(f))
	}
}

func functionName(f interface{}) string {
	funcValue := reflect.ValueOf(f)
	if funcValue.Kind() != reflect.Func {
		return fmt.Sprintf("%v is not a function", f)
	}
	// panics if f's Kind is not Chan, Func, Map, Ptr, Slice, or UnsafePointer.
	funcPointer := funcValue.Pointer()
	if runtimeFunc := runtime.FuncForPC(funcPointer); runtimeFunc != nil {
		return runtimeFunc.Name()
	}
	return "nil function"
}
