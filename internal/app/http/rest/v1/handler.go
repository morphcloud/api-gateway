package v1

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

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

func generateProxyURL(service string, r *http.Request) string {
	var proxyHost string
	switch service {
	case "oauth":
		proxyHost = os.Getenv("OAUTH_SERVER_HOST")
	case "customers":
		proxyHost = os.Getenv("CUSTOMER_SERVICE_HOST")
	case "orders":
		proxyHost = os.Getenv("ORDER_SERVICE_HOST")
	}

	return "http://" + proxyHost + r.URL.Path
}

func HandleRequestAndRedirect(w http.ResponseWriter, r *http.Request) {
	proxyURL := generateProxyURL(strings.Split(r.URL.Path, "/")[2], r)

	serveReverseProxy(proxyURL, w, r)
}
