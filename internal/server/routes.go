package server

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
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
	mux.HandleFunc("POST /micropub", s.MicropubHandler)
	mux.HandleFunc("GET /auth", s.AuthHandler)
	mux.HandleFunc("GET /auth/callback", s.AuthCallbackHandler)
	mux.HandleFunc("GET /auth/token", s.AuthTokenHandler)

	return mux
}

func (s *Server) MicropubHandler(w http.ResponseWriter, r *http.Request) {
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

func (s *Server) AuthHandler(w http.ResponseWriter, r *http.Request) {
	app_url := os.Getenv("APP_URL")

	me := "https://danielmcfarland.dev"
	redirect_uri := fmt.Sprintf("%v/auth/callback", app_url)
	client_id := app_url
	state := "abc123"
	scope := "create"
	response_type := "code"

	redirect_url := fmt.Sprintf("https://indieauth.com/auth?me=%v&redirect_uri=%v&client_id=%v&state=%v&scope=%v&response_type=%v", me, redirect_uri, client_id, state, scope, response_type)

	http.Redirect(w, r, redirect_url, http.StatusPermanentRedirect)
}

func (s *Server) AuthCallbackHandler(w http.ResponseWriter, r *http.Request) {
	params, _ := url.ParseQuery(r.URL.RawQuery)

	app_url := os.Getenv("APP_URL")

	grant_type := "authorization_code"
	me := params["me"][0]
	code := params["code"][0]

	redirect_uri := fmt.Sprintf("%v/auth/callback", app_url)
	client_id := app_url

	data := url.Values{}
	data.Set("grant_type", grant_type)
	data.Set("me", me)
	data.Set("code", code)
	data.Set("redirect_uri", redirect_uri)
	data.Set("client_id", client_id)

	token_url := "https://tokens.indieauth.com/token"
	method := "POST"

	client := &http.Client{}

	req, err := http.NewRequest(method, token_url, strings.NewReader(data.Encode()))

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	bodyString := string(bodyBytes)

	tokenParams, _ := url.ParseQuery(bodyString)

	access_token := tokenParams["access_token"][0]

	fmt.Println(access_token)

	resp := make(map[string]string)
	resp["access_token"] = access_token

	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("error handling JSON marshal. Err: %v", err)
	}

	w.Header().Add("Content-Type", "application/json")

	_, _ = w.Write(jsonResp)
}

func (s *Server) AuthTokenHandler(w http.ResponseWriter, r *http.Request) {
	params, _ := url.ParseQuery(r.URL.RawQuery)

	for k, v := range params {
		fmt.Println(fmt.Sprintf("key[%s] value[%s]\n", k, v))
	}

	//fmt.Println(params)
}
