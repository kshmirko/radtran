package librt3

// Пакет  - обертка для вызова программы расчета rt3
// RT3 - Radiation transfer code

import (
	//"bytes"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"os/exec"

	"github.com/kshmirko/radtran/libmath"
	"github.com/kshmirko/radtran/libmieadv"
	"github.com/kshmirko/radtran/libprepare"
)

const (
	n_layers = 49
	h_min    = 0.0
	h_max    = 48.0
)

type RT3Object struct {
	sdp           *libmieadv.SDistParams
	NStokes       int
	NGauss        int
	QuadType      string
	NAzimuth      int
	Layfile       string
	DeltaScaling  string
	Source        int
	SolFlux       float64
	Theta         float64
	GTemp         float64
	SurfType      string
	GRefl         float64
	STemp         float64
	RadUnits      string
	PolType       string
	NOutputLayers int
	OutLayIdx     int
	NOutAzi       int
	OutFname      string
}

type RT3Output struct {
	z, wl              float64
	Ang, phi, mu, I, Q *[]float64
}

func New() *RT3Object {
	return &RT3Object{
		sdp:           nil,
		NStokes:       2,
		NGauss:        32,
		QuadType:      "G",
		NAzimuth:      4,
		Layfile:       "testlay.lay",
		DeltaScaling:  "Y",
		Source:        1,
		SolFlux:       1.0,
		Theta:         20.0,
		GTemp:         0.0,
		SurfType:      "L",
		GRefl:         0.0,
		STemp:         0.0,
		RadUnits:      "W",
		PolType:       "IQ",
		NOutputLayers: 1,
		OutLayIdx:     n_layers,
		NOutAzi:       2,
		OutFname:      "tmpout.out",
	}
}

func (rt *RT3Object) GetSDistParams() *libmieadv.SDistParams {
	return rt.sdp
}

func (rt *RT3Object) SetSDistParams(sdp *libmieadv.SDistParams) {
	rt.sdp = sdp
}

func (rt *RT3Object) Prepare(taua, taum, hpbl float64) {

	alts := libmath.Linspace(h_max, h_min, n_layers)

	_ = libprepare.Prepare(rt.Layfile, rt.sdp, taua, taum, hpbl, alts)
}

func (rt *RT3Object) MakeScript() {
	f, err := os.Create("cmdfile.dat")
	if err != nil {
		log.Fatal(err)
	}

	io.WriteString(f, fmt.Sprintf("%d\n", rt.NStokes))
	io.WriteString(f, fmt.Sprintf("%d\n", rt.NGauss))
	io.WriteString(f, fmt.Sprintf("%s\n", rt.QuadType))
	io.WriteString(f, fmt.Sprintf("%d\n", rt.NAzimuth))
	io.WriteString(f, fmt.Sprintf("%s\n", rt.Layfile))
	io.WriteString(f, fmt.Sprintf("%s\n", rt.DeltaScaling))
	io.WriteString(f, fmt.Sprintf("%d\n", rt.Source))
	io.WriteString(f, fmt.Sprintf("%f\n", rt.SolFlux))
	io.WriteString(f, fmt.Sprintf("%f\n", rt.Theta))
	io.WriteString(f, fmt.Sprintf("%f\n", rt.GTemp))
	io.WriteString(f, fmt.Sprintf("%s\n", rt.SurfType))
	io.WriteString(f, fmt.Sprintf("%f\n", rt.GRefl))
	io.WriteString(f, fmt.Sprintf("%f\n", rt.STemp))
	io.WriteString(f, fmt.Sprintf("%f\n", rt.sdp.Wl))
	io.WriteString(f, fmt.Sprintf("%s\n", rt.RadUnits))
	io.WriteString(f, fmt.Sprintf("%s\n", rt.PolType))
	io.WriteString(f, fmt.Sprintf("%d\n", rt.NOutputLayers))
	io.WriteString(f, fmt.Sprintf("%d\n", rt.OutLayIdx))
	io.WriteString(f, fmt.Sprintf("%d\n", rt.NOutAzi))
	io.WriteString(f, fmt.Sprintf("%s\n", rt.OutFname))

	f.Close()

}

