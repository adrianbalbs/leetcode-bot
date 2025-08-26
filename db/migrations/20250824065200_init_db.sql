-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS difficulties (
    id SERIAL PRIMARY KEY,
    name TEXT UNIQUE
);


CREATE TABLE IF NOT EXISTS problems (
    id SERIAL PRIMARY KEY,
    slug TEXT UNIQUE,
    title TEXT,
    difficulty_id INT REFERENCES difficulties(id)
);

CREATE TABLE IF NOT EXISTS tags (
      id SERIAL PRIMARY KEY,
      name TEXT UNIQUE
);

CREATE TABLE IF NOT EXISTS problem_tags (
    problem_id SERIAL REFERENCES problems(id),
    tag_id BIGINT REFERENCES tags(id),
    PRIMARY KEY (problem_id, tag_id)
);

CREATE TABLE IF NOT EXISTS playlists (
    id SERIAL PRIMARY KEY,
    name TEXT UNIQUE,
    creator TEXT
);

CREATE TABLE IF NOT EXISTS playlist_entries (
  playlist_id SERIAL REFERENCES playlists(id),
  problem_id SERIAL REFERENCES problems(id),
  PRIMARY KEY (playlist_id, problem_id)
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE playlist_entries;
DROP TABLE playlists;
DROP TABLE problem_tags;
DROP TABLE tags;
DROP TABLE problems;
DROP TABLE difficulties;
-- +goose StatementEnd

