package handlers

import (
	"euanfblair/budgeting-app/internal/models"
	"github.com/labstack/echo/v4"
	"net/http"
	"time"
)

type tableData struct {
	Id              int
	Name            string
	TransactionType string
	Amount          models.Money
	Date            string
}

func (app *Application) GetTransactions(c echo.Context) error {
	data := TemplateData{
		Title: "Signup",
	}

	data.IsAuthenticated = app.SessionManager.Exists(c.Request().Context(), "authUserID")

	data.ActiveTab = "All"
	return c.Render(http.StatusOK, "transactions", data)
}

func (app *Application) FilteredTransactions(c echo.Context) error {

	data := TemplateData{}
	data.ActiveTab = c.QueryParam("tab")
	userID := app.getUserIdFromSession(c)
	TransactionData := app.Transactions.GetUserTransactions(userID)

	data.TableData = make([]tableData, len(TransactionData))

	data.TableData, data.TotalAmount = filterTransactions(TransactionData, data.ActiveTab)

	return c.Render(http.StatusOK, "transaction-table", data)
}

func convertDate(s string) string {
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return ""
	}
	return t.Format("2006-01-02")
}

func filterTransactions(transactions []models.Transactions, filter string) ([]tableData, models.Money) {
	var total int
	var filteredData []tableData

	for _, transaction := range transactions {
		// Check the filter condition
		if (filter == "Incoming" && transaction.TransactionType) ||
			(filter == "Outgoing" && !transaction.TransactionType) ||
			(filter == "All") {

			// Prepare tableData entry
			entry := tableData{
				Id:     transaction.TransactionId,
				Name:   transaction.Name,
				Amount: models.MoneyConvert(transaction.Amount),
				Date:   convertDate(transaction.TransactionDate),
			}

			if transaction.TransactionType {
				entry.TransactionType = "Income"
				total += transaction.Amount
			} else {
				entry.TransactionType = "Outgoing"
				total -= transaction.Amount
			}

			// Add entry to filteredData
			filteredData = append(filteredData, entry)
		}
	}

	// Convert total to Money
	totalAmount := models.MoneyConvert(total)

	return filteredData, totalAmount
}

func (app *Application) DeleteTransaction(c echo.Context) error {

	data := TemplateData{}
	transactionId := c.QueryParam("id")
	data.ActiveTab = c.QueryParam("tab")
	userID := app.getUserIdFromSession(c)

	err := app.Transactions.DeleteTransaction(transactionId)
	if err != nil {
		return err
	}

	TransactionData := app.Transactions.GetUserTransactions(userID)

	data.TableData = make([]tableData, len(TransactionData))

	data.TableData, data.TotalAmount = filterTransactions(TransactionData, data.ActiveTab)

	return c.Render(http.StatusOK, "transaction-table", data)

}
