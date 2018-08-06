package libplot

import (
	"fmt"
	"math"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	//"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
)

func VizualizePolarization(wl float64, xmu, I0, I1 *[]float64) {
	p, err := plot.New()
	if err != nil {
		panic(err)
	}

	p.Title.Text = "Polarization plot"
	p.X.Label.Text = "Theta, deg"
	p.Y.Label.Text = "Polarization, %"

	p.X.Min = 0.0
	p.X.Max = 180.0
	p.Y.Max = 100.0

	p.Legend.Top = true
	p.Add(plotter.NewGrid())

	lineData := polarizationPoints(xmu, I0, I1)
	line, err := plotter.NewLine(lineData)
	if err != nil {
		panic(err)
	}
	line.LineStyle.Width = vg.Points(0.5)

	p.Add(line)
	p.Legend.Add(fmt.Sprintf("wl=%7.3f", wl), line)

	// Save the plot to a PNG file.
	if err := p.Save(8*vg.Centimeter, 6*vg.Centimeter, "polarization.pdf"); err != nil {
		panic(err)
	}
}

func polarizationPoints(xmu, I0, I1 *[]float64) plotter.XYs {
	N := len(*I1)
	pts := make(plotter.XYs, N)

	for i := range pts {
		pts[i].X = math.Acos((*xmu)[i]) * 180.0 / math.Pi
		pts[i].Y = -(*I1)[i] / (*I0)[i] * 100.0
	}
	return pts
}
