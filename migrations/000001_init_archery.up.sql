create table users (
   id         serial primary key,
   username   text not null unique,
   password   text not null,
   role       text not null default 'user',
   avatar     text,
   locked     boolean default false,
   coin       integer default 0,
   created_at timestamptz default now()
);
CREATE TABLE images (
  id BIGSERIAL PRIMARY KEY,

  image_url TEXT NOT NULL,
  blur_url TEXT,
  tiny_blur_url TEXT,

  public_id TEXT NOT NULL,

  image_type TEXT NOT NULL,
  owner_id INTEGER,

  created_at TIMESTAMPTZ DEFAULT now(),

  CONSTRAINT fk_images_user
    FOREIGN KEY (owner_id)
    REFERENCES users(id)
    ON DELETE CASCADE
);
CREATE TABLE posts (
    id BIGSERIAL PRIMARY KEY,

    image_id BIGINT REFERENCES images(id) ON DELETE SET NULL,

    name VARCHAR(255) NOT NULL,
    description VARCHAR(255),
    topic VARCHAR(100) NOT NULL,
    prompt TEXT,
    hot_level SMALLINT DEFAULT 0 CHECK (hot_level BETWEEN 0 AND 9),
    hot_at TIMESTAMPTZ,

    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now()
);

CREATE INDEX idx_posts_topic_created
ON posts (topic, created_at DESC);

CREATE INDEX idx_posts_hot
ON posts (hot_level DESC, hot_at DESC);
