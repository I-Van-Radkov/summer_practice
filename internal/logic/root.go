package logic

func FindRoot(a, b, eps, yA float64) float64 {
	for (b - a) > eps {
		mid := (a + b) / 2
		if F(mid, yA)*F(a, yA) < 0 {
			b = mid
		} else {
			a = mid
		}
	}
	return (a + b) / 2
}
