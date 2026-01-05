package envutil

import (
	"os"
	"testing"
)

func TestGetEnv(t *testing.T) {
	// Test when environment variable exists
	key := "TEST_GET_ENV"
	expected := "test_value"
	os.Setenv(key, expected)
	defer os.Unsetenv(key)

	actual := GetEnv(key, "default_value")
	if actual != expected {
		t.Errorf("GetEnv(%s) = %s; want %s", key, actual, expected)
	}

	// Test when environment variable doesn't exist
	nonExistentKey := "NON_EXISTENT_KEY"
	defaultValue := "default_value"
	actual = GetEnv(nonExistentKey, defaultValue)
	if actual != defaultValue {
		t.Errorf("GetEnv(%s) = %s; want %s", nonExistentKey, actual, defaultValue)
	}

	// Test case conversion
	os.Setenv(key, "MiXeD_CaSe")
	actual = GetEnv(key, defaultValue)
	if actual != "mixed_case" {
		t.Errorf("GetEnv(%s) = %s; want %s", key, actual, "mixed_case")
	}
}

func TestGetEnvInt(t *testing.T) {
	// Test when environment variable exists with valid int
	key := "TEST_GET_ENV_INT"
	expected := 42
	os.Setenv(key, "42")
	defer os.Unsetenv(key)

	actual := GetEnvInt(key, 0)
	if actual != expected {
		t.Errorf("GetEnvInt(%s) = %d; want %d", key, actual, expected)
	}

	// Test when environment variable doesn't exist
	nonExistentKey := "NON_EXISTENT_INT_KEY"
	defaultValue := 99
	actual = GetEnvInt(nonExistentKey, defaultValue)
	if actual != defaultValue {
		t.Errorf("GetEnvInt(%s) = %d; want %d", nonExistentKey, actual, defaultValue)
	}

	// Test when environment variable exists with invalid int
	os.Setenv(key, "not_an_int")
	actual = GetEnvInt(key, defaultValue)
	if actual != defaultValue {
		t.Errorf("GetEnvInt(%s) with invalid value = %d; want %d", key, actual, defaultValue)
	}
}

func TestGetEnvBool(t *testing.T) {
	// Test when environment variable exists with valid bool (true)
	key := "TEST_GET_ENV_BOOL"

	// Test true values
	for _, trueVal := range []string{"true", "TRUE", "True", "1", "t", "T"} {
		os.Setenv(key, trueVal)
		actual := GetEnvBool(key, false)
		if actual != true {
			t.Errorf("GetEnvBool(%s) with value %s = %t; want true", key, trueVal, actual)
		}
	}

	// Test false values
	for _, falseVal := range []string{"false", "FALSE", "False", "0", "f", "F"} {
		os.Setenv(key, falseVal)
		actual := GetEnvBool(key, true)
		if actual != false {
			t.Errorf("GetEnvBool(%s) with value %s = %t; want false", key, falseVal, actual)
		}
	}

	defer os.Unsetenv(key)

	// Test when environment variable doesn't exist
	nonExistentKey := "NON_EXISTENT_BOOL_KEY"
	defaultValue := true
	actual := GetEnvBool(nonExistentKey, defaultValue)
	if actual != defaultValue {
		t.Errorf("GetEnvBool(%s) = %t; want %t", nonExistentKey, actual, defaultValue)
	}

	// Test when environment variable exists with invalid bool
	os.Setenv(key, "not_a_bool")
	actual = GetEnvBool(key, defaultValue)
	if actual != defaultValue {
		t.Errorf("GetEnvBool(%s) with invalid value = %t; want %t", key, actual, defaultValue)
	}
}
