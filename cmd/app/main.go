package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"markitos-it-svc-goldens/internal/application/services"
	"markitos-it-svc-goldens/internal/infrastructure/persistence/postgres"
	"net"
	"os"
	"os/signal"
	"syscall"

	grpcserver "markitos-it-svc-goldens/internal/infrastructure/grpc"
	pb "markitos-it-svc-goldens/proto"

	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func getEnvRequired(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("❌ Required environment variable %s is not set", key)
	}
	return value
}

func main() {
	log.Println("🚀 Starting Goldens gRPC Service...")
	db, repo := loadDatabase()
	defer db.Close()

	ctx := context.Background()
	if err := repo.InitSchema(ctx); err != nil {
		log.Fatalf("❌ Failed to initialize schema: %v", err)
	}

	if err := repo.SeedData(ctx); err != nil {
		log.Printf("⚠️  Failed to seed data: %v", err)
	}
	docService := services.NewGoldenService(repo)

	grpcPort := getEnvRequired("GRPC_PORT")
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", grpcPort))
	if err != nil {
		log.Fatalf("❌ Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterGoldenServiceServer(grpcServer, grpcserver.NewGoldenServer(docService))
	reflection.Register(grpcServer)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Printf("🎯 gRPC server listening on :%s", grpcPort)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("❌ Failed to serve: %v", err)
		}
	}()

	<-sigChan
	log.Println("\n🛑 Shutting down gracefully...")
	grpcServer.GracefulStop()
	log.Println("👋 Service stopped")
}

func loadDatabase() (*sql.DB, *postgres.GoldenRepository) {
	log.Println("🚀 loading database")
	dbHost := getEnvRequired("DB_HOST")
	dbPort := getEnvRequired("DB_PORT")
	dbUser := getEnvRequired("DB_USER")
	dbPass := getEnvRequired("DB_PASS")
	dbName := getEnvRequired("DB_NAME")
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPass, dbName)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("❌ Failed to connect to database: %v", err)
	}
	if err := db.Ping(); err != nil {
		log.Fatalf("❌ Failed to ping database: %v", err)
	}
	log.Println("✅ Connected to PostgreSQL")

	return db, postgres.NewGoldenRepository(db)
}
