package handlers

import (
	"encoding/json"
	"net/http"
	"sync"

	"github.com/I-Van-Radkov/summer_practice/internal/logic"
	"github.com/I-Van-Radkov/summer_practice/internal/models"
)

func SolveHandler(w http.ResponseWriter, r *http.Request) {
	var input models.Input
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "invalid input", http.StatusBadRequest)
		return
	}

	var wg sync.WaitGroup
	wg.Add(4)

	yA := logic.RungeKutta(input.A, 0, 1, input.E)

	zMax, fZMax := logic.FindMaximumParallel(input.C, input.D, input.E, yA)

	areaSimpson := logic.IntegrateSimpsonParallel(input.C, zMax, input.E, yA)
	areaTrapezoid := logic.IntegrateTrapezoidParallel(input.C, zMax, input.E, yA)

	root := logic.FindRoot(input.C, input.D, input.E, yA)

	response := models.Output{
		YA:        yA,
		ZMax:      zMax,
		FZMax:     fZMax,
		AreaSimp:  areaSimpson,
		AreaTrap:  areaTrapezoid,
		ZeroPoint: root,
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(response)
}
