package logic

import (
	"math"
	"sync"
)

func F(z, yA float64) float64 {
	return 15 + math.Pow(math.Sin(z), 3) + yA*z - z*z/yA
}

func FindMaximumParallel(c, d, e, yA float64) (float64, float64) {
	const numWorkers = 4
	step := (d - c) / float64(numWorkers)
	type result struct {
		z   float64
		val float64
	}

	var wg sync.WaitGroup
	ch := make(chan result, numWorkers)

	for i := 0; i < numWorkers; i++ {
		a := c + float64(i)*step
		b := a + step
		wg.Add(1)
		go func(start, end float64) {
			defer wg.Done()
			zLocal := start
			fLocal := F(zLocal, yA)
			for z := start; z <= end; z += e {
				fz := F(z, yA)
				if fz > fLocal {
					zLocal = z
					fLocal = fz
				}
			}
			ch <- result{z: zLocal, val: fLocal}
		}(a, b)
	}

	wg.Wait()
	close(ch)

	zMax, fMax := c, F(c, yA)
	for r := range ch {
		if r.val > fMax {
			zMax = r.z
			fMax = r.val
		}
	}

	return zMax, fMax
}

func adaptiveSimpson(f func(float64) float64, a, b, eps float64) float64 {
	var wg sync.WaitGroup
	resultChan := make(chan float64, 4)

	adaptive := func(a, b, eps float64) float64 {
		c := (a + b) / 2
		fa, fb, fc := f(a), f(b), f(c)
		S := (fa + 4*fc + fb) * (b - a) / 6
		return adaptiveSimpsonStep(f, a, b, eps, S, fa, fb, fc)
	}

	step := (b - a) / 4
	for i := 0; i < 4; i++ {
		subA := a + float64(i)*step
		subB := subA + step
		wg.Add(1)
		go func(start, end float64) {
			defer wg.Done()
			resultChan <- adaptive(start, end, eps/4)
		}(subA, subB)
	}

	wg.Wait()
	close(resultChan)
	total := 0.0
	for r := range resultChan {
		total += r
	}
	return total
}

func adaptiveSimpsonStep(f func(float64) float64, a, b, eps, S, fa, fb, fc float64) float64 {
	c := (a + b) / 2
	fd := f((a + c) / 2)
	fe := f((c + b) / 2)
	S1 := (fa + 4*fd + fc) * (c - a) / 6
	S2 := (fc + 4*fe + fb) * (b - c) / 6
	if math.Abs(S1+S2-S) <= 15*eps {
		return S1 + S2 + (S1+S2-S)/15
	}
	return adaptiveSimpsonStep(f, a, c, eps/2, S1, fa, fc, fd) + adaptiveSimpsonStep(f, c, b, eps/2, S2, fc, fb, fe)
}

func adaptiveTrapezoid(f func(float64) float64, a, b, eps float64) float64 {
	var wg sync.WaitGroup
	resultChan := make(chan float64, 4)

	adaptive := func(a, b, eps float64) float64 {
		fa := f(a)
		fb := f(b)
		S := (fa + fb) * (b - a) / 2
		return adaptiveTrapezoidStep(f, a, b, eps, S, fa, fb)
	}

	step := (b - a) / 4
	for i := 0; i < 4; i++ {
		subA := a + float64(i)*step
		subB := subA + step
		wg.Add(1)
		go func(start, end float64) {
			defer wg.Done()
			resultChan <- adaptive(start, end, eps/4)
		}(subA, subB)
	}

	wg.Wait()
	close(resultChan)
	total := 0.0
	for r := range resultChan {
		total += r
	}
	return total
}

func adaptiveTrapezoidStep(f func(float64) float64, a, b, eps, S, fa, fb float64) float64 {
	c := (a + b) / 2
	fc := f(c)
	S1 := (fa + fc) * (c - a) / 2
	S2 := (fc + fb) * (b - c) / 2
	if math.Abs(S1+S2-S) <= 3*eps {
		return S1 + S2 + (S1+S2-S)/3
	}
	return adaptiveTrapezoidStep(f, a, c, eps/2, S1, fa, fc) + adaptiveTrapezoidStep(f, c, b, eps/2, S2, fc, fb)
}

func IntegrateTrapezoidParallel(a, b, eps, yA float64) float64 {
	return adaptiveTrapezoid(func(x float64) float64 {
		return F(x, yA)
	}, a, b, eps)
}

func IntegrateSimpsonParallel(a, b, eps, yA float64) float64 {
	return adaptiveSimpson(func(x float64) float64 {
		return F(x, yA)
	}, a, b, eps)
}
