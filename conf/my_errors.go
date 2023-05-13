package conf

import (
	"fmt"
	"runtime"

	"github.com/pkg/errors"
)

var (
	ErrNotFound = errors.New("Not Found")
)

func AddCallerInfo(err error) error {
	if err != nil {
		// Get the caller's PC (program counter)
		pc, file, line, _ := runtime.Caller(1)

		// Get the function name
		funcName := runtime.FuncForPC(pc).Name()

		// Wrap the error and add caller information
		err = errors.Wrap(err, fmt.Sprintf(" >> [%s %s:%d]", funcName, file, line))
	}

	return err
}
