package logic

import "fmt"

// FindRoot находит корень F(z) = 0 на [a, b] методом бисекции
func FindRoot(a, b, eps, yA float64) (float64, error) {
	fa := F(a, yA)
	fb := F(b, yA)

	if fa*fb >= 0 {
		return 0, fmt.Errorf("F(a) and F(b) must have opposite signs")
	}

	for (b - a) > eps {
		mid := (a + b) / 2
		fmid := F(mid, yA)

		if fmid == 0 {
			return mid, nil
		}

		if fa*fmid < 0 {
			b = mid
			fb = fmid
		} else {
			a = mid
			fa = fmid
		}
	}
	return (a + b) / 2, nil
}
