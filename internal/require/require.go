package require

import (
	"errors"
	"reflect"
	"testing"
	"time"
)

func NoError(t *testing.T, err error, msg string) {
	t.Helper()
	if err != nil {
		t.Fatalf("%s: unexpected error: %s", msg, err)
	}
}

func Error(t *testing.T, err error, msg string) {
	t.Helper()
	if err == nil {
		t.Fatalf("%s: got nil, want error", msg)
	}
}

func ErrorIs(t *testing.T, err, target error, msg string) {
	t.Helper()
	if !errors.Is(err, target) {
		t.Fatalf("%s: error should be in err chain:\n got: %s\n want: %s",
			msg, err, target)
	}
}

func Equal(t *testing.T, got, want any, msg string) {
	t.Helper()
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("%s: values not equal:\n got %v\n want %v", msg, got, want)
	}
}

func Bool(t *testing.T, got, want bool, msg string) {
	t.Helper()
	if got != want {
		t.Fatalf("%s: unexpected bool value: got %t, want: %t", msg, got, want)
	}
}

func Int(t *testing.T, got, want int, msg string) {
	t.Helper()
	if got != want {
		t.Fatalf("%s: unexpected int value: got %d, want: %d", msg, got, want)
	}
}

func Int64(t *testing.T, got, want int64, msg string) {
	t.Helper()
	if got != want {
		t.Fatalf("%s: unexpected int64 value: got %d, want: %d", msg, got, want)
	}
}

func String(t *testing.T, got, want string, msg string) {
	t.Helper()
	if got != want {
		t.Fatalf("%s: unexpected string value: got %s, want: %s", msg, got, want)
	}
}

func Time(t *testing.T, got, want time.Time, msg string) {
	t.Helper()
	if got != want {
		t.Fatalf("%s: unexpected string value: got %s, want: %s", msg, got, want)
	}
}

func Nil(t *testing.T, got any, msg string) {
	t.Helper()
	if !isNil(got) {
		t.Fatalf("%s: got %v, want nil", msg, got)
	}
}

// isNil checks if a specified object is nil or not, without Failing.
//
//nolint:exhaustive
func isNil(object interface{}) bool {
	if object == nil {
		return true
	}

	value := reflect.ValueOf(object)
	switch value.Kind() {
	case
		reflect.Chan, reflect.Func,
		reflect.Interface, reflect.Map,
		reflect.Ptr, reflect.Slice, reflect.UnsafePointer:

		return value.IsNil()
	default:
		return false
	}
}
