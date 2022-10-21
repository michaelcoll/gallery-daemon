-- name: GetPhoto :one
SELECT *
FROM photos
WHERE hash = ?;

-- name: CreateOrReplacePhoto :exec
REPLACE INTO photos (hash, path, date_time, iso, exposure_time, x_dimension, y_dimension, model, f_number)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?);

-- name: CountPhotoByHash :one
SELECT COUNT(*)
FROM photos
WHERE hash = ?;

-- name: CountPhotos :one
SELECT COUNT(*)
FROM photos;

-- name: List :many
SELECT *
FROM photos;

-- name: DeletePhotoByPath :exec
DELETE FROM photos
WHERE path = ?;

-- name: DeleteAllPhotoInPath :exec
DELETE FROM photos
WHERE path LIKE ?;

-- name: DeleteAllPhotos :exec
DELETE FROM photos
WHERE 1
