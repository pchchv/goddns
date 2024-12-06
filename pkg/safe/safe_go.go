package safe

import (
	"errors"
	"log"
)

var ErrRecoverFromPanic = errors.New("recover from panic")

// Try invokes #{fn} with panic handler.
// If #{fn} causes a panic, #{ErrRecoverFromPanic} will be returned.
// Otherwise, nil will be returned.
func Try(fn func()) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = ErrRecoverFromPanic
		}
	}()
	fn()
	return nil
}

// SafeGo spawns a go-routine to run {#fn} with panic handler.
func SafeGo(fn func()) {
	go func() {
		if err := Try(fn); err != nil {
			log.Panicf("panic in go-routine: %s", err.Error())
		}
	}()
}
