package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"toanpm0510/soniclabs/internal/cache"
	"toanpm0510/soniclabs/internal/config"
	httpx "toanpm0510/soniclabs/internal/http"
	"toanpm0510/soniclabs/internal/obs"
	"toanpm0510/soniclabs/internal/store/pg"
)

var (
	doMigrate = flag.Bool("migrate", false, "run migrations then exit")
	doSeed    = flag.Bool("seed", false, "seed minimal data then exit")
)

func main() {
	flag.Parse()
	cfg := config.Load()
	log.Printf("starting api on %s", cfg.HTTPAddr)

	logger, cleanup := obs.NewLogger()
	defer cleanup()

	ctx := context.Background()
	// DB
	db, err := pg.Connect(ctx, cfg.PGURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Pool.Close()

	// Redis
	rc := cache.NewRedis(cfg.RedisAddr, cfg.RedisDB)
	if err := cache.Ping(ctx, rc); err != nil {
		log.Fatal(err)
	}
	defer rc.Close()

	// migrations (exec simple SQL files)
	if *doMigrate {
		pg.MustExecFile(ctx, db, "migrations/0001_init.sql")
		pg.MustExecFile(ctx, db, "migrations/0002_indexes.sql")
		log.Println("migrated")
		return
	}
	if *doSeed {
		pg.MustExecFile(ctx, db, "migrations/seed.sql")
		log.Println("seeded")
		return
	}

	// HTTP server
	router := httpx.NewRouter(logger)
	srv := &http.Server{
		Addr:              cfg.HTTPAddr,
		Handler:           router,
		ReadHeaderTimeout: cfg.ReadHeaderTimeout,
	}

	// graceful shutdown
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctxShut, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
	defer cancel()
	if err := srv.Shutdown(ctxShut); err != nil {
		log.Printf("shutdown err: %v", err)
	}
}
