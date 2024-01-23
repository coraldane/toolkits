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

func ScaleDecimal(inputValue float64, scale int) float64 {
	var delta float64
	if inputValue > 0 {
		delta = 0.5
	} else {
		delta = -0.5
	}
	return math.Trunc(inputValue*math.Pow10(scale)+delta) / math.Pow10(scale)
	//format := strings.Join([]string{"%.", "f"}, strconv.Itoa(scale))
	//value, _ := strconv.ParseFloat(fmt.Sprintf(format, inputValue), 64)
	//return value
}
