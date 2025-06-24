package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/I-Van-Radkov/summer_practice/internal/logic"
	"github.com/I-Van-Radkov/summer_practice/internal/models"
)

func SolveHandler(w http.ResponseWriter, r *http.Request) {
	var input models.Input
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "invalid input", http.StatusBadRequest)
		return
	}

	yA, err := logic.RungeKutta(input.A, 0, 1, input.E)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	zMax, fZMax, err := logic.FindMaximumParallel(input.C, input.D, input.E, yA)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	areaSimpson, err := logic.IntegrateSimpsonParallel(input.C, zMax, input.E, yA)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	areaTrapezoid, err := logic.IntegrateTrapezoidParallel(input.C, zMax, input.E, yA)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	root, err := logic.FindRoot(input.C, input.D, input.E, yA)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

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
