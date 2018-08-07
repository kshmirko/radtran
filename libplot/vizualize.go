package libplot

import (
	"fmt"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	//"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
)

func VizualizePolarization(fname string, wl float64, ang, I0, I1 *[]float64) {
	p, err := plot.New()
	if err != nil {
		panic(err)
	}

	p.Title.Text = fmt.Sprintf("Polarization plot @ wl=%7.3f", wl)
	p.X.Label.Text = "Theta, deg"
	p.Y.Label.Text = "Polarization, %"

	p.X.Min = -90.0
	p.X.Max = 90.0
	p.Y.Max = 100.0

	p.Legend.Top = true
	p.Add(plotter.NewGrid())

	lineData := polarizationPoints(ang, I0, I1)
	line, err := plotter.NewScatter(lineData)
	if err != nil {
		panic(err)
	}
	//line.LineStyle.Width = vg.Points(0.5)

	p.Add(line)

	// Save the plot to a PNG file.
	if err := p.Save(8*vg.Centimeter, 6*vg.Centimeter, fname); err != nil {
		panic(err)
	}
}

func polarizationPoints(ang, I0, I1 *[]float64) plotter.XYs {
	N := len(*I1)
	pts := make(plotter.XYs, N)

	for i := range pts {
		pts[i].X = (*ang)[i]
		pts[i].Y = -(*I1)[i] / (*I0)[i] * 100.0
	}
	return pts
}

func VizualizeIntensity(fname string, wl float64, ang, I0, I1 *[]float64) {
	p, err := plot.New()
	if err != nil {
		panic(err)
	}

	p.Title.Text = fmt.Sprintf("Intensity plot @ wl=%7.3f", wl)
	p.X.Label.Text = "Theta, deg"
	p.Y.Label.Text = "Intensity, a.u"

	p.X.Min = -90.0
	p.X.Max = 90.0
	//p.Y.Max = 100.0

	p.Legend.Top = true
	p.Add(plotter.NewGrid())

	Idata := intensityPoints(ang, I0)
	Qdata := intensityPoints(ang, I1)
	Iline, err := plotter.NewScatter(Idata)
	if err != nil {
		panic(err)
	}
	Qline, err := plotter.NewScatter(Qdata)
	if err != nil {
		panic(err)
	}
	//line.LineStyle.Width = vg.Points(0.5)

	p.Add(Iline)
	p.Add(Qline)

	p.Legend.Add("I", Iline)
	p.Legend.Add("Q", Qline)
	// Save the plot to a PNG file.
	if err := p.Save(8*vg.Centimeter, 6*vg.Centimeter, fname); err != nil {
		panic(err)
	}
}

func intensityPoints(ang, I *[]float64) plotter.XYs {
	N := len(*I)
	pts := make(plotter.XYs, N)

	for i := range pts {
		pts[i].X = (*ang)[i]
		pts[i].Y = (*I)[i]
	}
	return pts
}
