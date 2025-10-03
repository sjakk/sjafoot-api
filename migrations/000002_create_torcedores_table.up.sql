CREATE TABLE IF NOT EXISTS torcedores (
    id bigserial PRIMARY KEY,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    nome text NOT NULL,
    email citext UNIQUE NOT NULL,
    time_clube text NOT NULL
);
