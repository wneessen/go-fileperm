// SPDX-FileCopyrightText: 2022 Winni Neessen <winni@neessen.dev>
//
// SPDX-License-Identifier: MIT
//go:build !windows && !plan9
// +build !windows,!plan9

package fileperm

import (
	"os"
	"os/user"
	"strconv"
	"syscall"
)

// List of different OS permission bits
const (
	OsRead       = 0o4
	OsWrite      = 0o2
	OsEx         = 0o1
	OsUserShift  = 6
	OsGroupShift = 3
	OsOthShift   = 0

	OsUserR = OsRead << OsUserShift
	OsUserW = OsWrite << OsUserShift
	OsUserX = OsEx << OsUserShift

	OsGroupR = OsRead << OsGroupShift
	OsGroupW = OsWrite << OsGroupShift
	OsGroupX = OsEx << OsGroupShift

	OsOthR = OsRead << OsOthShift
	OsOthW = OsWrite << OsOthShift
	OsOthX = OsEx << OsOthShift
)

// PermUser implements the main struct of the fileperm packages. All methods are based on it
type PermUser struct {
	Path        string
	Stat        os.FileInfo
	UID         uint32
	GID         uint32
	CurUserUID  int64
	CurUserGIDs []int64
}

// New returns a new PermUser struct. NewFileUserPerm expects a file path string
// as input and will return an error if the initial operations failed
func New(f string) (PermUser, error) {
	fpuObj := PermUser{Path: f}
	fs, err := os.Lstat(fpuObj.Path)
	if fs == nil {
		return fpuObj, err
	}
	fpuObj.Stat = fs
	uid, gid := getUIDGID(fpuObj.Stat)
	fpuObj.UID = uid
	fpuObj.GID = gid

	if err := fpuObj.getCurUserIds(); err != nil {
		return fpuObj, err
	}
	return fpuObj, nil
}

// UserReadable returns true if the filepath is readable by the current user
func (p *PermUser) UserReadable() bool {
	if p.Stat.Mode().Perm()&OsOthR != 0 {
		return true
	}
	if p.isInGroup() && p.Stat.Mode().Perm()&OsGroupR != 0 {
		return true
	}
	if p.isOwner() && p.Stat.Mode().Perm()&OsUserR != 0 {
		return true
	}
	return false
}

// UserWritable returns true if the filepath is writable by the current user
func (p *PermUser) UserWritable() bool {
	if p.Stat.Mode().Perm()&OsOthW != 0 {
		return true
	}
	if p.isInGroup() && p.Stat.Mode().Perm()&OsGroupW != 0 {
		return true
	}
	if p.isOwner() && p.Stat.Mode().Perm()&OsUserW != 0 {
		return true
	}
	return false
}

// UserExecutable returns true if the filepath is executable by the current user
func (p *PermUser) UserExecutable() bool {
	if p.Stat.Mode().Perm()&OsOthX != 0 {
		return true
	}
	if p.isInGroup() && p.Stat.Mode().Perm()&OsGroupX != 0 {
		return true
	}
	if p.isOwner() && p.Stat.Mode().Perm()&OsUserX != 0 {
		return true
	}
	return false
}

// UserWriteReadable returns true if the filepath is write- and readable by the current user
func (p *PermUser) UserWriteReadable() bool {
	if p.Stat.Mode().Perm()&OsOthR != 0 && p.Stat.Mode().Perm()&OsOthW != 0 {
		return true
	}
	if p.isInGroup() && (p.Stat.Mode().Perm()&OsGroupR != 0 && p.Stat.Mode().Perm()&OsGroupW != 0) {
		return true
	}
	if p.isOwner() && (p.Stat.Mode().Perm()&OsUserR != 0 && p.Stat.Mode().Perm()&OsUserW != 0) {
		return true
	}
	return false
}

// UserWriteExecutable returns true if the filepath is write- and executable by the current user
func (p *PermUser) UserWriteExecutable() bool {
	if p.Stat.Mode().Perm()&OsOthX != 0 && p.Stat.Mode().Perm()&OsOthW != 0 {
		return true
	}
	if p.isInGroup() && (p.Stat.Mode().Perm()&OsGroupX != 0 && p.Stat.Mode().Perm()&OsGroupW != 0) {
		return true
	}
	if p.isOwner() && (p.Stat.Mode().Perm()&OsUserX != 0 && p.Stat.Mode().Perm()&OsUserW != 0) {
		return true
	}
	return false
}

// UserReadExecutable returns true if the filepath is read- and executable by the current user
func (p *PermUser) UserReadExecutable() bool {
	if p.Stat.Mode().Perm()&OsOthR != 0 && p.Stat.Mode().Perm()&OsOthX != 0 {
		return true
	}
	if p.isInGroup() && (p.Stat.Mode().Perm()&OsGroupR != 0 && p.Stat.Mode().Perm()&OsGroupX != 0) {
		return true
	}
	if p.isOwner() && (p.Stat.Mode().Perm()&OsUserR != 0 && p.Stat.Mode().Perm()&OsUserX != 0) {
		return true
	}
	return false
}

// UserWriteReadExecutable returns true if the filepath is write- and read- and executable by the
// current user
func (p *PermUser) UserWriteReadExecutable() bool {
	if p.Stat.Mode().Perm()&OsOthR != 0 && p.Stat.Mode().Perm()&OsOthW != 0 &&
		p.Stat.Mode().Perm()&OsOthX != 0 {
		return true
	}
	if p.isInGroup() && (p.Stat.Mode().Perm()&OsGroupR != 0 && p.Stat.Mode().Perm()&OsGroupW != 0 &&
		p.Stat.Mode().Perm()&OsGroupX != 0) {
		return true
	}
	if p.isOwner() && (p.Stat.Mode().Perm()&OsUserR != 0 && p.Stat.Mode().Perm()&OsUserW != 0 &&
		p.Stat.Mode().Perm()&OsUserX != 0) {
		return true
	}
	return false
}

// isOwner returns true if the file in question is owned by the current user
func (p *PermUser) isOwner() bool {
	return int64(p.UID) == p.CurUserUID
}

// isInGroup returns true if the current user is part the files group
func (p *PermUser) isInGroup() bool {
	for _, gID := range p.CurUserGIDs {
		if int64(p.GID) == gID {
			return true
		}
	}
	return false
}

// getCurUserIds retrieves the userid and groupids of the processes current user and stores them
// in the PermUser struct
func (p *PermUser) getCurUserIds() error {
	curUser, err := user.Current()
	if err != nil {
		return err
	}
	uidInt, err := strconv.ParseInt(curUser.Uid, 10, 32)
	if err != nil {
		return err
	}
	p.CurUserUID = uidInt

	userGroups, err := curUser.GroupIds()
	if err != nil {
		return err
	}
	for _, gid := range userGroups {
		gidInt, err := strconv.ParseInt(gid, 10, 32)
		if err != nil {
			return err
		}
		p.CurUserGIDs = append(p.CurUserGIDs, gidInt)
	}
	return nil
}

// getUIDGID retrieves the owner id and group id of the current PermUser structs file
// Please note: due to the nature of the syscall, this does not work on Windows and Plan9
func getUIDGID(s os.FileInfo) (uid, gid uint32) {
	fs, ok := s.Sys().(*syscall.Stat_t)
	if !ok {
		return 0, 0
	}
	return fs.Uid, fs.Gid
}
