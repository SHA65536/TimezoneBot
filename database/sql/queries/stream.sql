-- name: GetAlertChannel :one
SELECT alert_channel FROM streams WHERE guild_id = @guild_id;

-- name: IsStreamAlertEnabled :one
SELECT EXISTS (SELECT 1 FROM streams WHERE guild_id = @guild_id) AS enabled;

-- name: SetAlertChannel :exec
INSERT INTO streams (guild_id, alert_channel) VALUES (@guild_id, @alert_channel) ON CONFLICT (guild_id) DO UPDATE SET alert_channel = @alert_channel;  