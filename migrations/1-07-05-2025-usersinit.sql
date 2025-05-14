CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) NOT NULL,
    login VARCHAR(255),
    password VARCHAR(255)
)