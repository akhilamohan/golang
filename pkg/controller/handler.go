package controller

import (
	"encoding/json"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"io/ioutil"
	"net/http"
	"strings"
)

func sendErrorResponse(w http.ResponseWriter, err error, responseCode int) {
	fmt.Println(err.Error())
	w.WriteHeader(responseCode)
	fmt.Fprintf(w, "%s", err.Error())
}

// http handler to add a table
func (app *App) AddTableHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		sendErrorResponse(w, err, http.StatusInternalServerError)
		return
	}
	var table Table
	err = json.Unmarshal(body, &table)
	if err != nil {
		sendErrorResponse(w, err, http.StatusInternalServerError)
		return
	}
	errs, err := validateCapacity(table.Capacity)
	if err != nil {
		sendErrorResponse(w, err, http.StatusInternalServerError)
		return
	}
	if errs != "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "%s", errs)
		return
	}
	table, err, respCode := AddTable(app, table)
	if err != nil {
		sendErrorResponse(w, err, respCode)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(table)
}

// http handler to get guest list
func (app *App) GetGuestListHandler(w http.ResponseWriter, r *http.Request) {
	guestList, err, respCode := GetGuestList(app)
	if err != nil {
		sendErrorResponse(w, err, respCode)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(guestList)
}

// http handler to get arrived guests
func (app *App) GetGuestsHandler(w http.ResponseWriter, r *http.Request) {
	guests, err, respCode := GetGuests(app)
	if err != nil {
		sendErrorResponse(w, err, respCode)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(guests)
}

// http hander to allot a table to guest
func (app *App) AddGuestListHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		sendErrorResponse(w, err, http.StatusInternalServerError)
		return
	}
	params := mux.Vars(r)
	name := strings.ToLower(params["name"])
	errs, err := validateName(name)
	if err != nil {
		sendErrorResponse(w, err, http.StatusInternalServerError)
		return
	}
	if errs != "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "%s", errs)
		return
	}
	var guestList GuestList
	err = json.Unmarshal(body, &guestList)
	if err != nil {
		sendErrorResponse(w, err, http.StatusInternalServerError)
		return
	}
	errs, err = validateGuestList(guestList)
	if err != nil {
		sendErrorResponse(w, err, http.StatusInternalServerError)
		return
	}
	if errs != "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "%s", errs)
		return
	}
	guestList.Name = name
	guestName, err, respCode := AddGuestList(app, guestList)
	if err != nil {
		sendErrorResponse(w, err, respCode)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(guestName)
}

// http handler to check-in a guest
func (app *App) UpdateGuestHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		sendErrorResponse(w, err, http.StatusInternalServerError)
		return
	}
	params := mux.Vars(r)
	name := strings.ToLower(params["name"])
	errs, err := validateName(name)
	if err != nil {
		sendErrorResponse(w, err, http.StatusInternalServerError)
		return
	}
	if errs != "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "%s", errs)
		return
	}
	var guestList GuestList
	err = json.Unmarshal(body, &guestList)
	if err != nil {
		sendErrorResponse(w, err, http.StatusInternalServerError)
		return
	}
	errs, err = validateAccompanyingGuests(guestList.AccompanyingGuests)
	if err != nil {
		sendErrorResponse(w, err, http.StatusInternalServerError)
		return
	}
	if errs != "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "%s", errs)
		return
	}
	guestList.Name = name
	guestName, err, respCode := UpdateGuestList(app, guestList)
	if err != nil {
		sendErrorResponse(w, err, respCode)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(guestName)
}

// http handler to check-out a guest
func (app *App) DeleteGuestHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	name := strings.ToLower(params["name"])
	err, respCode := DeleteGuest(app, name)
	if err != nil {
		sendErrorResponse(w, err, respCode)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// http handler to get empty seats
func (app *App) GetEmptySeatsHandler(w http.ResponseWriter, r *http.Request) {
	emptySeats, err, respCode := GetEmptySeats(app)
	if err != nil {
		sendErrorResponse(w, err, respCode)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(emptySeats)
}
