package controller

import (
	"fmt"
	models "github.com/getground/tech-tasks/backend/pkg/models"
	_ "github.com/go-sql-driver/mysql"
	"net/http"
)

type App struct {
	Party interface {
		DbAddTable(int64) (int64, error)
		DbCheckTableExists(int64) (int64, error)
		DbAddGuestList(models.Guests) error
		DbUpdateGuestStatus(string, string) error
		DbUpdateGuestList(models.Guests) error
		DbGetGuestInTable(int64) (string, error)
		DbGetGuestStatus(string) (string, error)
		DbGetCapacitySum() (int64, error)
		DbGetAccompanyingGuestsSum(string) (int64, int64, error)
		DbGetGuestList() ([]models.Guests, error)
		DbGetArrivedGuests() ([]models.Guests, error)
		DbGetTableCapacity(int64) (int64, error)
		DbGetTableIdOfGuest(string) (int64, error)
		DbCheckGuestExists(string) (int64, error)
		DbIsTablesEmpty() (bool, error)
		DbIsGuestsEmpty(string) (bool, error)
	}
}

func AddTable(app *App, table Table) (Table, error, int) {
	id, err := app.Party.DbAddTable(table.Capacity)
	if err != nil {
		return table, err, http.StatusInternalServerError
	}
	table.ID = id
	return table, nil, http.StatusOK
}

func GetGuestList(app *App) ([]GuestList, error, int) {
	var guestList []GuestList
	guests, err := app.Party.DbGetGuestList()
	if err != nil {
		return guestList, err, http.StatusInternalServerError
	}
	for _, guest := range guests {
		var gl GuestList
		gl.Name = guest.Name
		gl.AccompanyingGuests = guest.AccompanyingGuests
		gl.Table = guest.Table
		guestList = append(guestList, gl)
	}
	return guestList, err, http.StatusOK
}

func GetGuests(app *App) ([]ArrivedGuests, error, int) {
	var arrGuests []ArrivedGuests
	guests, err := app.Party.DbGetArrivedGuests()
	if err != nil {
		return arrGuests, err, http.StatusInternalServerError
	}
	for _, guest := range guests {
		var ag ArrivedGuests
		ag.Name = guest.Name
		ag.AccompanyingGuests = guest.AccompanyingGuests
		ag.TimeArrived = guest.TimeArrived
		arrGuests = append(arrGuests, ag)
	}

	return arrGuests, nil, http.StatusOK
}

func AddGuestList(app *App, guestList GuestList) (GuestName, error, int) {
	var guestName GuestName

	exists, err := app.Party.DbCheckTableExists(guestList.Table)
	if err != nil {
		return guestName, err, http.StatusInternalServerError
	}
	if exists == 0 {
		return guestName, fmt.Errorf("Invalid table-id"), http.StatusBadRequest
	}

	exists, err = app.Party.DbCheckGuestExists(guestList.Name)
	if err != nil {
		return guestName, err, http.StatusInternalServerError
	}
	if exists > 0 {
		return guestName, fmt.Errorf("Guest %s already added", guestList.Name), http.StatusBadRequest
	}

	capacity, err := app.Party.DbGetTableCapacity(guestList.Table)
	if err != nil {
		return guestName, err, http.StatusInternalServerError
	}

	if guestList.AccompanyingGuests+1 > capacity {
		return guestName, fmt.Errorf("Cannot allot table. Table capacity is %d", capacity), http.StatusBadRequest
	}

	gname, err := app.Party.DbGetGuestInTable(guestList.Table)
	if err != nil {
		return guestName, err, http.StatusInternalServerError
	}

	if gname != "" {
		return guestName, fmt.Errorf("Table already allotted to %s", gname), http.StatusBadRequest
	}

	var guests models.Guests
	guests.Table = guestList.Table
	guests.Name = guestList.Name
	guests.AccompanyingGuests = guestList.AccompanyingGuests
	if err = app.Party.DbAddGuestList(guests); err != nil {
		return guestName, err, http.StatusInternalServerError
	}
	guestName.Name = guestList.Name
	return guestName, nil, http.StatusOK
}

// update status of guest in db to checked-in
// update accompanying guests if capacity is there for table
// update arrived time in db
func UpdateGuestList(app *App, guestList GuestList) (GuestName, error, int) {
	var guestName GuestName
	exists, err := app.Party.DbCheckGuestExists(guestList.Name)
	if exists == 0 {
		return guestName, fmt.Errorf("Guest %s is not present in Guestlist", guestList.Name), http.StatusBadRequest
	}

	id, err := app.Party.DbGetTableIdOfGuest(guestList.Name)
	if err != nil {
		return guestName, err, http.StatusInternalServerError
	}

	capacity, err := app.Party.DbGetTableCapacity(id)
	if err != nil {
		return guestName, err, http.StatusInternalServerError
	}

	if guestList.AccompanyingGuests+1 > capacity {
		return guestName, fmt.Errorf("Cannot update number of accompanying guests. Table capacity is %d", capacity), http.StatusBadRequest
	}

	var guest models.Guests
	guest.Name = guestList.Name
	guest.AccompanyingGuests = guestList.AccompanyingGuests
	guest.Status = models.CHECKEDIN
	if err = app.Party.DbUpdateGuestList(guest); err != nil {
		return guestName, err, http.StatusInternalServerError
	}
	guestName.Name = guestList.Name
	return guestName, nil, http.StatusOK

}

// Update status of guest in db to checked-out
func DeleteGuest(app *App, name string) (error, int) {
	exists, err := app.Party.DbCheckGuestExists(name)
	if err != nil {
		return err, http.StatusInternalServerError
	}
	if exists == 0 {
		return fmt.Errorf("Guest %s is not present in Guestlist", name), http.StatusBadRequest
	}

	status, err := app.Party.DbGetGuestStatus(name)
	if err != nil {
		return err, http.StatusInternalServerError
	}

	if status == models.ALLOTTED {
		return fmt.Errorf("Request failed, guest not checked-in"), http.StatusBadRequest
	} else if status == models.CHECKEDOUT {
		return fmt.Errorf("Request failed, guest already checked-out"), http.StatusBadRequest
	}

	if err := app.Party.DbUpdateGuestStatus(name, models.CHECKEDOUT); err != nil {
		return err, http.StatusInternalServerError
	}
	return nil, http.StatusOK
}

// sum of capacity of all tables -
// (no.of guests arrived + sum of accompanying guests of arrived guests)
func GetEmptySeats(app *App) (EmptySeats, error, int) {
	var emptySeats EmptySeats
	empty, err := app.Party.DbIsTablesEmpty()
	if err != nil {
		return emptySeats, err, http.StatusInternalServerError
	}
	if empty {
		return emptySeats, nil, http.StatusOK
	}
	capacity, err := app.Party.DbGetCapacitySum()
	if err != nil {
		return emptySeats, err, http.StatusInternalServerError
	}
	emptySeats.SeatsEmpty = capacity
	empty, err = app.Party.DbIsGuestsEmpty(models.CHECKEDIN)
	if err != nil {
		return emptySeats, err, http.StatusInternalServerError
	}
	if empty {
		return emptySeats, nil, http.StatusOK
	}
	accompGuests, guestsCount, err := app.Party.DbGetAccompanyingGuestsSum(models.CHECKEDIN)
	if err != nil {
		return emptySeats, err, http.StatusInternalServerError
	}
	emptySeats.SeatsEmpty -= (accompGuests + guestsCount)
	return emptySeats, nil, http.StatusOK
}
