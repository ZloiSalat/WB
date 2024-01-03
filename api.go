package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/nats-io/stan.go"
	"log"
	"net/http"
)

type APIServer struct {
	listenAddr string
	store      Storage
	cache      InMemoryCache
	stanConn   Stream
}

func NewAPIServer(listerAddr string, store Storage, c InMemoryCache, st Stream) *APIServer {
	return &APIServer{
		listenAddr: listerAddr,
		store:      store,
		cache:      c,
		stanConn:   st,
	}
}

func (s *APIServer) ConnectToNATSStreaming(clientID, clusterID, natsURL string) (stan.Conn, error) {
	sc, err := stan.Connect(clusterID, clientID, stan.NatsURL(natsURL))
	if err != nil {
		return nil, err
	}

	return sc, nil
}

func (s *APIServer) handleAccount(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "GET" {
		path := r.URL.Path[len("/user/"):]
		if path == "cache" {
			return s.PeriodicCacheCheck(w, r)
		} else {
			return s.handleGetAccount(w, r)
		}
	}
	if r.Method == "DELETE" {
		return s.handleDeleteSegment(w, r)
	}
	/*if r.Method == "PUT" {
		return s.handleAddUserToSegment(w, r)
	}*/
	return fmt.Errorf("method now allowed %s", r.Method)

}

func (s *APIServer) handleGetAccount(w http.ResponseWriter, r *http.Request) error {
	uid := mux.Vars(r)["id"]

	user, err := s.store.GetUserById(uid)
	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, user)
}

func (s *APIServer) handleDeleteSegment(w http.ResponseWriter, r *http.Request) error {
	uid := mux.Vars(r)["id"]

	if err := s.store.DeleteUser(uid); err != nil {
		return err
	}
	return WriteJSON(w, http.StatusOK, "segment deleted sucsessfully")
}

func (s *APIServer) PeriodicCacheCheck(w http.ResponseWriter, r *http.Request) error {
	//cacheStr := mux.Vars(r)["cache"]

	for {
		// Check if the cache is empty
		if s.cache.IsEmpty() {
			return WriteJSON(w, http.StatusOK, "Cache is empty")
		} else {
			//return WriteJSON(w, http.StatusOK, user)
		}
	}

}

func (s *APIServer) handleDataFromStream() error {
	s.store.CreateUserFromNATS(s.stanConn.GetMsgChannel())
	return nil
}

func (s *APIServer) Run() {
	router := mux.NewRouter()
	router.HandleFunc("/getUser/{id}", makeHTTPHandleFunc(s.handleGetAccount))
	router.HandleFunc("/delUser/{id}", makeHTTPHandleFunc(s.handleDeleteSegment))
	router.HandleFunc("/user/cache", makeHTTPHandleFunc(s.PeriodicCacheCheck))
	http.ListenAndServe(s.listenAddr, router)
	log.Panicln("API server running on port", s.listenAddr)

}

type ApiError struct {
	Error string `json:"error"`
}

type apiFunc func(http.ResponseWriter, *http.Request) error

func makeHTTPHandleFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			WriteJSON(w, http.StatusBadRequest, ApiError{Error: err.Error()})
		}
	}
}

func WriteJSON(w http.ResponseWriter, status int, v any) error {

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(v)
}
