db-setup:
	docker-compose up -d

run-backend:
	go run backend/main.go
	
db-clean:
	docker-compose down -v

db-start:
	docker-compose start

db-stop:
	docker-compose stop