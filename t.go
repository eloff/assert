package assert

import (
	"path/filepath"
	"runtime"
	"testing"
)

var ShortStringLength = 50

type tHelper interface {
	Helper()
}

type T struct {
	*testing.T
	name     string
	testName string
}

func new(t *testing.T, name string, parallel bool) *T {
	if parallel {
		t.Parallel()
	}
	callerName := getCallerName()
	if name == "" {
		name = callerName
	} else {
		name = callerName + ": " + name
	}
	return &T{
		T:    t,
		name: name,
	}
}

func New(t *testing.T, name string) *T {
	return new(t, name, true)
}

func NewSerial(t *testing.T, name string) *T {
	return new(t, name, false)
}

func (t *T) SetName(name string) {
	t.testName = name
}

func (t *T) Name() string {
	if t.testName == "" {
		return t.name
	}
	return t.name + ": " + t.testName
}

func getCallerName() string {
	fpcs := make([]uintptr, 1)
	// Skip 3 levels to get the test function
	n := runtime.Callers(3, fpcs)
	if n == 0 {
		return ""
	}

	caller := runtime.FuncForPC(fpcs[0])
	if caller == nil {
		return ""
	}

	qualifiedName := caller.Name()
	funcName := filepath.Ext(qualifiedName)
	if len(funcName) != 0 && funcName[0] == '.' {
		funcName = funcName[1:]
	} else {
		funcName = qualifiedName
	}
	return funcName
}

type noopHelper struct{}

func (noopHelper) Helper() {}

func helper(t *T) tHelper {
	var iface interface{}
	iface = t
	if h, ok := iface.(tHelper); ok {
		return h
	}
	return noopHelper{}
}

func checkPanics(f func()) (msg string, ok bool) {
	defer func() {
		message := recover()
		if message != nil {
			ok = true
			msg, _ = message.(string)
		}
	}()
	// Call function that may panic
	f()
	return
}

func shortStr(s string) string {
	if len(s) > ShortStringLength {
		cutoff := ShortStringLength / 2
		return s[:cutoff] + "..." + s[len(s)-cutoff:]
	}
	return s
}
