-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS difficulties (
    id INT PRIMARY KEY,
    name TEXT UNIQUE
);


CREATE TABLE IF NOT EXISTS problems (
    id INT PRIMARY KEY,
    slug TEXT UNIQUE,
    title TEXT,
    difficulty_id INT REFERENCES difficulties(id)
);

CREATE TABLE IF NOT EXISTS tags (
      id INT PRIMARY KEY,
      name TEXT UNIQUE
);

CREATE TABLE IF NOT EXISTS problem_tags (
    problem_id INT REFERENCES problems(id),
    tag_id INT REFERENCES tags(id),
    PRIMARY KEY (problem_id, tag_id)
);

CREATE TABLE IF NOT EXISTS playlists (
    id INT PRIMARY KEY,
    name TEXT,
    creator TEXT
);

CREATE TABLE IF NOT EXISTS playlist_entries (
  playlist_id INT REFERENCES playlists(id),
  problem_id INT REFERENCES problems(id),
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

