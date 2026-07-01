-- name: CreateReminder :one
INSERT INTO reminders (user_id, subscription_id, type, channel, scheduled_for)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: ListPendingReminders :many
SELECT * FROM reminders
WHERE status = 'pending'
  AND scheduled_for <= $1
ORDER BY scheduled_for;

-- name: MarkReminderSent :exec
UPDATE reminders
SET status = 'sent', sent_at = now()
WHERE id = $1;

-- name: ListRemindersBySubscription :many
SELECT * FROM reminders WHERE subscription_id = $1 ORDER BY scheduled_for;
