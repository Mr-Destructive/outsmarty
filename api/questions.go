package outsmarty_api

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	db "outsmarty.sqlc.dev/app/outsmarty"
)

func GenerateQuestionsHandler(w http.ResponseWriter, r *http.Request) {
	theme := r.URL.Query().Get("theme")
	rounds := r.URL.Query().Get("rounds")

	numOfRounds, err := strconv.Atoi(rounds)
	if err != nil {
		log.Fatal(err)
	}
	questions := getQuestions(theme, numOfRounds)
	json.NewEncoder(w).Encode(questions)

}

func getQuestions(theme string, numOfRounds int) []db.Question {
	database, err := sql.Open("sqlite3", "./outsmarty.db")
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()
	queries := db.New(database)
	themeID, err := queries.GetThemeByName(context.Background(), theme)
	if err != nil {
		log.Printf("Error: %v", err)
		return []db.Question{}
	}
	questions, err := queries.GetQuestions(context.Background(), themeID)
	if err != nil {
		log.Printf("Error: %v", err)
		return []db.Question{}
	}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	r.Shuffle(len(questions), func(i, j int) {
		questions[i], questions[j] = questions[j], questions[i]
	})
	questions = questions[:numOfRounds]
	return questions
}
