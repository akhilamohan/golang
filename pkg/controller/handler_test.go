package controller

import (
	models "github.com/getground/tech-tasks/backend/pkg/models"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

type mockPartyModel struct{}

var guests []models.Guests
var tables []models.Table

func addCheckedOutGuest() {
	guests = append(guests, models.Guests{3, 3, "checked-out", time.Time{}, "jack"})
}

func addGuests() {
	guests = append(guests, models.Guests{1, 1, "allotted", time.Time{}, "john"})
	guests = append(guests, models.Guests{2, 2, "checked-in", time.Time{}, "akhila"})
}

func addSingleGuest() {
	guests = append(guests, models.Guests{1, 1, "allotted", time.Time{}, "john"})
}

func addGuest(id int64, accGuest int64,
	status string, arrTime time.Time, name string) {
	guests = append(guests, models.Guests{id, accGuest,
		status, arrTime, name})
}

func addTables(count int64) {
	var start int64 = 1
	for i := start; i <= count; i++ {
		tables = append(tables, models.Table{i, i + 2})
	}
}

func cleanup() {
	guests = nil
	tables = nil
}

func (*mockPartyModel) DbGetTableIdOfGuest(name string) (int64, error) {
	var id int64
	for _, guest := range guests {
		if guest.Name == name {
			id = guest.Table
			break
		}
	}
	return id, nil
}

func (*mockPartyModel) DbGetGuestStatus(name string) (string, error) {
	var status string
	for _, guest := range guests {
		if guest.Name == name {
			status = guest.Status
			break
		}
	}
	return status, nil
}

func (*mockPartyModel) DbAddTable(int64) (int64, error) {
	var count int64 = 1
	var id int64 = 1
	addTables(count)
	return id, nil
}

func (*mockPartyModel) DbCheckTableExists(id int64) (int64, error) {
	var exists int64
	for _, table := range tables {
		if table.Id == id {
			exists = 1
			break
		}
	}
	return exists, nil
}

func (*mockPartyModel) DbGetGuestInTable(id int64) (string, error) {
	var name string
	for _, guest := range guests {
		if guest.Table == id {
			name = guest.Name
			break
		}
	}
	return name, nil
}

func (*mockPartyModel) DbAddGuestList(guest models.Guests) error {
	guests = append(guests, models.Guests{guest.Table, guest.AccompanyingGuests,
		"allotted", time.Time{}, guest.Name})
	return nil
}

func (*mockPartyModel) DbUpdateGuestStatus(name string, status string) error {
	for i, guest := range guests {
		if guest.Name == name {
			guests[i].Status = models.CHECKEDOUT
			break
		}
	}
	return nil
}

func (*mockPartyModel) DbUpdateGuestList(guest models.Guests) error {
	for i, g := range guests {
		if g.Name == guest.Name {
			guests[i].Status = guest.Status
			guests[i].AccompanyingGuests = guest.AccompanyingGuests
			guests[i].TimeArrived = guest.TimeArrived
			break
		}
	}
	return nil
}

func (*mockPartyModel) DbGetCapacitySum() (int64, error) {
	var sum int64
	for _, table := range tables {
		sum += table.Capacity
	}
	return sum, nil
}

func (*mockPartyModel) DbGetAccompanyingGuestsSum(status string) (int64, int64, error) {
	var sum, count int64
	for _, g := range guests {
		if g.Status == status {
			sum += g.AccompanyingGuests
			count++
		}
	}
	return sum, count, nil
}

func (*mockPartyModel) DbGetGuestList() ([]models.Guests, error) {
	return guests, nil
}

func (*mockPartyModel) DbGetArrivedGuests() ([]models.Guests, error) {
	var aguests []models.Guests
	for _, g := range guests {
		if g.Status == models.CHECKEDIN {
			aguests = append(aguests, g)
		}
	}
	return aguests, nil
}

func (*mockPartyModel) DbGetTableCapacity(id int64) (int64, error) {
	var capacity int64
	for _, table := range tables {
		if table.Id == id {
			capacity = table.Capacity
			break
		}
	}
	return capacity, nil
}

func (*mockPartyModel) DbCheckGuestExists(name string) (int64, error) {
	var exists int64
	for _, guest := range guests {
		if guest.Name == name {
			exists = 1
			break
		}
	}
	return exists, nil
}

func (*mockPartyModel) DbIsTablesEmpty() (bool, error) {
	if len(tables) > 1 {
		return false, nil
	}
	return true, nil
}

func (*mockPartyModel) DbIsGuestsEmpty(string) (bool, error) {
	if len(guests) > 1 {
		return false, nil
	}
	return true, nil
}

func TestAddTableHandler(t *testing.T) {
	var emptyTable, testTable []models.Table
	testTable = append(testTable, models.Table{1, 3})

	tt := []struct {
		name       string
		method     string
		body       string
		expTable   []models.Table
		want       string
		statusCode int
	}{
		{
			name:       "Add 0 capacity",
			method:     http.MethodPost,
			body:       `{"capacity": 0}`,
			expTable:   emptyTable,
			want:       `must be greater than 0`,
			statusCode: http.StatusBadRequest,
		},
		{
			name:       "Add negative capacity",
			method:     http.MethodPost,
			body:       `{"capacity": -1}`,
			expTable:   emptyTable,
			want:       `must be greater than 0`,
			statusCode: http.StatusBadRequest,
		},
		{
			name:       "Add out of limit capacity",
			method:     http.MethodPost,
			body:       `{"capacity": 99999999999999999}`,
			expTable:   emptyTable,
			want:       `must be less than 4,294,967,295`,
			statusCode: http.StatusBadRequest,
		},
		{
			name:       "Add valid capacity",
			method:     http.MethodPost,
			body:       `{"capacity":10}`,
			expTable:   testTable,
			want:       `{"id":1,"capacity":10}`,
			statusCode: http.StatusOK,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			defer cleanup()
			request := httptest.NewRequest(tc.method, "/tables", strings.NewReader(tc.body))
			responseRecorder := httptest.NewRecorder()

			app := App{Party: &mockPartyModel{}}

			handler := http.HandlerFunc(app.AddTableHandler)
			handler.ServeHTTP(responseRecorder, request)

			if responseRecorder.Code != tc.statusCode {
				t.Errorf("Want status '%d', got '%d'", tc.statusCode, responseRecorder.Code)
			}

			if strings.TrimSpace(responseRecorder.Body.String()) != tc.want {
				t.Errorf("Want '%s', got '%s'", tc.want, responseRecorder.Body)
			}

			assert.Equal(t, tc.expTable, tables, "both should be equal")
		})
	}
}

