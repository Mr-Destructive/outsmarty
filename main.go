package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/net/websocket"
	outsmarty_api "outsmarty.sqlc.dev/app/api"
	db "outsmarty.sqlc.dev/app/outsmarty"
)

var queries *db.Queries

type PublicUser struct {
	ID   int32  `json:"id"`
	Name string `json:"name"`
}

func main() {
	database, err := sql.Open("sqlite3", "./outsmarty.db")
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()

	queries = db.New(database)

	err = initializeDatabase(database)
	if err != nil {
		log.Fatal(err)
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/auth/register", registerHandler)
	mux.HandleFunc("/auth/login", loginHandler)
	mux.HandleFunc("/auth/logout", logoutHandler)

	mux.HandleFunc("/rooms/create", createRoomHandler)
	mux.HandleFunc("/rooms/join", joinRoomHandler)
	mux.HandleFunc("/players", createPlayerHandler)

	mux.HandleFunc("/games", createGameHandler)
	mux.HandleFunc("/games/rounds", startRoundHandler)
	mux.HandleFunc("/games/answers", submitAnswerHandler)
	mux.HandleFunc("/games/status", getGameStatusHandler)

	mux.HandleFunc("/app", appHandler)
	chatServer := outsmarty_api.NewServer()
	mux.Handle("/ws/room/", websocket.Handler(chatServer.HandleWS))
	corsHandler := outsmarty_api.CORSMiddleware(mux)

	log.Fatal(http.ListenAndServe(":8080", corsHandler))
}

func initializeDatabase(db *sql.DB) error {
	// open and read schema.sql
	ddl, err := os.ReadFile("./schema.sql")
	combinedDDL := string(ddl)
	if err != nil {
		return err
	}
	_, err = db.Exec(combinedDDL)
	return err
}

type CreatePlayerRequest struct {
	Name string `json:"name"`
}

func createPlayerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {

		var req CreatePlayerRequest
		log.Println(r.Body)
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

type CreateGameRequest struct {
	ThemeID   int `json:"theme_id"`
	NumRounds int `json:"num_rounds"`
}

func createGameHandler(w http.ResponseWriter, r *http.Request) {
	var req CreateGameRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err := queries.CreateGame(r.Context(), db.CreateGameParams{
		ThemeID:   int64(req.ThemeID),
		NumRounds: int64(req.NumRounds),
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	game, err := queries.GetLastInsertGame(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(game)
}

type StartRoundRequest struct {
	GameID      int `json:"game_id"`
	RoundNumber int `json:"round_number"`
}

func startRoundHandler(w http.ResponseWriter, r *http.Request) {
	var req StartRoundRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Implement logic to start a new round
}

type SubmitAnswerRequest struct {
	GameID     int    `json:"game_id"`
	PlayerID   int    `json:"player_id"`
	Answer     string `json:"answer"`
	QuestionID int    `json:"question_id"`
}

func submitAnswerHandler(w http.ResponseWriter, r *http.Request) {
	var req SubmitAnswerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Implement logic to submit an answer
}

func getGameStatusHandler(w http.ResponseWriter, r *http.Request) {
	gameIDStr := r.URL.Query().Get("game_id")
	gameID, err := strconv.Atoi(gameIDStr)
	if err != nil {
		http.Error(w, "Invalid game_id", http.StatusBadRequest)
		return
	}

	game, err := queries.GetGame(r.Context(), int64(gameID))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(game)
}

func createRoomHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		body := r.Body
		validatedBody, err := validateRoomPayload(body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		fmt.Println(validatedBody)
		roomSlug := generateRoomSlug(validatedBody.Name)
		roomObject := db.CreateRoomWithSlugParams{
			Name:       validatedBody.Name,
			MaxPlayers: validatedBody.MaxPlayers,
			GameRounds: validatedBody.GameRounds,
			Slug:       roomSlug,
		}
		err = queries.CreateRoomWithSlug(r.Context(), roomObject)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(roomObject)
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}

}

func validateRoomPayload(body io.ReadCloser) (db.CreateRoomParams, error) {
	createGameParams := db.CreateRoomParams{}
	err := json.NewDecoder(body).Decode(&createGameParams)
	if err != nil {
		return createGameParams, err
	}
	return createGameParams, nil

}

func generateRoomSlug(name string) string {
	// add a random string to the end of the slug
	slug := name + "_" + uuid.New().String()[0:8]
	return slug
}

func joinRoomHandler(w http.ResponseWriter, r *http.Request) {
}

func appHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		http.ServeFile(w, r, "mockup/index.html")
	}
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var user db.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Error hashing password", http.StatusInternalServerError)
		return
	}
	ctx := context.Background()
	newUser, err := queries.CreateUser(ctx, db.CreateUserParams{
		Name:     user.Name,
		Password: string(hashedPassword),
	})
	if err != nil {
		http.Error(w, "Error creating user", http.StatusInternalServerError)
		return
	}
	userID := sql.NullInt64{Int64: newUser.ID, Valid: true}
	_, err = queries.CreatePlayer(ctx, userID)
	if err != nil {
		http.Error(w, "Error creating player", http.StatusInternalServerError)
		return
	}

	response := PublicUser{
		Name: newUser.Name,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var name, password string

	contentType := r.Header.Get("Content-Type")

	switch contentType {
	case "application/json":
		var user db.User
		if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
			http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
			return
		}
		name = user.Name
		password = user.Password
	case "application/x-www-form-urlencoded":
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Invalid form request payload", http.StatusBadRequest)
			return
		}
		name = r.Form.Get("name")
		password = r.Form.Get("password")
	default:
		http.Error(w, "Unsupported content type", http.StatusUnsupportedMediaType)
		return
	}
	fmt.Println(name, password)
	ctx := context.Background()
	dbUser, err := queries.GetUserByName(ctx, name)
	if err != nil {
		http.Error(w, "Invalid name or password", http.StatusUnauthorized)
		return
	}
	name = dbUser.Name
	hashedPassword := dbUser.Password

	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)); err != nil {
		fmt.Println("Error comparing hash and password: ", err)
		return
	}
	sessionID, err := getSession(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	userID := strconv.Itoa(int(dbUser.ID))
	path := "/"
	if sessionID == nil {
		cookie := &http.Cookie{
			Domain:   "localhost",
			Name:     "outsmarty_uid",
			Value:    userID,
			Path:     path,
			MaxAge:   86400,
			SameSite: http.SameSiteNoneMode,
			Secure:   true,
		}
		userName := &http.Cookie{
			Domain:   "localhost",
			Name:     "outsmarty_name",
			Value:    name,
			Path:     path,
			MaxAge:   86400,
			SameSite: http.SameSiteNoneMode,
			Secure:   true,
		}
		http.SetCookie(w, cookie)
        http.SetCookie(w, userName)
	}

	response := PublicUser{
		Name: dbUser.Name,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	cookie := &http.Cookie{
		Name:   "outsmarty_uid",
		Value:  "",
		MaxAge: -1,
	}
	http.SetCookie(w, cookie)
	w.WriteHeader(http.StatusOK)
}

func getSession(r *http.Request) (*http.Cookie, error) {
	if cookie, err := r.Cookie("outsmarty_uid"); err == nil {
		return cookie, nil
	}
	return nil, nil
}
