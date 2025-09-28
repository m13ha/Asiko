CREATE TABLE IF NOT EXISTS ban_list_entries (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    banned_email TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(user_id, banned_email)
);
CREATE INDEX IF NOT EXISTS idx_ban_list_user_id ON ban_list_entries(user_id);
