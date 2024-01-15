package main

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"todo_list/models"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Database interface {
	Create(value interface{}) *gorm.DB
}

type MockDB struct {
	mock.Mock
}

func (m *MockDB) FindAll(out interface{}, conditions ...interface{}) *gorm.DB {
	args := m.Called(append([]interface{}{out}, conditions...)...)
	return args.Get(0).(*gorm.DB)
}

func (m *MockDB) Find(out interface{}, where ...interface{}) *gorm.DB {
	args := m.Called(append([]interface{}{out}, where...)...)
	return args.Get(0).(*gorm.DB)
}

func (m *MockDB) Create(value interface{}) *gorm.DB {
	args := m.Called(value)
	return args.Get(0).(*gorm.DB)
}

var mockDB *gorm.DB
var mockServer *echo.Echo

func TestMain(m *testing.M) {
	var errConn error
	mockDB, errConn = gorm.Open(postgres.Open("host=localhost user=postgres password=pgsql089 dbname=todo_list port=5432"), &gorm.Config{})

	if errConn != nil {
		log.Fatal("cannot connect to db:", errConn)
	}

	models.DB = mockDB
	mockServer = echo.New()

	os.Exit(m.Run())
}

func TestGetAllLists(t *testing.T) {
	mockDB := new(MockDB)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := mockServer.NewContext(req, rec)

	expectedLists := []models.List{
		{
			ID:          1,
			Title:       "title again",
			Description: "desc again",
			SubLists:    nil,
			Files:       "6726269950565090337.txt, 88281261534593740.txt",
		},
		{
			ID:          2,
			Title:       "post title again",
			Description: "post desc again",
			SubLists:    nil,
			Files:       "6960174124068648599.pdf, 7377041218994731760.pdf",
		},
	}

	mockDB.On("FindAll", mock.Anything, &[]models.List{}).Return(
		mockDB,
		nil,
	).Run(func(args mock.Arguments) {
		arg := args.Get(0).(*[]models.List)
		*arg = expectedLists
	})

	err := GetAllLists(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var responseLists []models.List
	err = json.Unmarshal(rec.Body.Bytes(), &responseLists)
	assert.NoError(t, err)
	assert.Equal(t, expectedLists, responseLists)
}

func TestGetList(t *testing.T) {
	mockDB := new(MockDB)

	req := httptest.NewRequest(http.MethodGet, "/list/2", nil)
	rec := httptest.NewRecorder()
	c := mockServer.NewContext(req, rec)
	c.SetPath("/list/:id")
	c.SetParamNames("id")
	c.SetParamValues("2")

	expectedList := models.List{
		ID:          2,
		Title:       "post title again",
		Description: "post desc again",
		SubLists:    nil,
		Files:       "6960174124068648599.pdf, 7377041218994731760.pdf",
	}

	mockDB.On("Find", mock.Anything, "id = ?", "2").Return(mockDB, gorm.ErrRecordNotFound).Run(func(args mock.Arguments) {
		arg := args.Get(0).(*models.List)
		*arg = expectedList
	})

	err := GetList(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var responseList models.List
	err = json.Unmarshal(rec.Body.Bytes(), &responseList)
	assert.NoError(t, err)
	assert.Equal(t, expectedList, responseList)
}

func TestGetOneSubList(t *testing.T) {
	mockDB := new(MockDB)

	req := httptest.NewRequest(http.MethodGet, "/subList/7", nil)
	rec := httptest.NewRecorder()
	c := mockServer.NewContext(req, rec)
	c.SetPath("/subList/:id")
	c.SetParamNames("id")
	c.SetParamValues("7")

	expectedList := models.SubList{
		ID:           7,
		ListID:       1,
		Title:        "title subList one",
		Description:  "description subList one",
		SubListFiles: "2460808849991462451.pdf, 430019192659676733.pdf",
	}

	mockDB.On("Find", mock.Anything, "id = ?", "7").Return(mockDB, gorm.ErrRecordNotFound).Run(func(args mock.Arguments) {
		arg := args.Get(0).(*models.SubList)
		*arg = expectedList
	})

	err := GetOneSubList(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var responseList models.SubList
	err = json.Unmarshal(rec.Body.Bytes(), &responseList)
	assert.NoError(t, err)
	assert.Equal(t, expectedList, responseList)
}
