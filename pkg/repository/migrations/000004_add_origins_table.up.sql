CREATE TABLE notification_service.origins (
    "name"        TEXT NOT NULL,
    "class"       TEXT NOT NULL,
    "service_id"   TEXT NOT NULL
);

CREATE INDEX idx_origins_service_id ON notification_service.origins(service_id);
