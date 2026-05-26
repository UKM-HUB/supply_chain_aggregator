CREATE TABLE sme_categories (
    sme_id      UUID         NOT NULL REFERENCES smes       (id)  ON DELETE CASCADE,
    category_id VARCHAR(100) NOT NULL REFERENCES categories (id)  ON DELETE CASCADE,
    PRIMARY KEY (sme_id, category_id)
);

CREATE INDEX idx_sme_categories_category_id ON sme_categories (category_id);
