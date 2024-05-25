package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type tokenResponse struct {
	Me       string `json:"me"`
	IssuedBy string `json:"issued_by"`
	ClientId string `json:"client_id"`
	IssuedAt int    `json:"issued_at"`
	Scope    string `json:"scope"`
	Nonce    int    `json:"nonce"`
}

func (s *Server) RegisterRoutes() http.Handler {

	mux := http.NewServeMux()
	mux.HandleFunc("POST /micropub", s.HelloWorldHandler)

	return mux
}

func (s *Server) HelloWorldHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()

	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var data = r.Form

	if !data.Has("access_token") {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	accessToken := data.Get("access_token")
	if !isTokenValid(accessToken) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// ACTUALLY CREATE POST HERE

	w.Header().Add("Location", fmt.Sprintf("https://danielmcfarland.dev/%v", ""))
	w.WriteHeader(http.StatusCreated)

	_, _ = w.Write(nil)
}

func isTokenValid(token string) bool {
	url := "https://tokens.indieauth.com/token"
	method := "GET"

	client := &http.Client{}

	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return false
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %v", token))
	req.Header.Add("Accept", "application/json")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return false
	}
	defer res.Body.Close()

	var data tokenResponse
	err = json.NewDecoder(res.Body).Decode(&data)
	if err != nil {
		fmt.Println(err)
		return false
	}

	return data.Me == "http://danielmcfarland.dev/"
}
