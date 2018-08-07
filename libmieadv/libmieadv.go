package libmieadv

/*
#cgo CFLAGS: -O3
#cgo LDFLAGS: -L. -lmiev0 -L /usr/local/Cellar/gcc/8.2.0/lib/gcc/8/ -lgfortran

void miev0easydriver(double xx, double mre, double mim,
                        double *qext, double *qsca, double *gqsc, double *pmom, int mxmdm);

void miesdist(double *xs, double *ys, int N, double mre, double mim, double wl, double *pmom,
                        int momdim, double *ext, double *sca, double *asy, double *vol, int *ierr);
*/
import (
	"C"
)

import (
	"errors"
	"fmt"
	"github.com/kshmirko/radtran/libmath"
	"math"
	"os"
	"unsafe"
)

type SDistParams struct {
	*AerosolDistrParams
	mre, mim, Wl float64
	momdim       int
}

func NewSDistParams(npts, momdim int,
	r0, r1, gamma, dens, mre, mim, wl float64) *SDistParams {

	return &SDistParams{
		AerosolDistrParams: &AerosolDistrParams{
			npts:  npts,
			r0:    r0,
			r1:    r1,
			gamma: gamma,
			Dens:  dens,
		},
		mre:    mre,
		mim:    mim,
		Wl:     wl,
		momdim: momdim,
	}
}

func (sdp *SDistParams) SetRefIdx(mre, mim float64) {
	sdp.mre = mre
	sdp.mim = mim
}

func (sdp *SDistParams) SetWl(wl float64) {
	sdp.Wl = wl
}

func MieSDist1(params *SDistParams) (ext, sca, asy, vol float64, pmom [][]float64, err error) {
	xs := libmath.Linspace(params.r0, params.r1, params.npts)
	ys := make([]float64, params.npts)

	for i := range xs {
		ys[i] = math.Pow(xs[i], params.gamma)
	}

	return MieSDist(xs, ys, params.mre, params.mim, params.Wl, params.momdim)
}

func MieSDist(xs, ys []float64, mre, mim, wl float64, momdim int) (ext, sca, asy, vol float64, pmom [][]float64, err error) {

	N := len(xs)

	x_tmp := make([]C.double, N)
	y_tmp := make([]C.double, N)
	Nelems := (momdim + 1) * 4
	pmom_tmp := make([]C.double, Nelems)
	for i := 0; i < N; i++ {
		x_tmp[i] = C.double(xs[i])
		y_tmp[i] = C.double(ys[i])
	}

	re := C.double(mre)
	im := C.double(mim)
	wla := C.double(wl)

	var exta, scaa, asya, vola C.double
	var ierr C.int

	C.miesdist(&x_tmp[0], &y_tmp[0], C.int(N),
		re, im, wla, &pmom_tmp[0], C.int(momdim), &exta, &scaa, &asya, &vola, &ierr)

	if ierr != 0 {
		err = errors.New("")
		return
	}

	pmom = make([][]float64, 4)

	for i := 0; i < 4; i++ {
		pmom[i] = make([]float64, momdim+1)
	}

	ext = float64(exta)
	sca = float64(scaa)
	asy = float64(asya)
	vol = float64(vola)
	momdim = momdim + 1
	/*
	   разбираем вывод программы на фортране
	 **/
	for n := 0; n < Nelems; n++ {
		i, j := n/momdim, n%momdim
		pmom[i][j] = float64(pmom_tmp[n])
	}
	err = nil
	return

}

func Miev0EasyDriver(xx, mre, mim float64, momdim int) (qext, qsca, gqsc float64, pmom [][]float64) {

	Nelems := (momdim + 1) * 4
	pmom_tmp := make([]C.double, Nelems)
	pmom = make([][]float64, 4)

	for i := 0; i < 4; i++ {
		pmom[i] = make([]float64, momdim+1)
	}

	var qext_tmp, qsca_tmp, gqsc_tmp C.double

	C.miev0easydriver(C.double(xx),
		C.double(mre),
		C.double(mim),
		&qext_tmp,
		&qsca_tmp,
		&gqsc_tmp,
		(*C.double)(unsafe.Pointer(&pmom_tmp[0])),
		C.int(momdim),
	)

	qext = float64(qext_tmp)
	qsca = float64(qsca_tmp)
	gqsc = float64(gqsc_tmp)

	momdim = momdim + 1
	/*
	   разбираем вывод программы на фортране
	 **/
	for n := 0; n < Nelems; n++ {
		i, j := n/momdim, n%momdim
		pmom[i][j] = float64(pmom_tmp[n])
	}

	return

}

func Wiscombe2Evans(wis *[][]float64) (evans [][]float64) {
	/*
	   матрицы wiscombe хранят моменты лежандра (pl), а марицы эванса - коэффициенты (fl)
	   fl = (2*l+1)pl
	*/

	evans = make([][]float64, 6)
	for i := 0; i < len(evans); i++ {
		evans[i] = make([]float64, len((*wis)[0]))
	}

	norm := (*wis)[0][0] + (*wis)[1][0]
	for i := 0; i < len(evans[0]); i++ {
		factor := float64(2*i + 1)
		evans[0][i] = ((*wis)[0][i] + (*wis)[1][i]) / norm * factor
		evans[1][i] = ((*wis)[1][i] - (*wis)[0][i]) / norm * factor
		evans[2][i] = 2.0 * (*wis)[2][i] / norm * factor
		evans[3][i] = 2.0 * (*wis)[3][i] / norm * factor
		evans[4][i] = evans[0][i]
		evans[5][i] = evans[2][i]
	}
	return
}

func PrintData(fname string, evans *[][]float64, taue, omega float64) {

	f, err := os.Create(fname)
	if err != nil {
		panic("Error!")
	}
	defer f.Close()
	Nl := len((*evans)[0])
	taus := taue * omega
	str := fmt.Sprintf("%12.5E\n%12.5E\n%12.5E\n%d\n", taue, taus, omega, Nl-1)
	f.WriteString(str)

	for i := 0; i < Nl; i++ {
		str = fmt.Sprintf("%4d%13.8f%13.8f%13.8f%13.8f%13.8f%13.8f\n", i,
			(*evans)[0][i], (*evans)[1][i], (*evans)[2][i],
			(*evans)[3][i], (*evans)[4][i], (*evans)[5][i])
		f.WriteString(str)
	}
	f.Sync()

}
