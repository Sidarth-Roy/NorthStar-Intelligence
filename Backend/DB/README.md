docker-compose up -d
go mod edit -replace github.com/Sidarth-Roy/NorthStar-Intelligence/Backend=../Backend
go mod download
go mod tidy
go run seed.go




2. How to Access the UI
Run docker-compose up -d again to pull and start the new UI container.

Open your browser and go to: http://localhost:8080

Login with:

Email: admin@admin.com

Password: admin


3. Connect pgAdmin to your DB
Once logged into pgAdmin, you need to register your specific database server:

Right-click "Servers" > Register > Server...

General Tab: Name it NorthStar-Local.

Connection Tab:

Host name/address: db (This is the service name from your docker-compose file).

Port: 5432

Maintenance database: northwind

Username: user

Password: password

Save.