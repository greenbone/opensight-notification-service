CREATE TABLE notification_service.notification_channel
(
    "id"                           UUID         NOT NULL DEFAULT gen_random_uuid() PRIMARY KEY,
    "created_at"                   TIMESTAMP    NOT NULL DEFAULT NOW(),
    "updated_at"                   TIMESTAMP,
    "channel_type"                 VARCHAR(100) NOT NULL CHECK (length(trim(channel_type)) > 0),
    "channel_name"                 VARCHAR(255) UNIQUE,
    "webhook_url"                  VARCHAR(2048),
    "description"                  VARCHAR(2048),
    "domain"                       VARCHAR(255),
    "port"                         INT,
    "is_authentication_required"   BOOLEAN,
    "is_tls_enforced"              BOOLEAN,
    "username"                     VARCHAR(255),
    "password"                     VARCHAR(1024),
    "max_email_attachment_size_mb" INT,
    "max_email_include_size_mb"    INT,
    "sender_email_address"         VARCHAR(255)
);
