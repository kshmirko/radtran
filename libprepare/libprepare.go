package libprepare

import (
	"fmt"
	"os"

	"github.com/kshmirko/radtran/libatmos"
	"github.com/kshmirko/radtran/libmath"
	"github.com/kshmirko/radtran/libmieadv"
)

func Prepare(layfname string,
	params *libmieadv.SDistParams,
	taua, taum, hpbl float64,
	alts []float64) (evans [][]float64) {

	ext, sca, _, _, pmom, err := libmieadv.MieSDist1(params)

	if err != nil {
		panic("Error calculating size distribution properties")
	}
	ext = ext * params.Dens
	sca = sca * params.Dens
	omega0 := sca / ext

	evans = libmieadv.Wiscombe2Evans(&pmom)

	f, err := os.Create(layfname)
	if err != nil {
		panic("Error")
	}
	defer f.Close()

	exta := libatmos.Exta_at_h(alts, taua, hpbl)
	extm := libatmos.Extm_at_h(alts, params.Wl)

	fmt.Println(-libmath.Trapz(&alts, &exta))
	fmt.Println(-libmath.Trapz(&alts, &extm))
	for i := range alts {
		//altinm := alts[i] // * 1000.0

		exta_at := exta[i] //libatmos.Ext_a(taua, altinm, hpbl)
		extm_at := extm[i] //libatmos.Ext_m(params.Wl, altinm)
		fmt.Printf("%10.7f %10.7f %10.7f\n", alts[i], exta_at, extm_at)
		evans_tmp := libatmos.MakeAtmosphereEvans(&evans, exta_at, extm_at)
		scat := exta_at*omega0 + extm_at
		extt := extm_at + exta_at
		omega := scat / extt
		tmpfname := fmt.Sprintf(".scat_file%03d", i)
		libmieadv.PrintData(tmpfname, &evans_tmp, extt, omega)
		strbuf := fmt.Sprintf("%6.2f%7.2f%12.6f  '%s'\n", alts[i],
			0.0, 0.0, tmpfname)
		f.WriteString(strbuf)
	}
	f.Sync()

	return

}
