package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"time"
	"log"
	"net/http"
	"encoding/json"
)

type Holiday struct {
	Holiday bool
}

func check_holiday() bool{
	var hdate string
	is_holiday := true
	now := time.Now()
	year, month, day := now.Date()
	timezone, err := time.LoadLocation("Asia/Kuala_Lumpur")

	if err != nil{
		log.Fatal(err)
	}
	// Because Holiday start on 0 hour on that day. And there is no date field
	test_time := time.Date(year, month, day, 0, 0, 0, 0, timezone)

	db, err := sql.Open("sqlite3", "./holiday.db")
	if err != nil{
		log.Fatal(err)
	}
	defer db.Close()

	err = db.QueryRow("select hdate from holiday where hdate=?", test_time.Format(time.RFC3339)).Scan(&hdate)

	if err != nil {
		if err == sql.ErrNoRows {
			is_holiday = false
		}else {
			log.Fatal(err)
		}
	}
	return is_holiday
}

func holiday_view(w http.ResponseWriter, r *http.Request){
	is_holiday := check_holiday()
	if is_holiday == true{
		fmt.Fprintf(w, "It's a holiday!")
	}else{
		fmt.Fprintf(w, "It's not a holiday!")
	}
}

func holiday_api(w http.ResponseWriter, r *http.Request){
	var holiday_struct Holiday
	if check_holiday() == true{
		holiday_struct.Holiday = true
	}else{
		holiday_struct.Holiday = false
	}

	js, err := json.Marshal(holiday_struct)
	if err != nil{
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func main(){
	http.HandleFunc("/", holiday_view)
	http.HandleFunc("/api", holiday_api)
	err := http.ListenAndServe(":9000", nil)
	if err != nil {
		log.Fatal(err)
	}
}
