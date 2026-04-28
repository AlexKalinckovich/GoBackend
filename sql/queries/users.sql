-- name: CreateUser :exec
INSERT INTO users (
    name,
    surname,
    age,
    country_code,
    role,
    is_premium,
    subscription_tier,
    account_balance,
    timezone
) VALUES (
             ?, ?, ?, ?, ?, ?, ?, ?, ?
         );

-- name: CreateUserWithID :execresult
INSERT INTO users (
    name,
    surname,
    age,
    country_code,
    role,
    is_premium,
    subscription_tier,
    account_balance,
    timezone
) VALUES (
             ?, ?, ?, ?, ?, ?, ?, ?, ?
         );

-- name: GetUser :one
SELECT * FROM users WHERE id = ? LIMIT 1;

-- name: GetUserByUsername :one
SELECT * FROM users WHERE name = ? AND surname = ? LIMIT 1;

-- name: UpdateUserRole :exec
UPDATE users SET role = ? WHERE id = ?;

-- name: UpdateUser :exec
UPDATE users
SET name = ?,
    surname = ?,
    age = ?,
    country_code = ?,
    timezone = ?
WHERE id = ?;

-- name: DeleteUser :exec
DELETE FROM users WHERE id = ?;