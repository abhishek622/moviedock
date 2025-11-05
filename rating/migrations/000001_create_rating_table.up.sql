CREATE TABLE IF NOT EXISTS ratings (
    record_id BIGINT NOT NULL,
    record_type TEXT NOT NULL,
    user_id TEXT NOT NULL,
    value INT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now()
);

-- trigger for updated_at
CREATE TRIGGER update_ratings_updated_at
BEFORE UPDATE ON ratings
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();