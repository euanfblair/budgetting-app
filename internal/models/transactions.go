package models

import (
	"database/sql"
	"math"
)

type Money float64

func MoneyConvert(value int) Money {
	m := Money(math.Round((float64(value)/100)*100)) / 100
	return m
}

type Transactions struct {
	TransactionId   int
	Name            string
	TransactionType bool
	Amount          int
	MoneyAmount     Money
	TransactionDate string
	UserId          int
	Category        string
}

type TransactionModel struct {
	DB *sql.DB
}

func (t *TransactionModel) GetUserTransactions(Id int) []Transactions {
	var transactions []Transactions
	var tx Transactions
	stmt := `SELECT * FROM transactions WHERE user_id = $1 ORDER BY transaction_date DESC`
	rows, err := t.DB.Query(stmt, Id)
	if err != nil {
		return nil
	}
	for rows.Next() {
		err := rows.Scan(&tx.TransactionId, &tx.Name, &tx.TransactionType, &tx.Amount, &tx.TransactionDate, &tx.UserId, &tx.Category)
		if err != nil {
			return nil
		}
		transactions = append(transactions, tx)
	}

	return transactions
}

func (t *TransactionModel) DeleteTransaction(id string) error {
	stmt := `DELETE FROM transactions WHERE transactionid = $1`
	_, err := t.DB.Exec(stmt, id)
	if err != nil {
		return err
	}
	return nil
}

func (t *TransactionModel) GetUniqueCategories(id int) []string {
	var categories = []string{"All"}
	stmt := `SELECT DISTINCT category FROM transactions WHERE user_id = $1 AND category <> 'All';`
	rows, err := t.DB.Query(stmt, id)
	if err != nil {
		return nil
	}
	var category string
	for rows.Next() {
		err = rows.Scan(&category)
		if err != nil {
			return nil
		}
		categories = append(categories, category)
	}

	return categories
}

func (t *TransactionModel) CreateTransaction(name, transactionDate, category string, amount, userId int, transactionType bool) error {
	stmt := `INSERT INTO transactions (name, type, amount, transaction_date, user_id, category ) 
			VALUES ($1, $2, $3, $4, $5, $6)`

	_, err := t.DB.Exec(stmt, name, transactionType, amount, transactionDate, userId, category)
	if err != nil {
		return err
	}

	return nil
}

func (t *TransactionModel) EditTransaction(name, transactionDate, category string, transactionId, amount, userId int, transactionType bool) error {
	stmt := `UPDATE transactions
			SET 
				name = $1,
				type = $2,
				amount = $3,
				transaction_date = $4,
				category = $5
			WHERE 
    			user_id = $6 AND transactionId = $7;`

	_, err := t.DB.Exec(stmt, name, transactionType, amount, transactionDate, category, userId, transactionId)
	if err != nil {
		return err
	}

	return nil
}
