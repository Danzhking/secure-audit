CREATE TABLE IF NOT EXISTS security_events (
    id          BIGSERIAL    PRIMARY KEY,
    service     VARCHAR(255) NOT NULL,
    event_type  VARCHAR(255) NOT NULL,
    severity    VARCHAR(20)  NOT NULL DEFAULT 'low',
    user_id     VARCHAR(255) NOT NULL DEFAULT '',
    ip          VARCHAR(45)  NOT NULL DEFAULT '',
    metadata    JSONB        DEFAULT '{}',
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_security_events_service    ON security_events (service);
CREATE INDEX IF NOT EXISTS idx_security_events_event_type ON security_events (event_type);
CREATE INDEX IF NOT EXISTS idx_security_events_severity   ON security_events (severity);
CREATE INDEX IF NOT EXISTS idx_security_events_user_id    ON security_events (user_id);
CREATE INDEX IF NOT EXISTS idx_security_events_ip         ON security_events (ip);
CREATE INDEX IF NOT EXISTS idx_security_events_created_at ON security_events (created_at);

CREATE TABLE IF NOT EXISTS alerts (
    id          BIGSERIAL    PRIMARY KEY,
    rule_name   VARCHAR(255) NOT NULL,
    severity    VARCHAR(20)  NOT NULL,
    message     TEXT         NOT NULL,
    user_id     VARCHAR(255) DEFAULT '',
    ip          VARCHAR(45)  DEFAULT '',
    event_count INT          NOT NULL DEFAULT 0,
    status      VARCHAR(20)  NOT NULL DEFAULT 'new',
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    resolved_at TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_alerts_rule_name  ON alerts (rule_name);
CREATE INDEX IF NOT EXISTS idx_alerts_severity   ON alerts (severity);
CREATE INDEX IF NOT EXISTS idx_alerts_status     ON alerts (status);
CREATE INDEX IF NOT EXISTS idx_alerts_user_id    ON alerts (user_id);
CREATE INDEX IF NOT EXISTS idx_alerts_ip         ON alerts (ip);
CREATE INDEX IF NOT EXISTS idx_alerts_created_at ON alerts (created_at);
