package main

import (
	"icssight/pcap-service/controller"
	"icssight/seeder/inc"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setup()  {}
func result() {}

func TestMain(m *testing.M) {
	setup()
	ret := m.Run()
	if ret == 0 {
		result()
	}
	os.Exit(ret)
}

func TestA(t *testing.T) {
}

func TestB(t *testing.T) {
}

func TestTopla(t *testing.T) {

	// assert.Equal(t, 3, Topla([]int{1, 2}), "they should be equal")

	/* var v int
	v = Topla([]int{1, 2})
	if v != 3 {
		t.Error("Expected 3, got ", v)
	} */
}

func TestGetHandlerSingleValue(t *testing.T) {
	db := inc.InitDB()
	defer db.Close()

	router := gin.Default()
	router.GET("/:id", controller.GetHandler)

	w := httptest.NewRecorder()
	reqEmpty, _ := http.NewRequest("GET", "/", nil)
	router.ServeHTTP(w, reqEmpty)
	assert.Equal(t, 404, w.Code)

	w = httptest.NewRecorder()
	reqInt, _ := http.NewRequest("GET", "/424", nil)
	router.ServeHTTP(w, reqInt)
	assert.Equal(t, 200, w.Code)

	w = httptest.NewRecorder()
	reqStr, _ := http.NewRequest("GET", "/asd", nil)
	router.ServeHTTP(w, reqStr)
	assert.Equal(t, 404, w.Code)

	// assert.Equal(t, "pong", w.Body.String())
}

func TestGetHandlerMultiValue(t *testing.T) {
	db := inc.InitDB()
	defer db.Close()

	router := gin.Default()
	router.GET("/:id/:count", controller.GetHandler)

	w := httptest.NewRecorder()
	reqPing, _ := http.NewRequest("GET", "/erhan/download", nil)
	router.ServeHTTP(w, reqPing)
	assert.Equal(t, 200, w.Code)
}
