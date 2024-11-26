-- name: GetSongs :many
SELECT * FROM songs 
WHERE group_name ILIKE COALESCE($1, '%') 
OR song_name ILIKE COALESCE($2, '%') 
LIMIT $3 OFFSET $4;

-- name: GetSongWithFiltersAndPagination :one
SELECT id, group_name, song_name, release_date, link
FROM songs
WHERE 
  (group_name ILIKE '%' || $1 || '%' OR $1 IS NULL) AND
  (song_name ILIKE '%' || $2 || '%' OR $2 IS NULL) AND
  (release_date = $3 OR $3 IS NULL)
ORDER BY release_date DESC
LIMIT $4 OFFSET $5;

-- name: GetSongVersesWithPagination :one
WITH verses AS (
  SELECT unnest(string_to_array(text, E'\n\n')) AS verse
  FROM songs
  WHERE id = $1
)
SELECT verse
FROM verses
LIMIT $2 OFFSET $3;

-- name: InsertSong :exec
INSERT INTO songs (id, group_name, song_name, release_date, text, link)
VALUES ($1, $2, $3, $4, $5, $6);

-- name: UpdateSong :exec
UPDATE songs SET group_name = $2, song_name = $3, text = $4, release_date = $5, link = $6 WHERE id = $1;

-- name: DeleteSong :exec
DELETE FROM songs WHERE id = $1;
