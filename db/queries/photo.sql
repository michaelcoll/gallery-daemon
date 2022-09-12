-- name: GetPhoto :one
SELECT *
FROM photos
WHERE hash = ?;

-- name: CreatePhoto :exec
INSERT INTO photos (hash, path, date_time, iso, exposure_time, x_dimension, y_dimension, model, focal_length)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
RETURNING *;
