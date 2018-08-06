package main

import (
	"fmt"
	//"github.com/kshmirko/radtran/libatmos"
	"github.com/kshmirko/radtran/libmath"
	"github.com/kshmirko/radtran/libmieadv"
	//"github.com/kshmirko/radtran/libplot"
	"github.com/kshmirko/radtran/libprepare"
	"math"
)

func main() {

	params := libmieadv.NewSDistParams(101, 60, 0.1, 1.0, -3.0, 1.0, 1.4, 0.01, 0.750)
	alts := libmath.Linspace(48, 0, 49)
	fmt.Println(alts)
	evans := libprepare.Prepare("testlay.lay", params, 0.1, 0, 3.0, alts)

	xmu := libmath.Linspace(0.0, math.Pi, 101)
	for i := range xmu {
		xmu[i] = math.Cos(xmu[i])
	}

	I0 := libmath.PhaseFunction(&xmu, &evans[0])
	fmt.Printf("∫F0(μ) dμ = %7.3f\n", -libmath.Trapz(&xmu, &I0))
	/*
		npts := 151
		nmoms := 60
		mre := 1.4
		mim := 0.01
		xs := libmath.Linspace(0.1, 1.0, npts)
		ys := make([]float64, npts)

		for i := range xs {
			ys[i] = math.Pow(xs[i], -3.5)
		}

		wl := 0.355
		ext, sca, asy, vol, pmom, err := libmieadv.MieSDist(xs, ys, mre, mim, wl, nmoms)
		if err != nil {
			panic("Error calculating size distribution properties")
		}

		omega := sca / ext

		fmt.Printf("ɷ=%7.3f\n", omega)
		fmt.Printf("asy=%12.3e\n", asy)
		fmt.Printf("vol=%12.3e\n", vol)

		evans := libmieadv.Wiscombe2Evans(&pmom)

		xmu := libmath.Linspace(0.0, math.Pi, npts)
		for i := range xmu {
			xmu[i] = math.Cos(xmu[i])
		}

		taua := 0.3
		taum := libatmos.Tau_m0(wl)

		evans = libatmos.MakeAtmosphereEvans(&evans, taua*omega, taum)
		omega = (taua*omega + taum) / (taua + taum)

		I0 := libmath.PhaseFunction(&xmu, &evans[0])
		I1 := libmath.PhaseFunction(&xmu, &evans[1])

		// Проверяем нормировку нашей индикатриссы
		// в нашем случае интеграл должен быть равен 2
		libplot.VizualizePolarization(wl, &xmu, &I0, &I1)
		fmt.Printf("∫F0(μ) dμ = %7.3f, ɷ=%7.3f\n", -libmath.Trapz(&xmu, &I0), omega)
		libmieadv.PrintData(&evans, (taua + taum), omega)
		libatmos.MakeLayFile("test.lay", taua, taum, wl, 3000.0, []float64{15.0, 5.0, 0.0})*/
}
