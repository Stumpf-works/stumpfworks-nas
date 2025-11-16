// Revision: 2025-11-16 | Author: Claude | Version: 1.1.1
package sysutil

import (
	"fmt"
	"os/user"
	"strconv"
)

// LookupUID looks up a user's UID by username
func LookupUID(username string) (int, error) {
	u, err := user.Lookup(username)
	if err != nil {
		return -1, fmt.Errorf("failed to lookup user %s: %w", username, err)
	}

	uid, err := strconv.Atoi(u.Uid)
	if err != nil {
		return -1, fmt.Errorf("invalid UID for user %s: %w", username, err)
	}

	return uid, nil
}

// LookupGID looks up a group's GID by group name
func LookupGID(groupname string) (int, error) {
	g, err := user.LookupGroup(groupname)
	if err != nil {
		return -1, fmt.Errorf("failed to lookup group %s: %w", groupname, err)
	}

	gid, err := strconv.Atoi(g.Gid)
	if err != nil {
		return -1, fmt.Errorf("invalid GID for group %s: %w", groupname, err)
	}

	return gid, nil
}

// LookupUsername looks up a username by UID
func LookupUsername(uid int) (string, error) {
	u, err := user.LookupId(strconv.Itoa(uid))
	if err != nil {
		return "", fmt.Errorf("failed to lookup UID %d: %w", uid, err)
	}

	return u.Username, nil
}

// LookupGroupname looks up a group name by GID
func LookupGroupname(gid int) (string, error) {
	g, err := user.LookupGroupId(strconv.Itoa(gid))
	if err != nil {
		return "", fmt.Errorf("failed to lookup GID %d: %w", gid, err)
	}

	return g.Name, nil
}

// ParseUIDOrUsername parses a string that could be either a numeric UID or username
// Returns the UID and username
func ParseUIDOrUsername(input string) (uid int, username string, err error) {
	// Try parsing as numeric UID first
	if parsedUID, parseErr := strconv.Atoi(input); parseErr == nil {
		// It's a numeric UID, lookup the username
		username, err = LookupUsername(parsedUID)
		if err != nil {
			return -1, "", err
		}
		return parsedUID, username, nil
	}

	// Treat as username
	uid, err = LookupUID(input)
	if err != nil {
		return -1, "", err
	}

	return uid, input, nil
}

// ParseGIDOrGroupname parses a string that could be either a numeric GID or group name
// Returns the GID and group name
func ParseGIDOrGroupname(input string) (gid int, groupname string, err error) {
	// Try parsing as numeric GID first
	if parsedGID, parseErr := strconv.Atoi(input); parseErr == nil {
		// It's a numeric GID, lookup the group name
		groupname, err = LookupGroupname(parsedGID)
		if err != nil {
			return -1, "", err
		}
		return parsedGID, groupname, nil
	}

	// Treat as group name
	gid, err = LookupGID(input)
	if err != nil {
		return -1, "", err
	}

	return gid, input, nil
}
