-- name: CreateLog :one
INSERT INTO logs (
  event_id,
  executed_on,
  status
) VALUES (
  $1, $2, $3
)
RETURNING *;

-- name: GetLog :one
SELECT * FROM logs
WHERE id = $1;

-- name: ListLogs :many
SELECT * FROM logs
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: DeleteLog :exec
DELETE FROM logs
WHERE id = $1;