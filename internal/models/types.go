package models

type Input struct {
	A float64 `json:"a"`
	C float64 `json:"c"`
	D float64 `json:"d"`
	E float64 `json:"e"`
}

type Output struct {
	YA        float64 `json:"y_a"`
	ZMax      float64 `json:"z_max"`
	FZMax     float64 `json:"f_z_max"`
	AreaTrap  float64 `json:"area_trapezoid"`
	AreaSimp  float64 `json:"area_simpson"`
	ZeroPoint float64 `json:"zero_point"`
}