func TestAddGuestListHandler(t *testing.T) {
	var testGuests1, testGuests2 []models.Guests

	testGuests1 = append(testGuests1, models.Guests{1, 1, "allotted", time.Time{}, "john"})

	testGuests2 = append(testGuests2, testGuests1...)
	testGuests2 = append(testGuests2, models.Guests{2, 2, "allotted", time.Time{}, "akhila"})

	tt := []struct {
		name       string
		method     string
		guestName  string
		body       string
		expGuests  []models.Guests
		want       string
		statusCode int
	}{
		{
			name:       "invalid table",
			method:     http.MethodPost,
			guestName:  "akhila",
			body:       `{"table": 3, "accompanying_guests": 10}`,
			expGuests:  testGuests1,
			want:       `Invalid table-id`,
			statusCode: http.StatusBadRequest,
		},
		{
			name:       "repeat name",
			method:     http.MethodPost,
			guestName:  "john",
			body:       `{"table": 2, "accompanying_guests": 2}`,
			expGuests:  testGuests1,
			want:       `Guest john already added`,
			statusCode: http.StatusBadRequest,
		},
		{
			name:       "repeat name case sensitive",
			method:     http.MethodPost,
			guestName:  "JOHN",
			body:       `{"table": 2, "accompanying_guests": 2}`,
			expGuests:  testGuests1,
			want:       `Guest john already added`,
			statusCode: http.StatusBadRequest,
		},
		{
			name:       "valid input",
			method:     http.MethodPost,
			guestName:  "akhila",
			body:       `{"table": 2, "accompanying_guests": 2}`,
			expGuests:  testGuests2,
			want:       `{"name":"akhila"}`,
			statusCode: http.StatusOK,
		},
		{
			name:       "negative accompanying guests",
			method:     http.MethodPost,
			guestName:  "akhila",
			body:       `{"table": 2, "accompanying_guests": -1}`,
			expGuests:  testGuests1,
			want:       `AccompanyingGuests must be 0 or greater`,
			statusCode: http.StatusBadRequest,
		},
		{
			name:       "more accompanying guests than capacity",
			method:     http.MethodPost,
			guestName:  "akhila",
			body:       `{"table": 1, "accompanying_guests": 10}`,
			expGuests:  testGuests1,
			want:       `Cannot allot table. Table capacity is 3`,
			statusCode: http.StatusBadRequest,
		},
		{
			name:       "already allocated table",
			method:     http.MethodPost,
			guestName:  "akhila",
			body:       `{"table": 1, "accompanying_guests": 0}`,
			expGuests:  testGuests1,
			want:       `Table already allotted to john`,
			statusCode: http.StatusBadRequest,
		},
	}

	app := App{Party: &mockPartyModel{}}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			addTables(2)
			addSingleGuest()
			defer cleanup()
			requrl := "/guest_list/" + tc.guestName
			request := httptest.NewRequest(tc.method, requrl, strings.NewReader(tc.body))
			param := map[string]string{"name": tc.guestName}
			request = mux.SetURLVars(request, param)
			responseRecorder := httptest.NewRecorder()
			handler := http.HandlerFunc(app.AddGuestListHandler)
			handler.ServeHTTP(responseRecorder, request)

			if responseRecorder.Code != tc.statusCode {
				t.Errorf("Want status '%d', got '%d'", tc.statusCode, responseRecorder.Code)
			}

			if strings.TrimSpace(responseRecorder.Body.String()) != tc.want {
				t.Errorf("Want '%s', got '%s'", tc.want, responseRecorder.Body)
			}

			assert.Equal(t, tc.expGuests, guests, "both should be equal")
		})
	}
}

