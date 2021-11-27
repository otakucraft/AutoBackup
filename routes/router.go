package routes

import (
	"backup/routes/endpoints"
	"backup/utils"
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func StartRouter() {
	port := "8462"
	log.Println("Starting a server on http://127.0.0.1:" + port)
	router := mux.NewRouter().StrictSlash(true)

	router.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)

		res := utils.Result{Code: http.StatusNotFound, Data: "Method not found!"}
		response, _ := json.Marshal(res)
		_, err := w.Write(response)
		if err != nil {
			return
		}
	})

	router.MethodNotAllowedHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)

		res := utils.Result{Code: http.StatusMethodNotAllowed, Data: "Method not allowed!"}
		response, _ := json.Marshal(res)
		_, err := w.Write(response)
		if err != nil {
			return
		}
	})

	router.HandleFunc("/healthcheck", func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(http.StatusOK)
	}).Methods("GET")

	router.HandleFunc("/reload", endpoints.ReloadConfigFile).Methods("GET")
	router.HandleFunc("/create/{name}", endpoints.CreateBackup).Methods("GET")

	log.Fatal(http.ListenAndServe(":" + port, router))
}
