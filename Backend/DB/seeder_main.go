package main

import (
    "log"
    "github.com/Sidarth-Roy/NorthStar-Intelligence/Backend/DB/seeder"
)

func main() {
    log.Println("🚀 Starting standalone seeder...")    
    seeder.RunDatabaseSetup()
    log.Println("✅ Seeding finished!")
}