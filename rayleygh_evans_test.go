package main

import (
	"fmt"
	"testing"

	"github.com/kshmirko/radtran/libatmos"
)

func TestRayleighEvans(t *testing.T) {
	for i := 0; i < 10; i++ {
		fmt.Printf("Em[%02d] = %7.3f %7.3f %7.3f %7.3f\n", i, libatmos.Evansm[0][i], libatmos.Evansm[1][i], libatmos.Evansm[2][i], libatmos.Evansm[3][i])
	}
}
