package er

import "errors"

var ErrUniqueValue = errors.New("provided URL already shorted and stored in DB")
