DROP TABLE IF EXISTS collections;

CREATE TABLE collections (
    id BIGSERIAL PRIMARY KEY,
    date TIMESTAMPTZ NOT NULL,
    handout_id BIGINT NOT NULL,
    amount DECIMAL(15,2) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    FOREIGN KEY (handout_id) REFERENCES handouts(id)
);

CREATE INDEX idx_collections_amount ON collections(amount);
CREATE INDEX idx_collections_date ON collections(date);
CREATE INDEX idx_collections_created_at ON collections(created_at);

CREATE TRIGGER update_collections_updated_at
BEFORE UPDATE ON collections
FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

ALTER TABLE handouts
DROP COLUMN nominee_id;
