-- +migrate Up
ALTER TABLE
    authors
ADD
    COLUMN phone varchar(20);

-- +migrate Down
ALTER TABLE
    authors DROP COLUMN phone;