-- +goose Up
CREATE TABLE feed_follows (
  id UUID PRIMARY KEY,
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL,
  user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  feed_id UUID NOT NULL REFERENCES feeds(id) ON DELETE CASCADE,
  UNIQUE (user_id, feed_id)
);

INSERT INTO feed_follows (id, created_at, updated_at, user_id, feed_id)
SELECT gen_random_uuid(), NOW(), NOW(), f.user_id, f.id
  FROM feeds f;

ALTER TABLE feeds
DROP COLUMN user_id;

-- +goose Down
ALTER TABLE feeds
ADD COLUMN user_id UUID;

UPDATE feeds SET user_id = ff.id
  FROM feed_follows ff
    WHERE ff.feed_id = feeds.id;

UPDATE feeds SET user_id = '00000000-0000-0000-0000-000000000000'
  WHERE user_id IS NULL;

ALTER TABLE feeds
ALTER COLUMN user_id SET NOT NULL;

ALTER TABLE feeds
ADD CONSTRAINT feeds_user_id_fk FOREIGN KEY (user_id)
REFERENCES users(id) ON DELETE CASCADE;

DROP TABLE feed_follows;

