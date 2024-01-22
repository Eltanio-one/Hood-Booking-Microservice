package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"bookings.com/m/handlers"
)

func main() {
	// instantiate a new logger
	l := log.New(os.Stdout, "booking-api", log.LstdFlags)

	// instantiate handlers
	regHandler := handlers.NewRegisterHandler(l)
	loginHandler := handlers.NewLoginHandler(l)
	userHandler := handlers.NewUserHandler(l)
	hoodHandler := handlers.NewHoodHandler(l)
	bookingHandler := handlers.NewBookingHandler(l)

	mux := http.NewServeMux()
	mux.Handle("/register", regHandler)
	mux.Handle("/login", loginHandler)
	mux.Handle("/user", userHandler)
	mux.Handle("/user/", userHandler)
	mux.Handle("/hood", hoodHandler)
	mux.Handle("/hood/", hoodHandler)
	mux.Handle("/booking", bookingHandler)
	mux.Handle("/booking/", bookingHandler)

	srvr := &http.Server{
		Addr:         "localhost:9090",
		Handler:      mux,
		ErrorLog:     l,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// set up Goroutine for when the server starts
	go func() {
		l.Println("Server starting on port 9090")

		err := srvr.ListenAndServe()
		if err != nil {
			l.Println("Error serving due to:", err)
			os.Exit(1)
		}
	}()

	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)

	sig := <-signalChannel
	l.Println("Server terminating, gracefully exiting", sig)

	timeoutContext, cancelCtx := context.WithTimeout(context.Background(), 30*time.Second)

	defer func() {
		cancelCtx()
	}()

	if e := srvr.Shutdown(timeoutContext); e != nil {
		log.Fatalf("server Shutdown Failed:%+s", e)
	}
}
