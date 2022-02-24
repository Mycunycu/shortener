package helpers

import "errors"

var ErrUnique = errors.New("duplicate unique field")
var ErrDeletedItem = errors.New("requested thing is deleted")
