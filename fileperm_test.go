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
		name string
		perm os.FileMode
		r    bool
		w    bool
		x    bool
	}{
		{"Fileperm: 000", 0o000, false, false, false},
		{"Fileperm: 001", 0o001, false, false, true},
		{"Fileperm: 010", 0o010, false, false, true},
		{"Fileperm: 100", 0o100, false, false, true},
		{"Fileperm: 011", 0o011, false, false, true},
		{"Fileperm: 101", 0o101, false, false, true},
		{"Fileperm: 110", 0o110, false, false, true},
		{"Fileperm: 111", 0o111, false, false, true},
		{"Fileperm: 002", 0o002, false, true, false},
		{"Fileperm: 020", 0o020, false, true, false},
		{"Fileperm: 200", 0o200, false, true, false},
		{"Fileperm: 022", 0o022, false, true, false},
		{"Fileperm: 202", 0o202, false, true, false},
		{"Fileperm: 220", 0o220, false, true, false},
		{"Fileperm: 222", 0o222, false, true, false},
		{"Fileperm: 003", 0o003, false, true, true},
		{"Fileperm: 030", 0o030, false, true, true},
		{"Fileperm: 300", 0o300, false, true, true},
		{"Fileperm: 033", 0o033, false, true, true},
		{"Fileperm: 303", 0o303, false, true, true},
		{"Fileperm: 330", 0o330, false, true, true},
		{"Fileperm: 333", 0o333, false, true, true},
		{"Fileperm: 004", 0o004, true, false, false},
		{"Fileperm: 040", 0o040, true, false, false},
		{"Fileperm: 400", 0o400, true, false, false},
		{"Fileperm: 044", 0o044, true, false, false},
		{"Fileperm: 404", 0o404, true, false, false},
		{"Fileperm: 440", 0o440, true, false, false},
		{"Fileperm: 444", 0o444, true, false, false},
		{"Fileperm: 005", 0o005, true, false, true},
		{"Fileperm: 050", 0o050, true, false, true},
		{"Fileperm: 500", 0o500, true, false, true},
		{"Fileperm: 055", 0o055, true, false, true},
		{"Fileperm: 505", 0o505, true, false, true},
		{"Fileperm: 550", 0o550, true, false, true},
		{"Fileperm: 555", 0o555, true, false, true},
		{"Fileperm: 006", 0o006, true, true, false},
		{"Fileperm: 060", 0o060, true, true, false},
		{"Fileperm: 600", 0o600, true, true, false},
		{"Fileperm: 066", 0o066, true, true, false},
		{"Fileperm: 606", 0o606, true, true, false},
		{"Fileperm: 660", 0o660, true, true, false},
		{"Fileperm: 666", 0o666, true, true, false},
		{"Fileperm: 007", 0o007, true, true, true},
		{"Fileperm: 070", 0o070, true, true, true},
		{"Fileperm: 700", 0o700, true, true, true},
		{"Fileperm: 077", 0o077, true, true, true},
		{"Fileperm: 707", 0o707, true, true, true},
		{"Fileperm: 770", 0o770, true, true, true},
		{"Fileperm: 777", 0o777, true, true, true},
	}
	for _, tc := range testTable {
		t.Run(tc.name, func(t *testing.T) {
			testFile, err := os.CreateTemp("", "go-fileperm_testing")
			defer func() { _ = os.Remove(testFile.Name()) }()
			if err != nil {
				t.Errorf("Could not create temporary file: %s", err)
			}
			if err := testFile.Chmod(tc.perm); err != nil {
				t.Errorf("Failed to change filemod of temporary file: %s", err)
			}
			fp, err := New(testFile.Name())
			if err != nil {
				t.Errorf("New() on temporary file failed: %s", err)
			}
			if !tc.r && fp.UserReadable() {
				t.Errorf("file with perms %o was supposed to be not readable but UserReadable returned true",
					tc.perm)
			}
			if tc.r && !fp.UserReadable() {
				t.Errorf("file with perms %o was supposed to be readable but UserReadable returned false",
					tc.perm)
			}
			if tc.w && !fp.UserWritable() {
				t.Errorf("file with perms %o was supposed to be writable but UserWritable returned false",
					tc.perm)
			}
			if !tc.w && fp.UserWritable() {
				t.Errorf("file with perms %o was supposed to be not writable but UserWritable returned true",
					tc.perm)
			}
			if tc.x && !fp.UserExecutable() {
				t.Errorf("file with perms %o was supposed to be executable but UserExecutable returned false",
					tc.perm)
			}
			if !tc.x && fp.UserExecutable() {
				t.Errorf("file with perms %o was supposed to be not executable but UserExecutable returned true",
					tc.perm)
			}
			if (tc.r && tc.w) && !fp.UserWriteReadable() {
				t.Errorf("file with perms %o was supposed to be write-/readable but UserWriteReadable "+
					"returned false", tc.perm)
			}
			if (!tc.r && !tc.w) && fp.UserWriteReadable() {
				t.Errorf("file with perms %o was supposed to be not write-/readable but UserWriteReadable "+
					"returned true", tc.perm)
			}
			if (tc.r && tc.x) && !fp.UserReadExecutable() {
				t.Errorf("file with perms %o was supposed to be read-/executable but UserReadExecutable "+
					"returned false", tc.perm)
			}
			if (!tc.r && !tc.x) && fp.UserReadExecutable() {
				t.Errorf("file with perms %o was supposed to be not read-/executable but UserReadExecutable "+
					"returned true", tc.perm)
			}
			if (tc.w && tc.x) && !fp.UserWriteExecutable() {
				t.Errorf("file with perms %o was supposed to be write-/executable but UserWriteExecutable "+
					"returned false", tc.perm)
			}
			if (!tc.w && !tc.x) && fp.UserWriteExecutable() {
				t.Errorf("file with perms %o was supposed to be not write-/executable but UserWriteExecutable "+
					"returned true", tc.perm)
			}
			if (tc.r && tc.w && tc.x) && !fp.UserWriteReadExecutable() {
				t.Errorf("file with perms %o was supposed to be write-/read-/executable but "+
					"UserWriteReadExecutable returned false", tc.perm)
			}
			if (!tc.r && !tc.w && !tc.x) && fp.UserWriteReadExecutable() {
				t.Errorf("file with perms %o was supposed to be not write-/read-/executable but "+
					"UserWriteReadExecutable returned true", tc.perm)
			}
		})
	}
}

