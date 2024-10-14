package xenvironment

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

const (
	base10                    = 10
	defaultEnvironment string = ".env"
	tagName            string = "env"
)

var (
	// ErrMissingEnv - missing environment variable.
	ErrMissingEnv = errors.New("missing environment variable")
	// ErrParsingEnv - error parsing environment variable.
	ErrParsingEnv = errors.New("error parsing environment variable")
	// ErrUnsupportedEnvironmentType - unsupported environment type.
	ErrUnsupportedEnvironmentType = errors.New("unsupported environment type")
	// ErrEnvironmentMustBePointer - environment has to be a pointer and cannot be nil.
	ErrEnvironmentMustBePointer = errors.New("environment has to be a pointer and cannot be nil")
)

const (
	intSize   int = 0
	int8Size  int = 8
	int16Size int = 16
	int32Size int = 32
	int64Size int = 64
)

func getIntBitSize(kind reflect.Kind) (int, error) {
	switch kind {
	case reflect.Int, reflect.Uint:
		return intSize, nil
	case reflect.Int8, reflect.Uint8:
		return int8Size, nil
	case reflect.Int16, reflect.Uint16:
		return int16Size, nil
	case reflect.Int32, reflect.Uint32:
		return int32Size, nil
	case reflect.Int64, reflect.Uint64:
		return int64Size, nil
	default:
		return intSize, ErrUnsupportedEnvironmentType
	}
}

func loadEnvironment() {
	env := os.Getenv("ENV")
	if env == "" {
		env = defaultEnvironment
	} else {
		env = ".env." + env
	}

	_ = godotenv.Load(env)
}

func getEnvironmentVariables(environment interface{}) error {
	loadEnvironment()

	value := reflect.ValueOf(environment)

	if value.Kind() != reflect.Ptr || value.IsNil() {
		return ErrEnvironmentMustBePointer
	}

	element := value.Elem()
	typeOfT := element.Type()

	for i := 0; i < typeOfT.NumField(); i++ {
		structField := typeOfT.Field(i)
		tag := structField.Tag.Get(tagName)

		if tag == "" {
			continue
		}

		envValue := os.Getenv(tag)
		if envValue == "" {
			continue
		}

		kind := structField.Type.Kind()

		valueField := element.Field(i)

		err := setStructFieldValue(&valueField, kind, envValue)
		if err != nil {
			return err
		}
	}

	return nil
}

func setStructFieldValue(valueField *reflect.Value, kind reflect.Kind, envValue string) error {
	switch kind {
	case reflect.String:
		valueField.SetString(envValue)
	case reflect.Bool:
		boolValue, err := strconv.ParseBool(envValue)
		if err != nil {
			return fmt.Errorf("%s, %w", err.Error(), ErrParsingEnv)
		}

		valueField.SetBool(boolValue)
	case reflect.Float32, reflect.Float64:
		bitSize := 32
		if kind == reflect.Float64 {
			bitSize = 64
		}

		floatValue, err := strconv.ParseFloat(envValue, bitSize)
		if err != nil {
			return fmt.Errorf("%s, %w", err.Error(), ErrParsingEnv)
		}

		valueField.SetFloat(floatValue)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		bitSize, err := getIntBitSize(kind)
		if err != nil {
			return fmt.Errorf("%w, %s", ErrUnsupportedEnvironmentType, kind)
		}

		intValue, err := strconv.ParseInt(envValue, base10, bitSize)
		if err != nil {
			return fmt.Errorf("%s, %w", err.Error(), ErrParsingEnv)
		}

		valueField.SetInt(intValue)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		bitSize, err := getIntBitSize(kind)
		if err != nil {
			return fmt.Errorf("%w, %s", ErrUnsupportedEnvironmentType, kind)
		}

		uintValue, err := strconv.ParseUint(envValue, base10, bitSize)
		if err != nil {
			return fmt.Errorf("%s, %w", err.Error(), ErrParsingEnv)
		}

		valueField.SetUint(uintValue)
	case reflect.Slice:
		sliceType := valueField.Type().Elem()
		if sliceType.Kind() == reflect.String {
			values := strings.Split(envValue, "|")
			valueField.Set(reflect.ValueOf(values))
		} else {
			return fmt.Errorf("%w, %s of type %s", ErrUnsupportedEnvironmentType, kind, sliceType.Kind())
		}

	default:
		return fmt.Errorf("%w, %s", ErrUnsupportedEnvironmentType, kind)
	}

	return nil
}

// GetEnvironment - returns the application environment.
func GetEnvironment(environment interface{}) error {
	err := getEnvironmentVariables(environment)
	if err != nil {
		return err
	}

	return nil
}
