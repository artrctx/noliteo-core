-- name: ValidateToken :one
SELECT
    validate_token_key (sqlc.arg (key));