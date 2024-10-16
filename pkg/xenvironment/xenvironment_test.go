//nolint:paralleltest // disabled paralleltest due to the nature of os environment handling.
package xenvironment_test

import (
	"os"
	"reflect"
	"testing"

	"github.com/nicce/go-grpc-lab/pkg/xenvironment"

	"github.com/google/go-cmp/cmp"
)

type Environment struct {
	Foo string `env:"FOO"`
	Bar string `env:"BAR"`
	Baz string
}

func TestShouldRunWithNoEnvFile(t *testing.T) {
	e := &Environment{}
	err := xenvironment.GetEnvironment(e)

	if err != nil {
		t.Errorf("Failed running GetEnvironment without .env file: %v", err)
	}
}

func TestShouldKeepNonMatchingValues(t *testing.T) {
	e := &Environment{
		Foo: "foo",
		Bar: "bar",
		Baz: "",
	}

	_ = os.Setenv("FOO2", "foo2")
	_ = os.Setenv("BAR2", "bar2")

	_ = xenvironment.GetEnvironment(e)

	if e.Bar != "bar" || e.Foo != "foo" {
		t.Errorf("Expected foo and bar in environment but got %v", e)
	}
}

func TestShouldUseEnvironmentVariables(t *testing.T) {
	e := &Environment{}

	expected := Environment{
		Foo: "123",
		Bar: "456",
		Baz: "",
	}

	_ = os.Setenv("FOO", "123")
	_ = os.Setenv("BAR", "456")

	_ = xenvironment.GetEnvironment(e)

	if !cmp.Equal(*e, expected) {
		t.Errorf("Expected %v to match %v", *e, expected)
	}
}

func TestShouldNotModifyUntaggedField(t *testing.T) {
	expected := "Ignored"

	e := &Environment{
		Foo: "",
		Bar: "",
		Baz: expected,
	}

	_ = xenvironment.GetEnvironment(e)

	if e.Baz != expected {
		t.Errorf("Expected field Baz to be '%s' but got '%s'", expected, e.Baz)
	}
}

type boolValuesTestcase struct {
	name     string
	value    string
	expected bool
}

func TestShouldParseBoolValues(t *testing.T) {
	testcases := []boolValuesTestcase{
		{
			name:     "testcase false",
			value:    "false",
			expected: false,
		},
		{
			name:     "testcase true",
			value:    "true",
			expected: true,
		},
	}

	for _, testcase := range testcases {
		testcase := testcase
		t.Run(testcase.name, func(t *testing.T) {
			type env struct {
				TestBool bool `env:"TESTBOOL"`
			}

			e := &env{}

			_ = os.Setenv("TESTBOOL", testcase.value)

			_ = xenvironment.GetEnvironment(e)

			if e.TestBool != testcase.expected {
				t.Errorf("Expected field TestBool to be '%v' but got '%v'", testcase.expected, e.TestBool)
			}
		})
	}
}

type floatValuesTestcase struct {
	name       string
	value      string
	expected32 float32
	expected64 float64
}

func TestShouldParseFloatValues(t *testing.T) {
	testcases := []floatValuesTestcase{
		{
			name:       "testcase float parse 10.65",
			value:      "10.65",
			expected32: 10.65,
			expected64: 10.65,
		},
		{
			name:       "testcase float parse -0.2234234",
			value:      "-0.2234234",
			expected32: -0.2234234,
			expected64: -0.2234234,
		},
	}

	for _, testcase := range testcases {
		testcase := testcase

		t.Run(testcase.name, func(t *testing.T) {
			type env struct {
				Test32 float32 `env:"TESTVAL"`
				Test64 float64 `env:"TESTVAL"`
			}

			e := &env{}

			_ = os.Setenv("TESTVAL", testcase.value)

			_ = xenvironment.GetEnvironment(e)

			if e.Test32 != testcase.expected32 {
				t.Errorf("Expected field Test32 to be '%v' but got '%v'", testcase.expected32, e.Test32)
			}

			if e.Test64 != testcase.expected64 {
				t.Errorf("Expected field Test64 to be '%v' but got '%v'", testcase.expected64, e.Test64)
			}
		})
	}
}

type intValuesTestcase struct {
	name             string
	envVal           string
	expectedIntVal   int
	expectedInt8Val  int8
	expectedInt16Val int16
	expectedInt32Val int32
	expectedInt64Val int64
}