func TestUpateGuestHandler(t *testing.T) {
	var testGuests1, testGuests2, testGuests3, testGuests4 []models.Guests
	testGuests1 = append(testGuests1, models.Guests{1, 1, "allotted", time.Time{}, "john"})
	testGuests2 = append(testGuests2, models.Guests{1, 1, "checked-in", time.Time{}, "john"})
	testGuests3 = append(testGuests3, models.Guests{1, 0, "checked-in", time.Time{}, "john"})
	testGuests4 = append(testGuests4, models.Guests{1, 2, "checked-in", time.Time{}, "john"})

	tt := []struct {
		name       string
		method     string
		guestName  string
		body       string
		expGuests  []models.Guests
		want       string
		statusCode int
	}{
		{
			name:       "invalid guest",
			method:     http.MethodPut,
			guestName:  "akhila",
			body:       `{"accompanying_guests": 10}`,
			expGuests:  testGuests1,
			want:       `Guest akhila is not present in Guestlist`,
			statusCode: http.StatusBadRequest,
		},
		{
			name:       "more acc guest than capacity",
			method:     http.MethodPut,
			guestName:  "john",
			body:       `{"accompanying_guests": 4}`,
			expGuests:  testGuests1,
			want:       `Cannot update number of accompanying guests. Table capacity is 3`,
			statusCode: http.StatusBadRequest,
		},
		{
			name:       "same acc guest as capacity",
			method:     http.MethodPut,
			guestName:  "john",
			expGuests:  testGuests1,
			body:       `{"accompanying_guests": 3}`,
			want:       `Cannot update number of accompanying guests. Table capacity is 3`,
			statusCode: http.StatusBadRequest,
		},
		{
			name:       "no update in acc guest",
			method:     http.MethodPut,
			guestName:  "john",
			body:       `{"accompanying_guests": 1}`,
			expGuests:  testGuests2,
			want:       `{"name":"john"}`,
			statusCode: http.StatusOK,
		},
		{
			name:       "less acc guest than previous",
			method:     http.MethodPut,
			guestName:  "john",
			body:       `{"accompanying_guests": 0}`,
			expGuests:  testGuests3,
			want:       `{"name":"john"}`,
			statusCode: http.StatusOK,
		},
		{
			name:       "greater valid acc guest than prev",
			method:     http.MethodPut,
			guestName:  "john",
			expGuests:  testGuests4,
			body:       `{"accompanying_guests": 2}`,
			want:       `{"name":"john"}`,
			statusCode: http.StatusOK,
		},
	}

	app := App{Party: &mockPartyModel{}}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			addTables(2)
			addSingleGuest()
			defer cleanup()
			requrl := "/guests/" + tc.guestName
			request := httptest.NewRequest(tc.method, requrl, strings.NewReader(tc.body))
			param := map[string]string{"name": tc.guestName}
			request = mux.SetURLVars(request, param)
			responseRecorder := httptest.NewRecorder()
			handler := http.HandlerFunc(app.UpdateGuestHandler)
			handler.ServeHTTP(responseRecorder, request)

			if responseRecorder.Code != tc.statusCode {
				t.Errorf("Want status '%d', got '%d'", tc.statusCode, responseRecorder.Code)
			}

			if strings.TrimSpace(responseRecorder.Body.String()) != tc.want {
				t.Errorf("Want '%s', got '%s'", tc.want, responseRecorder.Body)
			}

			assert.Equal(t, tc.expGuests, guests, "both should be equal")
		})
	}
}

