package main

import (
	"context"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

var (
	appName, lisAddr string
)

func configureEnv() {
	if err := godotenv.Load(); err != nil {
		log.Fatalln(err.Error())
	}

	appName = os.Getenv("APP_NAME")
	if appName == "" {
		appName = "api-gateway"
	}

	lisAddr = os.Getenv("PORT")
	if lisAddr == "" {
		lisAddr = ":8080"
	} else {
		lisAddr = ":" + lisAddr
	}
}

func serveReverseProxy(target string, w http.ResponseWriter, r *http.Request) {
	resp, err := http.Get(target)
	if err != nil {
		log.Fatalln(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	w.Write(body)
	return
}

func handleRequestAndRedirect(w http.ResponseWriter, r *http.Request) {
	orderServicePort := os.Getenv("ORDER_SERVICE_PORT")
	if orderServicePort == "" {
		orderServicePort = "8080"
	}
	proxyURL := "http://localhost:" + orderServicePort + "/v1/orders"
	serveReverseProxy(proxyURL, w, r)
}

func waitForShutdown(srv http.Server, l *log.Logger) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	sig := <-sigChan
	l.Println("Graceful shutdown:", sig)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	srv.Shutdown(ctx)
}

func main() {
	configureEnv()

	l := log.New(os.Stdout, strings.ToUpper(appName)+" ", log.LstdFlags)

	router := mux.NewRouter()
	router.HandleFunc("/v1/orders", handleRequestAndRedirect)

	srv := http.Server{
		Addr:         lisAddr,
		Handler:      router,
		ErrorLog:     l,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 20 * time.Second,
		IdleTimeout:  30 * time.Second,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			l.Fatalln(err)
		}
	}()
	l.Printf("%s has been started on %s\n", appName, lisAddr)

	waitForShutdown(srv, l)
}
