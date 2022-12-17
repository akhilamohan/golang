package main

import (
	"fmt"
	"log"
	"net/http"

	controller "github.com/getground/tech-tasks/backend/pkg/controller"
	db "github.com/getground/tech-tasks/backend/pkg/db"
	models "github.com/getground/tech-tasks/backend/pkg/models"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

func main() {
	// init mysql.
	sqlDB, err := db.ConnectToDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.CloseConnection(sqlDB)
	app := &controller.App{
		Party: models.PartyModel{DB: sqlDB},
	}
	router := mux.NewRouter()
	router.HandleFunc("/tables", app.AddTableHandler).Methods("POST")
	router.HandleFunc("/guest_list/{name}", app.AddGuestListHandler).Methods("POST")
	router.HandleFunc("/guest_list", app.GetGuestListHandler).Methods("GET")
	router.HandleFunc("/guests", app.GetGuestsHandler).Methods("GET")
	router.HandleFunc("/guests/{name}", app.UpdateGuestHandler).Methods("PUT")
	router.HandleFunc("/guests/{name}", app.DeleteGuestHandler).Methods("DELETE")
	router.HandleFunc("/seats_empty", app.GetEmptySeatsHandler).Methods("GET")
	router.HandleFunc("/ping", handlerPing).Methods("GET")
	if err = http.ListenAndServe(":3000", router); err != nil {
		log.Fatal(err)
	}
}

func handlerPing(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "pong\n")
}
