-- name: GetPhoto :one
SELECT hash,
       path,
       date_time,
       iso,
       exposure_time,
       x_dimension,
       y_dimension,
       model,
       f_number,
       orientation
FROM photos
WHERE hash = ?;

-- name: GetThumbnail :one
SELECT thumbnail
FROM thumbnails
WHERE hash = ?
  AND width = ?
  AND height = ?;

-- name: CreateOrReplaceThumbnail :exec
REPLACE INTO thumbnails (hash, height, width, thumbnail)
VALUES (?, ?, ?, ?);

-- name: CreateOrReplacePhoto :exec
REPLACE INTO photos (hash, path, date_time, iso, exposure_time, x_dimension, y_dimension, model, f_number, orientation)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?);

-- name: CountPhotoByHash :one
SELECT COUNT(*)
FROM photos
WHERE hash = ?;

-- name: CountPhotos :one
SELECT COUNT(*)
FROM photos;

-- name: List :many
SELECT hash,
       path,
       date_time,
       iso,
       exposure_time,
       x_dimension,
       y_dimension,
       model,
       f_number,
       orientation
FROM photos
ORDER BY date_time DESC
LIMIT ? OFFSET ?;

-- name: DeletePhotoByPath :exec
DELETE
FROM photos
WHERE path = ?;

-- name: DeleteAllPhotoInPath :exec
DELETE
FROM photos
WHERE path LIKE ?;

-- name: DeleteAllPhotos :exec
DELETE
FROM photos
WHERE 1
