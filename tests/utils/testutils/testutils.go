package testutils

import (
	"os"
	"strings"
	"testing"
)

// SkipUnlessEnvFlag lets us enable normally-skipped test cases by specifying an environment flag.
// If running locally just prepend `TEST_INCLUDE="Test reading & writing of module block without tags"`.
func SkipUnlessEnvFlag(t *testing.T) {
	testName := t.Name()

	_, ok := os.LookupEnv("TEST_ALL")
	if ok {
		return
	}

	noSkip, ok := os.LookupEnv("TEST_INCLUDE")
	if !ok {
		t.Skipf("Skipping %s - TEST_INCLUDE doesn't exist", testName)
	}

	noSkipSplit := strings.Split(noSkip, ",")
	for _, s := range noSkipSplit {
		if s == testName {
			return
		}
	}

	// We didn't find this test name in TEST_INCLUDE - so just skip
	t.Skipf("Skipping %s - not specified in TEST_INCLUDE", testName)
}
