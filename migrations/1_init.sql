-- +migrate Up

CREATE TABLE color
(
    id SERIAL,
    color INTEGER
);

INSERT INTO color (color) VALUES (1);


-- +migrate Down
DROP TABLE color;
