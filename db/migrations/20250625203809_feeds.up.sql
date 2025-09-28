create table if not exists feeds (
    id uuid not null primary key,
    feeder_id uuid references feeders (id) on delete cascade not null,
    title text not null,
    description text not null,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now()
);