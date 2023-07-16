DELETE FROM members
WHERE "last_seen" < NOW() - $1::INTERVAL;