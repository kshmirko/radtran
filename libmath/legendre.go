package libmath

func Coef2phase(f *[]float64, x float64) (phase float64) {
	var l int
	var p0, p1 float64
	l = 0
	p0 = 1.0
	p1 = x

	phase = p0*(*f)[0] + p1*(*f)[1]

	for l = 2; l < len(*f); l++ {
		p2 := (float64(2*l-1)*x*p1 - float64(l-1)*p0) / float64(l)
		phase += (*f)[l] * p2
		p0 = p1
		p1 = p2
	}
	return
}

func mom2phase(f *[]float64, x float64) (phase float64) {
	var l int
	var p0, p1 float64
	l = 0
	p0 = 1.0
	p1 = x

	phase = p0*(*f)[0] + 3.0*p1*(*f)[1]

	for l = 2; l < len(*f); l++ {
		p2 := (float64(2*l-1)*x*p1 - float64(l-1)*p0) / float64(l)
		phase += float64(2*l+1) * (*f)[l] * p2
		p0 = p1
		p1 = p2
	}
	return
}

func PhaseFunction(xmu, coef *[]float64) (phase []float64) {
	phase = make([]float64, len(*xmu))

	for i := 0; i < len(*xmu); i++ {
		phase[i] = Coef2phase(coef, (*xmu)[i])
	}
	return
}

func PhaseFunction1(xmu, moms *[]float64) (phase []float64) {
	phase = make([]float64, len(*xmu))

	for i := 0; i < len(*xmu); i++ {
		phase[i] = mom2phase(moms, (*xmu)[i])
	}
	return
}
