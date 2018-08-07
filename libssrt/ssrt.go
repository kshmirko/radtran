package libssrt

import (
	"fmt"
	"math"
	"os"

	"github.com/kshmirko/radtran/libmath"
)

type SSRTData struct {
	Sza, Omega, Taue, Wl float64
	P1, P2               *[]float64
	NTheta               int
}

func New(fname string, taue, omega, sza, wl float64, ntheta int) *SSRTData {

	f, err := os.Open(fname)
	if err != nil {
		panic("Error opening scattering file")
	}

	defer f.Close()

	var ext, sca, omega1 float64
	var NL int

	// читаем коэффициенты лежандра из файла
	fmt.Fscanf(f, "%f\n", &ext)
	fmt.Fscanf(f, "%f\n", &sca)
	fmt.Fscanf(f, "%f\n", &omega1)
	fmt.Fscanf(f, "%d\n", &NL)

	var P1, P2 []float64
	P1 = make([]float64, NL+1)
	P2 = make([]float64, NL+1)

	var j int
	var P3, P4, P5, P6 float64
	for i := 0; i <= NL; i++ {

		fmt.Fscanf(f, "%d %f %f %f %f %f %f\n", &j, &P1[i], &P2[i], &P3, &P4, &P5, &P6)
	}

	return &SSRTData{
		Sza:    sza,
		Omega:  omega,
		Taue:   taue,
		Wl:     wl,
		NTheta: ntheta,
		P1:     &P1,
		P2:     &P2,
	}
}

func (v *SSRTData) L0(tau, theta, phi float64) float64 {
	//Прямое солнечное излучение
	delta := theta - v.Sza + phi
	if math.Abs(delta) > 0.01 {
		return 0.0
	}

	theta = theta * math.Pi / 180.0

	return math.Exp(-tau / math.Cos(theta))
}

func (v *SSRTData) L1(tau, theta, phi float64) float64 {
	delta := math.Abs(v.Sza-theta) + phi
	if delta <= 0.01 {
		return 0.0
	}

	theta = theta * math.Pi / 180.0
	theta0 := v.Sza * math.Pi / 180.0
	phi = phi * math.Pi / 180.0

	u0 := math.Cos(theta0)
	u := math.Cos(theta)
	uphi := math.Cos(phi)

	ua := u*u0 + math.Sqrt(1.0-u*u)*math.Sqrt(1.0-u0*u0)*uphi
	F1_i := libmath.Coef2phase(v.P1, ua) //Scale by 0.5 beacause ||F1_i|| = 2

	Ld := F1_i * u0 / (u0 - u) * (math.Exp(-tau/u0) - math.Exp(-tau/u)) / u0

	return Ld
}

func (v *SSRTData) L(tau, theta, phi float64) float64 {
	return v.L0(tau, theta, phi) + v.Omega*v.L1(tau, theta, phi)
}
