package main

import (
	"log"
	"net/http"
	"strings"

	"github.com/casbin/casbin/v2"
	"github.com/tg123/go-htpasswd"
)

var (
	enforcer *casbin.Enforcer
	users    *htpasswd.File
)

func initCasbin() {
	var err error
	// Load the model and policy from a file
	enforcer, err = casbin.NewEnforcer("/conf/model.conf", "/conf/policy.csv")
	if err != nil {
		log.Fatalf("Failed to create enforcer: %v", err)
	}
}

func readUsers() {
	var err error
	// Load the users from a .htpasswd formatted file
	users, err = htpasswd.New("/conf/users.htpasswd", htpasswd.DefaultSystems, nil)
	if err != nil {
		log.Fatalf("Failed to load users: %v", err)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	method := r.Header.Get("X-Forwarded-Method")
	uri := r.Header.Get("X-Forwarded-Uri")
	user, password, basicAuthOK := r.BasicAuth()

	log.Printf("user, uri, method: %s %s %s", user, uri, method)

	// docker-cli sends a request to /v2/ to determine whether or not to send the `Authorization` header
	if uri == "/v2/" {
		if !basicAuthOK || !users.Match(user, password) {
			w.Header().Add("WWW-Authenticate", "Basic realm=\"Restricted\"")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		w.WriteHeader(http.StatusOK)
		return
	}

	if !strings.HasPrefix(uri, "/v2/") {
		log.Println("not handling this request")
		w.WriteHeader(http.StatusOK)
		return
	}

	// Check permission
	enforcerOK, err := enforcer.Enforce(user, uri, method)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if !enforcerOK {
		unauthorized(w)
		return
	}

	// If allowed, proceed with the request
	w.WriteHeader(http.StatusOK)
}

func main() {
	initCasbin()
	readUsers()

	http.HandleFunc("/", handler)

	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Could not start server: %s\n", err)
	}
}

func unauthorized(w http.ResponseWriter) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusForbidden)
	w.Write([]byte(`{
	"errors": [
		{
			"code": "UNAUTHORIZED",
			"message": "authentication required",
			"detail": []
		}
	]
}`))
}
