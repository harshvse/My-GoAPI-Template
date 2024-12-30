CREATE TABLE IF NOT EXISTS user_invitation (
    token bytea NOT NULL,
    user_id bigint NOT NULL,
    PRIMARY KEY (token, user_id)
);
