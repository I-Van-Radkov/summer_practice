package logic

import "strconv"

// FindRoot находит корень F(z) = 0 на [a, b] методом бисекции
func FindRoot(a, b, eps, yA float64) string {
	fa := F(a, yA)
	fb := F(b, yA)

	if fa*fb >= 0 {
		//return 0, fmt.Errorf("F(a) and F(b) must have opposite signs, the root was not found")
		return "не существует"
	}

	for (b - a) > eps {
		mid := (a + b) / 2
		fmid := F(mid, yA)

		if fmid == 0 {
			return strconv.FormatFloat(mid, 'f', 6, 64)
		}

		if fa*fmid < 0 {
			b = mid
			fb = fmid
		} else {
			a = mid
			fa = fmid
		}
	}
	return strconv.FormatFloat((a+b)/2, 'f', 6, 64)
}
