ALTER TABLE posts
ADD COLUMN tags varchar(100)[];

ALTER TABLE posts
ADD COLUMN updated_at timestamp(0) with time zone DEFAULT NOW();
