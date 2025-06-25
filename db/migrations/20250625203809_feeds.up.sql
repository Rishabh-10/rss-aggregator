create table if not exists feeds (
    id uuid not null primary key,
    title text not null,
    description text not null,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now()
);