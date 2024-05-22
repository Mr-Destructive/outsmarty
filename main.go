package main

import (
	"database/sql"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

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
	http.HandleFunc("/rooms/create", createRoomHandler)
	http.HandleFunc("/rooms/join", joinRoomHandler)

	http.HandleFunc("/players", createPlayerHandler)
	http.HandleFunc("/games", createGameHandler)
	http.HandleFunc("/games/rounds", startRoundHandler)
	http.HandleFunc("/games/answers", submitAnswerHandler)
	http.HandleFunc("/games/status", getGameStatusHandler)

	log.Fatal(http.ListenAndServe(":8080", nil))
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

func createRoomHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		body := r.Body
		validatedBody, err := validateRoomPayload(body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
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
	slug := name + strconv.FormatInt(time.Now().UnixNano(), 10)
	return slug
}

func joinRoomHandler(w http.ResponseWriter, r *http.Request) {
}
