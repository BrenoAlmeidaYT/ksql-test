-- +migrate Up
CREATE TABLE authors (
    id serial primary key,
    name varchar(100)
);

CREATE TABLE videos (
    id serial primary key,
    title varchar(100),
    description text,
    like_count integer,
    view_count integer,
    author_id integer references authors(id)
);

-- +migrate Down
DROP TABLE videos;

DROP TABLE authors;