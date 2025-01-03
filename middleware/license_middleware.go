package middleware

import (
	"errors"
	"fmt"
	"reflect"
)

type LicenseMiddleware struct {
	licenseService LicenseValidator
	licenseKey     string
}

type LicenseValidator interface {
	ValidateLicense(key string) (bool, error)
}

func NewLicenseMiddleware(validator LicenseValidator, key string) *LicenseMiddleware {
	return &LicenseMiddleware{
		licenseService: validator,
		licenseKey:     key,
	}
}

func (m *LicenseMiddleware) WrapFunction(fn interface{}) interface{} {
	return func(args ...interface{}) (interface{}, error) {
		// Check license first
		valid, err := m.licenseService.ValidateLicense(m.licenseKey)
		if err != nil {
			return nil, fmt.Errorf("license validation error: %w", err)
		}
		if !valid {
			return nil, errors.New("invalid license")
		}

		// Call original function if license valid
		result := reflect.ValueOf(fn).Call(makeValueArgs(args))
		if len(result) == 0 {
			return nil, nil
		}
		return result[0].Interface(), nil
	}
}

func makeValueArgs(args []interface{}) []reflect.Value {
	vals := make([]reflect.Value, len(args))
	for i, arg := range args {
		vals[i] = reflect.ValueOf(arg)
	}
	return vals
}
