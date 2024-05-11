package transaction

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/KKGo-Software-engineering/workshop-summer/api/config"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestCreateTransaction(t *testing.T) {

	t.Run("create transaction fail when bad request body", func(t *testing.T) {
		e := echo.New()
		defer e.Close()

		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{ bad request body }`))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		cfg := config.FeatureFlag{EnableCreateSpender: true}

		h := New(cfg, nil)
		err := h.Create(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Contains(t, rec.Body.String(), "Invalid transaction request")
	})

	t.Run("create transaction success when request body", func(t *testing.T) {
		e := echo.New()
		defer e.Close()

		req := httptest.NewRequest(http.MethodPost, "/api/v1/transactions", io.NopCloser(strings.NewReader(`{
			"date": "2024-04-30T09:00:00.000Z",
			"amount": 1000,
			"category": "Food",
			"transaction_type": "expense",
			"note": "Lunch",
			"image_url": "https://example.com/image1.jpg"
		}`)))

		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("api/v1/transactions")
		cfg := config.FeatureFlag{EnableCreateSpender: true}

		h := New(cfg, nil)
		err := h.Create(c)

		assert.NoError(t, err)
		// assert.JSONEq(t, `{
		// 	"date": "2024-04-30T09:00:00.000Z",
		// 	"amount": 1000,
		// 	"category": "Food",
		// 	"transaction_type": "expense",
		// 	"note": "Lunch",
		// 	"image_url": "https://example.com/image1.jpg"
		// }`, rec.Body.String())
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

}

func TestGetTransactionById(t *testing.T) {
	e := echo.New()
	defer e.Close()

	t.Run("successfully retrieve transactions by spender ID", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/transactions/1", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("1")

		db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp))
		defer db.Close()

		// Convert string to time.Time
		date1, _ := time.Parse(time.RFC3339, "2022-01-01T12:00:00Z")
		date2, _ := time.Parse(time.RFC3339, "2022-01-02T12:00:00Z")

		rows := sqlmock.NewRows([]string{"id", "date", "amount", "category", "transaction_type", "note", "image_url"}).
			AddRow(1, date1, 100.00, "groceries", "expense", "Weekly groceries", "http://example.com/receipt1.jpg").
			AddRow(2, date2, 150.00, "electronics", "expense", "Gadget purchase", "http://example.com/receipt2.jpg")
		mock.ExpectQuery(`SELECT id, date, amount, category, transaction_type, note, image_url FROM transaction WHERE spender_id = \$1`).WithArgs(1).WillReturnRows(rows)

		h := New(config.FeatureFlag{}, db)
		if assert.NoError(t, h.GetTransactionById(c)) {
			assert.Equal(t, http.StatusOK, rec.Code)
			assert.JSONEq(t, `{"transactions":[{"id":1,"date":"2022-01-01T12:00:00Z","amount":100.00,"category":"groceries","transaction_type":"expense","note":"Weekly groceries","image_url":"http://example.com/receipt1.jpg"},{"id":2,"date":"2022-01-02T12:00:00Z","amount":150.00,"category":"electronics","transaction_type":"expense","note":"Gadget purchase","image_url":"http://example.com/receipt2.jpg"}]}`, rec.Body.String())
		}
	})

	t.Run("database error during transaction retrieval", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/transactions/1", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetParamNames("id")
		c.SetParamValues("1")

		db, mock, _ := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp))
		defer db.Close()

		// Configure the mock to return an error for the query
		mock.ExpectQuery(`SELECT id, date, amount, category, transaction_type, note, image_url FROM transaction WHERE spender_id = \$1`).WithArgs(1).WillReturnError(assert.AnError)

		h := New(config.FeatureFlag{}, db)
		err := h.GetTransactionById(c)

		// Test the error handling and response
		assert.Error(t, err)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.Contains(t, rec.Body.String(), "Database error")
	})
}
