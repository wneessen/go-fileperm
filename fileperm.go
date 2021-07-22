package go_fileperm

import (
	"os"
)

const (
	OS_READ        = 04
	OS_WRITE       = 02
	OS_EX          = 01
	OS_USER_SHIFT  = 6
	OS_GROUP_SHIFT = 3
	OS_OTH_SHIFT   = 0

	OS_USER_R   = OS_READ << OS_USER_SHIFT
	OS_USER_W   = OS_WRITE << OS_USER_SHIFT
	OS_USER_X   = OS_EX << OS_USER_SHIFT
	OS_USER_RW  = OS_USER_R | OS_USER_W
	OS_USER_RWX = OS_USER_RW | OS_USER_X

	OS_GROUP_R   = OS_READ << OS_GROUP_SHIFT
	OS_GROUP_W   = OS_WRITE << OS_GROUP_SHIFT
	OS_GROUP_X   = OS_EX << OS_GROUP_SHIFT
	OS_GROUP_RW  = OS_GROUP_R | OS_GROUP_W
	OS_GROUP_RWX = OS_GROUP_RW | OS_GROUP_X

	OS_OTH_R   = OS_READ << OS_OTH_SHIFT
	OS_OTH_W   = OS_WRITE << OS_OTH_SHIFT
	OS_OTH_X   = OS_EX << OS_OTH_SHIFT
	OS_OTH_RW  = OS_OTH_R | OS_OTH_W
	OS_OTH_RWX = OS_OTH_RW | OS_OTH_X

	OS_ALL_R   = OS_USER_R | OS_GROUP_R | OS_OTH_R
	OS_ALL_W   = OS_USER_W | OS_GROUP_W | OS_OTH_W
	OS_ALL_X   = OS_USER_X | OS_GROUP_X | OS_OTH_X
	OS_ALL_RW  = OS_ALL_R | OS_ALL_W
	OS_ALL_RWX = OS_ALL_RW | OS_GROUP_X
)

// UserWritable returns true if the filepath is writable by the current user
func UserWritable(f string) bool {
	fs, err := os.Lstat(f)
	if err != nil {
		return false
	}
	return fs.Mode().Perm()&OS_USER_W != 0
}

// UserReadable returns true if the filepath is readable by the current user
func UserReadable(f string) bool {
	fs, err := os.Lstat(f)
	if err != nil {
		return false
	}
	return fs.Mode().Perm()&OS_USER_R != 0
}

// UserExecutable returns true if the filepath is executable by the current user
func UserExecutable(f string) bool {
	fs, err := os.Lstat(f)
	if err != nil {
		return false
	}
	return fs.Mode().Perm()&OS_USER_X != 0
}

// UserWriteReadable returns true if the filepath is write- and readable by the current user
func UserWriteReadable(f string) bool {
	fs, err := os.Lstat(f)
	if err != nil {
		return false
	}
	return fs.Mode().Perm()&OS_USER_RW != 0
}

// UserWriteReadExecutable returns true if the filepath is write- and read- and executable by the
// current user
func UserWriteReadExecutable(f string) bool {
	fs, err := os.Lstat(f)
	if err != nil {
		return false
	}
	return fs.Mode().Perm()&OS_USER_RWX != 0
}
