package glue

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"os"

	"encoding/json"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
)

type Fruit struct {
	ID    uint   `gorm:"primary_key" json:"id"`
	Name  string `json:"name"`
	Color string `json:"color"`
	Taste string `json:"taste"`
}

func beforeTest() (e *echo.Echo) {
	e = echo.New()

	db, _ := gorm.Open("sqlite3", "test.db")
	db.LogMode(true)
	db.AutoMigrate(&Fruit{})

	Glue(e.Group("/api"), db, &Fruit{})

	return
}

var Apple = &Fruit{
	Name:  "Apple",
	Color: "Red",
	Taste: "Sweet",
}

var Pear = &Fruit{
	Name:  "Pear",
	Color: "Green",
	Taste: "Sweet",
}

func TestNew(t *testing.T) {
	defer cleanup()
	e := beforeTest()
	code, body := request(GET, "/api/fruits", nil, e)
	assert.Equal(t, http.StatusOK, code)
	assert.Equal(t, "[]", body)
}

func TestPost(t *testing.T) {
	defer cleanup()
	e := beforeTest()
	var (
		code int
		body string
	)

	code, body = request(POST, "/api/fruits", Apple, e)
	assert.Equal(t, http.StatusOK, code)
	assert.Equal(t, "{\"id\":1,\"name\":\"Apple\",\"color\":\"Red\",\"taste\":\"Sweet\"}", body)

	code, body = request(POST, "/api/fruits", Pear, e)
	assert.Equal(t, http.StatusOK, code)
	assert.Equal(t, "{\"id\":2,\"name\":\"Pear\",\"color\":\"Green\",\"taste\":\"Sweet\"}", body)

	code, body = request(GET, "/api/fruits", nil, e)
	assert.Equal(t, http.StatusOK, code)
	assert.Equal(t, "[{\"id\":1,\"name\":\"Apple\",\"color\":\"Red\",\"taste\":\"Sweet\"},{\"id\":2,\"name\":\"Pear\",\"color\":\"Green\",\"taste\":\"Sweet\"}]", body)
}

func TestPut(t *testing.T) {
	defer cleanup()
	e := beforeTest()
	var (
		code int
		body string
	)

	code, body = request(POST, "/api/fruits", Apple, e)
	assert.Equal(t, http.StatusOK, code)
	assert.Equal(t, "{\"id\":1,\"name\":\"Apple\",\"color\":\"Red\",\"taste\":\"Sweet\"}", body)

	code, body = request(PUT, "/api/fruits/1", Pear, e)
	assert.Equal(t, http.StatusOK, code)
	assert.Equal(t, "{\"id\":1,\"name\":\"Pear\",\"color\":\"Green\",\"taste\":\"Sweet\"}", body)

	code, body = request(GET, "/api/fruits", nil, e)
	assert.Equal(t, http.StatusOK, code)
	assert.Equal(t, "[{\"id\":1,\"name\":\"Pear\",\"color\":\"Green\",\"taste\":\"Sweet\"}]", body)
}

func TestPatch(t *testing.T) {
	defer cleanup()
	e := beforeTest()
	var (
		code int
		body string
	)

	code, body = request(POST, "/api/fruits", Apple, e)
	assert.Equal(t, http.StatusOK, code)
	assert.Equal(t, "{\"id\":1,\"name\":\"Apple\",\"color\":\"Red\",\"taste\":\"Sweet\"}", body)

	code, body = request(POST, "/api/fruits", Pear, e)
	assert.Equal(t, http.StatusOK, code)
	assert.Equal(t, "{\"id\":2,\"name\":\"Pear\",\"color\":\"Green\",\"taste\":\"Sweet\"}", body)

	code, body = request(PATCH, "/api/fruits/2", &Fruit{Taste: "Bitter"}, e)
	assert.Equal(t, http.StatusOK, code)
	assert.Equal(t, "{\"id\":2,\"name\":\"Pear\",\"color\":\"Green\",\"taste\":\"Bitter\"}", body)

	code, body = request(GET, "/api/fruits", nil, e)
	assert.Equal(t, http.StatusOK, code)
	assert.Equal(t, "[{\"id\":1,\"name\":\"Apple\",\"color\":\"Red\",\"taste\":\"Sweet\"},{\"id\":2,\"name\":\"Pear\",\"color\":\"Green\",\"taste\":\"Bitter\"}]", body)
}

const (
	GET   = "GET"
	POST  = "POST"
	PUT   = "PUT"
	PATCH = "PATCH"
)

func request(method string, path string, requestBody interface{}, e *echo.Echo) (statusCode int, body string) {
	var req *http.Request
	if requestBody == nil {
		req = httptest.NewRequest(method, path, nil)
	} else {
		jsondata, _ := json.Marshal(requestBody)
		req = httptest.NewRequest(method, path, bytes.NewReader(jsondata))
		req.Header.Add("Content-Type", "application/json")
	}
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	return rec.Code, rec.Body.String()
}

func cleanup() {
	fmt.Println("Removing db")
	os.Remove("test.db")

}
