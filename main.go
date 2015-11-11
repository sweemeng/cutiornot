package main

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"time"
	"log"
)

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

func main(){
	is_holiday := check_holiday()
	if is_holiday == true{
		fmt.Printf("Is holiday\n")
	}else{
		fmt.Printf("Not holiday\n")
	}
}
