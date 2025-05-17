package getdate

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"salesAnalysisLumel/dbconnect"
	"time"
)

type CustomerResp struct {
	Status         string  `json:"status"`
	ErrMsg         string  `json:"errmsg"`
	TotalCustomers int     `json:"totalcustomers"`
	TotalOrders    int     `json:"totalorders"`
	AverageValue   float64 `json:"averagevalue"`
}

func CustomerAnalysis(w http.ResponseWriter, r *http.Request) {
	(w).Header().Set("Access-Control-Allow-Origin", "*")
	(w).Header().Set("Access-Control-Allow-Credentials", "true")
	(w).Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	(w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, credentials")
	if r.Method == "GET" {
		log.Println("CustomerAnalysis(+)")
		var lResponse CustomerResp

		lResponse.Status = "S"

		start := r.URL.Query().Get("start")
		end := r.URL.Query().Get("end")
		startDate, lErr := time.Parse("2006-01-02", start)
		if lErr != nil {
			log.Println("CA01: ", lErr)
			lResponse.Status = "E"
			lResponse.ErrMsg = "CA01: " + lErr.Error()
		}

		endDate, lErr := time.Parse("2006-01-02", end)
		if lErr != nil {
			log.Println("CA02: ", lErr)
			lResponse.Status = "E"
			lResponse.ErrMsg = "CA02: " + lErr.Error()
		}
		fromstr := startDate.Format("2006-01-02")
		endstr := endDate.Format("2006-01-02")
		lResponse.TotalCustomers, lErr = TotalCustomers(fromstr, endstr)
		if lErr != nil {
			log.Println("CA03: ", lErr)
			lResponse.Status = "E"
			lResponse.ErrMsg = "CA03: " + lErr.Error()
		} else {
			lResponse.TotalOrders, lErr = TotalOrders(fromstr, endstr)
			if lErr != nil {
				log.Println("CA04: ", lErr)
				lResponse.Status = "E"
				lResponse.ErrMsg = "CA04: " + lErr.Error()
			} else {
				lResponse.AverageValue, lErr = AverageValue(fromstr, endstr)
				if lErr != nil {
					log.Println("CA05: ", lErr)
					lResponse.Status = "E"
					lResponse.ErrMsg = "CA05: " + lErr.Error()
				}
			}
		}
		lData, lErr := json.Marshal(lResponse)
		if lErr != nil {
			fmt.Fprintf(w, "Error taking data"+lErr.Error())
		} else {
			fmt.Fprint(w, string(lData))
		}

	}
}

func TotalCustomers(lFromdate, lTodate string) (int, error) {
	log.Println("TotalCustomers(+)")
	var lCount int
	lQuery := `SELECT COUNT(*)
        FROM orders
        WHERE saledate BETWEEN ? AND ?
		group by customer_id`
	lRows, lErr := dbconnect.GDBConnection.Query(lQuery, lFromdate, lTodate)
	if lErr != nil {
		log.Println("TC01: ", lErr)
		return lCount, lErr
	}
	lErr = lRows.Err()
	if lErr != nil {
		log.Println("TC02: ", lErr)
		return lCount, lErr
	}
	defer lRows.Close()
	for lRows.Next() {
		lErr = lRows.Scan(&lCount)
		if lErr != nil {
			log.Println("TC03: ", lErr)
			return lCount, lErr
		}
	}
	log.Println("TotalCustomers(-)")
	return lCount, nil
}

func TotalOrders(lFromdate, lTodate string) (int, error) {
	log.Println("TotalOrder(+)")
	var lCount int
	lQuery := `SELECT nvl(COUNT(1),0)
        FROM orders
        WHERE saledate BETWEEN ? AND ?`
	lRows, lErr := dbconnect.GDBConnection.Query(lQuery, lFromdate, lTodate)
	if lErr != nil {
		log.Println("TO01: ", lErr)
		return lCount, lErr
	}
	lErr = lRows.Err()
	if lErr != nil {
		log.Println("TO02: ", lErr)
		return lCount, lErr
	}
	defer lRows.Close()
	for lRows.Next() {
		lErr = lRows.Scan(&lCount)
		if lErr != nil {
			log.Println("TO03: ", lErr)
			return lCount, lErr
		}
	}
	log.Println("TotalOrder(-)")
	return lCount, nil
}

func AverageValue(lFromdate, lTodate string) (float64, error) {
	log.Println("AverageValue(+)")
	var lCount float64
	lQuery := `SELECT nvl(AVG(order_total), 0)
        FROM (
            SELECT SUM(s.unit_price * s.quantity_sold * (1 - s.discount)) AS order_total
            FROM saleitems s
            JOIN orders o ON s.order_id = o.id
            WHERE o.date_of_sale BETWEEN ? AND ?
            GROUP BY o.id
        ) AS order_totals`
	lRows, lErr := dbconnect.GDBConnection.Query(lQuery, lFromdate, lTodate)
	if lErr != nil {
		log.Println("AV01: ", lErr)
		return lCount, lErr
	}
	lErr = lRows.Err()
	if lErr != nil {
		log.Println("AV02: ", lErr)
		return lCount, lErr
	}
	defer lRows.Close()
	for lRows.Next() {
		lErr = lRows.Scan(&lCount)
		if lErr != nil {
			log.Println("AV03: ", lErr)
			return lCount, lErr
		}
	}
	log.Println("AverageValue(-)")
	return lCount, nil
}
