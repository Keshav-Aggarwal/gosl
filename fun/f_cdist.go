// Copyright 2015 Dorival Pedroso and Raul Durand. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fun

import (
	"math"

	"github.com/cpmech/gosl/chk"
)

// Cdist implements the distance from point to a circle (2D) or a sphere (3D)
// where the circle/sphere is implicitly defined by means of F(x) = 0.
//  where
//    F(x) = sqrt((x-xc) dot (x-xc)) - r
//  with r being the radius and xc the coordinates of the centre.
//  Thus F > 0 is outside and F < 0 is inside the circle/sphere
type Cdist struct {
	xc []float64 // centre; len(xc) = 2 or 3 => must enter "xc", "yc" and "zc" as parameters
	r  float64   // radius
}

// set allocators database
func init() {
	allocators["cdist"] = func() Func { return new(Cdist) }
}

// Init initialises the function
func (o *Cdist) Init(prms Prms) (err error) {
	var xc, yc, zc float64
	is3d := false
	for _, p := range prms {
		switch p.N {
		case "xc":
			xc = p.V
		case "yc":
			yc = p.V
		case "zc":
			zc = p.V
			is3d = true
		case "r":
			o.r = p.V
		default:
			return chk.Err("cdist: parameter named %q is invalid", p.N)
		}
	}
	rtol := 1e-10
	if o.r < rtol {
		return chk.Err("cdist: radius must be greater than %g", rtol)
	}
	if is3d {
		o.xc = []float64{xc, yc, zc}
	} else {
		o.xc = []float64{xc, yc}
	}
	return
}

// F returns y = F(t, x)
func (o Cdist) F(t float64, x []float64) float64 {
	f := (x[0]-o.xc[0])*(x[0]-o.xc[0]) + (x[1]-o.xc[1])*(x[1]-o.xc[1])
	if len(x) == 3 {
		f += (x[2] - o.xc[2]) * (x[2] - o.xc[2])
	}
	return math.Sqrt(f) - o.r
}

// G returns ∂y/∂t_cteX = G(t, x)
func (o Cdist) G(t float64, x []float64) float64 {
	return 0
}

// H returns ∂²y/∂t²_cteX = H(t, x)
func (o Cdist) H(t float64, x []float64) float64 {
	return 0
}

// Grad returns ∇F = ∂y/∂x = Grad(t, x)
func (o Cdist) Grad(v []float64, t float64, x []float64) {
	d := (x[0]-o.xc[0])*(x[0]-o.xc[0]) + (x[1]-o.xc[1])*(x[1]-o.xc[1])
	if len(x) == 3 {
		d += (x[2] - o.xc[2]) * (x[2] - o.xc[2])
	}
	d = math.Sqrt(d)
	v[0] = (x[0] - o.xc[0]) / d
	v[1] = (x[1] - o.xc[1]) / d
	if len(x) == 3 {
		v[2] = (x[2] - o.xc[1]) / d
	}
	return
}
