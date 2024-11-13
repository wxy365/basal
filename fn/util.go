package fn

func Call(f func(), e *error) {
	if *e != nil {
		return
	}
	f()
}
