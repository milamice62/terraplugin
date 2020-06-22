package provider

import (
	"fmt"
	"regexp"
)

type Model struct {
	ID string `json:"_id"`
}

func validateName(v interface{}, k string) (ws []string, es []error) {
	var errs []error
	var warns []string
	value, ok := v.(string)
	if !ok {
		errs = append(errs, fmt.Errorf("Expected value to be string"))
		return warns, errs
	}
	whiteSpace := regexp.MustCompile(`\s+`)
	if whiteSpace.Match([]byte(value)) {
		errs = append(errs, fmt.Errorf("String value cannot contain whitespace. Got %s", value))
		return warns, errs
	}
	return warns, errs
}

func validateInt(v interface{}, k string) (ws []string, es []error) {
	var errs []error
	var warns []string
	value, ok := v.(int)
	if !ok {
		errs = append(errs, fmt.Errorf("Expected positive integer number. Got %v", value))
		return warns, errs
	}
	return warns, errs
}

func validateFloat(v interface{}, k string) (ws []string, es []error) {
	var errs []error
	var warns []string
	value, ok := v.(float64)
	if !ok {
		errs = append(errs, fmt.Errorf("Expected float number"))
		return warns, errs
	}
	if value < 0 {
		errs = append(errs, fmt.Errorf("Should be positive value"))
		return warns, errs
	}
	return warns, errs
}
