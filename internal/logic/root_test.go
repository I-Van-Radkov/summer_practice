package logic

import (
	"math"
	"strconv"
	"testing"
)

func TestFindRoot(t *testing.T) {
	tests := []struct {
		name     string
		a        float64
		b        float64
		eps      float64
		yA       float64
		expected string
	}{
		{
			name:     "корень существует при yA=1",
			a:        -20,
			b:        0,
			eps:      1e-6,
			yA:       1,
			expected: "-3.4074468",
		},
		{
			name:     "корень существует при yA=2",
			a:        -10,
			b:        0,
			eps:      1e-6,
			yA:       2,
			expected: "-3.883646",
		},
		{
			name:     "корень не существует (положительные значения)",
			a:        1,
			b:        2,
			eps:      1e-6,
			yA:       1,
			expected: "не существует",
		},
		{
			name:     "корень на границе a",
			a:        -3.791288,
			b:        0,
			eps:      1e-6,
			yA:       1,
			expected: "-3.407446",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FindRoot(tt.a, tt.b, tt.eps, tt.yA)

			if tt.expected == "не существует" {
				if got != tt.expected {
					t.Errorf("FindRoot() = %v, ожидалось %v", got, tt.expected)
				}
				return
			}

			gotFloat, err := strconv.ParseFloat(got, 64)
			if err != nil {
				t.Errorf("Ошибка парсинга результата: %v", err)
				return
			}

			expectedFloat, err := strconv.ParseFloat(tt.expected, 64)
			if err != nil {
				t.Errorf("Ошибка парсинга эталонного значения: %v", err)
				return
			}

			fValue := F(gotFloat, tt.yA)
			if math.Abs(fValue) > tt.eps*10 {
				t.Errorf("F(найденный корень) = %v (должно быть близко к 0)", fValue)
			}

			if math.Abs(gotFloat-expectedFloat) > tt.eps {
				t.Errorf("FindRoot() = %v, ожидалось %v с точностью %v", got, tt.expected, tt.eps)
			}
		})
	}
}
