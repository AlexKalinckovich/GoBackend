-- +goose Up
CREATE TABLE users (
                       id BIGINT PRIMARY KEY AUTO_INCREMENT,
                       name VARCHAR(100) NOT NULL,
                       surname VARCHAR(100) NOT NULL,
                       age INT CHECK (age >= 0 AND age <= 150),
                       country_code CHAR(2),
                       account_balance DECIMAL(10,2) DEFAULT 0.00
);

-- +goose Down
DROP TABLE users;