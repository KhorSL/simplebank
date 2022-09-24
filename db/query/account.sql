-- name: CreateAccount :one
INSERT INTO ACCOUNTS (
  OWNER,
  BALANCE,
  CURRENCY
) VALUES (
  $1, $2, $3
) RETURNING *;

-- name: GetAccount :one
SELECT * FROM ACCOUNTS
WHERE ID = $1 LIMIT 1;

-- name: GetAccountForUpdate :one
SELECT * FROM ACCOUNTS
WHERE ID = $1 LIMIT 1
FOR NO KEY UPDATE;

-- name: ListAccounts :many
SELECT * FROM ACCOUNTS
WHERE OWNER = $1
ORDER BY ID
LIMIT $2
OFFSET $3;

-- name: UpdateAccount :one
UPDATE ACCOUNTS 
SET BALANCE = $2
WHERE ID = $1
RETURNING *;

-- name: AddAccountBalance :one
UPDATE ACCOUNTS 
SET BALANCE = BALANCE + sqlc.arg(amount)
WHERE ID = sqlc.arg(id)
RETURNING *;

-- name: DeleteAccount :exec
DELETE FROM ACCOUNTS
WHERE ID = $1;