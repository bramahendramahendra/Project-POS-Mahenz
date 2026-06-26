CREATE TABLE IF NOT EXISTS log_schedulers (
    id             VARCHAR(36)              NOT NULL PRIMARY KEY,
    scheduler_name VARCHAR(100)             NOT NULL,
    status         ENUM('success','failed') NOT NULL,
    message        TEXT                     NULL,
    duration_ms    INT                      NULL,
    executed_at    DATETIME                 NOT NULL DEFAULT CURRENT_TIMESTAMP
) DEFAULT CHARSET=utf8mb4;
