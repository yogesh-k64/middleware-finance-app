
ALTER TABLE handouts
ADD COLUMN user_id BIGINT,
ADD COLUMN nominee_id BIGINT;

ALTER TABLE handouts
DROP COLUMN name,
DROP COLUMN nominee;

ALTER TABLE handouts
DROP COLUMN mobile,
DROP COLUMN address;

ALTER TABLE handouts
ADD CONSTRAINT fk_handouts_user
FOREIGN KEY (user_id) REFERENCES users(id);

ALTER TABLE handouts
ADD CONSTRAINT fk_handouts_nominee
FOREIGN KEY (nominee_id) REFERENCES users(id);

ALTER TABLE users
DROP COLUMN phone_number;

ALTER TABLE users
ADD COLUMN mobile BIGINT;

UPDATE handouts SET user_id = 0 WHERE user_id = 1;

ALTER TABLE handouts ALTER COLUMN user_id SET NOT NULL;