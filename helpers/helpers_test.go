package helpers

import (
	"os"
	"testing"
)

func TestSetValueFromEnv(t *testing.T) {

	key := "_TATANKA_TEST_HELPER"
	value := "_test_value"

	os.Setenv(key, value)

	var test string
	SetValueFromEnv(&test, key)

	if test != value {
		t.Error("Failed to set variable to environment setting, value is: ", test)
	}

}
