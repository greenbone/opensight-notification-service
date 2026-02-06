CREATE TABLE notification_service.origins (
    "name"        TEXT NOT NULL,
    "class"       TEXT NOT NULL,
    "namespace"   TEXT NOT NULL
);

CREATE INDEX idx_origins_namespace ON notification_service.origins(namespace);
