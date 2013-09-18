package config

import (
	"testing"
)

var (
	expectedList = []string{"one", "two", "three", "global_extra", "override_extra"}
)

func init() {
	globalFile = "testConfigs/global.cfg"
	configFile = "testConfigs/config.cfg"
	overrideFile = "testConfigs/override.cfg"
	LoadConfigs()
}

func TestAppCount(t *testing.T) {
	actual := List()
	// Assert whether a regular option was merged from source -> target
	if len(expectedList) != len(actual) {
		t.Errorf("Expected App count to be %d but instead it was '%d'", len(actual), len(expectedList))
	}
}

/*
func TestGlobal(t *testing.T) {
	actual := Process("one").Kill
	expected := false
	// Assert whether a regular option was merged from source -> target
	if actual != expected {
		t.Errorf("'one' Kill command expected to be %s but instead it was '%d'", actual, expected)
	}

	actual = Process("two").Kill
	expected = true
	// Assert whether a regular option was merged from source -> target
	if actual != expected {
		t.Errorf("'two' Kill command expected to be % but instead it was '%'", actual, expected)
	}

	actual = Process("three").Kill
	expected = false
	// Assert whether a regular option was merged from source -> target
	if actual != expected {
		t.Errorf("'three' Kill command expected to be % but instead it was '%'", actual, expected)
	}
}

func TestOverride(t *testing.T) {
	actual := Process("one").Port
	expected := 1111
	// Assert whether a regular option was merged from source -> target
	if actual != expected {
		t.Errorf("'one' Port expected to be %d but instead it was '%d'", actual, expected)
	}

	actual = Process("two").Port
	expected = 3000
	// Assert whether a regular option was merged from source -> target
	if actual != expected {
		t.Errorf("'two' Port expected to be %d but instead it was '%d'", actual, expected)
	}

	actual = Process("three").Port
	expected = 2
	// Assert whether a regular option was merged from source -> target
	if actual != expected {
		t.Errorf("'three' Port expected to be %d but instead it was '%d'", actual, expected)
	}
}
*/
