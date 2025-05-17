# Sales Analytics API (Go)

A high-performance backend system built in Go to manage and analyze large-scale sales data from CSV files. This service supports periodic data refresh, automatic scheduled tasks, and RESTful APIs for customer analysis.

---

## Features

- Efficient CSV ingestion for millions of records
- Normalized Mysql schema for orders, products, customers, etc.
- Daily and on-demand data refresh mechanism
- REST API for customer metrics (total customers, orders, AOV)
- Scheduled daily data refresh at **6:00 AM**
- Logging for data refresh activity


## Tech Stack

- **Go (Golang)**: Backend API
- **PostgreSQL**: Relational database
- **CSV**: File ingestion


###  Prerequisites

- Go 1.18+
- Mariadb 11.7.2


### calling Method

http://localhost:29095//api/CustomerAnalysis

### Running Servers
go run main.go



### Data Refresh
Automatic: Every day at 6:00 AM

### Api List

Manual: See API below
Method	-  POST
Route  - /api/refresh
Description	 -  Refresh data from CSV 
Body / Params	-  { "path": "/path/to/file.csv" }
Sample Response -  { "status": "S" ,"errmsg":"" }
			


Method - GET	
Route - /api/CustomerAnalysis
Description - 	Customer insights by date range	
Params  - ?start=2024-01-01&end=2024-04-01
Sample Response -  { "status": "S" ,"errmsg":"" , "totalcustomers": 100, "totalorders": 150, "averagevalue": 83.33 }

