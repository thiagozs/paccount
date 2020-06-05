package server

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"paccount/database"
	"paccount/pkg/account"
	"paccount/pkg/oprtype"
	"paccount/pkg/transaction"
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
	d.AutoMigrate(&account.Entity{},
		&transaction.Entity{},
		&oprtype.Entity{},
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
		payload  string
		code     int
	}{
		{
			endpoint: "/accounts",
			method:   "POST",
			code:     201,
			payload:  `{"document_number":12312,"limit":100}`,
			expected: `{"created_at":1590984657,"document_number":12312,"id":1,"limit":100,"updated_at":1590984657}`,
		},

		{
			endpoint: "/accounts/1",
			method:   "GET",
			code:     200,
			payload:  "",
			expected: `{"created_at":1590984657,"document_number":12312,"id":1,"limit":100,"updated_at":1590984657}`,
		},

		{
			endpoint: "/transactions",
			method:   "POST",
			code:     201,
			payload:  `{"account_id":1,"amount":-1,"operation_id":1}`,
			expected: `{"account_id":1,"amount":-1,"balance":-1,"created_at":1590981025,"id":1,"operation_id":1}`,
		},

		{
			endpoint: "/transactions/account/1",
			method:   "GET",
			code:     200,
			payload:  "",
			expected: `[{"account_id":1,"amount":1,"created_at":1590981025,"id":1,"operation_id":1}]`,
		},
	}

	// options server...
	opts := func(s *Server) {
		s.DB = db
	}

	server := New(opts)
	server.PublicRoutes()

	for _, test := range tests {

		t.Run(test.endpoint, func(t *testing.T) {

			var w *httptest.ResponseRecorder
			if test.payload == "" {
				w = performRequest(server.Engine, test.method, test.endpoint, nil)
			} else {
				w = performRequest(server.Engine, test.method, test.endpoint, strings.NewReader(test.payload))
			}

			var response map[string]interface{}
			respErr := json.Unmarshal(w.Body.Bytes(), &response)

			var testErr error
			var rtest map[string]interface{}
			//if test.payload == "" {
			testErr = json.Unmarshal([]byte(test.expected), &rtest)
			//} else {
			//testErr = json.Unmarshal([]byte(test.body), &rtest)
			//}

			switch test.endpoint {
			case "/transactions":
				delete(response, "updated_at")
				delete(response, "created_at")
				delete(response, "id")

				delete(rtest, "updated_at")
				delete(rtest, "created_at")
				delete(rtest, "id")

				// t.Logf("%#v\n", response)

				assert.Nil(t, respErr)
				assert.Nil(t, testErr)
				assert.Equal(t, rtest, response)
				assert.Equal(t, test.code, w.Code)

			case "/transactions/account/1":
				body, errMarshal := json.Marshal(response)
				bstr := strings.TrimSpace(string(body))

				t.Logf("%#v\n", response)

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
