-- +goose Up
ALTER TABLE users
    ADD COLUMN role ENUM('user', 'admin', 'moderator') NOT NULL DEFAULT 'user',
    ADD COLUMN is_premium BOOLEAN DEFAULT FALSE,
    ADD COLUMN subscription_tier ENUM('free', 'basic', 'pro') NOT NULL DEFAULT 'free',
    ADD COLUMN timezone VARCHAR(50) DEFAULT 'UTC';

-- +goose Down
ALTER TABLE users
    DROP COLUMN role,
    DROP COLUMN is_premium,
    DROP COLUMN subscription_tier,
    DROP COLUMN timezone;