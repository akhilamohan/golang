package models

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

const (
	CHECKEDIN  = "checked-in"
	CHECKEDOUT = "checked-out"
	ALLOTTED   = "allotted"
)

type Table struct {
	Id       int64
	Capacity int64
}

type Guests struct {
	Table              int64
	AccompanyingGuests int64
	Status             string
	TimeArrived        time.Time
	Name               string
}

type PartyModel struct {
	DB *sql.DB
}

func (p PartyModel) DbAddTable(capacity int64) (int64, error) {
	var resId int64
	res, err := p.DB.Exec(
		`INSERT INTO tables(capacity) VALUES (?)`,
		capacity)
	if err != nil {
		fmt.Println(err)
		return resId, err
	}
	resId, err = res.LastInsertId()
	if err != nil {
		return resId, err
	}

	return resId, nil
}

func (p PartyModel) DbGetTableIdOfGuest(name string) (int64, error) {
	var id int64
	res := p.DB.QueryRow("SELECT id FROM guests WHERE name = ?",
		name)
	err := res.Scan(&id)
	if err != nil {
		fmt.Println(err)
		return id, err
	}
	return id, err

}

func (p PartyModel) DbAddGuestList(guest Guests) error {
	_, err := p.DB.Exec(
		`INSERT INTO guests(id, accompanying_guests, name) VALUES (?, ?, ?)`,
		guest.Table, guest.AccompanyingGuests, guest.Name)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func (p PartyModel) DbUpdateGuestStatus(name string, status string) error {
	_, err := p.DB.Exec(
		`UPDATE guests SET 
		status = ?
		WHERE name = ?`,
		status, name)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func (p PartyModel) DbUpdateGuestList(guest Guests) error {
	_, err := p.DB.Exec(
		`UPDATE guests SET 
		status = ?,
		accompanying_guests = ?,
		time_arrived = ?
		WHERE name = ?`,
		guest.Status, guest.AccompanyingGuests,
		time.Now(), guest.Name)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func (p PartyModel) DbGetGuestInTable(id int64) (string, error) {
	var name string
	res := p.DB.QueryRow("SELECT name FROM guests WHERE id = ?",
		id)
	err := res.Scan(&name)
	if err != nil && err != sql.ErrNoRows {
		fmt.Println(err)
		return name, err
	}
	return name, nil
}

func (p PartyModel) DbGetGuestStatus(name string) (string, error) {
	var status string
	res := p.DB.QueryRow("SELECT status FROM guests WHERE name = ?",
		name)
	err := res.Scan(&status)
	if err != nil {
		fmt.Println(err)
		return status, err
	}
	return status, nil
}

func (p PartyModel) DbGetCapacitySum() (int64, error) {
	var totalCapacity int64
	res := p.DB.QueryRow("SELECT SUM(capacity) FROM tables")
	err := res.Scan(&totalCapacity)
	if err != nil && err != sql.ErrNoRows {
		fmt.Println(err)
		return totalCapacity, err
	}
	return totalCapacity, nil
}

func (p PartyModel) DbGetAccompanyingGuestsSum(status string) (int64, int64, error) {
	var accompanyingGuests int64
	var guestsCount int64
	res := p.DB.QueryRow("SELECT SUM(accompanying_guests), COUNT(name) FROM guests WHERE status = ?",
		status)
	err := res.Scan(&accompanyingGuests, &guestsCount)
	if err != nil && err != sql.ErrNoRows {
		fmt.Println(err)
		return accompanyingGuests, guestsCount, err
	}
	return accompanyingGuests, guestsCount, nil

}

func (p PartyModel) DbGetGuestList() ([]Guests, error) {
	var guestList []Guests
	res, err := p.DB.Query("SELECT id, name, accompanying_guests FROM guests")
	if err != nil {
		fmt.Println(err)
		return guestList, err
	}
	defer res.Close()
	for res.Next() {
		var gl Guests
		err := res.Scan(&gl.Table, &gl.Name, &gl.AccompanyingGuests)
		if err != nil {
			fmt.Println(err)
			return guestList, err
		}
		guestList = append(guestList, gl)
	}
	return guestList, nil
}

func (p PartyModel) DbGetArrivedGuests() ([]Guests, error) {
	var arrivedGuests []Guests
	res, err := p.DB.Query(`SELECT name, accompanying_guests, time_arrived 
				   FROM guests 
				   WHERE (status = ?) OR (status = ?)`,
		CHECKEDIN, CHECKEDOUT)
	if err != nil {
		fmt.Println(err)
		return arrivedGuests, err
	}
	defer res.Close()
	for res.Next() {
		var ag Guests
		err := res.Scan(&ag.Name, &ag.AccompanyingGuests, &ag.TimeArrived)
		if err != nil {
			fmt.Println(err)
			return arrivedGuests, err
		}
		arrivedGuests = append(arrivedGuests, ag)
	}
	return arrivedGuests, nil
}

func (p PartyModel) DbGetTableCapacity(id int64) (int64, error) {
	var capacity int64
	res := p.DB.QueryRow("SELECT capacity FROM tables WHERE id = ?", id)
	err := res.Scan(&capacity)
	if err != nil {
		fmt.Println(err)
		return capacity, err
	}
	return capacity, nil
}

func (p PartyModel) DbCheckTableExists(id int64) (int64, error) {
	var exists int64
	res := p.DB.QueryRow("SELECT EXISTS(SELECT * FROM tables WHERE id = ?)", id)
	err := res.Scan(&exists)
	if err != nil && err != sql.ErrNoRows {
		fmt.Println(err)
		return exists, err
	}
	return exists, nil
}

func (p PartyModel) DbCheckGuestExists(name string) (int64, error) {
	var exists int64
	res := p.DB.QueryRow("SELECT EXISTS(SELECT * FROM guests WHERE name = ?)", name)
	err := res.Scan(&exists)
	if err != nil && err != sql.ErrNoRows {
		fmt.Println(err)
		return exists, err
	}
	return exists, nil
}

func (p PartyModel) DbIsTablesEmpty() (bool, error) {
	var count int64
	res := p.DB.QueryRow("SELECT COUNT(*) from tables")
	err := res.Scan(&count)
	if err != nil {
		fmt.Println(err)
		return true, err
	}
	if count > 0 {
		return false, nil
	}
	return true, nil
}

func (p PartyModel) DbIsGuestsEmpty(status string) (bool, error) {
	var count int64
	res := p.DB.QueryRow("SELECT COUNT(*) from guests WHERE status = ?", status)
	err := res.Scan(&count)
	if err != nil {
		fmt.Println(err)
		return true, err
	}
	if count > 0 {
		return false, nil
	}
	return true, nil
}
