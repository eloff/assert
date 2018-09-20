package assert

import (
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

var ShortStringLength = 50

type tHelper interface {
	Helper()
}

type T struct {
	*testing.T
	assertions *assert.Assertions
	Name       string
}

func New(t *testing.T, name string) *T {
	callerName := getCallerName()
	return &T{
		T:          t,
		assertions: assert.New(t),
		Name:       callerName + ": " + name,
	}
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
