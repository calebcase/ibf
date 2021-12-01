package ibf

import "github.com/zeebo/errs"

// Errors for this package.
var (
	Error = errs.Class("ibf")

	ErrNoPureCell = Error.New("no pure cell")
	ErrEmptySet   = Error.New("empty set")
)
