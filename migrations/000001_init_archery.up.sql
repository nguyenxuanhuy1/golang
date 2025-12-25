create table users (
   id         serial primary key,
   username   text not null unique,
   password   text not null,
   role       text not null default 'user',
   avatar     text,
   elo        integer default 1000,
   created_at timestamptz default now()
);