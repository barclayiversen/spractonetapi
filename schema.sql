--CREATE DATABASE spractonet;

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255),
    password VARCHAR(255),
    username VARCHAR(255),
    created_at BIGINT,
    activated BOOLEAN,
    activation_key VARCHAR(255)
);

CREATE TABLE posts (
    id SERIAL PRIMARY KEY,
    title VARCHAR(64),
    post VARCHAR(2000),
    user_id INT REFERENCES users (id),
    created_at BIGINT
);

CREATE TABLE comments (
    id SERIAL PRIMARY KEY,
    comment VARCHAR(300),
    post_id INT REFERENCES posts(id)
);