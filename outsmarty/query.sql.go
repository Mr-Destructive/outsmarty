// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: query.sql

package outsmarty

import (
	"context"
	"database/sql"
)

const createGame = `-- name: CreateGame :exec
INSERT INTO games (theme_id, num_rounds) VALUES (?, ?)
`

type CreateGameParams struct {
	ThemeID   int64 `json:"theme_id"`
	NumRounds int64 `json:"num_rounds"`
}

func (q *Queries) CreateGame(ctx context.Context, arg CreateGameParams) error {
	_, err := q.db.ExecContext(ctx, createGame, arg.ThemeID, arg.NumRounds)
	return err
}

const createPlayer = `-- name: CreatePlayer :one
INSERT INTO players (user_id)
VALUES (?)
RETURNING id, user_id, score, game_history
`

// Player queries
func (q *Queries) CreatePlayer(ctx context.Context, userID sql.NullInt64) (Player, error) {
	row := q.db.QueryRowContext(ctx, createPlayer, userID)
	var i Player
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Score,
		&i.GameHistory,
	)
	return i, err
}

const createRoom = `-- name: CreateRoom :exec
INSERT INTO rooms (name, max_players, game_rounds) VALUES (?, ?, ?)
`

type CreateRoomParams struct {
	Name       string `json:"name"`
	MaxPlayers int64  `json:"max_players"`
	GameRounds int64  `json:"game_rounds"`
}

func (q *Queries) CreateRoom(ctx context.Context, arg CreateRoomParams) error {
	_, err := q.db.ExecContext(ctx, createRoom, arg.Name, arg.MaxPlayers, arg.GameRounds)
	return err
}

const createRoomWithSlug = `-- name: CreateRoomWithSlug :exec
INSERT INTO rooms (name, slug, max_players, game_rounds) VALUES (?, ?, ?, ?)
`

type CreateRoomWithSlugParams struct {
	Name       string `json:"name"`
	Slug       string `json:"slug"`
	MaxPlayers int64  `json:"max_players"`
	GameRounds int64  `json:"game_rounds"`
}

func (q *Queries) CreateRoomWithSlug(ctx context.Context, arg CreateRoomWithSlugParams) error {
	_, err := q.db.ExecContext(ctx, createRoomWithSlug,
		arg.Name,
		arg.Slug,
		arg.MaxPlayers,
		arg.GameRounds,
	)
	return err
}

const createUser = `-- name: CreateUser :one
INSERT INTO users (name, password)
VALUES (?, ?)
RETURNING id, name, password
`

type CreateUserParams struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

// User queries
func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.db.QueryRowContext(ctx, createUser, arg.Name, arg.Password)
	var i User
	err := row.Scan(&i.ID, &i.Name, &i.Password)
	return i, err
}

const getGame = `-- name: GetGame :one
SELECT id, theme_id, num_rounds, current_round
FROM games
WHERE id = ?
`

func (q *Queries) GetGame(ctx context.Context, id int64) (Game, error) {
	row := q.db.QueryRowContext(ctx, getGame, id)
	var i Game
	err := row.Scan(
		&i.ID,
		&i.ThemeID,
		&i.NumRounds,
		&i.CurrentRound,
	)
	return i, err
}

const getLastInsertGame = `-- name: GetLastInsertGame :one
SELECT id, theme_id, num_rounds, current_round FROM games WHERE id = last_insert_rowid()
`

func (q *Queries) GetLastInsertGame(ctx context.Context) (Game, error) {
	row := q.db.QueryRowContext(ctx, getLastInsertGame)
	var i Game
	err := row.Scan(
		&i.ID,
		&i.ThemeID,
		&i.NumRounds,
		&i.CurrentRound,
	)
	return i, err
}

const getPlayerByUserID = `-- name: GetPlayerByUserID :one
SELECT id, user_id, score, game_history
FROM players
WHERE user_id = ?
`

func (q *Queries) GetPlayerByUserID(ctx context.Context, userID sql.NullInt64) (Player, error) {
	row := q.db.QueryRowContext(ctx, getPlayerByUserID, userID)
	var i Player
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Score,
		&i.GameHistory,
	)
	return i, err
}

const getQuestions = `-- name: GetQuestions :many
SELECT id, theme_id, question_text, correct_answer
FROM questions
WHERE theme_id = ?
`

func (q *Queries) GetQuestions(ctx context.Context, themeID int64) ([]Question, error) {
	rows, err := q.db.QueryContext(ctx, getQuestions, themeID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Question
	for rows.Next() {
		var i Question
		if err := rows.Scan(
			&i.ID,
			&i.ThemeID,
			&i.QuestionText,
			&i.CorrectAnswer,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getRoomPlayers = `-- name: GetRoomPlayers :many
SELECT id, player_id, room_id
FROM room_players
WHERE room_id = ?
`

type GetRoomPlayersRow struct {
	ID       int64 `json:"id"`
	PlayerID int64 `json:"player_id"`
	RoomID   int64 `json:"room_id"`
}

func (q *Queries) GetRoomPlayers(ctx context.Context, roomID int64) ([]GetRoomPlayersRow, error) {
	rows, err := q.db.QueryContext(ctx, getRoomPlayers, roomID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetRoomPlayersRow
	for rows.Next() {
		var i GetRoomPlayersRow
		if err := rows.Scan(&i.ID, &i.PlayerID, &i.RoomID); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getThemeByName = `-- name: GetThemeByName :one
SELECT id
FROM themes
WHERE name = ?
`

func (q *Queries) GetThemeByName(ctx context.Context, name string) (int64, error) {
	row := q.db.QueryRowContext(ctx, getThemeByName, name)
	var id int64
	err := row.Scan(&id)
	return id, err
}

const getUserByID = `-- name: GetUserByID :one
SELECT id, name, password
FROM users
WHERE id = ?
`

func (q *Queries) GetUserByID(ctx context.Context, id int64) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserByID, id)
	var i User
	err := row.Scan(&i.ID, &i.Name, &i.Password)
	return i, err
}

const getUserByName = `-- name: GetUserByName :one
SELECT id, name, password
FROM users
WHERE name = ?
`

func (q *Queries) GetUserByName(ctx context.Context, name string) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserByName, name)
	var i User
	err := row.Scan(&i.ID, &i.Name, &i.Password)
	return i, err
}

const listThemes = `-- name: ListThemes :many
SELECT id, name
FROM themes
`

func (q *Queries) ListThemes(ctx context.Context) ([]Theme, error) {
	rows, err := q.db.QueryContext(ctx, listThemes)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Theme
	for rows.Next() {
		var i Theme
		if err := rows.Scan(&i.ID, &i.Name); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
