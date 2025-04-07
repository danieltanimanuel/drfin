insert into transaction_types values (1, 'Expense',NOW());
insert into transaction_types values (2, 'Income',NOW());
insert into transaction_types values (3, 'Transfer',NOW());
insert into transaction_types values (4, 'Lend',NOW());
insert into transaction_types values (5, 'Borrow',NOW());
insert into transaction_types values (6, 'Adjustment',NOW());

insert into categories values (1,'Food',NOW());
insert into categories values (2,'Transport',NOW());
insert into categories values (3,'Utilities',NOW());
insert into categories values (4,'Gifts',NOW());
insert into categories values (5,'Debts',NOW());
insert into categories values (6,'Misc',NOW());
insert into categories values (7,'Fee',NOW());
insert into categories values (8,'Salary',NOW());

insert into accounts values (1, 'Cash', 0, 1, 'O',NOW(),null);
insert into accounts values (2, 'Tahapan Gold', 0, 2, 'O',NOW(),null);
insert into accounts values (3, 'Tahapan Xpresi', 0, 2, 'O',NOW(),null);
insert into accounts values (4, 'CC BCA', 0, 3, 'O',NOW(),null);
insert into accounts values (5, 'Go-Pay', 0, 5, 'O',NOW(),null);
insert into accounts values (6, 'Badmin 2025/03/12', 0, 7, 'O',NOW(),null);
insert into accounts values (7, 'Blu', 0, 2, 'O',NOW(),null);
insert into accounts values (8, 'Flazz', 0, 4, 'O',NOW(),null);

insert into account_groups values (1, 'Cash',NOW());
insert into account_groups values (2, 'Debit',NOW());
insert into account_groups values (3, 'Credit',NOW());
insert into account_groups values (4, 'E-Money',NOW());
insert into account_groups values (5, 'E-Wallet',NOW());
insert into account_groups values (6, 'Debt',NOW());
insert into account_groups values (7, 'Loan',NOW());
insert into account_groups values (8, 'Saving',NOW());