-- name: CreateTransfers :one
INSERT INTO transfers (
  from_account_id,
  to_account_id,
  amount
) VALUES (
  $1, $2, $3
)
RETURNING *;

-- name: GetTransfers :one
SELECT * FROM transfers
WHERE id = $1 
LIMIT 1;

-- name: ListTransfers :many
SELECT * FROM transfers
WHERE from_account_id = sqlc.arg(from_account_id)
   OR to_account_id = sqlc.arg(to_account_id)
ORDER BY id
LIMIT $1
OFFSET $2;


-- name: UpdateTransfers :exec
UPDATE transfers
set amount = $2
WHERE id = $1;

-- name: DeleteTransfers :exec
DELETE FROM transfers
WHERE id = $1;