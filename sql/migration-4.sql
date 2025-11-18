DROP TABLE IF EXISTS collections;

CREATE TABLE collections (
    id BIGINT PRIMARY KEY,
    date TIMESTAMPTZ NOT NULL,
    handout_id BIGINT NOT NULL,
    amount DECIMAL(15,2) NOT NULL,
    total_paid DECIMAL(15,2),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    FOREIGN KEY (handout_id) REFERENCES handouts(id)
);

