DROP TABLE IF EXISTS handouts;

CREATE TABLE handouts (
  id BIGSERIAL PRIMARY KEY,
  address TEXT,
  amount DECIMAL(15,2) NOT NULL,
  date TIMESTAMPTZ NOT NULL,
  name TEXT NOT NULL,
  nominee TEXT NOT NULL,
  mobile BIGINT,
  created_at TIMESTAMPTZ DEFAULT NOW(),
  updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_handouts_name ON handouts(name);
CREATE INDEX idx_handouts_date ON handouts(date);
CREATE INDEX idx_handouts_mobile ON handouts(mobile);
CREATE INDEX idx_handouts_created_at ON handouts(created_at);

CREATE OR REPLACE FUNCTION update_updated_at_column() 
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = NOW();
  RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_handouts_updated_at 
BEFORE UPDATE ON  handouts
FOR EACH ROW
EXECUTE FUNCTION update_updated_at_column();


