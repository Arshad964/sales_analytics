package reader

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"salesAnalysisLumel/dbconnect"
)

type FileData struct {
	OrderId      string
	ProductId    string
	CustomerId   string
	ProductName  string
	Category     string
	Region       string
	SaleDate     string
	SaleQuantity string
	UnitPrice    string
	Discount     string
	ShippingCost string
	PaymentType  string
	CustmrName   string
	CustmrEmail  string
	CustmrAddr   string
}

func LoadCsv(pFilePath string) error {
	log.Println("LoadCsv(+)")
	lFile, lErr := os.Open(pFilePath)
	if lErr != nil {
		return lErr
	}
	defer lFile.Close()

	var lGetFileData FileData
	var lGetFileDataArr []FileData

	reader := csv.NewReader(lFile)
	lCount := 0
	for {
		row, lErr := reader.Read()
		if lErr == io.EOF {
			break
		}
		if lErr != nil {
			return lErr
		}
		if lCount > 1 {
			lGetFileData.OrderId = row[0]
			lGetFileData.ProductId = row[1]
			lGetFileData.CustomerId = row[2]
			lGetFileData.ProductName = row[3]
			lGetFileData.Category = row[4]
			lGetFileData.Region = row[5]
			lGetFileData.SaleDate = row[6]
			lGetFileData.SaleQuantity = row[7]
			lGetFileData.UnitPrice = row[8]
			lGetFileData.Discount = row[9]
			lGetFileData.ShippingCost = row[10]
			lGetFileData.PaymentType = row[11]
			lGetFileData.CustmrName = row[12]
			lGetFileData.CustmrEmail = row[13]
			lGetFileData.CustmrAddr = row[14]

			lGetFileDataArr = append(lGetFileDataArr, lGetFileData)
		}

		// Extract data and insert with prepared statements
		// Convert data types, validate, and normalize
		lCount++
	}
	lInsCount := 0

	CustomerSql := `insert into customers(customer_id,name,email,address) values `
	OrderSql := `insert into Orders (Order_id,customer_id,saledate,region,payment_method) values `
	ProductSql := `insert into Products(product_id,name,category)`
	SaleItem := `insert into SaleItems(order_id,product_id,quantity_sold,unit_price,discount,shipping_cost) values `
	CustomerStringValues := ``
	OrderStringValues := ``
	ProductStringValues := ``
	SaleItemStringValues := ``

	for _, lFileDate := range lGetFileDataArr {
		lInsCount++
		CustomerStringValues += `(` + lFileDate.CustomerId + `,'` + lFileDate.CustmrName + `','` + lFileDate.CustmrEmail + `','` + lFileDate.CustmrAddr + `'),`
		OrderStringValues += `(` + lFileDate.OrderId + `,` + lFileDate.CustomerId + `,` + lFileDate.SaleDate + `,'` + lFileDate.Region + `','` + lFileDate.PaymentType + `'),`
		ProductStringValues += `(` + lFileDate.ProductId + `,'` + lFileDate.ProductName + `','` + lFileDate.Category + `'),`
		SaleItemStringValues += `(` + lFileDate.OrderId + `,` + lFileDate.ProductId + `,` + lFileDate.SaleQuantity + `,` + lFileDate.UnitPrice + `,` + lFileDate.Discount + `,` + lFileDate.ShippingCost + `),`
		// log.Println(CustomerSql + CustomerStringValues)
		if lInsCount == 1000 {
			lErr = dbconnect.ExecuteBulkStatement(dbconnect.GDBConnection, CustomerStringValues, CustomerSql)
			if lErr != nil {
				log.Println("LCSV01 Error inserting file into database for Customer table ", lErr)
				return lErr
			}
			lInsCount = 0
			CustomerStringValues = ""
			lErr = dbconnect.ExecuteBulkStatement(dbconnect.GDBConnection, OrderStringValues, OrderSql)
			if lErr != nil {
				log.Println("LCSV02 Error inserting file into database for Orders table ", lErr)
				return lErr
			}
			OrderStringValues = ""
			lErr = dbconnect.ExecuteBulkStatement(dbconnect.GDBConnection, ProductStringValues, ProductSql)
			if lErr != nil {
				log.Println("LCSV03 Error inserting file into database for Products table ", lErr)
				return lErr
			}
			ProductStringValues = ""
			lErr = dbconnect.ExecuteBulkStatement(dbconnect.GDBConnection, SaleItemStringValues, SaleItem)
			if lErr != nil {
				log.Println("LCSV04 Error inserting file into database for SaleItem table ", lErr)
				return lErr

			}
			SaleItemStringValues = ""
		}
	}
	if lInsCount > 0 {
		lErr = dbconnect.ExecuteBulkStatement(dbconnect.GDBConnection, CustomerStringValues, CustomerSql)
		if lErr != nil {
			log.Println("LCSV05 Error inserting file into database for Customer table ", lErr)
			return lErr
		}
		lInsCount = 0
		CustomerStringValues = ""
		lErr = dbconnect.ExecuteBulkStatement(dbconnect.GDBConnection, OrderStringValues, OrderSql)
		if lErr != nil {
			log.Println("LCSV06 Error inserting file into database for Orders table ", lErr)
			return lErr
		}
		OrderStringValues = ""
		lErr = dbconnect.ExecuteBulkStatement(dbconnect.GDBConnection, ProductStringValues, ProductSql)
		if lErr != nil {
			log.Println("LCSV07 Error inserting file into database for Products table ", lErr)
			return lErr
		}
		ProductStringValues = ""
		lErr = dbconnect.ExecuteBulkStatement(dbconnect.GDBConnection, SaleItemStringValues, SaleItem)
		if lErr != nil {
			log.Println("LCSV08 Error inserting file into database for SaleItem table ", lErr)
			return lErr

		}
		SaleItemStringValues = ""
	}
	log.Println("LoadCsv(-)")
	return nil
}

type RefreshResp struct {
	Status string `json:"status"`
	ErrMsg string `json:"errmsg"`
}

func RefreshData(w http.ResponseWriter, r *http.Request) {
	(w).Header().Set("Access-Control-Allow-Origin", "*")
	(w).Header().Set("Access-Control-Allow-Credentials", "true")
	(w).Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	(w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, credentials")
	if r.Method == "POST" {
		log.Println("RefreshData(+)")
		var lResponse RefreshResp
		lResponse.Status = "S"
		csvPath := "./data/sales.csv"
		lErr := LoadCsv(csvPath)
		logMessage := "Refresh successful"
		success := true
		if lErr != nil {
			logMessage = lErr.Error()
			success = false
		}

		_, lErr = dbconnect.GDBConnection.Exec(`INSERT INTO refresh_logs (success, message) VALUES ($1, $2)`, success, logMessage)
		if lErr != nil {
			log.Println("Error while refresh logs", lErr)
			lResponse.Status = "E"

		} else {
			log.Println("inserted Sucessfully")
		}
		lData, lErr := json.Marshal(lResponse)
		if lErr != nil {
			fmt.Fprintf(w, "Error taking data"+lErr.Error())
		} else {
			fmt.Fprint(w, string(lData))
		}
	}
	log.Println("RefreshData(-)")
}
