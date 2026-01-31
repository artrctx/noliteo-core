-- name: GetUserCount :one
SELECT COUNT(*) FROM usr WHERE id=$1;