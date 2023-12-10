package http

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"io"
	"log"
	"net/http"
	"otus_msa_docker/api"
	"otus_msa_docker/internal/users"
)

type userServer struct {
	usersStore users.Store
}

func NewServer(store users.Store) (api.ServerInterface, error) {
	return &userServer{
		usersStore: store,
	}, nil
}

func RunServer(s api.ServerInterface) {
	h := api.Handler(s)
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Mount("/", h)
	r.HandleFunc("/health", health)

	http.ListenAndServe(":8000", r)
}

func (s *userServer) CreateUser(w http.ResponseWriter, r *http.Request) {
	var u users.User
	ctx := r.Context()
	b, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(500)
		return
	}
	err = json.Unmarshal(b, &u)
	if err != nil {
		w.WriteHeader(500)
		return
	}

	err = s.usersStore.CreateUser(ctx, &u)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(500)
		return
	}
}

func (s *userServer) DeleteUser(w http.ResponseWriter, r *http.Request, userId int64) {
	ctx := r.Context()
	err := s.usersStore.DeleteUser(ctx, userId)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(500)
		return
	}
}

func (s *userServer) FindUserById(w http.ResponseWriter, r *http.Request, userId int64) {
	ctx := r.Context()
	u, err := s.usersStore.GetUser(ctx, userId)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}
	d, err := json.Marshal(u)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}
	w.Write(d)
}

func (s *userServer) UpdateUser(w http.ResponseWriter, r *http.Request, userId int64) {
	var u users.User
	ctx := r.Context()
	b, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(500)
		return
	}
	err = json.Unmarshal(b, &u)
	if err != nil {
		w.WriteHeader(500)
		return
	}

	err = s.usersStore.UpdateUser(ctx, &u)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(500)
		return
	}
}

func health(w http.ResponseWriter, req *http.Request) {
	resp := struct {
		Status string `json:"status"`
	}{
		Status: "OK",
	}

	w.Header().Set("Content-Type", "application/json")
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	}
	w.Write(jsonResp)
}
