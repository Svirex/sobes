package main

import (
	"encoding/json"
	"flag"
	"io"
	"log"
	"log/slog"
	"math/rand"
	"net/http"
	"os"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

var errorFile *os.File
var logger *slog.Logger

var index = 0

var firstList = []string{
	"Adele Ingram",
	"Imran Bryan",
	"Cecilia Odom",
	"Jaydon Gould",
	"Elodie Hendrix",
}

var secondList = []string{
	"Sumaiya Cruz",
	"Ernest Stafford",
	"Lorraine House",
	"Gregory O'Doherty",
	"Mikolaj Dale",
}

func main() {
	var addr string
	flag.StringVar(&addr, "a", "localhost:8080", "host:port")
	errorsFile := flag.String("e", "errors.log", "errors file")
	flag.Parse()
	var err error
	errorFile, err = os.OpenFile(*errorsFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0664)
	if err != nil {
		log.Fatalf("open error file: %v", err)
	}

	logger = slog.New(slog.NewJSONHandler(errorFile, nil))
	Main(addr)
	defer os.Remove(*errorsFile)
}

func Main(addr string) {

	router := chi.NewRouter()
	router.Use(middleware.Recoverer)
	router.Post("/api/auth", Auth)
	router.Get("/api/servers", ListServers)

	http.ListenAndServe(addr, router)
}

type AuthData struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type Token struct {
	Value string `json:"token"`
}

func Auth(response http.ResponseWriter, request *http.Request) {
	if request.Header.Get("Content-Type") != "application/json" {
		response.WriteHeader(http.StatusBadRequest)
		logger.Error("invalid Content-Type header")
		return
	}
	bodyBytes, err := io.ReadAll(request.Body)
	if err != nil {
		response.WriteHeader(http.StatusBadRequest)
		logger.Error("couldn't read body", "err", err)
		return
	}
	if len(bodyBytes) == 0 {
		response.WriteHeader(http.StatusBadRequest)
		logger.Error("empty body")
		return
	}
	var authData AuthData
	err = json.Unmarshal(bodyBytes, &authData)
	if err != nil {
		response.WriteHeader(http.StatusBadRequest)
		logger.Error("couldn't unmarshal body", "err", err)
		return
	}
	if authData.Login != "test" || authData.Password != "test" {
		e := ErrorResponse{
			Error: "invalid login or password",
		}
		body, err := json.Marshal(&e)
		if err != nil {
			response.WriteHeader(http.StatusBadRequest)
			logger.Error("couldn't marshal error response", "err", err)
			return
		}
		response.Header().Add("Content-Type", "application/json")
		response.WriteHeader(http.StatusBadRequest)
		response.Write(body)
		return
	}
	token := Token{
		Value: "this_is_simple_token",
	}
	body, err := json.Marshal(&token)
	if err != nil {
		response.WriteHeader(http.StatusBadRequest)
		logger.Error("couldn't marshal token", "err", err)
		return
	}
	response.Write(body)
}

func ListServers(response http.ResponseWriter, request *http.Request) {
	authHeader := request.Header.Get("Authorization")
	if authHeader == "" {
		response.WriteHeader(http.StatusBadRequest)
		logger.Error("not found Authorization header")
		return
	}
	splitted := strings.Split(authHeader, " ")
	if len(splitted) != 2 {
		response.WriteHeader(http.StatusBadRequest)
		logger.Error("invalid authorization header")
		return
	}
	if splitted[0] != "Bearer" {
		response.WriteHeader(http.StatusBadRequest)
		logger.Error("invalid scheme, need Bearer")
		return
	}
	if splitted[1] != "this_is_simple_token" {
		response.WriteHeader(http.StatusBadRequest)
		logger.Error("invalid token")
		return
	}
	if index == 0 {
		index = 1
		rand.Shuffle(len(firstList), func(i, j int) { firstList[i], firstList[j] = firstList[j], firstList[i] })
		body, err := json.Marshal(&firstList)
		if err != nil {
			response.WriteHeader(http.StatusBadRequest)
			logger.Error("marshal forst list")
			return
		}
		response.Write(body)
	} else {
		index = 0
		rand.Shuffle(len(secondList), func(i, j int) { secondList[i], secondList[j] = secondList[j], secondList[i] })
		body, err := json.Marshal(&secondList)
		if err != nil {
			response.WriteHeader(http.StatusBadRequest)
			logger.Error("marshal forst list")
			return
		}
		response.Write(body)
	}
}