func BenchmarkPermUser_UserReadable(b *testing.B) {
	testFile, err := os.CreateTemp("", "go-fileperm_testing")
	defer func() { _ = os.Remove(testFile.Name()) }()
	if err != nil {
		b.Errorf("Could not create temporary file: %s", err)
	}
	if err := testFile.Chmod(0o777); err != nil {
		b.Errorf("Failed to change filemod of temporary file: %s", err)
	}
	fp, err := New(testFile.Name())
	if err != nil {
		b.Errorf("failed to create permuser instance: %s", err)
		return
	}
	for i := 0; i < b.N; i++ {
		_ = fp.UserReadable()
	}
}

func BenchmarkPermUser_UserWritable(b *testing.B) {
	testFile, err := os.CreateTemp("", "go-fileperm_testing")
	defer func() { _ = os.Remove(testFile.Name()) }()
	if err != nil {
		b.Errorf("Could not create temporary file: %s", err)
	}
	if err := testFile.Chmod(0o777); err != nil {
		b.Errorf("Failed to change filemod of temporary file: %s", err)
	}
	fp, err := New(testFile.Name())
	if err != nil {
		b.Errorf("failed to create permuser instance: %s", err)
		return
	}
	for i := 0; i < b.N; i++ {
		_ = fp.UserWritable()
	}
}

func BenchmarkPermUser_UserExecutable(b *testing.B) {
	testFile, err := os.CreateTemp("", "go-fileperm_testing")
	defer func() { _ = os.Remove(testFile.Name()) }()
	if err != nil {
		b.Errorf("Could not create temporary file: %s", err)
	}
	if err := testFile.Chmod(0o777); err != nil {
		b.Errorf("Failed to change filemod of temporary file: %s", err)
	}
	fp, err := New(testFile.Name())
	if err != nil {
		b.Errorf("failed to create permuser instance: %s", err)
		return
	}
	for i := 0; i < b.N; i++ {
		_ = fp.UserExecutable()
	}
}

func BenchmarkPermUser_UserWriteReadable(b *testing.B) {
	testFile, err := os.CreateTemp("", "go-fileperm_testing")
	defer func() { _ = os.Remove(testFile.Name()) }()
	if err != nil {
		b.Errorf("Could not create temporary file: %s", err)
	}
	if err := testFile.Chmod(0o777); err != nil {
		b.Errorf("Failed to change filemod of temporary file: %s", err)
	}
	fp, err := New(testFile.Name())
	if err != nil {
		b.Errorf("failed to create permuser instance: %s", err)
		return
	}
	for i := 0; i < b.N; i++ {
		_ = fp.UserWriteReadable()
	}
}

func BenchmarkPermUser_UserWriteExecutable(b *testing.B) {
	testFile, err := os.CreateTemp("", "go-fileperm_testing")
	defer func() { _ = os.Remove(testFile.Name()) }()
	if err != nil {
		b.Errorf("Could not create temporary file: %s", err)
	}
	if err := testFile.Chmod(0o777); err != nil {
		b.Errorf("Failed to change filemod of temporary file: %s", err)
	}
	fp, err := New(testFile.Name())
	if err != nil {
		b.Errorf("failed to create permuser instance: %s", err)
		return
	}
	for i := 0; i < b.N; i++ {
		_ = fp.UserWriteExecutable()
	}
}

func BenchmarkPermUser_UserReadExecutable(b *testing.B) {
	testFile, err := os.CreateTemp("", "go-fileperm_testing")
	defer func() { _ = os.Remove(testFile.Name()) }()
	if err != nil {
		b.Errorf("Could not create temporary file: %s", err)
	}
	if err := testFile.Chmod(0o777); err != nil {
		b.Errorf("Failed to change filemod of temporary file: %s", err)
	}
	fp, err := New(testFile.Name())
	if err != nil {
		b.Errorf("failed to create permuser instance: %s", err)
		return
	}
	for i := 0; i < b.N; i++ {
		_ = fp.UserReadExecutable()
	}
}

func BenchmarkPermUser_UserWriteReadExecutable(b *testing.B) {
	testFile, err := os.CreateTemp("", "go-fileperm_testing")
	defer func() { _ = os.Remove(testFile.Name()) }()
	if err != nil {
		b.Errorf("Could not create temporary file: %s", err)
	}
	if err := testFile.Chmod(0o777); err != nil {
		b.Errorf("Failed to change filemod of temporary file: %s", err)
	}
	fp, err := New(testFile.Name())
	if err != nil {
		b.Errorf("failed to create permuser instance: %s", err)
		return
	}
	for i := 0; i < b.N; i++ {
		_ = fp.UserWriteReadExecutable()
	}
}
