
-- name: CreatePlayer :exec
INSERT INTO players (name) VALUES (?);

-- name: GetLastInsertPlayer :one
SELECT id, name FROM players WHERE id = last_insert_rowid();

-- name: CreateGame :exec
INSERT INTO games (theme_id, num_rounds) VALUES (?, ?);

-- name: GetLastInsertGame :one
SELECT id, theme_id, num_rounds, current_round FROM games WHERE id = last_insert_rowid();

-- name: GetGame :one
SELECT id, theme_id, num_rounds, current_round
FROM games
WHERE id = ?;
