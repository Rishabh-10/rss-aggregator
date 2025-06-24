-- name: CreateFeeder :one
insert into
    feeders (id, name, link)
values ($1, $2, $3)
returning
    *;