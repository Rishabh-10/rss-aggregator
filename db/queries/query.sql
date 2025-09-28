-- name: CreateFeeder :one
insert into
    feeders (id, name, link)
values ($1, $2, $3) returning *;

-- name: GetFeeders :many
select * from feeders;

-- name: CreateFeed :one
insert into
    feeds (
        id,
        feeder_id,
        title,
        description
    )
values ($1, $2, $3, $4) returning *;

-- name: GetFeedsByFeederID :many
select * from feeds where feeder_id = $1 LIMIT $2 OFFSET $3;