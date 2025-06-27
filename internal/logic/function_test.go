package logic

import (
	"math"
	"testing"
)

func TestF(t *testing.T) {
	tests := []struct {
		name string
		z    float64
		yA   float64
		want float64
	}{
		{
			name: "simple case 1",
			z:    0,
			yA:   1,
			want: 15,
		},
		{
			name: "simple case 2",
			z:    math.Pi / 2,
			yA:   2,
			want: 15 + 1 + 2*(math.Pi/2) - math.Pow(math.Pi/2, 2)/2,
		},
		{
			name: "negative z",
			z:    -1,
			yA:   1,
			want: 15 + math.Pow(math.Sin(-1), 3) + 1*(-1) - math.Pow(-1, 2)/1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := F(tt.z, tt.yA)
			if math.Abs(got-tt.want) > 1e-9 {
				t.Errorf("F() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFindMaximumParallel(t *testing.T) {
	tests := []struct {
		name    string
		c       float64
		d       float64
		eps     float64
		yA      float64
		wantZ   float64
		wantVal float64
		wantErr bool
	}{
		{
			name:    "simple maximum search",
			c:       0,
			d:       math.Pi,
			eps:     1e-6,
			yA:      1,
			wantZ:   math.Pi / 2, // Ожидаемый максимум sin³(z) в [0, π]
			wantErr: false,
		},
		{
			name:    "invalid interval",
			c:       1,
			d:       0,
			eps:     1e-6,
			yA:      1,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			z, val, err := FindMaximumParallel(tt.c, tt.d, tt.eps, tt.yA)
			if (err != nil) != tt.wantErr {
				t.Errorf("FindMaximumParallel() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				// Проверяем, что найденное значение действительно максимум
				testVal := F(z, tt.yA)
				if math.Abs(testVal-val) > tt.eps {
					t.Errorf("FindMaximumParallel() returned value %v doesn't match F(%v) = %v", val, z, testVal)
				}
				// Проверяем, что значение в соседних точках не больше
				left := F(z-tt.eps*10, tt.yA)
				right := F(z+tt.eps*10, tt.yA)
				if left > val || right > val {
					t.Errorf("FindMaximumParallel() found value %v at z=%v, but nearby values are larger: left=%v, right=%v", val, z, left, right)
				}
			}
		})
	}
}

func TestIntegrateSimpsonParallel(t *testing.T) {
	tests := []struct {
		name    string
		a       float64
		b       float64
		eps     float64
		yA      float64
		want    float64
		wantErr bool
	}{
		{
			name:    "integrate constant",
			a:       0,
			b:       1,
			eps:     1e-6,
			yA:      0, // F(z) = 15 + sin³(z)
			want:    15,
			wantErr: false,
		},
		{
			name:    "invalid interval",
			a:       1,
			b:       0,
			eps:     1e-6,
			yA:      1,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := IntegrateSimpsonParallel(tt.a, tt.b, tt.eps, tt.yA)
			if (err != nil) != tt.wantErr {
				t.Errorf("IntegrateSimpsonParallel() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && math.Abs(got-tt.want) > tt.eps {
				t.Errorf("IntegrateSimpsonParallel() = %v, want %v within epsilon %v", got, tt.want, tt.eps)
			}
		})
	}
}

func TestIntegrateTrapezoidParallel(t *testing.T) {
	tests := []struct {
		name    string
		a       float64
		b       float64
		eps     float64
		yA      float64
		want    float64
		wantErr bool
	}{
		{
			name:    "integrate constant",
			a:       0,
			b:       1,
			eps:     1e-6,
			yA:      0, // F(z) = 15 + sin³(z)
			want:    15,
			wantErr: false,
		},
		{
			name:    "invalid interval",
			a:       1,
			b:       0,
			eps:     1e-6,
			yA:      1,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := IntegrateTrapezoidParallel(tt.a, tt.b, tt.eps, tt.yA)
			if (err != nil) != tt.wantErr {
				t.Errorf("IntegrateTrapezoidParallel() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && math.Abs(got-tt.want) > tt.eps {
				t.Errorf("IntegrateTrapezoidParallel() = %v, want %v within epsilon %v", got, tt.want, tt.eps)
			}
		})
	}
}

func BenchmarkFindMaximumParallel(b *testing.B) {
	for i := 0; i < b.N; i++ {
		FindMaximumParallel(0, math.Pi, 1e-6, 1)
	}
}

func BenchmarkIntegrateSimpsonParallel(b *testing.B) {
	for i := 0; i < b.N; i++ {
		IntegrateSimpsonParallel(0, math.Pi, 1e-6, 1)
	}
}

func BenchmarkIntegrateTrapezoidParallel(b *testing.B) {
	for i := 0; i < b.N; i++ {
		IntegrateTrapezoidParallel(0, math.Pi, 1e-6, 1)
	}
}
