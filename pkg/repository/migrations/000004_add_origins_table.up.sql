CREATE TABLE notification_service.origins (
    "name"        TEXT NOT NULL,
    "class"       TEXT NOT NULL CONSTRAINT origins_class_unique UNIQUE,
    "service_id"  TEXT NOT NULL
);

CREATE INDEX idx_origins_service_id ON notification_service.origins(service_id);
