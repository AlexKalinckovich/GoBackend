package main

import (
	"database/sql"
	handler "github.com/brota/gobackend/internal/shared/common_handler"
	"github.com/brota/gobackend/internal/shared/config"
	"github.com/brota/gobackend/internal/shared/db"
	"github.com/brota/gobackend/internal/shared/redis"
	userhandler "github.com/brota/gobackend/internal/user/handler"
	repository2 "github.com/brota/gobackend/internal/user/repository"
	redis2 "github.com/redis/go-redis/v9"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	portString := os.Getenv("PORT")
	if portString == "" {
		log.Fatal("$PORT must be set in .env")
	}

	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		dbHost = "localhost"
	}

	dbPort := os.Getenv("DB_PORT")
	if dbPort == "" {
		dbPort = "3306"
	}

	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		log.Fatal("$DB_NAME must be set in .env")
	}

	dbUser := os.Getenv("DB_USER")
	if dbUser == "" {
		log.Fatal("$DB_USER must be set in .env")
	}

	dbPassword := os.Getenv("DB_PASSWORD")
	if dbPassword == "" {
		log.Fatal("$DB_PASSWORD must be set in .env")
	}

	dbURL := dbUser + ":" + dbPassword + "@tcp(" + dbHost + ":" + dbPort + ")/" + dbName + "?parseTime=true"

	conn, err := sql.Open("mysql", dbURL)
	if err != nil {
		log.Fatal("Error connecting to database: " + err.Error())
	}
	defer func(conn *sql.DB) {
		err := conn.Close()
		if err != nil {
			log.Printf("Error closing database connection: " + err.Error())
		}
	}(conn)

	err = conn.Ping()
	if err != nil {
		log.Fatal("Error pinging database: " + err.Error())
	}
	log.Println("Successfully connected to database")

	redisCfg := config.RedisConfig{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	}
	if redisCfg.Addr == "" {
		redisCfg.Addr = "localhost:6379"
	}
	rdb, err := redis.NewClient(redisCfg)
	if err != nil {
		log.Fatalf("Cannot connect to Redis: %v", err)
	}
	defer func(rdb redis2.Client) {
		_ = rdb.Close()
	}(*rdb)

	queries := db.New(conn)

	baseUserRepo := repository2.NewUserRepositoryWithQueriesAndConn(queries, conn)

	cachedUserRepo := repository2.NewCachedUserRepository(baseUserRepo, rdb, 5*time.Minute)

	userHandler := userhandler.NewUserHandler(cachedUserRepo)

	readinessHandler := handler.NewReadinessHandler()
	testErrorHandler := handler.NewTestErrorHandler()

	router := chi.NewRouter()
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	v1Router := chi.NewRouter()
	v1Router.Get("/health", readinessHandler.ServeHTTP)
	v1Router.Get("/testError", testErrorHandler.ServeHTTP)

	v1Router.Post("/users", userHandler.CreateUser)
	v1Router.Patch("/users/{id}", userHandler.PatchUser)
	v1Router.Put("/users/{id}", userHandler.UpdateUser)
	v1Router.Get("/users/{id}", userHandler.GetUser)

	router.Mount("/v1", v1Router)

	srv := http.Server{
		Handler: router,
		Addr:    ":" + portString,
	}

	log.Printf("Starting server on port %s", portString)
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
