CREATE TABLE event_tasks (
    task_id VARCHAR NOT NULL,
    event_id BIGINT NOT NULL,
    is_canceled BOOLEAN DEFAULT false,
    UNIQUE (task_id, event_id)
);