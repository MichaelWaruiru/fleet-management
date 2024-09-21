package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type Driver struct {
	ID         int
	DriverName string
	IDNumber   string
	Contact    string
	CarID      int
	SaccoID    int
	Car        string
	SaccoName  string
}

// GetAllDrivers retrieves all drivers from the database
func getAllDrivers(db *sql.DB) ([]Driver, error) {
	rows, err := db.Query(`SELECT drivers.id, driver_name, drivers.id_number, drivers.contact, car_id, drivers.sacco_id, cars.number_plate, sacco_name 
						   FROM drivers 
						   JOIN cars ON drivers.car_id = cars.id
						   JOIN saccos ON drivers.sacco_id = saccos.id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var drivers []Driver
	for rows.Next() {
		var driver Driver
		if err := rows.Scan(&driver.ID, &driver.DriverName, &driver.IDNumber, &driver.Contact, &driver.CarID, &driver.SaccoID, &driver.Car, &driver.SaccoName); err != nil {
			return nil, err
		}
		drivers = append(drivers, driver)
	}

	return drivers, nil
}

// CreateDriver inserts a new driver into the database
func addDriver(db *sql.DB, driver Driver) error {
	_, err := db.Exec("INSERT INTO drivers (driver_name, id_number, contact, car_id, sacco_id) VALUES (?, ?, ?, ?, ?)", driver.DriverName, driver.IDNumber, driver.Contact, driver.CarID, driver.SaccoID)
	return err
}

// driverHandler handles requests to the /drivers route
func driverHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		drivers, err := getAllDrivers(db)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Fetch cars and saccos for dropdowns
		cars, err := getAllCars(db)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		saccos, err := getAllSaccos(db)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		data := struct {
			Drivers []Driver
			Cars    []Car
			Saccos  []Sacco
		}{
			Drivers: drivers,
			Cars:    cars,
			Saccos:  saccos,
		}

		if err := tmpl.ExecuteTemplate(w, "drivers", data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

	case "POST":
		r.ParseForm()
		carID, err := strconv.Atoi(r.FormValue("assignedCar"))
		if err != nil {
			http.Error(w, "Invalid car ID", http.StatusBadRequest)
			return
		}

		saccoID, err := strconv.Atoi(r.FormValue("saccoID"))
		if err != nil {
			http.Error(w, "Invalid sacco ID", http.StatusBadRequest)
			return
		}

		driver := Driver{
			DriverName: r.FormValue("driverName"),
			IDNumber:   r.FormValue("idNumber"),
			Contact:    r.FormValue("contact"),
			CarID:      carID,
			SaccoID:    saccoID,
		}
		if err := addDriver(db, driver); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/drivers", http.StatusSeeOther)
	}
}

func getSaccoByCarHandler(w http.ResponseWriter, r *http.Request) {
	carID := r.URL.Query().Get("car_id")
	if carID == "" {
		http.Error(w, "Car ID is required", http.StatusBadRequest)
		return
	}

	sacco, err := getSaccoByCarID(db, carID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(sacco)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

func getSaccoByCarID(db *sql.DB, carID string) (Sacco, error) {
	var sacco Sacco
	query := `SELECT s.id, s.sacco_name FROM saccos s
              JOIN cars c ON s.id = c.sacco_id
              WHERE c.id = ?`
	row := db.QueryRow(query, carID)
	err := row.Scan(&sacco.ID, &sacco.SaccoName)
	if err != nil {
		return sacco, err
	}
	return sacco, nil
}

func editDriverHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Parse form values
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	carID, err := strconv.Atoi(r.FormValue("editAssignedCar"))
	if err != nil {
		http.Error(w, "Invalid car ID", http.StatusBadRequest)
		return
	}

	saccoID, err := strconv.Atoi(r.FormValue("editSaccoID"))
	if err != nil {
		http.Error(w, "Invalid sacco ID", http.StatusBadRequest)
		return
	}

	driver := Driver{
		ID:         id,
		DriverName: r.FormValue("editDriverName"),
		IDNumber:   r.FormValue("editIDNumber"),
		Contact:    r.FormValue("editContact"),
		CarID:      carID,
		SaccoID:    saccoID,
	}

	fmt.Println(driver)

	if err := updateDriver(db, driver); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/drivers", http.StatusSeeOther)
}

func updateDriver(db *sql.DB, driver Driver) error {
	_, err := db.Exec(`
	UPDATE drivers
	SET driver_name = ?, id_number = ?, contact = ?, car_id = ?, sacco_id = ?
	WHERE id = ?`,
		driver.DriverName, driver.IDNumber, driver.Contact, driver.CarID, driver.SaccoID, driver.ID) // Corrected: sacco_id and car_id
	return err
}