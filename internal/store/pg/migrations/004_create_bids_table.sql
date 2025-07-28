-- Write your migrate up statements here
CREATE TABLE IF NOT EXISTS bids (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  product_id UUID NOT NULL REFERENCES products(id),
  bidder_id UUID NOT NULL REFERENCES users(id),
  amount FLOAT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

---- create above / drop below ----
DROP TABLE IF EXISTS bids;
-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
