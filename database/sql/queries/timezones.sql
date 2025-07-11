-- name: GetTimezone :one
SELECT timezone FROM timezones WHERE user_id = @user_id;

-- name: SetTimezone :exec
INSERT INTO timezones (user_id, timezone) VALUES (@user_id, @timezone) ON CONFLICT (user_id) DO UPDATE SET timezone = @timezone;  