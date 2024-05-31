package server

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

type tokenResponse struct {
	Me       string `json:"me"`
	IssuedBy string `json:"issued_by"`
	ClientId string `json:"client_id"`
	IssuedAt int    `json:"issued_at"`
	Scope    string `json:"scope"`
	Nonce    int    `json:"nonce"`
}

type GitHubCommitter struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type GitHubCommit struct {
	Message   string          `json:"message"`
	Committer GitHubCommitter `json:"committer"`
	Content   string          `json:"content"`
}

type DataEntryProperties struct {
	Published  []string `json:"published"`
	BookmarkOf []string `json:"bookmark-of"`
	PostStatus []string `json:"post-status"`
}

type DataEntry struct {
	Date       string              `json:"date"`
	Deleted    bool                `json:"deleted"`
	Draft      bool                `json:"draft"`
	H          string              `json:"h"`
	Properties DataEntryProperties `json:"properties"`
	Type       string              `json:"type"`
	Slug       string              `json:"slug"`
	ClientId   string              `json:"client_id"`
}

func (s *Server) RegisterRoutes() http.Handler {

	mux := http.NewServeMux()
	mux.HandleFunc("POST /micropub", s.MicropubHandler)
	mux.HandleFunc("GET /auth", s.AuthHandler)
	mux.HandleFunc("GET /auth/callback", s.AuthCallbackHandler)

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

	if !data.Has("h") {
		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	switch data.Get("h") {
	case "bookmark":
		if !data.Has("bookmark-of") {
			w.WriteHeader(http.StatusUnprocessableEntity)
			return
		}

		w.Header().Add("Location", createBookmark(data.Get("bookmark-of")))
		w.WriteHeader(http.StatusCreated)
	}

	_, _ = w.Write(nil)
}

func isTokenValid(token string) bool {
	tokenUrl := "https://tokens.indieauth.com/token"
	method := "GET"

	client := &http.Client{}

	req, err := http.NewRequest(method, tokenUrl, nil)

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

	return data.Me == "https://danielmcfarland.dev/"
}

func (s *Server) AuthHandler(w http.ResponseWriter, r *http.Request) {
	app_url := os.Getenv("APP_URL")

	me := "https://danielmcfarland.dev"
	state := randSeq(5)
	redirect_uri := fmt.Sprintf("%v/auth/callback", app_url)
	client_id := app_url
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

	resp := make(map[string]string)
	resp["access_token"] = access_token

	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("error handling JSON marshal. Err: %v", err)
	}

	w.Header().Add("Content-Type", "application/json")

	_, _ = w.Write(jsonResp)
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func createBookmark(bookmarkUrl string) string {
	slug := randSeq(5)
	date := time.Now()
	month := date.Month()
	year := date.Year()

	githubUrl := fmt.Sprintf("https://api.github.com/repos/danielmcfarland/danielmcfarland.dev/contents/content/data/%04d/%02d/%v.md", year, month, slug)
	method := "PUT"

	entryContent := DataEntry{
		Date:    date.Format(time.RFC3339),
		Deleted: false,
		Draft:   false,
		H:       "h-entry",
		Properties: DataEntryProperties{
			Published:  []string{date.Format(time.RFC3339)},
			BookmarkOf: []string{bookmarkUrl},
			PostStatus: []string{"published"},
		},
		Type:     "bookmarks",
		Slug:     fmt.Sprintf("2024/05/%v", slug),
		ClientId: "https://app.danielmcfarland.dev",
	}

	entryContentString, err := json.Marshal(entryContent)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	commitPayload := GitHubCommit{
		Message: "Adding Bookmark",
		Committer: GitHubCommitter{
			Name:  "Daniel McFarland",
			Email: "micropub@danielmcfarland.dev",
		},
		Content: base64.StdEncoding.EncodeToString(entryContentString),
	}

	commitPayloadString, err := json.Marshal(commitPayload)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	payload := strings.NewReader(string(commitPayloadString))

	client := &http.Client{}
	req, err := http.NewRequest(method, githubUrl, payload)

	if err != nil {
		fmt.Println(err)
		return ""
	}
	req.Header.Add("Accept", "application/vnd.github+json")
	req.Header.Add("X-GitHub-Api-Version", "2022-11-28")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %v", os.Getenv("GITHUB_API_KEY")))

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	defer res.Body.Close()

	return fmt.Sprintf("https://danielmcfarland.dev/data/%04d/%02d/%v", year, month, slug)
}
