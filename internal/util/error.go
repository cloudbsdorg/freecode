package util

import (
	"fmt"
	"reflect"
	"runtime"
	"strings"
)

// ErrorFormat formats an error to its stack trace or "name: message" representation.
// Returns empty string for nil errors.
func ErrorFormat(err error) string {
	if err == nil {
		return ""
	}

	if e, ok := err.(interface{ StackTrace() string }); ok && e.StackTrace() != "" {
		return fmt.Sprintf("%s\n%s", err.Error(), e.StackTrace())
	}

	if stack := getStackTrace(err); stack != "" {
		return fmt.Sprintf("%s\n%s", err.Error(), stack)
	}

	return err.Error()
}

// getStackTrace attempts to extract stack trace from error interface.
func getStackTrace(err error) string {
	t := reflect.TypeOf(err)
	if t.Kind() != reflect.Ptr {
		return ""
	}

	errVal := reflect.ValueOf(err)
	stackField := errVal.Elem().FieldByName("Stack")
	if stackField.IsValid() && stackField.CanInterface() {
		if stack, ok := stackField.Interface().([]uintptr); ok && len(stack) > 0 {
			var frames []string
			for _, pc := range stack {
				if fn := runtimeFrame(pc); fn != "" {
					frames = append(frames, fn)
				}
			}
			if len(frames) > 0 {
				return strings.Join(frames, "\n")
			}
		}
	}

	return ""
}

func runtimeFrame(pc uintptr) string {
	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return ""
	}
	name := fn.Name()
	if idx := strings.LastIndex(name, "."); idx >= 0 {
		name = name[idx+1:]
	}
	file, line := fn.FileLine(pc)
	return fmt.Sprintf("%s (%s:%d)", name, file, line)
}

// ErrorMessage extracts a clean error message from various error types.
// Returns empty string for nil errors.
func ErrorMessage(err error) string {
	if err == nil {
		return ""
	}

	if e, ok := err.(interface{ Message() string }); ok {
		if msg := e.Message(); msg != "" {
			return msg
		}
	}

	if e, ok := err.(interface{ error }); ok {
		if msg := e.Error(); msg != "" {
			return msg
		}
	}

	return err.Error()
}

// ErrorData extracts structured error data including type, message, stack, and cause.
// Returns nil for nil errors.
func ErrorData(err error) map[string]interface{} {
	if err == nil {
		return nil
	}

	data := make(map[string]interface{})

	data["type"] = getErrorType(err)
	data["message"] = ErrorMessage(err)
	data["formatted"] = ErrorFormat(err)

	if stack := getStackTrace(err); stack != "" {
		data["stack"] = stack
	}

	if cause := extractCause(err); cause != nil {
		data["cause"] = ErrorFormat(cause)
	}

	extractCustomData(err, data)

	return data
}

// getErrorType returns the error type name.
func getErrorType(err error) string {
	t := reflect.TypeOf(err)
	if t == nil {
		return "error"
	}
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return t.Name()
}

// extractCause extracts the cause from an error if it implements the causer interface
// or has a Cause field.
func extractCause(err error) error {
	if err == nil {
		return nil
	}

	if c, ok := err.(interface{ Cause() error }); ok {
		return c.Cause()
	}

	errVal := reflect.ValueOf(err)
	if errVal.Kind() != reflect.Ptr {
		return nil
	}

	causeField := errVal.Elem().FieldByName("Cause")
	if causeField.IsValid() && causeField.CanInterface() {
		if cause, ok := causeField.Interface().(error); ok && cause != nil {
			return cause
		}
	}

	return nil
}

// extractCustomData extracts additional fields from plain objects.
func extractCustomData(err error, data map[string]interface{}) {
	if err == nil {
		return
	}

	errVal := reflect.ValueOf(err)
	if errVal.Kind() != reflect.Ptr {
		return
	}

	elem := errVal.Elem()
	if elem.Kind() != reflect.Struct {
		return
	}

	t := elem.Type()

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if field.PkgPath != "" && !field.Anonymous {
			continue
		}

		fieldVal := elem.Field(i)
		if !fieldVal.CanInterface() {
			continue
		}

		name := field.Name
		value := fieldVal.Interface()

		switch value := value.(type) {
		case string:
			if name != "stack" && name != "Cause" {
				data[name] = value
			}
		case int, int8, int16, int32, int64:
			data[name] = value
		case uint, uint8, uint16, uint32, uint64:
			data[name] = value
		case float32, float64:
			data[name] = value
		case bool:
			data[name] = value
		case error:
			if name != "stack" && name != "Cause" && value != nil {
				data[name] = value.Error()
			}
		}
	}
}
