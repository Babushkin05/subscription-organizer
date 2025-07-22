CREATE TABLE IF NOT EXISTS subscriptions (
    id UUID PRIMARY KEY,
    service_name TEXT NOT NULL,
    price INTEGER NOT NULL CHECK (price >= 0),
    user_id UUID NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE,
    is_deleted BOOLEAN NOT NULL DEFAULT FALSE
);

CREATE INDEX IF NOT EXISTS idx_subscriptions_user_id ON subscriptions(user_id);

CREATE INDEX IF NOT EXISTS idx_subscriptions_service_name ON subscriptions(service_name);

CREATE INDEX IF NOT EXISTS idx_subscriptions_start_date ON subscriptions(start_date);
CREATE INDEX IF NOT EXISTS idx_subscriptions_end_date ON subscriptions(end_date);

CREATE INDEX IF NOT EXISTS idx_subscriptions_filtering ON subscriptions(user_id, service_name, start_date, end_date);

CREATE INDEX IF NOT EXISTS idx_subscriptions_is_deleted ON subscriptions(is_deleted);
