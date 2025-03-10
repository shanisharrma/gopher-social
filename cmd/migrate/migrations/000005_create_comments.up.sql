CREATE TABLE IF NOT EXISTS comments (
id bigserial PRIMARY KEY,
  user_id bigserial NOT NULL,
  post_id bigserial NOT NULL,
  content TEXt NOT NULL,
  created_at timestamp(0) with time zone DEFAULT NOW()
);
