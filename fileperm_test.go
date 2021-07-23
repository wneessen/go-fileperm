package fileperm

import (
	"fmt"
	"os"
	"testing"
)

// TestNew tests the New() function
func TestNew(t *testing.T) {
	testFile, err := os.CreateTemp("", "go-fileperm_testing")
	defer func() { _ = os.Remove(testFile.Name()) }()
	if err != nil {
		t.Errorf("Could not create temporary file: %s", err)
	}

	// Successful New() call
	_, err = New(testFile.Name())
	if err != nil {
		t.Errorf("New() on temporary file failed: %s", err)
	}

	// Unsuccessful New() call
	_, err = New(fmt.Sprintf("%s.notexistsing", testFile.Name()))
	if err == nil {
		t.Errorf("New() on temporary file failed: %s", err)
	}
}

// TestFuncs tests the different functions of the library
func TestFuncs(t *testing.T) {
	testTable := []struct {
		testName    string
		filePerms   os.FileMode
		shouldRead  bool
		shouldWrite bool
		shouldExec  bool
		shouldFail  bool
	}{
		{"Expect: 777, has: 777", 0777, true, true, true, false},
		{"Expect: 100, has: 100", 0100, false, false, true, false},
		{"Expect: 100, has: 000", 0000, false, false, true, true},
		{"Expect: 200, has: 200", 0200, false, true, false, false},
		{"Expect: 200, has: 000", 0000, false, true, false, true},
		{"Expect: 400, has: 400", 0400, true, false, false, false},
		{"Expect: 400, has: 000", 0000, true, false, false, true},
		{"Expect: 500, has: 500", 0500, true, false, true, false},
		{"Expect: 500, has: 000", 0000, true, false, true, true},
		{"Expect: 600, has: 600", 0600, true, true, false, false},
		{"Expect: 600, has: 000", 0000, true, true, false, true},
		{"Expect: 700, has: 700", 0700, true, true, true, false},
		{"Expect: 700, has: 000", 0000, true, true, true, true},
		{"Expect: 500, has: 700", 0700, true, false, true, true},
		{"Expect: 700, has: 500", 0500, true, true, true, true},
	}
	for _, testCase := range testTable {
		t.Run(testCase.testName, func(t *testing.T) {
			testFile, err := os.CreateTemp("", "go-fileperm_testing")
			defer func() { _ = os.Remove(testFile.Name()) }()
			if err != nil {
				t.Errorf("Could not create temporary file: %s", err)
			}
			if err := testFile.Chmod(testCase.filePerms); err != nil {
				t.Errorf("Failed to change filemod of temporary file: %s", err)
			}
			filePerm, err := New(testFile.Name())
			if err != nil {
				t.Errorf("New() on temporary file failed: %s", err)
			}
			if filePerm.UserReadable() != testCase.shouldRead {
				if !testCase.shouldFail {
					t.Errorf("File is supposed to be user-readable but isn't.")
				}
			}
			if filePerm.UserWritable() != testCase.shouldWrite {
				if !testCase.shouldFail {
					t.Errorf("File is supposed to be user-writable but isn't.")
				}
			}
			if filePerm.UserExecutable() != testCase.shouldExec {
				if !testCase.shouldFail {
					t.Errorf("File is supposed to be user-executable but isn't.")
				}
			}
			if filePerm.UserWriteReadable() != (testCase.shouldRead && testCase.shouldWrite) {
				if !testCase.shouldFail {
					t.Errorf("File is supposed to be user-readable and -writable but isn't.")
				}
			}
			if filePerm.UserWriteReadExecutable() != (testCase.shouldRead && testCase.shouldWrite &&
				testCase.shouldExec) {
				if !testCase.shouldFail {
					t.Errorf("File is supposed to be user-readable, -writable and -executable but isn't.")
				}
			}
		})
	}
}
