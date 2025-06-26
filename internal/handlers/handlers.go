package handlers

import (
	"encoding/csv"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/I-Van-Radkov/summer_practice/internal/logic"
	"github.com/I-Van-Radkov/summer_practice/internal/models"
)

const fileName = "./data/output.csv"

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

	root := logic.FindRoot(input.C, input.D, input.E, yA)

	log.Printf("Вычислено root: %v\n", root)

	response := models.Output{
		YA:        yA,
		ZMax:      zMax,
		FZMax:     fZMax,
		AreaSimp:  areaSimpson,
		AreaTrap:  areaTrapezoid,
		ZeroPoint: root,
	}

	err = saveToCSV(response, fileName)
	if err != nil {
		log.Printf("Ошибка: %v\n", err)
		http.Error(w, "Ошибка при сохранении CSV: "+err.Error(), http.StatusInternalServerError)
		return
	}

	_ = json.NewEncoder(w).Encode(response)
}

func EnableCORS(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		if r.Method == "OPTIONS" {
			return
		}

		next.ServeHTTP(w, r)
	})
}

func DownloadHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Запрос на скачивание файла")

	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		log.Println("Файл не найден:", fileName)
		http.Error(w, "Файл не найден", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Disposition", "attachment; filename="+fileName)
	w.Header().Set("Content-Type", "text/csv")

	http.ServeFile(w, r, fileName)
}

func saveToCSV(data models.Output, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	writer.Write([]string{
		"y(a)", "z_max", "F(z_max)", "Area (Trapezoid)", "Area (Simpson)", "F(z) = 0",
	})

	writer.Write([]string{
		formatFloat(data.YA),
		formatFloat(data.ZMax),
		formatFloat(data.FZMax),
		formatFloat(data.AreaTrap),
		formatFloat(data.AreaSimp),
		data.ZeroPoint,
	})

	return nil
}

func formatFloat(v float64) string {
	return strconv.FormatFloat(v, 'f', 6, 64)
}
