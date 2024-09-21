package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type Car struct {
	ID                 int
	NumberPlate        string
	Make               string
	Model              string
	NumberOfPassengers int
	SaccoID            int
	SaccoName          string
}

// Get all cars in the database
func getAllCars(db *sql.DB) ([]Car, error) {
	rows, err := db.Query(`SELECT cars.id, number_plate, make, model, no_of_passengers, saccos.id AS sacco_id, sacco_name
						   FROM cars
						   JOIN saccos ON cars.sacco_id = saccos.id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cars []Car
	for rows.Next() {
		var car Car
		if err := rows.Scan(&car.ID, &car.NumberPlate, &car.Make, &car.Model, &car.NumberOfPassengers, &car.SaccoID, &car.SaccoName); err != nil {
			return nil, err
		}
		cars = append(cars, car)
	}

	return cars, nil
}

// Add cars to the database
func addCar(db *sql.DB, car Car) error {
	_, err := db.Exec("INSERT INTO cars (number_plate, make, model, no_of_passengers, sacco_id) VALUES (?, ?, ?, ?, ?)", car.NumberPlate, car.Make, car.Model, car.NumberOfPassengers, car.SaccoID)
	return err
}

// Handles requests to the /car route
func carHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
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
			Cars   []Car
			Saccos []Sacco
		}{
			Cars:   cars,
			Saccos: saccos,
		}

		err = tmpl.ExecuteTemplate(w, "cars", data)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			fmt.Println("Error executing template")
			return
		}

	case "POST":
		noOfpassengers, err := strconv.Atoi(r.FormValue("numberOfPassengers"))
		if err != nil {
			http.Error(w, "Invalid number of passengers", http.StatusBadRequest)
			return
		}

		saccoID, err := strconv.Atoi(r.FormValue("saccoID"))
		if err != nil {
			http.Error(w, "Invalid sacco ID", http.StatusBadRequest)
			return
		}

		car := Car{
			NumberPlate:        r.FormValue("numberPlate"),
			Make:               r.FormValue("make"),
			Model:              r.FormValue("model"),
			NumberOfPassengers: noOfpassengers,
			SaccoID:            saccoID,
		}
		if err := addCar(db, car); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/cars", http.StatusSeeOther)
	}
}

func editCarHandler(w http.ResponseWriter, r *http.Request) {
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

	noOfPassengers, err := strconv.Atoi(r.FormValue("numberOfPassengers"))
	if err != nil {
		http.Error(w, "Invalid number of passengers", http.StatusBadRequest)
		return
	}

	saccoID, err := strconv.Atoi(r.FormValue("saccoID"))
	if err != nil {
		http.Error(w, "Invalid sacco ID", http.StatusBadRequest)
		return
	}

	car := Car{
		ID:                 id,
		NumberPlate:        r.FormValue("numberPlate"),
		Make:               r.FormValue("make"),
		Model:              r.FormValue("model"),
		NumberOfPassengers: noOfPassengers,
		SaccoID:            saccoID,
	}

	if err := updateCar(db, car); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/cars", http.StatusSeeOther)
}

// Handles the GET request to fetch a car's details for editing
func getCarHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Path[len("/cars/"):]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	if r.Method == http.MethodGet {
		car, err := getCarByID(db, id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		response, err := json.Marshal(car)
		if err != nil {
			http.Error(w, "Error marshaling JSON", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(response)
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}

// Fetch a single car's details from the database
func getCarByID(db *sql.DB, id int) (Car, error) {
	var car Car
	row := db.QueryRow("SELECT cars.id, number_plate, make, model, no_of_passengers, saccos.id AS sacco_id, sacco_name FROM cars JOIN saccos ON cars.sacco_id = saccos.id WHERE cars.id = ?", id)
	err := row.Scan(&car.ID, &car.NumberPlate, &car.Make, &car.Model, &car.NumberOfPassengers, &car.SaccoID, &car.SaccoName)
	if err != nil {
		return car, err
	}
	return car, nil
}

// Update car
func updateCar(db *sql.DB, car Car) error {
	_, err := db.Exec(`
        UPDATE cars
        SET number_plate = ?, make = ?, model = ?, no_of_passengers = ?, sacco_id = ?
        WHERE id = ?`,
		car.NumberPlate, car.Make, car.Model, car.NumberOfPassengers, car.SaccoID, car.ID)
	return err
}

// Handles the DELETE request to delete a car
func deleteCarHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	if err := deleteCar(db, id); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/cars", http.StatusSeeOther)
}

// Delete car from the database
func deleteCar(db *sql.DB, id int) error {
	_, err := db.Exec(`DELETE FROM cars WHERE id = ?`, id)
	return err
}
