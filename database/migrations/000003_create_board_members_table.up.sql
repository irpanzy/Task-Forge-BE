CREATE TABLE IF NOT EXISTS board_members (
    board_id  BIGINT NOT NULL,
    user_id   BIGINT NOT NULL,
    joined_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
    PRIMARY KEY (board_id, user_id),
    CONSTRAINT fk_board_members_board_id FOREIGN KEY (board_id) REFERENCES boards(internal_id) ON DELETE CASCADE,
    CONSTRAINT fk_board_members_user_id FOREIGN KEY (user_id) REFERENCES users(internal_id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_board_members_board_id ON board_members(board_id);
CREATE INDEX IF NOT EXISTS idx_board_members_user_id ON board_members(user_id);
