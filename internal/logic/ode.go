package logic

import (
	"fmt"
	"math"
)

// RungeKutta решает ОДУ y' = 3x*y методом Рунге-Кутты 4-го порядка
func RungeKutta(xEnd, x0, y0, eps float64) (float64, error) {
	if eps <= 0 {
		return 0, fmt.Errorf("epsilon must be positive")
	}

	y := y0
	x := x0
	h := 0.1 // Начальный шаг

	for x < xEnd {
		if x+h > xEnd {
			h = xEnd - x
		}

		k1 := h * fODE(x, y)
		k2 := h * fODE(x+h/2, y+k1/2)
		k3 := h * fODE(x+h/2, y+k2/2)
		k4 := h * fODE(x+h, y+k3)

		yNext := y + (k1+2*k2+2*k3+k4)/6
		errorEstimate := math.Abs(yNext - y)

		if errorEstimate < eps {
			y = yNext
			x += h
			// Увеличиваем шаг если ошибка мала
			if errorEstimate < eps/10 {
				h *= 1.5
			}
		} else {
			// Уменьшаем шаг если ошибка велика
			h /= 2
		}
	}
	return y, nil
}

func fODE(x, y float64) float64 {
	return 3 * x * y
}
