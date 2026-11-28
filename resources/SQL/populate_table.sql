-- multi-owner: seed order is now users -> account_groups -> accounts, so the
-- accounts.account_group_id and accounts.owner_id foreign keys are satisfied.

INSERT INTO transaction_types VALUES
(1, 'Expense', NOW()),
(2, 'Income', NOW()),
(3, 'Transfer', NOW()),
(4, 'Lend', NOW()),
(5, 'Borrow', NOW()),
(6, 'Adjustment', NOW());

INSERT INTO categories VALUES
(DEFAULT, 'Transport', NOW()),
(DEFAULT, 'Salary', NOW()),
(DEFAULT, 'Unknown', NOW()),
(DEFAULT, 'Mami', NOW()),
(DEFAULT, 'Tagihan', NOW()),
(DEFAULT, 'Kebutuhan Rumah', NOW()),
(DEFAULT, 'Mama', NOW()),
(DEFAULT, 'Makan', NOW()),
(DEFAULT, 'Entertainment', NOW()),
(DEFAULT, 'Skin Care', NOW()),
(DEFAULT, 'Occasional', NOW());

-- multi-owner: moved before accounts (accounts references account_groups)
INSERT INTO account_groups VALUES
(1, 'Cash', NOW()),
(2, 'Debit', NOW()),
(3, 'Credit', NOW()),
(4, 'E-Money', NOW()),
(5, 'E-Wallet', NOW()),
(6, 'Debt', NOW()),
(7, 'Loan', NOW()),
(8, 'Saving', NOW());

-- multi-owner: owners. Placeholders only; this repository is public.
-- After installing, set the real values locally, for example:
--   UPDATE users SET name = '...', email = '...', bank_holder_name = '...' WHERE id = 1;
-- bank_holder_name must match the name printed in that owner's blu emails exactly.
INSERT INTO users (id, email, name, bank_holder_name) VALUES
(1, 'owner1@example.com', 'OWNER_1_NAME', 'HOLDER_NAME_1'),
(2, 'owner2@example.com', 'OWNER_2_NAME', 'HOLDER_NAME_2');

-- multi-owner: explicit column list (owner_id is NOT NULL and last in the table,
-- so the old positional insert no longer works). The starter accounts belong to
-- owner 1; add owner 2's accounts from the dashboard.
INSERT INTO accounts (name, balance, account_group_id, status, inserted_at, owner_id) VALUES
('Cash', 0, 1, 'O', NOW(), 1),
('Tahapan Gold', 0, 1, 'O', NOW(), 1),
('Tahapan Xpresi', 0, 1, 'O', NOW(), 1),
('CC BCA', 0, 1, 'O', NOW(), 1),
('Blu', 0, 1, 'O', NOW(), 1),
('Blu Transport', 0, 1, 'O', NOW(), 1),
('Blu Makan', 0, 1, 'O', NOW(), 1),
('Blu Entertainment', 0, 1, 'O', NOW(), 1),
('Mami', 0, 1, 'O', NOW(), 1);

-- multi-owner: one Gmail checkpoint row per owner. internal_date is Unix
-- milliseconds: the first sync only imports emails after this moment. Using
-- install time means no old mailbox history is imported; replace it with an
-- earlier timestamp if you want history. 0 would import the whole mailbox.
INSERT INTO last_email (message_id, internal_date, owner_id) VALUES
('', (EXTRACT(EPOCH FROM NOW()) * 1000)::bigint, 1),
('', (EXTRACT(EPOCH FROM NOW()) * 1000)::bigint, 2);

INSERT INTO owned_account VALUES
('BCA Tahapan Gold IDR/Forex', '0840945740'),
('BCA Tahapan Gold IDR', '2060524290'),
('BCA Tahapan Xpresi IDR', '0840945740'),
('BLU BY BCA DIGITAL', '005493616619');

-- multi-owner: rows with explicit IDs do not advance their sequences. Without
-- this, the next insert from the app fails with a duplicate key error.
SELECT setval('transaction_types_id_seq', (SELECT MAX(id) FROM transaction_types));
SELECT setval('account_groups_id_seq', (SELECT MAX(id) FROM account_groups));
SELECT setval('users_id_seq', (SELECT MAX(id) FROM users));