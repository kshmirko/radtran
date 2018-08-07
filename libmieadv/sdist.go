package libmieadv

type AerosolDistrParams struct {
	npts   int
	r0, r1 float64
	gamma  float64
	Dens   float64
}

func (adp *AerosolDistrParams) SetRadiusRange(r0, r1 float64) {
	adp.r0 = r0
	adp.r1 = r1
}

func (adp *AerosolDistrParams) SetGamma(gamma float64) {
	adp.gamma = gamma
}
