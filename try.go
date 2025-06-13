package onerr

func Try(e error) {
	if e != nil {
		panic(e)
	}
}

func Try1[T any](val T, e error) T {
	if e == nil {
		return val
	} else {
		panic(e)
	}
}

func Try2[T1 any, T2 any](v1 T1, v2 T2, e error) (T1, T2) {
	if e == nil {
		return v1, v2
	} else {
		panic(e)
	}
}

func Try3[T1 any, T2 any, T3 any](v1 T1, v2 T2, v3 T3, e error) (T1, T2, T3) {
	if e == nil {
		return v1, v2, v3
	} else {
		panic(e)
	}
}

func Enter(f func()) error {
	return runAndCatch(f)
}

func Scope(f func()) { f() }

func Eval[T any](f func() T) T                   { return f() }
func Eval2[T any, U any](f func() (T, U)) (T, U) { return f() }
