package main

import (
	"fmt"
	"testing"

	"github.com/kshmirko/radtran/libmath"
	"github.com/kshmirko/radtran/libssrt"
)

func TestSsrt(t *testing.T) {
	a := libssrt.New(".scat_file048", 0.1, 0.9, 9.9, 0.532, 10)

	theta := libmath.Linspace(0, 90, 101)

	for i := range theta {
		fmt.Printf("%7.3f %7.3f %10.3f %10.3f %10.3f\n",
			theta[i],
			0.0,
			a.L0(0.13, theta[i], 0.0),
			a.L1(0.13, theta[i], 0.0),
			a.L(0.13, theta[i], 0.0),
		)
	}

	for i := range theta {
		fmt.Printf("%7.3f %7.3f %10.3f %10.3f %10.3f\n",
			theta[i],
			180.0,
			a.L0(0.13, theta[i], 180.0),
			a.L1(0.13, theta[i], 180.0),
			a.L(0.13, theta[i], 180.0),
		)
	}
}
