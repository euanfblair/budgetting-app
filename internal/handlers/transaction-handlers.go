package handlers

import (
	"euanfblair/budgeting-app/internal/models"
	"fmt"
	"github.com/labstack/echo/v4"
	"math"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode/utf8"
)

var transactionCache = struct {
	sync.RWMutex
	cache map[int][]models.Transactions
}{cache: make(map[int][]models.Transactions)}

var CategoriesCache = struct {
	sync.RWMutex
	cache map[int][]string
}{cache: make(map[int][]string)}

type tableData struct {
	Id              int
	Name            string
	TransactionType string
	Amount          models.Money
	Date            string
	Category        string
}

func (app *Application) getUserTransactions(userId int) ([]models.Transactions, error) {
	transactionCache.RLock()
	transactions, found := transactionCache.cache[userId]
	transactionCache.RUnlock()

	if found {
		return transactions, nil
	}

	// Data not in cache, fetch from DB
	transactions = app.Transactions.GetUserTransactions(userId)

	// Update cache
	transactionCache.Lock()
	transactionCache.cache[userId] = transactions
	transactionCache.Unlock()

	return transactions, nil
}

func (app *Application) getUserCategories(userId int) ([]string, error) {
	transactionCache.RLock()
	categories, found := CategoriesCache.cache[userId]
	transactionCache.RUnlock()

	if found {
		return categories, nil
	}

	// Data not in cache, fetch from DB
	categories = app.Transactions.GetUniqueCategories(userId)

	// Update cache
	CategoriesCache.Lock()
	CategoriesCache.cache[userId] = categories
	CategoriesCache.Unlock()

	return categories, nil
}

func (app *Application) GetTransactions(c echo.Context) error {
	data := TemplateData{
		Title: "Transactions",
	}

	IsAuthenticated := app.SessionManager.Exists(c.Request().Context(), "authUserID")
	if !IsAuthenticated {
		return c.Render(http.StatusOK, "not-auth", data)
	}

	data.ActiveTab = "All"
	return c.Render(http.StatusOK, "transactions", data)
}