func TestDeleteGuestHandler(t *testing.T) {
	var testGuests1, testGuests2, testGuests3 []models.Guests

	testGuests1 = append(testGuests1, models.Guests{1, 1, "allotted", time.Time{}, "john"})
	testGuests1 = append(testGuests1, models.Guests{2, 2, "checked-in", time.Time{}, "akhila"})

	testGuests2 = append(testGuests2, models.Guests{1, 1, "allotted", time.Time{}, "john"})
	testGuests2 = append(testGuests2, models.Guests{2, 2, "checked-out", time.Time{}, "akhila"})

	testGuests3 = append(testGuests3, testGuests1...)
	testGuests3 = append(testGuests3, models.Guests{3, 3, "checked-out", time.Time{}, "jack"})

	tt := []struct {
		name       string
		method     string
		guestName  string
		checkout   bool
		expGuests  []models.Guests
		want       string
		statusCode int
	}{
		{
			name:       "invalid guest",
			method:     http.MethodDelete,
			guestName:  "prasob",
			checkout:   false,
			expGuests:  testGuests1,
			want:       `Guest prasob is not present in Guestlist`,
			statusCode: http.StatusBadRequest,
		},
		{
			name:       "guest not arrived",
			method:     http.MethodDelete,
			guestName:  "john",
			checkout:   false,
			expGuests:  testGuests1,
			want:       `Request failed, guest not checked-in`,
			statusCode: http.StatusBadRequest,
		},
		{
			name:       "guest already checked-out",
			method:     http.MethodDelete,
			guestName:  "john",
			checkout:   true,
			expGuests:  testGuests3,
			want:       `Request failed, guest not checked-in`,
			statusCode: http.StatusBadRequest,
		},
		{
			name:       "valid input",
			method:     http.MethodDelete,
			guestName:  "akhila",
			checkout:   false,
			expGuests:  testGuests2,
			want:       ``,
			statusCode: http.StatusNoContent,
		},
	}

	app := App{Party: &mockPartyModel{}}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			addTables(2)
			addGuests()
			if tc.checkout {
				addCheckedOutGuest()
			}
			defer cleanup()
			requrl := "/guests/" + tc.guestName
			request := httptest.NewRequest(tc.method, requrl, nil)
			param := map[string]string{"name": tc.guestName}
			request = mux.SetURLVars(request, param)
			responseRecorder := httptest.NewRecorder()
			handler := http.HandlerFunc(app.DeleteGuestHandler)
			handler.ServeHTTP(responseRecorder, request)

			if responseRecorder.Code != tc.statusCode {
				t.Errorf("Want status '%d', got '%d'", tc.statusCode, responseRecorder.Code)
			}

			if strings.TrimSpace(responseRecorder.Body.String()) != tc.want {
				t.Errorf("Want '%s', got '%s'", tc.want, responseRecorder.Body)
			}

			assert.Equal(t, tc.expGuests, guests, "both should be equal")
		})
	}
}

