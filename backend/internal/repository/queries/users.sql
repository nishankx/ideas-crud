-- name: CreateUser :one
INSERT INTO users (email, password_hash, display_name, is_demo)
VALUES ($1, $2, $3, $4)
RETURNING id, email, password_hash, display_name, is_demo, created_at;

-- name: GetUserByEmail :one
SELECT id, email, password_hash, display_name, is_demo, created_at
FROM users
WHERE lower(email) = lower($1)
LIMIT 1;

-- name: GetUserByID :one
SELECT id, email, password_hash, display_name, is_demo, created_at
FROM users
WHERE id = $1
LIMIT 1;

-- name: GetDemoUser :one
SELECT id, email, password_hash, display_name, is_demo, created_at
FROM users
WHERE is_demo = TRUE
LIMIT 1;
