-- +goose Up
CREATE TABLE songs (
  id UUID PRIMARY KEY,
  group_name TEXT NOT NULL,
  song_name TEXT NOT NULL,
  release_date DATE,
  text TEXT,
  link TEXT
);

-- +goose Down
DROP TABLE IF EXISTS songs;
