package server

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"paccount/database"
	"paccount/models"
	"path/filepath"
	"strings"
	"testing"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"

	"github.com/stretchr/testify/assert"
)

func setupTestCase(t *testing.T) (func(t *testing.T), database.IGormRepo) {
	dbFile := filepath.Base("test_acc.db")
	d, err := gorm.Open("sqlite3", dbFile)
	if err != nil {
		t.Error(err)
	}
	d.AutoMigrate(&models.Account{},
		&models.OprType{},
		&models.Transaction{},
	)

	db := database.NewGormRepo(d)

	return func(t *testing.T) {
		err := os.Remove(dbFile)
		if err != nil {
			t.Error(err)
		}
		d.Close()
	}, db
}

func performRequest(r http.Handler, method,
	path string, body io.Reader) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, body)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func TestServerHandlerEndpoint(t *testing.T) {

	teardownTestCase, db := setupTestCase(t)
	defer teardownTestCase(t)

	tests := []struct {
		endpoint string
		expected string
		method   string
		body     string
		code     int
	}{
		{
			endpoint: "/transactions",
			method:   "POST",
			code:     201,
			body:     `{"account_id":1,"amount":-1.005,"operation_id":4}`,
			expected: `{"account_id":1,"amount":-1.005,"created_at":1590981025,"id":4,"operation_id":4}`,
		},

		{
			endpoint: "/transactions/account/1",
			method:   "GET",
			code:     200,
			body:     "",
			expected: `[{"account_id":1,"amount":-1.005,"created_at":1590981025,"id":4,"operation_id":4}]`,
		},

		{
			endpoint: "/accounts",
			method:   "POST",
			code:     201,
			body:     `{"document_number":12312}`,
			expected: `{"created_at":1590984657,"document_number":12312,"id":1,"updated_at":1590984657}`,
		},

		{
			endpoint: "/accounts/1",
			method:   "GET",
			code:     200,
			body:     "",
			expected: `{"created_at":1590984657,"document_number":12312,"id":1,"updated_at":1590984657}`,
		},
	}

	// options server...
	opts := func(s *Server) {
		s.DB = db
	}

	server := New(opts)
	server.Seeds()
	server.PublicRoutes()

	for _, test := range tests {

		t.Run(test.endpoint, func(t *testing.T) {

			var w *httptest.ResponseRecorder
			if test.body == "" {
				w = performRequest(server.Engine, test.method, test.endpoint, nil)
			} else {
				w = performRequest(server.Engine, test.method, test.endpoint, strings.NewReader(test.body))
			}

			var response map[string]interface{}
			respErr := json.Unmarshal(w.Body.Bytes(), &response)

			var testErr error
			var rtest map[string]interface{}
			if test.body == "" {
				testErr = json.Unmarshal([]byte(test.expected), &rtest)
			} else {
				testErr = json.Unmarshal([]byte(test.body), &rtest)
			}

			switch test.endpoint {
			case "/transactions":
				delete(response, "created_at")
				delete(response, "id")

				assert.Nil(t, respErr)
				assert.Nil(t, testErr)
				assert.Equal(t, response, rtest)
				assert.Equal(t, w.Code, test.code)

			case "/transactions/1":
				body, errMarshal := json.Marshal(response)
				bstr := strings.TrimSpace(string(body))

				assert.Nil(t, respErr)
				assert.Nil(t, errMarshal)
				assert.Equal(t, body, bstr)
				assert.Equal(t, w.Code, test.code)

			case "/accounts", "/accounts/1":
				delete(response, "created_at")
				delete(response, "updated_at")
				delete(response, "id")
				delete(rtest, "created_at")
				delete(rtest, "updated_at")
				delete(rtest, "id")

				assert.Nil(t, respErr)
				assert.Nil(t, testErr)
				assert.Equal(t, response, rtest)
				assert.Equal(t, w.Code, test.code)

			}

		})
	}

}
