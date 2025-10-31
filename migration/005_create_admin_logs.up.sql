CREATE TABLE IF NOT EXISTS account_audit (
  id SERIAL PRIMARY KEY,
  account_id INT NOT NULL,
  admin_id INT NOT NULL,
  action VARCHAR(26), -- "block"/"unblock"
  reason TEXT NULL,
  created_at TIMESTAMPTZ DEFAULT NOW()
);
