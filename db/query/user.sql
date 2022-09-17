-- name: CreateUser :one
INSERT INTO USERS (
  USERNAME,
  HASHED_PASSWORD,
  FULL_NAME,
  EMAIL
) VALUES (
  $1, $2, $3, $4
) RETURNING *;

-- name: GetUserById :one
SELECT * FROM USERS
WHERE ID = $1 LIMIT 1;

-- name: GetUserByUsername :one
SELECT * FROM USERS
WHERE USERNAME = $1 LIMIT 1;
