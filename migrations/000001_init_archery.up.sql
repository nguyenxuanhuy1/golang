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
create table images (
   id            bigserial primary key,
   image_url     text not null,
   blur_url      text,
   tiny_blur_url text,
   public_id     text not null,
   image_type    text not null,
   owner_id      integer,
   created_at    timestamptz default now(),
   constraint fk_images_user foreign key ( owner_id )
      references users ( id )
         on delete cascade
);
create table posts (
   id          bigserial primary key,
   image_id    bigint
      references images ( id )
         on delete set null,
   name        varchar(255) not null,
   description varchar(255),
   topic       varchar(100) not null,
   prompt      text,
   hot_level   smallint default 0 check ( hot_level between 0 and 9 ),
   hot_at      timestamptz,
   created_at  timestamptz default now(),
   updated_at  timestamptz default now()
);

create index idx_posts_topic_created on
   posts (
      topic,
      created_at
   desc );

create index idx_posts_hot on
   posts (
      hot_level
   desc,
      hot_at
   desc );