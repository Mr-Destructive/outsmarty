
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

-- name: GetRoomPlayers :many
SELECT id, player_id, room_id
FROM room_players
WHERE room_id = ?;

-- name: CreateRoom :exec
INSERT INTO rooms (name, max_players, game_rounds) VALUES (?, ?, ?);

-- name: CreateRoomWithSlug :exec
INSERT INTO rooms (name, slug, max_players, game_rounds) VALUES (?, ?, ?, ?);
