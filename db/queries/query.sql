-- name: CreateFeeder :one
insert into
    feeders (id, name, link)
values ($1, $2, $3)
returning
    *;

-- name: GetFeeders :many
select * from feeders;

-- name: CreateFeed :one
insert into
    feeds (id, title, description)
values ($1, $2, $3)
returning
    *;