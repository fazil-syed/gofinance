package main

import (
	_ "fazil-syed/gofinance/docs"
	"fazil-syed/gofinance/internal/config"
	"fazil-syed/gofinance/internal/handlers"
	"fazil-syed/gofinance/internal/middleware"
	"fmt"
	"log"
	"net/http"

	httpSwagger "github.com/swaggo/http-swagger"
)

// @title						Go Finance Server
// @version					1.0
// @host						localhost:8080
// @BasePath					/
// @description				Server to answer Stock related queries
//
// @securitydefinitions.basic	BasicAuth
func main() {
	cfg, err := config.LoadConfig(".")
	port := 8080
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}
	handler := handlers.NewHandler(cfg)
	authMiddleWare := middleware.NewAuthMiddleWare(cfg)

	http.Handle("/swagger/", httpSwagger.WrapHandler)

	http.HandleFunc("/price", authMiddleWare.BasicAuthMiddleWare(handler.GetPriceAtDayHandler))
	log.Printf("server running on :%d", port)
	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}
