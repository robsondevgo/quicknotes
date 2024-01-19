package main

import (
	"errors"
	"fmt"
)

var CustomErr100 = &CustomError{msg: "aconteceu um erro 100", code: 100}
var CustomErr200 = &CustomError{msg: "aconteceu um erro 200", code: 200}

type CustomError struct {
	msg  string
	code int
}

func NewCustomError(msg string, code int) error {
	return &CustomError{msg: msg, code: code}
}

func (ce *CustomError) Error() string {
	return ce.msg
}

func subOperation() error {
	return CustomErr100
}

func execute() error {
	err := subOperation()
	if err != nil {
		return fmt.Errorf("%s: (%w)", err, err)
	}
	return nil
}

// wrap e unrap
func main() {
	err := execute()
	if err != nil {
		fmt.Println(err)
		err = errors.Unwrap(err)
		fmt.Println(err)
		err = errors.Unwrap(err)
		fmt.Println(err)
	}
	// var err100 *CustomError
	// if errors.As(err, &err100) {
	// 	fmt.Println(err100.code)
	// } else {
	// 	fmt.Println(err)
	// }
}
