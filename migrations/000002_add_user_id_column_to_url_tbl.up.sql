ALTER TABLE url 
    ADD COLUMN user_id BIGINT;
CREATE INDEX links_user_id_idx ON url(user_id);