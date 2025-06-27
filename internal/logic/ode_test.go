package logic

import (
	"math"
	"testing"
)

func TestRungeKutta(t *testing.T) {
	tests := []struct {
		name    string
		xEnd    float64
		x0      float64
		y0      float64
		eps     float64
		want    float64
		wantErr bool
	}{
		// Основные базовые случаи
		{
			name:    "basic case: x from 0 to 1, y0=1",
			xEnd:    1.0,
			x0:      0.0,
			y0:      1.0,
			eps:     1e-6,
			want:    math.Exp(1.5), // y(1) = e^(1.5 * 1^2) ≈ 4.481689
			wantErr: false,
		},
		{
			name:    "basic case: x from 0 to 0.5, y0=1",
			xEnd:    0.5,
			x0:      0.0,
			y0:      1.0,
			eps:     1e-6,
			want:    math.Exp(1.5 * 0.25), // y(0.5) = e^(1.5 * 0.5^2) ≈ e^0.375 ≈ 1.454991
			wantErr: false,
		},
		{
			name:    "basic case: x from 0 to 0.1, y0=2",
			xEnd:    0.1,
			x0:      0.0,
			y0:      2.0,
			eps:     1e-6,
			want:    2 * math.Exp(1.5*0.01), // y(0.1) = 2 * e^(1.5 * 0.1^2) ≈ 2 * e^0.015 ≈ 2.030151
			wantErr: false,
		},

		// Граничные и ошибочные случаи
		{
			name:    "zero step case",
			xEnd:    0.0,
			x0:      0.0,
			y0:      1.0,
			eps:     1e-6,
			want:    1.0,
			wantErr: false,
		},
		{
			name:    "invalid epsilon",
			xEnd:    1.0,
			x0:      0.0,
			y0:      1.0,
			eps:     -1e-6,
			want:    0.0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := RungeKutta(tt.xEnd, tt.x0, tt.y0, tt.eps)
			if (err != nil) != tt.wantErr {
				t.Errorf("RungeKutta() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && math.Abs(got-tt.want) > tt.eps {
				t.Errorf("RungeKutta() = %v, want %v within epsilon %v", got, tt.want, tt.eps)
			}
		})
	}
}
