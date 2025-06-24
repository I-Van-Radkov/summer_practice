package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/I-Van-Radkov/summer_practice/internal/logic"
	"github.com/I-Van-Radkov/summer_practice/internal/models"
)

func SolveHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Получены данные")

	var input models.Input
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		log.Println("Ошибка при парсинге данных")
		http.Error(w, "invalid input", http.StatusBadRequest)
		return
	}

	yA, err := logic.RungeKutta(input.A, 0, 1, input.E)
	if err != nil {
		log.Printf("Ошибка: %v\n", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("Вычислено yA: %v\n", yA)

	zMax, fZMax, err := logic.FindMaximumParallel(input.C, input.D, input.E, yA)
	if err != nil {
		log.Printf("Ошибка: %v\n", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("Вычислено zMax: %v; fZMax: %v\n", zMax, fZMax)

	areaSimpson, err := logic.IntegrateSimpsonParallel(input.C, zMax, input.E, yA)
	if err != nil {
		log.Printf("Ошибка: %v\n", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	areaTrapezoid, err := logic.IntegrateTrapezoidParallel(input.C, zMax, input.E, yA)
	if err != nil {
		log.Printf("Ошибка: %v\n", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("Вычислено areaSimpson: %v; areaTrapezoid: %v\n", areaSimpson, areaTrapezoid)

	root, err := logic.FindRoot(input.C, input.D, input.E, yA)
	if err != nil {
		log.Printf("Ошибка: %v\n", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("Вычислено root: %v\n", root)

	response := models.Output{
		YA:        yA,
		ZMax:      zMax,
		FZMax:     fZMax,
		AreaSimp:  areaSimpson,
		AreaTrap:  areaTrapezoid,
		ZeroPoint: root,
	}

	_ = json.NewEncoder(w).Encode(response)
}
