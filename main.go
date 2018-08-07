package main

import (
	"fmt"
	//"github.com/kshmirko/radtran/libatmos"
	//"github.com/kshmirko/radtran/libmath"
	"github.com/kshmirko/radtran/libmieadv"
	//"github.com/kshmirko/radtran/libplot"
	//"github.com/kshmirko/radtran/libprepare"
	"github.com/kshmirko/radtran/librt3"
	"runtime"
	//"math"
)

func main() {
	runtime.GOMAXPROCS(1)
	// func NewSDistParams(npts, momdim, r0, r1, gamma, dens, mre, mim, wl) float64)
	params := libmieadv.NewSDistParams(101,
		60,
		0.1,
		1.0,
		-3.0,
		1.0,
		1.4, 0.01,
		0.532)

	rt3par := librt3.New()
	rt3par.SetSDistParams(params)
	rt3par.Prepare(0.1, 0.0, 3.0)
	rt3par.Theta = 30.0
	ret := rt3par.Run()

	ret.PlotPolarization("pol.pdf")
	ret.PlotIntensities("Int.pdf")

	for i := 0; i < len(*(ret.Ang)); i++ {
		fmt.Printf("%10.3f%10.3f\n", (*(ret.Ang))[i], (*(ret.I))[i])
	}

	//	alts := libmath.Linspace(48, 0, 49)

	//	evans := libprepare.Prepare("testlay.lay", params, 0.1, 0, 3.0, alts)

	//	xmu := libmath.Linspace(0.0, math.Pi, 101)
	//	for i := range xmu {
	//		xmu[i] = math.Cos(xmu[i])
	//	}

	//	I0 := libmath.PhaseFunction(&xmu, &evans[0])
	//	fmt.Printf("∫F0(μ) dμ = %7.3f\n", -libmath.Trapz(&xmu, &I0))

	//	fmt.Println(rt3par)

}