func TestGetGuestListHandler(t *testing.T) {
	tt := []struct {
		name       string
		method     string
		addGuests  bool
		want       string
		statusCode int
	}{
		{
			name:       "Empty guests",
			method:     http.MethodGet,
			addGuests:  false,
			want:       `null`,
			statusCode: http.StatusOK,
		},
		{
			name:       "Get guests",
			method:     http.MethodGet,
			addGuests:  true,
			want:       `[{"table":1,"accompanying_guests":1,"name":"john"},{"table":2,"accompanying_guests":2,"name":"akhila"}]`,
			statusCode: http.StatusOK,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			if tc.addGuests {
				addGuests()
			}
			defer cleanup()
			request := httptest.NewRequest(tc.method, "/guest_list", nil)
			responseRecorder := httptest.NewRecorder()

			app := App{Party: &mockPartyModel{}}

			handler := http.HandlerFunc(app.GetGuestListHandler)
			handler.ServeHTTP(responseRecorder, request)

			if responseRecorder.Code != tc.statusCode {
				t.Errorf("Want status '%d', got '%d'", tc.statusCode, responseRecorder.Code)
			}

			if strings.TrimSpace(responseRecorder.Body.String()) != tc.want {
				t.Errorf("Want '%s', got '%s'", tc.want, responseRecorder.Body)
			}
		})
	}
}

func TestGetGuestsHandler(t *testing.T) {
	tt := []struct {
		name       string
		method     string
		addGuests  bool
		want       string
		statusCode int
	}{
		{
			name:       "Empty arrived guests",
			method:     http.MethodGet,
			addGuests:  false,
			want:       `null`,
			statusCode: http.StatusOK,
		},
		{
			name:       "Get arrived guests",
			method:     http.MethodGet,
			addGuests:  true,
			want:       `[{"time_arrived":"0001-01-01T00:00:00Z","accompanying_guests":2,"name":"akhila"}]`,
			statusCode: http.StatusOK,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			if tc.addGuests {
				addGuests()
				addCheckedOutGuest()
			} else {
				addCheckedOutGuest()
			}
			defer cleanup()
			request := httptest.NewRequest(tc.method, "/guests", nil)
			responseRecorder := httptest.NewRecorder()

			app := App{Party: &mockPartyModel{}}

			handler := http.HandlerFunc(app.GetGuestsHandler)
			handler.ServeHTTP(responseRecorder, request)

			if responseRecorder.Code != tc.statusCode {
				t.Errorf("Want status '%d', got '%d'", tc.statusCode, responseRecorder.Code)
			}

			if strings.TrimSpace(responseRecorder.Body.String()) != tc.want {
				t.Errorf("Want '%s', got '%s'", tc.want, responseRecorder.Body)
			}
		})
	}
}

func TestEmptySeatsHandler(t *testing.T) {
	tt := []struct {
		name       string
		method     string
		empty      bool
		want       string
		statusCode int
	}{
		{
			name:       "Empty tables",
			method:     http.MethodGet,
			empty:      true,
			want:       `{"seats_empty":0}`,
			statusCode: http.StatusOK,
		},
		{
			name:       "Get arrived guests",
			method:     http.MethodGet,
			empty:      false,
			want:       `{"seats_empty":9}`,
			statusCode: http.StatusOK,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			if !tc.empty {
				addTables(3)
				addGuests()
				addCheckedOutGuest()
			}
			defer cleanup()
			request := httptest.NewRequest(tc.method, "/seats_empty", nil)
			responseRecorder := httptest.NewRecorder()

			app := App{Party: &mockPartyModel{}}

			handler := http.HandlerFunc(app.GetEmptySeatsHandler)
			handler.ServeHTTP(responseRecorder, request)

			if responseRecorder.Code != tc.statusCode {
				t.Errorf("Want status '%d', got '%d'", tc.statusCode, responseRecorder.Code)
			}

			if strings.TrimSpace(responseRecorder.Body.String()) != tc.want {
				t.Errorf("Want '%s', got '%s'", tc.want, responseRecorder.Body)
			}
		})
	}
}
