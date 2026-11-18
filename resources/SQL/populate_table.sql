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

INSERT INTO accounts VALUES
(DEFAULT, 'Cash', 0, 1, 'O', NOW(), NULL),
(DEFAULT, 'Tahapan Gold', 0, 1, 'O', NOW(), NULL),
(DEFAULT, 'Tahapan Xpresi', 0, 1, 'O', NOW(), NULL),
(DEFAULT, 'CC BCA', 0, 1, 'O', NOW(), NULL),
(DEFAULT, 'Blu', 0, 1, 'O', NOW(), NULL),
(DEFAULT, 'Blu Transport', 0, 1, 'O', NOW(), NULL),
(DEFAULT, 'Blu Makan', 0, 1, 'O', NOW(), NULL),
(DEFAULT, 'Blu Entertainment', 0, 1, 'O', NOW(), NULL),
(DEFAULT, 'Mami', 0, 1, 'O', NOW(), NULL);

INSERT INTO account_groups VALUES
(1, 'Cash', NOW()),
(2, 'Debit', NOW()),
(3, 'Credit', NOW()),
(4, 'E-Money', NOW()),
(5, 'E-Wallet', NOW()),
(6, 'Debt', NOW()),
(7, 'Loan', NOW()),
(8, 'Saving', NOW());

INSERT INTO owned_account VALUES
('BCA Tahapan Gold IDR/Forex', '0840945740'),
('BCA Tahapan Gold IDR', '2060524290'),
('BCA Tahapan Xpresi IDR', '0840945740'),
('BLU BY BCA DIGITAL', '005493616619');