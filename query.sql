-- name: CreateGame :exec
INSERT INTO games (theme_id, num_rounds) VALUES (?, ?);

-- name: GetLastInsertGame :one
SELECT id, theme_id, num_rounds, current_round FROM games WHERE id = last_insert_rowid();

-- name: GetGame :one
SELECT id, theme_id, num_rounds, current_round
FROM games
WHERE id = ?;

-- name: GetRoomPlayers :many
SELECT id, player_id, room_id
FROM room_players
WHERE room_id = ?;

-- name: CreateRoom :exec
INSERT INTO rooms (name, max_players, game_rounds) VALUES (?, ?, ?);

-- name: CreateRoomWithSlug :exec
INSERT INTO rooms (name, slug, max_players, game_rounds) VALUES (?, ?, ?, ?);

-- User queries
-- name: CreateUser :one
INSERT INTO users (name, password)
VALUES (?, ?)
RETURNING id, name, password;

-- name: GetUserByName :one
SELECT id, name, password
FROM users
WHERE name = ?;

-- name: GetUserByID :one
SELECT id, name, password
FROM users
WHERE id = ?;

-- Player queries
-- name: CreatePlayer :one
INSERT INTO players (user_id)
VALUES (?)
RETURNING id, user_id, score, game_history;

-- name: GetPlayerByUserID :one
SELECT id, user_id, score, game_history
FROM players
WHERE user_id = ?;


-- name: GetQuestions :many
SELECT id, theme_id, question_text, correct_answer
FROM questions
WHERE theme_id = ?;

-- name: GetThemeByName :one
SELECT id
FROM themes
WHERE name = ?;

-- name: ListThemes :many
SELECT id, name
FROM themes;
