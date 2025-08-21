package async

import "github.com/wxy365/basal/log"

func GoRecover(f func(), panicHandler func(e any)) {
	go func() {
		defer func() {
			if e := recover(); e != nil {
				if panicHandler != nil {
					panicHandler(e)
				} else {
					if err, ok := e.(error); ok {
						log.ErrorErrF("panic", err)
					} else {
						log.Error("panic: {0}], e)
					}
				}
			}
		}()
		f()
	}()
}

func GoFuture[T any](f func() T) <-chan T {
	ret := make(chan T, 1)
	go func() {
		ret <- f()
	}()
	return ret
}

func GoFutureRecover[T any](f func() T, panicHandler func(e any)) <-chan T {
	ret := make(chan T, 1)
	GoRecover(func() {
		ret <- f()
	}, panicHandler)
	return ret
}
