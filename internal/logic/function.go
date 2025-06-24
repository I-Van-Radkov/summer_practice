package logic

import (
	"fmt"
	"math"
	"sync"
)

// F вычисляет значение функции F(z) = 15 + sin³(z) + yA*z - z²/yA
func F(z, yA float64) float64 {
	return 15 + math.Pow(math.Sin(z), 3) + yA*z - math.Pow(z, 2)/yA
}

// FindMaximumParallel находит максимум F(z) на [c, d] с точностью eps
func FindMaximumParallel(c, d, eps, yA float64) (float64, float64, error) {
	if c >= d {
		return 0, 0, fmt.Errorf("c must be less than d")
	}

	const numWorkers = 4
	step := (d - c) / numWorkers
	results := make(chan struct {
		z   float64
		val float64
	}, numWorkers)

	var wg sync.WaitGroup
	wg.Add(numWorkers)

	for i := 0; i < numWorkers; i++ {
		go func(start, end float64) {
			defer wg.Done()
			zMax := start

			// Используем золотое сечение для эффективного поиска
			phi := (1 + math.Sqrt(5)) / 2
			a, b := start, end
			for (b - a) > eps {
				x1 := b - (b-a)/phi
				x2 := a + (b-a)/phi
				if F(x1, yA) > F(x2, yA) {
					b = x2
				} else {
					a = x1
				}
			}
			zMax = (a + b) / 2
			fMax := F(zMax, yA)

			results <- struct {
				z   float64
				val float64
			}{z: zMax, val: fMax}
		}(c+float64(i)*step, c+float64(i+1)*step)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	globalMax := -math.MaxFloat64
	var globalZ float64
	for res := range results {
		if res.val > globalMax {
			globalMax = res.val
			globalZ = res.z
		}
	}

	return globalZ, globalMax, nil
}

// IntegrateSimpsonParallel вычисляет интеграл методом Симпсона
func IntegrateSimpsonParallel(a, b, eps, yA float64) (float64, error) {
	if a >= b {
		return 0, fmt.Errorf("a must be less than b")
	}

	f := func(x float64) float64 { return F(x, yA) }
	return adaptiveSimpson(f, a, b, eps), nil
}

// IntegrateTrapezoidParallel вычисляет интеграл методом трапеций
func IntegrateTrapezoidParallel(a, b, eps, yA float64) (float64, error) {
	if a >= b {
		return 0, fmt.Errorf("a must be less than b")
	}

	f := func(x float64) float64 { return F(x, yA) }
	return adaptiveTrapezoid(f, a, b, eps), nil
}

func adaptiveSimpson(f func(float64) float64, a, b, eps float64) float64 {
	const maxDepth = 20
	var worker func(a, b, eps float64, depth int) float64
	worker = func(a, b, eps float64, depth int) float64 {
		if depth > maxDepth {
			return (f(a) + f(b)) * (b - a) / 2
		}

		c := (a + b) / 2
		h := b - a
		fa := f(a)
		fb := f(b)
		fc := f(c)

		S := (fa + 4*fc + fb) * h / 6
		d := (a + c) / 2
		e := (c + b) / 2
		fd := f(d)
		fe := f(e)
		S2 := (fa + 4*fd + 2*fc + 4*fe + fb) * h / 12

		if math.Abs(S2-S) < 15*eps {
			return S2 + (S2-S)/15
		}
		return worker(a, c, eps/2, depth+1) + worker(c, b, eps/2, depth+1)
	}

	return worker(a, b, eps, 0)
}

func adaptiveTrapezoid(f func(float64) float64, a, b, eps float64) float64 {
	const maxDepth = 20
	var worker func(a, b, eps float64, depth int) float64
	worker = func(a, b, eps float64, depth int) float64 {
		if depth > maxDepth {
			return (f(a) + f(b)) * (b - a) / 2
		}

		h := b - a
		fa := f(a)
		fb := f(b)
		S := (fa + fb) * h / 2

		c := (a + b) / 2
		fc := f(c)
		S2 := (fa + 2*fc + fb) * h / 4

		if math.Abs(S2-S) < 3*eps {
			return S2 + (S2-S)/3
		}
		return worker(a, c, eps/2, depth+1) + worker(c, b, eps/2, depth+1)
	}

	return worker(a, b, eps, 0)
}