func (app *Application) FilteredTransactions(c echo.Context) error {

	data := TemplateData{
		PageIndex: "0",
	}

	data.ActiveTab = c.QueryParam("tab")
	userID := app.getUserIdFromSession(c)
	data.AllCategories, _ = app.getUserCategories(userID)

	var pages = make([][]tableData, 1)
	pages, data.TotalAmount = app.getPages(userID, data.ActiveTab, data.ActiveCategory, data.ActiveMonth)
	if pages != nil {
		data.PageData = make([]tableData, 0)
		data.PageData = pages[0]
	}

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

func filterTransactions(transactions []models.Transactions, filter, category, timeframe string) ([]tableData, models.Money) {
	var total int
	var filteredData []tableData

	baseDate := time.Now()

	// Calculate the date range (one month before and after the baseDate)
	oneMonthBefore := baseDate.AddDate(0, -1, 0)
	threeMonthBefore := baseDate.AddDate(0, -3, 0)
	oneYearBefore := baseDate.AddDate(0, -12, 0)

	for _, transaction := range transactions {
		// Check the filter condition
		if (filter == "Incoming" && transaction.TransactionType) ||
			(filter == "Outgoing" && !transaction.TransactionType) ||
			(filter == "All") {

			date, err := time.Parse(time.RFC3339, transaction.TransactionDate)
			if err != nil {
				fmt.Println("Error parsing transaction date:", err)
				continue
			}

			if (timeframe == "1" && date.After(oneMonthBefore)) ||
				(timeframe == "3" && date.After(threeMonthBefore)) ||
				(timeframe == "12" && date.After(oneYearBefore)) ||
				(timeframe == "All" || timeframe == "") {

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
						entry.TransactionType = "Incoming"
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
	}

	// Convert total to Money
	totalAmount := models.MoneyConvert(total)

	return filteredData, totalAmount
}

func (app *Application) DeleteTransaction(c echo.Context) error {

	data := TemplateData{
		PageIndex: "0",
	}
	transactionId := c.QueryParam("id")
	data.ActiveTab = c.QueryParam("tab")
	data.ActiveMonth = c.QueryParam("month")
	data.ActiveCategory = c.QueryParam("category")
	userID := app.getUserIdFromSession(c)

	err := app.Transactions.DeleteTransaction(transactionId)
	if err != nil {
		return err
	}

	transactionCache.Lock()
	delete(transactionCache.cache, userID)
	transactionCache.Unlock()

	CategoriesCache.Lock()
	delete(CategoriesCache.cache, userID)
	CategoriesCache.Unlock()

	data.AllCategories, _ = app.getUserCategories(userID)

	var pages = make([][]tableData, 1)

	pages, data.TotalAmount = app.getPages(userID, data.ActiveTab, data.ActiveCategory, data.ActiveMonth)
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

func (app *Application) FilterCategory(c echo.Context) error {
	data := TemplateData{
		PageIndex: "0",
	}
	data.ActiveCategory = c.QueryParam("categories")
	data.ActiveTab = c.QueryParam("tab")
	data.ActiveMonth = c.QueryParam("month")
	userID := app.getUserIdFromSession(c)

	data.AllCategories, _ = app.getUserCategories(userID)

	var pages = make([][]tableData, 1)

	pages, data.TotalAmount = app.getPages(userID, data.ActiveTab, data.ActiveCategory, data.ActiveMonth)
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
	data.ActiveMonth = c.QueryParam("month")
	userID := app.getUserIdFromSession(c)

	data.AllCategories, _ = app.getUserCategories(userID)

	pageNumber, err := strconv.Atoi(data.PageIndex)
	pageNumber += 1

	if err != nil {
		fmt.Println("Error parsing page number:", err)
		return err
	}

	var pages = make([][]tableData, 1)
	pages, data.TotalAmount = app.getPages(userID, data.ActiveTab, data.ActiveCategory, data.ActiveMonth)

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
	data.ActiveMonth = c.QueryParam("month")
	userID := app.getUserIdFromSession(c)

	data.AllCategories, _ = app.getUserCategories(userID)

	pageNumber, err := strconv.Atoi(data.PageIndex)
	pageNumber -= 1

	if err != nil {
		return err
	}

	var pages = make([][]tableData, 1)
	pages, data.TotalAmount = app.getPages(userID, data.ActiveTab, data.ActiveCategory, data.ActiveMonth)

	data.PageData = pages[pageNumber]

	data.PageIndex = strconv.Itoa(pageNumber)
	pageCount := len(pages)
	data.PageCount = strconv.Itoa(pageCount - 1)

	return c.Render(http.StatusOK, "table-body", data)

}

func (app *Application) FilterTimeFrame(c echo.Context) error {
	data := TemplateData{
		PageIndex: "0",
	}
	data.ActiveMonth = c.QueryParam("time")
	data.ActiveTab = c.QueryParam("tab")

	userID := app.getUserIdFromSession(c)
	data.AllCategories, _ = app.getUserCategories(userID)

	var pages = make([][]tableData, 1)
	pages, data.TotalAmount = app.getPages(userID, data.ActiveTab, data.ActiveCategory, data.ActiveMonth)

	if pages != nil {
		data.PageData = pages[0]
	}
	data.PageIndex = "0"
	pageCount := len(pages)
	data.PageCount = strconv.Itoa(pageCount - 1)

	return c.Render(http.StatusOK, "time-frame-header-table", data)

}

func (app *Application) getPages(userID int, activeTab, activeCategory, timeframe string) ([][]tableData, models.Money) {

	const transactionsPerPage = 5

	TransactionData, err := app.getUserTransactions(userID)
	if err != nil {
		return nil, 0
	}

	var pageData = make([]tableData, 5)

	pageData, totalAmount := filterTransactions(TransactionData, activeTab, activeCategory, timeframe)
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

func (app *Application) CreateTransaction(c echo.Context) error {
	data := TemplateData{}
	userID := app.getUserIdFromSession(c)

	data.ActiveTab = c.QueryParam("tab")
	data.ActiveMonth = c.QueryParam("month")

	name := c.FormValue("name")
	ticked := c.FormValue("incoming")
	amount := c.FormValue("amount")
	date := c.FormValue("date")
	category := c.FormValue("category")

	var incoming = false

	value, err := strconv.ParseFloat(amount, 64)
	if err != nil {
		fmt.Println(data.ErrorMessage)
		err := c.Render(http.StatusBadRequest, "error-message", "Error adding this amount, please try again in xx.xx format")
		fmt.Println(err)
	}

	switch {
	case utf8.RuneCountInString(name) > 16:
		data.ErrorMessage = "Name too long, max 16 characters"
	case utf8.RuneCountInString(amount) > 15:
		data.ErrorMessage = "Transaction amount too large, max 999,999,999,999.99"
	case utf8.RuneCountInString(category) > 12:
		data.ErrorMessage = "Category must be 12 characters or less"
	}

	if strings.Compare(data.ErrorMessage, "") != 0 {
		err := c.Render(http.StatusBadRequest, "error-message", data)
		return err
	}

	moneyAmount := int(value * 100)
	if ticked == "on" {
		incoming = true
	}

	err = app.Transactions.CreateTransaction(name, date, category, moneyAmount, userID, incoming)
	if err != nil {
		err := c.Render(http.StatusBadRequest, "error-message", "Error adding transaction, please try again later")
		return err
	}

	transactionCache.Lock()
	delete(transactionCache.cache, userID)
	transactionCache.Unlock()

	CategoriesCache.Lock()
	delete(CategoriesCache.cache, userID)
	CategoriesCache.Unlock()

	data.AllCategories, _ = app.getUserCategories(userID)

	fmt.Println(data.AllCategories)

	var pages = make([][]tableData, 1)

	pages, data.TotalAmount = app.getPages(userID, data.ActiveTab, data.ActiveCategory, data.ActiveMonth)
	if pages == nil {
		pageCount := len(pages)
		data.PageCount = strconv.Itoa(pageCount)
		return c.Render(http.StatusOK, "table-head", data)
	}

	data.PageData = pages[0]
	pageCount := len(pages)
	data.PageCount = strconv.Itoa(pageCount - 1)

	return c.Render(http.StatusOK, "table-head", data)

}

func (app *Application) EditTransaction(c echo.Context) error {

	data := TemplateData{}

	data.ActiveTab = c.QueryParam("tab")
	data.ActiveMonth = c.QueryParam("month")

	userID := app.getUserIdFromSession(c)
	id := c.FormValue("id")
	name := c.FormValue("name")
	amount := c.FormValue("amount")
	date := c.FormValue("date")
	category := c.FormValue("category")
	incoming := c.FormValue("incoming")

	switch {
	case utf8.RuneCountInString(name) > 16:
		data.ErrorMessage = "Name too long, max 16 characters"
	case utf8.RuneCountInString(amount) > 15:
		data.ErrorMessage = "Transaction amount too large, max 999,999,999,999.99"
	case utf8.RuneCountInString(category) > 12:
		data.ErrorMessage = "Category must be 12 characters or less"
	}

	transactionType := false

	if incoming == "on" {
		transactionType = true
	}

	value, err := strconv.ParseFloat(amount, 64)
	if err != nil {
		fmt.Println(data.ErrorMessage)
		err = c.Render(http.StatusBadRequest, "error-message", "Error adding this amount, please try again in xx.xx format")
		if err != nil {
			return err
		}
	}

	transactionId, err := strconv.Atoi(id)
	if err != nil {
		err := c.Render(http.StatusBadRequest, "error-message", "Error try again later")
		if err != nil {
			return err
		}
	}

	moneyAmount := int(value * 100)

	err = app.Transactions.EditTransaction(name, date, category, transactionId, moneyAmount, userID, transactionType)
	fmt.Println(err)

	transactionCache.Lock()
	delete(transactionCache.cache, userID)
	transactionCache.Unlock()

	CategoriesCache.Lock()
	delete(CategoriesCache.cache, userID)
	CategoriesCache.Unlock()

	data.AllCategories, _ = app.getUserCategories(userID)

	fmt.Println(data.AllCategories)

	var pages = make([][]tableData, 1)

	pages, data.TotalAmount = app.getPages(userID, data.ActiveTab, data.ActiveCategory, data.ActiveMonth)
	if pages == nil {
		pageCount := len(pages)
		data.PageCount = strconv.Itoa(pageCount)
		return c.Render(http.StatusOK, "table-head", data)
	}

	data.PageData = pages[0]
	pageCount := len(pages)
	data.PageCount = strconv.Itoa(pageCount - 1)

	return c.Render(http.StatusOK, "table-head", data)

}
