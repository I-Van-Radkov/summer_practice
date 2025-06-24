package logic

import "math"

func RungeKutta(xEnd, x0, y0 float64, eps float64) float64 {
	y := y0
	x := x0
	h := 0.01

	for x < xEnd {
		if x+h > xEnd {
			h = xEnd - x
		}

		k1 := h * fODE(x, y)
		k2 := h * fODE(x+h/2, y+k1/2)
		k3 := h * fODE(x+h/2, y+k2/2)
		k4 := h * fODE(x+h, y+k3)

		yNext := y + (k1+2*k2+2*k3+k4)/6
		if math.Abs(yNext-y) < eps {
			y = yNext
			x += h
		} else {
			h /= 2
		}
	}
	return y
}

func fODE(x, y float64) float64 {
	return 3 * x * y
}
