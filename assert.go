package assert

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/davecgh/go-spew/spew"

	"github.com/ansel1/merry"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func (t *T) CheckError(actual error, expected interface{}) bool {
	helper(t).Helper()
	hadError := actual != nil
	switch val := expected.(type) {
	case nil:
		if hadError {
			t.T.Errorf("%s: unexpected error %s", t.Name, merry.Details(actual))
		}
	case bool:
		if val {
			if !hadError {
				t.T.Errorf("%s: expected an error, but call succeeded", t.Name)
			}
		} else if hadError {
			t.T.Errorf("%s: unexpected error %s", t.Name, merry.Details(actual))
		}
	case error:
		if !merry.Is(actual, val) && (actual == nil || actual.Error() != val.Error()) {
			details := "<nil>"
			if actual != nil {
				details = merry.Details(actual)
			}
			t.T.Errorf("%s: expected error %s but got %s", t.Name, val, details)
		}
	case string:
		if val != "" {
			if !hadError {
				t.T.Errorf("%s: expected an error containing %s, but call succeeded", t.Name, val)
			} else if !strings.Contains(actual.Error(), val) {
				t.T.Errorf("%s: expected an error containing %s, but got %s", t.Name, val, merry.Details(actual))
			}
		} else if hadError {
			t.T.Errorf("%s: unexpected error %s", t.Name, merry.Details(actual))
		}
	default:
		panic("fatal error: expected must be an error or a string or nil")
	}
	return !hadError
}

func (t *T) Nil(actual interface{}) bool {
	helper(t).Helper()
	if actual != nil {
		t.Errorf("%s: should be nil, not %v", t.Name, actual)
		return false
	}
	return true
}

func (t *T) NotNil(actual interface{}) bool {
	helper(t).Helper()
	if actual == nil {
		t.Errorf("%s: should not be nil", t.Name)
		return false
	}
	return true
}

func (t *T) Len(actual interface{}, length int) bool {
	helper(t).Helper()
	return t.Equal(length, reflect.ValueOf(actual).Len())
}

func (t *T) NotEqual(expected, actual interface{}, opts ...cmp.Option) bool {
	helper(t).Helper()
	// Add special flags for comparing structs from this service
	opts = append(opts,
		cmpopts.EquateEmpty(),
	)
	if cmp.Equal(actual, expected, opts...) {
		t.Errorf("%s: actual equals expected:\n%s", t.Name, spew.Sdump(expected))
		return false
	}
	return true
}

func (t *T) Equal(expected, actual interface{}, opts ...cmp.Option) bool {
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
		t.Errorf("%s: differs: (-got +want)\n%s", t.Name, diff)
		return false
	}
	return true
}

func (t *T) Contains(haystack, needle string) bool {
	helper(t).Helper()
	if !strings.Contains(haystack, needle) {
		t.Errorf("%q does not contain %q", shortStr(haystack), shortStr(needle))
		return false
	}
	return true
}

func (t *T) Panics(f func(), msgContains string) bool {
	helper(t).Helper()

	if msg, ok := checkPanics(f); ok {
		if !strings.Contains(msg, msgContains) {
			t.Errorf("panic message %q does not contain %q", shortStr(msg), shortStr(msgContains))
			return false
		}
		return true
	} else {
		t.Error("expected function to panic")
		return false
	}
}
