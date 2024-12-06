package safe

import "errors"

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
