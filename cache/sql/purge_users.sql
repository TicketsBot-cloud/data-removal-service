DELETE FROM users
WHERE "last_seen" < NOW() - $1::INTERVAL;