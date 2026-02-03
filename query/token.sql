-- name: ValidateToken :one
SELECT
    *
FROM
    validate_token_key (sqlc.arg (key));