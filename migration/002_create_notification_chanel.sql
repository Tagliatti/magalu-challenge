CREATE OR REPLACE FUNCTION notify_pending_notifications() RETURNS trigger AS $$
DECLARE
    payload json;
BEGIN
    IF (TG_OP = 'INSERT') THEN
        SELECT row_to_json(NEW) INTO payload;
        PERFORM pg_notify('pending_notifications', payload::text);
    END IF;

    RETURN NULL;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER pending_notifications_trigger
    AFTER INSERT ON notifications
    FOR EACH ROW EXECUTE FUNCTION notify_pending_notifications();
