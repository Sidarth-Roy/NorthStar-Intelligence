module github.com/Sidarth-Roy/NorthStar-Intelligence/DB

go 1.21

// issues due remote model fetching, will be resolved in the future
replace github.com/Sidarth-Roy/NorthStar-Intelligence/Backend => ../Backend

require (
	github.com/Sidarth-Roy/NorthStar-Intelligence/Backend v0.0.0
	golang.org/x/text v0.20.0
	gorm.io/driver/postgres v1.5.4
	gorm.io/gorm v1.31.1
)

require (
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20221227161230-091c0ba34f0a // indirect
	github.com/jackc/pgx/v5 v5.4.3 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	golang.org/x/crypto v0.14.0 // indirect
)
