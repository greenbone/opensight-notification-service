-- split origin_uri into origin_class and origin_resource_id, separated by the last `/`

ALTER TABLE notification_service.notifications
    ADD COLUMN origin_class TEXT,
    ADD COLUMN origin_resource_id TEXT;

UPDATE notification_service.notifications
SET
    origin_class = LEFT(origin_uri, LENGTH(origin_uri) - POSITION('/' IN REVERSE(origin_uri))),
    origin_resource_id = CASE
        WHEN POSITION('/' IN origin_uri) > 0
            THEN split_part(origin_uri, '/', -1)
        ELSE ''
    END;

ALTER TABLE notification_service.notifications
    DROP COLUMN origin_uri;

