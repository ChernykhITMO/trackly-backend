CREATE TABLE users
(
    id            SERIAL PRIMARY KEY,
    email         TEXT UNIQUE NOT NULL,
    username      VARCHAR,
    date_of_birth DATE,
    password      TEXT        NOT NULL,
    country       VARCHAR,
    city          VARCHAR,
    avatar_id     VARCHAR DEFAULT NULL
);

