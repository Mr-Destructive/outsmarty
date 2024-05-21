package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
	db "outsmarty.sqlc.dev/app/outsmarty"
)

var queries *db.Queries

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

	http.HandleFunc("/players", createPlayerHandler)
	http.HandleFunc("/games", createGameHandler)
	http.HandleFunc("/games/rounds", startRoundHandler)
	http.HandleFunc("/games/answers", submitAnswerHandler)
	http.HandleFunc("/games/status", getGameStatusHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func initializeDatabase(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS players (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL
		);
		CREATE TABLE IF NOT EXISTS themes (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL
		);
		CREATE TABLE IF NOT EXISTS questions (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			theme_id INTEGER NOT NULL,
			question_text TEXT NOT NULL,
			correct_answer TEXT NOT NULL,
			FOREIGN KEY (theme_id) REFERENCES themes(id)
		);
		CREATE TABLE IF NOT EXISTS answers (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			question_id INTEGER NOT NULL,
			player_id INTEGER NOT NULL,
			answer_text TEXT NOT NULL,
			is_correct BOOLEAN NOT NULL DEFAULT FALSE,
			FOREIGN KEY (question_id) REFERENCES questions(id),
			FOREIGN KEY (player_id) REFERENCES players(id)
		);
		CREATE TABLE IF NOT EXISTS games (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			theme_id INTEGER NOT NULL,
			num_rounds INTEGER NOT NULL,
			current_round INTEGER NOT NULL DEFAULT 0,
			FOREIGN KEY (theme_id) REFERENCES themes(id)
		);
		CREATE TABLE IF NOT EXISTS game_rounds (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			game_id INTEGER NOT NULL,
			round_number INTEGER NOT NULL,
			question_id INTEGER NOT NULL,
			FOREIGN KEY (game_id) REFERENCES games(id),
			FOREIGN KEY (question_id) REFERENCES questions(id)
		);
		CREATE TABLE IF NOT EXISTS game_players (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			game_id INTEGER NOT NULL,
			player_id INTEGER NOT NULL,
			points INTEGER NOT NULL DEFAULT 0,
			FOREIGN KEY (game_id) REFERENCES games(id),
			FOREIGN KEY (player_id) REFERENCES players(id)
		);
	`)
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

		err := queries.CreatePlayer(r.Context(), req.Name)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		player, err := queries.GetLastInsertPlayer(r.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(player)
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
