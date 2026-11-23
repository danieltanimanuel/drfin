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
        WHERE ($1::int IS NULL OR a.owner_id = $1)
        ORDER BY 
            t.id DESC
    `
    queryGetTransactionByID = 
    `
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
        WHERE t.id = $1
    `

    

	queryInsertTransactionManual = `INSERT INTO transactions (parent_id, trade_date, type_id, category_id, account_id, account_dest_id, amount, description, status, source, inserted_at, last_updated_at, confirmed_at) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, NOW(), NOW(), NOW()) RETURNING id`

	queryInsertTransactionAuto = `INSERT INTO transactions (parent_id, trade_date, type_id, category_id, account_id, account_dest_id, amount, description, status, source, inserted_at, last_updated_at) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, NOW(), NOW()) RETURNING id`

    queryInsertTransactionAdjustment = `INSERT INTO transactions (parent_id, trade_date, type_id, category_id, account_id, account_dest_id, amount, description, status, source, inserted_at, last_updated_at, confirmed_at)
		VALUES (0, $1, $2, $3, $4, NULL, $5::numeric, $6, 'Confirmed', 'Manual', NOW(), NOW(), NOW()) RETURNING id`

	// Cancel path: lock the row so concurrent cancel/update serialize.
	queryGetTransactionForUpdate = `
		SELECT id, type_id, account_id, account_dest_id, amount::text, status
		FROM transactions
		WHERE id = $1
		FOR UPDATE`

	queryMarkTransactionCancelled = `
		UPDATE transactions
		SET status = 'Cancelled', last_updated_at = NOW()
		WHERE id = $1 AND status <> 'Cancelled'`
)