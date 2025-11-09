
CREATE TABLE IF NOT EXISTS movies (
  metadata_id SERIAL PRIMARY KEY,
  title TEXT NOT NULL,
  description TEXT,
  director TEXT,
  runtime INT,
  created_at TIMESTAMPTZ DEFAULT now(),
  updated_at TIMESTAMPTZ DEFAULT now()
);

CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = now();
    RETURN NEW;
END;
$$ language 'plpgsql';

-- trigger for updated_at
CREATE TRIGGER update_movies_updated_at
BEFORE UPDATE ON movies
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();