func TestShouldParseIntValues(t *testing.T) {
	type env struct {
		IntVal   int   `env:"TESTINT"`
		Int8Val  int8  `env:"TESTINT"`
		Int16Val int16 `env:"TESTINT"`
		Int32Val int32 `env:"TESTINT"`
		Int64Val int64 `env:"TESTINT"`
	}

	testcases := []intValuesTestcase{
		{
			name:             "testcase int parse 123",
			envVal:           "123",
			expectedIntVal:   123,
			expectedInt8Val:  123,
			expectedInt16Val: 123,
			expectedInt32Val: 123,
			expectedInt64Val: 123,
		},
		{
			name:             "testcase int parse -12",
			envVal:           "-12",
			expectedIntVal:   -12,
			expectedInt8Val:  -12,
			expectedInt16Val: -12,
			expectedInt32Val: -12,
			expectedInt64Val: -12,
		},
	}

	for _, testcase := range testcases {
		testcase := testcase

		t.Run(testcase.name, func(t *testing.T) {
			_ = os.Setenv("TESTINT", testcase.envVal)

			testEnv := &env{}
			_ = xenvironment.GetEnvironment(testEnv)

			if testEnv.IntVal != testcase.expectedIntVal {
				t.Errorf("Expected field IntVal to be '%v' but got '%v'", testcase.expectedIntVal, testEnv.IntVal)
			}

			if testEnv.Int8Val != testcase.expectedInt8Val {
				t.Errorf("Expected field Int8Val to be '%v' but got '%v'", testcase.expectedInt8Val, testEnv.Int8Val)
			}

			if testEnv.Int16Val != testcase.expectedInt16Val {
				t.Errorf("Expected field Int16Val to be '%v' but got '%v'", testcase.expectedInt16Val, testEnv.Int16Val)
			}

			if testEnv.Int32Val != testcase.expectedInt32Val {
				t.Errorf("Expected field Int32Val to be '%v' but got '%v'", testcase.expectedInt32Val, testEnv.Int32Val)
			}

			if testEnv.Int64Val != testcase.expectedInt64Val {
				t.Errorf("Expected field Int64Val to be '%v' but got '%v'", testcase.expectedInt64Val, testEnv.Int64Val)
			}
		})
	}
}

type uintValuesTestcase struct {
	name              string
	envVal            string
	expectedUintVal   uint
	expectedUint8Val  uint8
	expectedUint16Val uint16
	expectedUint32Val uint32
	expectedUint64Val uint64
}

func TestShouldParseUintValues(t *testing.T) {
	type env struct {
		UintVal   uint   `env:"TESTUINT"`
		Uint8Val  uint8  `env:"TESTUINT"`
		Uint16Val uint16 `env:"TESTUINT"`
		Uint32Val uint32 `env:"TESTUINT"`
		Uint64Val uint64 `env:"TESTUINT"`
	}

	testcases := []uintValuesTestcase{
		{
			name:              "testcase uint parse 123",
			envVal:            "123",
			expectedUintVal:   123,
			expectedUint8Val:  123,
			expectedUint16Val: 123,
			expectedUint32Val: 123,
			expectedUint64Val: 123,
		},
		{
			name:              "testcase uint parse 1",
			envVal:            "1",
			expectedUintVal:   1,
			expectedUint8Val:  1,
			expectedUint16Val: 1,
			expectedUint32Val: 1,
			expectedUint64Val: 1,
		},
	}

	for _, testcase := range testcases {
		testcase := testcase

		t.Run(testcase.name, func(t *testing.T) {
			_ = os.Setenv("TESTUINT", testcase.envVal)

			testEnv := &env{}
			_ = xenvironment.GetEnvironment(testEnv)

			if testEnv.UintVal != testcase.expectedUintVal {
				t.Errorf("Expected field UintVal to be '%v' but got '%v'", testcase.expectedUintVal, testEnv.UintVal)
			}

			if testEnv.Uint8Val != testcase.expectedUint8Val {
				t.Errorf("Expected field Uint8Val to be '%v' but got '%v'", testcase.expectedUint8Val, testEnv.Uint8Val)
			}

			if testEnv.Uint16Val != testcase.expectedUint16Val {
				t.Errorf("Expected field Uint16Val to be '%v' but got '%v'", testcase.expectedUint16Val, testEnv.Uint16Val)
			}

			if testEnv.Uint32Val != testcase.expectedUint32Val {
				t.Errorf("Expected field Uint32Val to be '%v' but got '%v'", testcase.expectedUint32Val, testEnv.Uint32Val)
			}

			if testEnv.Uint64Val != testcase.expectedUint64Val {
				t.Errorf("Expected field Uint64Val to be '%v' but got '%v'", testcase.expectedUint64Val, testEnv.Uint64Val)
			}
		})
	}
}

type unsupportedValuesTestcase struct {
	name  string
	value interface{}
}

func TestShouldReturnErrorOnUnsupportedValueType(t *testing.T) {
	testcases := []unsupportedValuesTestcase{
		{
			name: "testcase unsupported struct",
			value: &struct {
				Stuffs struct{} `env:"STUFFS"`
			}{
				struct{}{},
			},
		},
		{
			name: "testcase unsupported interface",
			value: &struct {
				Stuffs interface{} `env:"STUFFS"`
			}{
				"this is a triumph",
			},
		},
	}

	for _, testcase := range testcases {
		testcase := testcase

		t.Run(testcase.name, func(t *testing.T) {
			_ = os.Setenv("STUFFS", testcase.name)

			err := xenvironment.GetEnvironment(testcase.value)
			if err == nil {
				t.Errorf("Expected the returned error to be ErrUnsupportedEnvironmentType but got nil")
			}
		})
	}
}

type stringSlice struct {
	Values []string `env:"VALUES"`
}

func TestShouldParseStringSlice(t *testing.T) {
	// Arrange
	expected := []string{"a", "B", "c", "D"}
	_ = os.Setenv("VALUES", "a|B|c|D")

	// Act
	env := &stringSlice{}
	err := xenvironment.GetEnvironment(env)

	// Assert
	if err != nil {
		t.Errorf("unable to parse string slice: %v", err)
	}

	equal := reflect.DeepEqual(expected, env.Values)
	if !equal {
		t.Errorf("the parsed string slice did not match the expected values")
	}
}
