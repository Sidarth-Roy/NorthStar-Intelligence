db-setup:
	docker-compose up -d

run-backend:
	go run backend/main.go

db-close:
	docker-compose down
	
db-clean:
	docker-compose down -v