package routers

import (
	"apiTest/config"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
)

func TestGetUser(t *testing.T) {
	config.InitConfig()
	r := gin.Default()
	Register(r)

	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			w := httptest.NewRecorder()

			req, _ := http.NewRequest("GET", "/user/", nil)
			r.ServeHTTP(w, req)

			assert.Equal(t, http.StatusOK, w.Code)
			//data, err := io.ReadAll(c.Wr.Body)
			//assert.Equal(t, nil)
			assert.Equal(t, "{\"code\":1,\"data\":{\"id\":1,\"name\":\"张三\"},\"message\":\"ok\"}", w.Body.String())
		}()
	}
	wg.Wait()
}
