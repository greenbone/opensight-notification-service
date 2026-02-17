CREATE TABLE notification_service.rules (
    "id"                UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    "name"              TEXT NOT NULL CONSTRAINT rules_name_unique UNIQUE,
    "trigger_origins"   TEXT[] NOT NULL,
    "trigger_levels"    TEXT[] NOT NULL,
    "action_channel_id" UUID REFERENCES notification_service.notification_channel(id) ON DELETE CASCADE,
    "action_recipient"  TEXT,
    "active"            BOOLEAN NOT NULL
);

CREATE INDEX idx_rules_id ON notification_service.rules(id);
