CREATE DATABASE IF NOT EXISTS devbook_db;
USE devbook_db;

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

DROP TABLE IF EXISTS users;

CREATE TABLE users (
    id UUID PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(100) NOT NULL unique,
    password VARCHAR(100) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NULL
);

DROP TABLE IF EXISTS posts;

CREATE TABLE posts (
    id UUID PRIMARY KEY,
    title VARCHAR(100) NOT NULL,
    content TEXT NOT NULL,
    user_id UUID NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NULL
);

ALTER TABLE posts       
ADD CONSTRAINT fk_posts_user_id
FOREIGN KEY (user_id)
REFERENCES users (id);