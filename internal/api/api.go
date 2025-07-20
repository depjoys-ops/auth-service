package api

import (
	"log"
	"net/http"

	"github.com/depjoys-ops/auth-service/internal/config"
	apihttp "github.com/depjoys-ops/auth-service/internal/controller/http"
)

func Run(cfg *config.Config) {

	log.Println("Starting service on:", cfg.HTTPServer.Addr)

	srv := &http.Server{
		Addr:         cfg.HTTPServer.Addr,
		Handler:      apihttp.NewRouter(),
		ReadTimeout:  cfg.HTTPServer.ReadTimeout,
		WriteTimeout: cfg.HTTPServer.WriteTimeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
