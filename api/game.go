package outsmarty_api

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	db "outsmarty.sqlc.dev/app/outsmarty"
)

func GetGameHandler(w http.ResponseWriter, r *http.Request) {
	gameIDStr := r.URL.Query().Get("game_id")
	gameID, err := strconv.Atoi(gameIDStr)
	if err != nil {
		http.Error(w, "Invalid game_id", http.StatusBadRequest)
		return
	}
	database, err := sql.Open("sqlite3", "./outsmarty.db")
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()
	queries := db.New(database)

	game, err := queries.GetGame(r.Context(), int64(gameID))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(game)
}

func GetThemesHandler(w http.ResponseWriter, r *http.Request) {
	database, err := sql.Open("sqlite3", "./outsmarty.db")
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()
	queries := db.New(database)
	themes, err := queries.ListThemes(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(themes)
}
