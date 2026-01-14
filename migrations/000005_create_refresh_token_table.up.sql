BEGIN;

CREATE TABLE IF NOT EXISTS REFRESH_TOKENS (
    id BIGSERIAL PRIMARY KEY NOT NULL,
    user_id BIGINT NOT NULL,
    token_hash CHAR(64) NOT NULL UNIQUE,
    expires_at TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    revoked_at TIMESTAMPTZ NULL,
    replaced_by_token_id BIGINT NULL,
    CONSTRAINT fk_user_refresh FOREIGN KEY(user_id) REFERENCES USERS(user_id),
    CONSTRAINT fk_replaced_refresh FOREIGN KEY(replaced_by_token_id) REFERENCES REFRESH_TOKENS(id)
);

CREATE INDEX idx_refresh_tokens_user_id ON REFRESH_TOKENS(user_id);
CREATE INDEX idx_refresh_tokens_token_hash ON REFRESH_TOKENS(token_hash);

COMMIT;