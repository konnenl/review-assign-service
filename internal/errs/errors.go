package errs

import (
	"errors"
)

var (
	ErrTeamExists   = errors.New("team already exists")
	ErrPRExists     = errors.New("pull request already exists")
	ErrTeamNotFound = errors.New("team not found")
	ErrUserNotFound = errors.New("user not found")
	ErrPRNotFound   = errors.New("pull request not found")
)
