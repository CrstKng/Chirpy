-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email, hashed_password)
VALUES (
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1,
    $2
)
RETURNING *;

-- name: DeleteUsers :exec
DELETE FROM users;

-- name: CreateChirp :one
INSERT INTO chirps (id, created_at, updated_at, body, user_id)
VALUES (
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1,
    $2
)
RETURNING *;

-- name: GetChirps :many
SELECT * FROM chirps
ORDER BY created_at;

-- name: GetChirp :one
SELECT * FROM chirps
WHERE id = $1;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1;

-- name: GetUserByID :one
SELECT * FROM users
WHERE id = $1;

-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens (token, created_at, updated_at, user_id, expires_at, revoked_at)
VALUES (
    $1,
    NOW(),
    NOW(),
    $2,
    NOW() + INTERVAL '60 days',
    NULL
)
RETURNING *;

-- name: GetUserByRefreshToken :one
SELECT users.* FROM users
INNER JOIN refresh_tokens
ON users.id = refresh_tokens.user_id
WHERE refresh_tokens.token = $1
AND refresh_tokens.expires_at > NOW()
AND refresh_tokens.revoked_at IS NULL;

-- name: RevokeRefreshToken :exec
UPDATE refresh_tokens
SET
  updated_at = NOW(),
  revoked_at = NOW()
WHERE
  token = $1;

-- name: UpdateUserEmailPassword :exec
UPDATE users
SET
  updated_at = NOW(),
  email = $1,
  hashed_password = $2
WHERE
  id = $3;

-- name: DeleteChirpByID :exec
DELETE FROM chirps
WHERE
  id = $1;

-- name: ValidateChirpOwner :one
SELECT users.* FROM users
INNER JOIN chirps
ON users.id = chirps.user_id
WHERE users.id = $1 AND chirps.id = $2;

-- name: UpdateUserMembership :exec
UPDATE users
SET
  is_chirpy_red = TRUE
WHERE
  id = $1;