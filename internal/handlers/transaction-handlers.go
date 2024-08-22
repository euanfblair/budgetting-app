package handlers

import (
	"euanfblair/budgeting-app/internal/models"
	"fmt"
	"github.com/labstack/echo/v4"
	"math"
	"net/http"
	"strconv"
	"time"
)

type tableData struct {
	Id              int
	Name            string
	TransactionType string
	Amount          models.Money
	Date            string
	Category        string
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

	data := TemplateData{
		PageIndex: "0",
		PageData:  make([]tableData, 5),
	}
	data.ActiveTab = c.QueryParam("tab")
	userID := app.getUserIdFromSession(c)
	data.AllCategories = app.Transactions.GetUniqueCategories(userID)

	var pages = make([][]tableData, 1)
	pages, data.TotalAmount = app.getPages(userID, data.ActiveTab, data.ActiveCategory)

	data.PageData = pages[0]
	pageCount := len(pages)
	data.PageCount = strconv.Itoa(pageCount - 1)

	return c.Render(http.StatusOK, "transaction-table", data)
}

func convertDate(s string) string {
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return ""
	}
	return t.Format("2006-01-02")
}

func filterTransactions(transactions []models.Transactions, filter, category string) ([]tableData, models.Money) {
	var total int
	var filteredData []tableData

	for _, transaction := range transactions {
		// Check the filter condition
		if (filter == "Incoming" && transaction.TransactionType) ||
			(filter == "Outgoing" && !transaction.TransactionType) ||
			(filter == "All") {

			if (category == "All" || category == "") || category == transaction.Category {

				// Prepare tableData entry
				entry := tableData{
					Id:       transaction.TransactionId,
					Name:     transaction.Name,
					Amount:   models.MoneyConvert(transaction.Amount),
					Date:     convertDate(transaction.TransactionDate),
					Category: transaction.Category,
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
	data.AllCategories = app.Transactions.GetUniqueCategories(userID)

	data.PageData = make([]tableData, len(TransactionData))

	data.PageData, data.TotalAmount = filterTransactions(TransactionData, data.ActiveTab, "")

	return c.Render(http.StatusOK, "table-body", data)

}

func (app *Application) FilterCategory(c echo.Context) error {
	data := TemplateData{}
	data.ActiveCategory = c.QueryParam("categories")
	data.ActiveTab = c.QueryParam("tab")
	userID := app.getUserIdFromSession(c)

	data.AllCategories = app.Transactions.GetUniqueCategories(userID)

	var pages = make([][]tableData, 1)

	pages, data.TotalAmount = app.getPages(userID, data.ActiveTab, data.ActiveCategory)
	if pages == nil {
		pageCount := len(pages)
		data.PageCount = strconv.Itoa(pageCount)
		return c.Render(http.StatusOK, "table-body", data)
	}

	data.PageData = pages[0]
	pageCount := len(pages)
	data.PageCount = strconv.Itoa(pageCount - 1)

	return c.Render(http.StatusOK, "table-body", data)

}

func (app *Application) NextPage(c echo.Context) error {
	data := TemplateData{}
	data.PageIndex = c.QueryParam("page")
	data.ActiveTab = c.QueryParam("tab")
	data.ActiveCategory = c.QueryParam("category")
	fmt.Println(data.ActiveCategory)
	userID := app.getUserIdFromSession(c)

	data.AllCategories = app.Transactions.GetUniqueCategories(userID)

	pageNumber, err := strconv.Atoi(data.PageIndex)
	pageNumber += 1

	if err != nil {
		return err
	}

	var pages = make([][]tableData, 1)
	pages, data.TotalAmount = app.getPages(userID, data.ActiveTab, data.ActiveCategory)

	data.PageData = pages[pageNumber]

	data.PageIndex = strconv.Itoa(pageNumber)
	pageCount := len(pages)
	data.PageCount = strconv.Itoa(pageCount - 1)

	return c.Render(http.StatusOK, "table-body", data)

}

func (app *Application) PrevPage(c echo.Context) error {
	data := TemplateData{}
	data.PageIndex = c.QueryParam("page")
	data.ActiveTab = c.QueryParam("tab")
	data.ActiveCategory = c.QueryParam("category")
	userID := app.getUserIdFromSession(c)

	data.AllCategories = app.Transactions.GetUniqueCategories(userID)

	pageNumber, err := strconv.Atoi(data.PageIndex)
	pageNumber -= 1

	if err != nil {
		return err
	}

	var pages = make([][]tableData, 1)
	pages, data.TotalAmount = app.getPages(userID, data.ActiveTab, data.ActiveCategory)

	data.PageData = pages[pageNumber]

	data.PageIndex = strconv.Itoa(pageNumber)
	pageCount := len(pages)
	data.PageCount = strconv.Itoa(pageCount - 1)

	return c.Render(http.StatusOK, "table-body", data)

}

func (app *Application) getPages(userID int, activeTab, activeCategory string) ([][]tableData, models.Money) {

	const transactionsPerPage = 5

	TransactionData := app.Transactions.GetUserTransactions(userID)

	var pageData = make([]tableData, 5)

	pageData, totalAmount := filterTransactions(TransactionData, activeTab, activeCategory)
	if pageData == nil {
		return nil, 0
	}

	numPages := int(math.Ceil(float64(len(pageData)) / float64(transactionsPerPage)))

	Pages := make([][]tableData, numPages)

	for i := range Pages {
		Pages[i] = make([]tableData, 0)
	}

	for i, transaction := range pageData {
		pageIndex := i / transactionsPerPage
		Pages[pageIndex] = append(Pages[pageIndex], transaction)
	}

	return Pages, totalAmount
}