func (rt *RT3Object) Run() *RT3Output {
	cmd := exec.Command("./rt3")
	stdin, err := cmd.StdinPipe()

	if err != nil {
		log.Fatal(err)
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}

	cmd.Start()
	if err != nil {
		log.Fatal(err)
	}

	defer stdin.Close()
	defer stdout.Close()

	//outC := make(chan string)
	// copy the output in a separate goroutine so printing can't block indefinitely
	//go func() {
	//	var buf bytes.Buffer
	//	io.Copy(&buf, stdout)
	//	outC <- buf.String()
	//}()

	io.WriteString(stdin, fmt.Sprintf("%d\n", rt.NStokes))
	io.WriteString(stdin, fmt.Sprintf("%d\n", rt.NGauss))
	io.WriteString(stdin, fmt.Sprintf("%s\n", rt.QuadType))
	io.WriteString(stdin, fmt.Sprintf("%d\n", rt.NAzimuth))
	io.WriteString(stdin, fmt.Sprintf("%s\n", rt.Layfile))
	io.WriteString(stdin, fmt.Sprintf("%s\n", rt.DeltaScaling))
	io.WriteString(stdin, fmt.Sprintf("%d\n", rt.Source))
	io.WriteString(stdin, fmt.Sprintf("%f\n", rt.SolFlux))
	io.WriteString(stdin, fmt.Sprintf("%f\n", rt.Theta))
	io.WriteString(stdin, fmt.Sprintf("%f\n", rt.GTemp))
	io.WriteString(stdin, fmt.Sprintf("%s\n", rt.SurfType))
	io.WriteString(stdin, fmt.Sprintf("%f\n", rt.GRefl))
	io.WriteString(stdin, fmt.Sprintf("%f\n", rt.STemp))
	io.WriteString(stdin, fmt.Sprintf("%f\n", rt.sdp.Wl))
	io.WriteString(stdin, fmt.Sprintf("%s\n", rt.RadUnits))
	io.WriteString(stdin, fmt.Sprintf("%s\n", rt.PolType))
	io.WriteString(stdin, fmt.Sprintf("%d\n", rt.NOutputLayers))
	io.WriteString(stdin, fmt.Sprintf("%d\n", rt.OutLayIdx))
	io.WriteString(stdin, fmt.Sprintf("%d\n", rt.NOutAzi))
	io.WriteString(stdin, fmt.Sprintf("%s\n", rt.OutFname))

	//out := <-outC
	//fmt.Print(out)

	cmd.Wait()

	f, err := os.Open(rt.OutFname)
	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	for i := 0; i < 11; i++ {
		_ = ReadLine(f)
	}

	NLines := rt.NGauss * rt.NOutAzi
	z := 0.0
	phi := make([]float64, NLines)
	mu := make([]float64, NLines)
	I := make([]float64, NLines)
	Q := make([]float64, NLines)
	ang := make([]float64, NLines)
	deg2rad := math.Pi / 180.0
	j := 0
	for i := 0; i < NLines*2; i++ {
		fmt.Fscanf(f, "%f %f %f %f %f\n", &z, &phi[j], &mu[j], &I[j], &Q[j])
		if mu[j] > 0.0 {
			fmt.Printf("%f %f %f %f %f\n", z, phi[j], mu[j], I[j], Q[j])
			ang[j] = math.Cos(phi[j]*deg2rad) * math.Acos(mu[j]) / deg2rad
			j++
		}
	}

	return &RT3Output{
		z:   z,
		phi: &phi,
		mu:  &mu,
		I:   &I,
		Q:   &Q,
		Ang: &ang,
		wl:  rt.sdp.Wl,
	}
}

func ReadLine(r io.Reader) string {
	var ret []byte
	var tmp []byte
	tmp = make([]byte, 1)
	for tmp[0] != 10 {
		io.ReadFull(r, tmp)

		ret = append(ret, tmp[0])
	}

	return string(ret[:])

}
