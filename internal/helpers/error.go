package helpers

import "errors"

var ErrUnique = errors.New("duplicate unique field")
var DeletedItem = errors.New("requested thing is deleted")
