/*
Copyright 2021 The Vitess Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package flagutil

import (
	"errors"
	"fmt"
	"strconv"
)

// OptionalFlag augements the pflag.Value interface with a method to determine
// if a flag was set explicitly on the comand-line.
type OptionalFlag[T any] struct {
	val    T
	parser func(string) (T, error)
	set    bool
}

// NewOptionalFlag returns an OptionalFlag with the specified value as its
// starting value. The value is parsed using the parser provided
func NewOptionalFlag[T any](val T, parser func(string) (T, error)) *OptionalFlag[T] {
	return &OptionalFlag[T]{
		val:    val,
		parser: parser,
		set:    false,
	}
}

// Set is part of the pflag.Value interface.
func (f *OptionalFlag[T]) Set(arg string) error {
	v, err := f.parser(arg)
	if err != nil {
		return err
	}

	f.val = v
	f.set = true

	return nil
}

// String is part of the pflag.Value interface.
func (f *OptionalFlag[T]) String() string {
	return fmt.Sprintf("%v", f.val)
}

// Type is part of the pflag.Value interface.
func (f *OptionalFlag[T]) Type() string {
	return fmt.Sprintf("%T", f.val)
}

// Get returns the underlying value of this flag. If the flag was not
// explicitly set, this will be the initial value passed to the constructor.
func (f *OptionalFlag[T]) Get() T {
	return f.val
}

// IsSet is part of the OptionalFlag interface.
func (f *OptionalFlag[T]) IsSet() bool {
	return f.set
}

// OptionalFloat64 implements OptionalFlag for float64 values.
type OptionalFloat64 struct {
	val float64
	set bool
}

// NewOptionalFloat64 returns an *OptionalFlag[float64] with the specified value as its
// starting value.
func NewOptionalFloat64(val float64) *OptionalFlag[float64] {
	return &OptionalFlag[float64]{
		val:    val,
		parser: float64Parser,
		set:    false,
	}
}

// parses a float from a string as float64
func float64Parser(arg string) (v float64, err error) {
	v, err = strconv.ParseFloat(arg, 64)
	if err != nil {
		return v, numError(err)
	}

	return v, nil
}

// NewOptionalString returns an OptionalFlag[string] with the specified value as its
// starting value.
func NewOptionalString(val string) *OptionalFlag[string] {
	return &OptionalFlag[string]{
		val:    val,
		parser: func(s string) (string, error) { return s, nil },
		set:    false,
	}
}

// lifted directly from package flag to make the behavior of numeric parsing
// consistent with the standard library for our custom optional types.
var (
	errParse = errors.New("parse error")
	errRange = errors.New("value out of range")
)

// lifted directly from package flag to make the behavior of numeric parsing
// consistent with the standard library for our custom optional types.
func numError(err error) error {
	ne, ok := err.(*strconv.NumError)
	if !ok {
		return err
	}

	switch ne.Err {
	case strconv.ErrSyntax:
		return errParse
	case strconv.ErrRange:
		return errRange
	default:
		return err
	}
}
