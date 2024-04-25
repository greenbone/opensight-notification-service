CREATE SCHEMA IF NOT EXISTS notification_service;

COMMENT ON SCHEMA notification_service IS 'Notification Service schema';

CREATE TABLE notification_service.notifications (
    "id"            UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    "origin"        TEXT NOT NULL,
    "originUri"     TEXT NOT NULL,
    "timestamp"     TIMESTAMP,
    "title"         TEXT NOT NULL,
    "detail"        TEXT NOT NULL,
    "level"         VARCHAR(255) NOT NULL,
    "custom_fields" JSONB
);
