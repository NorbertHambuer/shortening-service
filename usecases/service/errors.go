package service

import "fmt"

var ErrCodeAlreadyExists = fmt.Errorf("code already exists in the database")
var ErrCheckCode = fmt.Errorf("unable to check if the code already exists in the database")
