package transactions

const (
	queryGetAllTransaction = `
        SELECT 
            t.id, t.parent_id, t.trade_date, t.type_id, t.category_id, t.account_id, 
            t.account_dest_id, t.amount, t.description, t.status, t.source, 
            t.inserted_at, t.last_updated_at, t.confirmed_at, 
            c.name, 
			tt.name,
			a.name,
			ad.name
        FROM 
            transactions t
        LEFT JOIN 
            transaction_types tt ON t.type_id = tt.id
        LEFT JOIN 
            accounts a ON t.account_id = a.id
        LEFT JOIN 
            accounts ad ON t.account_dest_id = ad.id
        LEFT JOIN 
            categories c ON t.category_id = c.id
        ORDER BY 
            t.id ASC
    `
	queryInsertTransactionManual = `INSERT INTO transactions (parent_id, trade_date, type_id, category_id, account_id, account_dest_id, amount, description, status, source, inserted_at, last_updated_at, confirmed_at) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, NOW(), NOW(), NOW()) RETURNING id`

	queryInsertTransactionAuto = `INSERT INTO transactions (parent_id, trade_date, type_id, category_id, account_id, account_dest_id, amount, description, status, source, inserted_at, last_updated_at) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, NOW(), NOW()) RETURNING id`
)
