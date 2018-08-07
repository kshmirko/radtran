package librt3

import (
	"github.com/kshmirko/radtran/libplot"
)

func (ret *RT3Output) PlotIntensities(fname string) {
	libplot.VizualizeIntensity(fname, ret.wl, ret.Ang, ret.I, ret.Q)
}

func (ret *RT3Output) PlotPolarization(fname string) {
	libplot.VizualizePolarization(fname, ret.wl, ret.Ang, ret.I, ret.Q)
}
