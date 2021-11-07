package customErrors

import (
	"errors"
	"strconv"
)

func Success() string {
	return "OK"
}

func InternalServerError() error {
	return errors.New("internal server error")
}

func NonSuccessStatusCode(statusCode int) error {
	return errors.New("Request failed with status code " + strconv.Itoa(statusCode))
}
