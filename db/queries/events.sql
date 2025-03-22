-- name: CreateEvent :one
INSERT INTO events (
  name,
  type,
  api_end_point,
  api_method,
  api_request_body,
  created_by 
) VALUES (
  $1, $2, $3, $4, $5, $6
)
RETURNING *;

-- name: GetEvent :one
SELECT * FROM events
WHERE id = $1;

-- name: ListEvents :many
SELECT * FROM events
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: ListUserEvents :many
SELECT * FROM events
WHERE created_by = $1
ORDER BY id
LIMIT $2
OFFSET $3;

-- name: DeleteEvent :exec
DELETE FROM events
WHERE id = $1;