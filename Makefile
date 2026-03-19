setup:
	docker-compose up -d

run-backend:
	go run backend/main.go

clean:
# 	docker-compose down -v
	docker-compose down