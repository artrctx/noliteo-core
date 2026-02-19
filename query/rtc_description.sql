-- name: GetRTCDescription :one
SELECT
    *
FROM
    rtc_description
WHERE
    token_id = $1
    AND type = $2;

-- name: CreateRTCDescription :exec
INSERT INTO
    rtc_description (token_id, sdp, type)
VALUES
    ($1, $2, $3);

-- name: DeleteRTCDescription :exec
DELETE FROM rtc_description
WHERE
    token_id = $1;