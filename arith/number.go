package arith

type Number interface {
	int | uint | int8 | uint8 | int16 | uint16 | int32 | uint32 | int64 | uint64 | float32 | float64
}

func Max[N Number](left, right N) N {
	if left > right {
		return left
	} else {
		return right
	}
}

func Min[N Number](left, right N) N {
	if left > right {
		return right
	} else {
		return left
	}
}
