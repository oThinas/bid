-- name: CreateBid :one
INSERT INTO bids (product_id, bidder_id, amount)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetBidsByProductID :many
SELECT * FROM bids
WHERE product_id = $1
ORDER BY amount DESC;

-- name: GetHighestBidByProductID :one
SELECT * FROM bids
WHERE product_id = $1
ORDER BY amount DESC
LIMIT 1;
