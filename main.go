package main

import (
	"log"
	"net/http"
	"os"
	"salesAnalysisLumel/dbconnect"
	"salesAnalysisLumel/getdate"
	"salesAnalysisLumel/reader"
	"time"
)

func autoloadcsv(pFilePath string) {
	log.Println("autoloadcsv(+)")
	for {
		_ = reader.LoadCsv(pFilePath)
		now := time.Now()
		next := time.Date(now.Year(), now.Month(), now.Day(), 6, 0, 0, 0, now.Location())
		if now.After(next) {
			next = next.Add(24 * time.Hour)
		}
		time.Sleep(time.Until(next))
		lErr := reader.LoadCsv(pFilePath)
		if lErr != nil {
			log.Println("Error while loading csv: ", lErr)
		}
	}
}

func main() {
	log.Println("Server started(+)")
	f, err := os.OpenFile("./log/logfile"+time.Now().Format("02012006.15.04.05.000000000")+".txt", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()
	filepath := "./data/sales.csv"
	log.SetOutput(f)
	lErr := dbconnect.DbConnect()
	if lErr != nil {
		log.Println("Error while connection global database connection")
	} else {
		defer dbconnect.GDBConnection.Close()
		go autoloadcsv(filepath)

		http.HandleFunc("/api/CustomerAnalysis", getdate.CustomerAnalysis)
		http.HandleFunc("/api/refresh", reader.RefreshData)
		http.ListenAndServe(":29095", nil)
	}
	log.SetOutput(f)
}
