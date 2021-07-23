// +build !windows,!plan9

package fileperm

import (
	"os"
	"os/user"
	"strconv"
	"syscall"
)

// VERSION of go-fileperm, follows Semantic Versioning. (http://semver.org/)
const VERSION = "0.1.2"

// Bitmask to represent different access level of the file in question
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

// FilePermUser implements the main struct of the fileperm packages. All methods are based on it
type FilePermUser struct {
	Path        string
	Stat        os.FileInfo
	Uid         uint32
	Gid         uint32
	CurUserUid  int64
	CurUserGids []int64
}

// New returns a new FilePermUser struct. NewFileUserPerm expects a file path string
// as input and will return an error if the initial operations failed
func New(f string) (FilePermUser, error) {
	fpuObj := FilePermUser{Path: f}
	fs, err := os.Lstat(fpuObj.Path)
	if fs == nil {
		return fpuObj, err
	}
	fpuObj.Stat = fs

	if err := fpuObj.getUidGid(); err != nil {
		return fpuObj, err
	}

	if err := fpuObj.getCurUserIds(); err != nil {
		return fpuObj, err
	}
	return fpuObj, nil
}

// UserReadable returns true if the filepath is readable by the current user
func (f *FilePermUser) UserReadable() bool {
	if int64(f.Uid) == f.CurUserUid {
		return f.Stat.Mode().Perm()&OS_USER_R != 0
	}
	for _, gId := range f.CurUserGids {
		if int64(f.Gid) == gId {
			return f.Stat.Mode().Perm()&OS_GROUP_R != 0
		}
	}
	return f.Stat.Mode().Perm()&OS_OTH_R != 0
}

// UserWritable returns true if the filepath is writable by the current user
func (f *FilePermUser) UserWritable() bool {
	if int64(f.Uid) == f.CurUserUid {
		return f.Stat.Mode().Perm()&OS_USER_W != 0
	}
	for _, gId := range f.CurUserGids {
		if int64(f.Gid) == gId {
			return f.Stat.Mode().Perm()&OS_GROUP_W != 0
		}
	}
	return f.Stat.Mode().Perm()&OS_OTH_W != 0
}

// UserExecutable returns true if the filepath is executable by the current user
func (f *FilePermUser) UserExecutable() bool {
	if int64(f.Uid) == f.CurUserUid {
		return f.Stat.Mode().Perm()&OS_USER_X != 0
	}
	for _, gId := range f.CurUserGids {
		if int64(f.Gid) == gId {
			return f.Stat.Mode().Perm()&OS_GROUP_X != 0
		}
	}
	return f.Stat.Mode().Perm()&OS_OTH_X != 0
}

// UserWriteReadable returns true if the filepath is write- and readable by the current user
func (f *FilePermUser) UserWriteReadable() bool {
	return f.UserReadable() && f.UserWritable()
}

// UserWriteReadExecutable returns true if the filepath is write- and read- and executable by the
// current user
func (f *FilePermUser) UserWriteReadExecutable() bool {
	return f.UserReadable() && f.UserWritable() && f.UserExecutable()
}

// getUidGid retrieves the owner id and group id of the current FilePermUser structs file
func (f *FilePermUser) getUidGid() error {
	fsSys := f.Stat.Sys().(*syscall.Stat_t)
	if fsSys != nil {
		f.Uid = fsSys.Uid
		f.Gid = fsSys.Gid
	}
	return nil
}

// getCurUserIds retrieves the userid and groupids of the processes current user and stores them
// in the FilePermUser struct
func (f *FilePermUser) getCurUserIds() error {
	curUser, err := user.Current()
	if err != nil {
		return err
	}
	uidInt, err := strconv.ParseInt(curUser.Uid, 10, 32)
	if err != nil {
		return err
	}
	f.CurUserUid = uidInt

	userGroups, err := curUser.GroupIds()
	if err != nil {
		return err
	}
	for _, gid := range userGroups {
		gidInt, err := strconv.ParseInt(gid, 10, 32)
		if err != nil {
			return err
		}
		f.CurUserGids = append(f.CurUserGids, gidInt)
	}
	return nil
}
