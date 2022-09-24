package main

import (
	"communautowatcher/pkg/api"
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/cron", func(w http.ResponseWriter, r *http.Request) {
		cars, err := api.GetAvailableCars(api.CarQuery{StartDate: "03/10/2022 11:00", EndDate: "03/10/2022 11:30"})

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Error while trying to get cars %v", err)
			return
		}

		fmt.Fprintf(w, "Cars, %+v", cars)
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
