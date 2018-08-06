package libatmos

import (
	//"bufio"
	"fmt"
	"math"
	"os"

	"github.com/kshmirko/radtran/libmath"
	"github.com/kshmirko/radtran/libmieadv"
)

const (
	h_mol = 8.4273  //km
	h_max = 45.0000 // km
)

var Evansm [][]float64

func init() {
	// Инициализируем матрицу Эванса для молекулярного расеяния
	// так как рна идентична относительно длин волн
	// не заморачиваемся относительно длин волн
	// и показателя преломления
	npts := 101
	wl := 0.355
	xs := libmath.Linspace(0.0001, 0.0005, npts)
	ys := make([]float64, npts)
	for i := range xs {
		ys[i] = math.Pow(xs[i], -7.0)
	}
	_, _, _, _, pmomm, _ := libmieadv.MieSDist(xs, ys, 1.4, 0.00, wl, 60)
	Evansm = libmieadv.Wiscombe2Evans(&pmomm)
}

func MakeAtmosphereEvans(evansa *[][]float64,
	taua, taum float64) (ret [][]float64) {

	N := len(Evansm[0]) //int(math.Min(float64(len(Evansm[0])), float64(len((*evansa)[0]))))
	ret = make([][]float64, 6)
	for i := 0; i < 6; i++ {
		ret[i] = make([]float64, len((*evansa)[0]))
	}

	for i := 0; i < N; i++ {
		ret[0][i] = (taua*(*evansa)[0][i] + taum*Evansm[0][i]) / (taua + taum)
		ret[1][i] = (taua*(*evansa)[1][i] + taum*Evansm[1][i]) / (taua + taum)
		ret[2][i] = (taua*(*evansa)[2][i] + taum*Evansm[2][i]) / (taua + taum)
		ret[3][i] = (taua*(*evansa)[3][i] + taum*Evansm[3][i]) / (taua + taum)
		ret[4][i] = (taua*(*evansa)[4][i] + taum*Evansm[4][i]) / (taua + taum)
		ret[5][i] = (taua*(*evansa)[5][i] + taum*Evansm[5][i]) / (taua + taum)
	}

	return
}

func Extm_at_h(h []float64, wl float64) []float64 {
	extm0 := Tau_m0(wl) / h_mol
	N := len(h)
	ret := make([]float64, N)

	for i := range h {
		ret[i] = extm0 * math.Exp(-h[i]/h_mol)
	}

	return ret
}

func Exta_at_h(h []float64, taua, hpbl float64) []float64 {
	exta0 := taua / hpbl
	N := len(h)
	ret := make([]float64, N)

	for i := range h {
		ret[i] = exta0 * math.Exp(-h[i]/hpbl)
	}
	return ret
}

func Tau_m0(wl float64) float64 {
	return 0.008569 / math.Pow(wl, 4.0) *
		(1.0 +
			0.0113/math.Pow(wl, 2.0) +
			0.00013/math.Pow(wl, 4.0))
}

func Ext_m0(wl float64) float64 {
	return Tau_m0(wl) / h_mol
}

func Ext_m(wl, h float64) float64 {
	// Вычисляет extinction  в интервале от 0 до h.
	return Ext_m0(wl) * math.Exp(-h/h_mol)
}

func Tau_m(wl, h float64) float64 {
	return Ext_m0(wl) * h_mol * (1.0 - math.Exp(-h/h_mol))
}

func Ext_a(taua, h, hpbl float64) float64 {
	return taua / hpbl * math.Exp(-h/hpbl)
}

func Tau(taua, h, hpbl float64) float64 {
	return taua * (1.0 - math.Exp(-h/hpbl))
}

func MakeLayFile(fname string, taua, taum, wl, hpbl float64, alts []float64) {

	f, err := os.Create(fname)
	if err != nil {
		panic("")
	}
	defer f.Close()

	for i := range alts {

		//taua_tmp := taua - Tau(taua, alts[i], hpbl)
		//taum_tmp := taum - Tau(taum, alts[i], h_mol)
		strbuf := fmt.Sprintf("%6.2f%7.2f%12.6f  '%s'\n", alts[i],
			0.0, 0.0, ".scat_file_001")
		f.WriteString(strbuf)
		f.Sync()
	}
}
