package libmath

func Linspace1(x0, x1 float64, x *[]float64) {
	N := len(*x)
	dx := (x1 - x0) / float64(N-1)

	for i := 0; i < N; i++ {
		(*x)[i] = x0 + float64(i)*dx
	}
}

func Linspace(x0, x1 float64, N int) (res []float64) {
	res = make([]float64, N)

	dx := (x1 - x0) / float64(N-1)

	for i := 0; i < N; i++ {
		res[i] = x0 + float64(i)*dx
	}

	return
}

func Trapz(x, y *[]float64) (sum float64) {
	/*
	   интегрирование методом трапеций одномерной функции, заданной таблично
	*/
	sum = 0.0
	if len(*x) != len(*y) {
		return sum
	}

	for i := 1; i < len(*x); i++ {
		sum += 0.5 * ((*y)[i] + (*y)[i-1]) * ((*x)[i] - (*x)[i-1])
	}

	return sum
}
