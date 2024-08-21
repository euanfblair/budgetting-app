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
	Recurring       bool
	TransactionDate string
	UserId          int
}

type TransactionModel struct {
	DB *sql.DB
}

func (t *TransactionModel) GetUserTransactions(Id int) []Transactions {
	var transactions []Transactions
	var tx Transactions
	stmt := `SELECT * FROM transactions WHERE user_id = $1 ORDER BY transactionid ASC`
	rows, err := t.DB.Query(stmt, Id)
	if err != nil {
		return nil
	}
	for rows.Next() {
		err := rows.Scan(&tx.TransactionId, &tx.Name, &tx.TransactionType, &tx.Amount, &tx.Recurring, &tx.TransactionDate, &tx.UserId)
		if err != nil {
			return nil
		}
		transactions = append(transactions, tx)
	}

	return transactions
}
