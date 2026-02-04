-- name: ValidateToken :one
SELECT
    id,
    ident,
    created_at
FROM
    validate_token_key (sqlc.arg (key));