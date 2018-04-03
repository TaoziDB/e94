// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package gin

import (
	"crypto/tls"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func formatAsDate(t time.Time) string {
	year, month, day := t.Date()
	return fmt.Sprintf("%d/%02d/%02d", year, month, day)
}

func setupHTMLFiles(t *testing.T, mode string, tls bool) func() {
	go func() {
		SetMode(mode)
		router := New()
		router.Delims("{[{", "}]}")
		router.SetFuncMap(template.FuncMap{
			"formatAsDate": formatAsDate,
		})
		router.LoadHTMLFiles("./fixtures/basic/hello.tmpl", "./fixtures/basic/raw.tmpl")
		router.GET("/test", func(c *Context) {
			c.HTML(http.StatusOK, "hello.tmpl", map[string]string{"name": "world"})
		})
		router.GET("/raw", func(c *Context) {
			c.HTML(http.StatusOK, "raw.tmpl", map[string]interface{}{
				"now": time.Date(2017, 07, 01, 0, 0, 0, 0, time.UTC),
			})
		})
		if tls {
			// these files generated by `go run $GOROOT/src/crypto/tls/generate_cert.go --host 127.0.0.1`
			router.RunTLS(":9999", "./fixtures/testdata/cert.pem", "./fixtures/testdata/key.pem")
		} else {
			router.Run(":8888")
		}
	}()
	t.Log("waiting 1 second for server startup")
	time.Sleep(1 * time.Second)
	return func() {}
}

func setupHTMLGlob(t *testing.T, mode string, tls bool) func() {
	go func() {
		SetMode(mode)
		router := New()
		router.Delims("{[{", "}]}")
		router.SetFuncMap(template.FuncMap{
			"formatAsDate": formatAsDate,
		})
		router.LoadHTMLGlob("./fixtures/basic/*")
		router.GET("/test", func(c *Context) {
			c.HTML(http.StatusOK, "hello.tmpl", map[string]string{"name": "world"})
		})
		router.GET("/raw", func(c *Context) {
			c.HTML(http.StatusOK, "raw.tmpl", map[string]interface{}{
				"now": time.Date(2017, 07, 01, 0, 0, 0, 0, time.UTC),
			})
		})
		if tls {
			// these files generated by `go run $GOROOT/src/crypto/tls/generate_cert.go --host 127.0.0.1`
			router.RunTLS(":9999", "./fixtures/testdata/cert.pem", "./fixtures/testdata/key.pem")
		} else {
			router.Run(":8888")
		}
	}()
	t.Log("waiting 1 second for server startup")
	time.Sleep(1 * time.Second)
	return func() {}
}

func TestLoadHTMLGlob(t *testing.T) {
	td := setupHTMLGlob(t, DebugMode, false)
	res, err := http.Get("http://127.0.0.1:8888/test")
	if err != nil {
		fmt.Println(err)
	}

	resp, _ := ioutil.ReadAll(res.Body)
	assert.Equal(t, "<h1>Hello world</h1>", string(resp[:]))

	td()
}

func TestLoadHTMLGlob2(t *testing.T) {
	td := setupHTMLGlob(t, TestMode, false)
	res, err := http.Get("http://127.0.0.1:8888/test")
	if err != nil {
		fmt.Println(err)
	}

	resp, _ := ioutil.ReadAll(res.Body)
	assert.Equal(t, "<h1>Hello world</h1>", string(resp[:]))

	td()
}

func TestLoadHTMLGlob3(t *testing.T) {
	td := setupHTMLGlob(t, ReleaseMode, false)
	res, err := http.Get("http://127.0.0.1:8888/test")
	if err != nil {
		fmt.Println(err)
	}

	resp, _ := ioutil.ReadAll(res.Body)
	assert.Equal(t, "<h1>Hello world</h1>", string(resp[:]))

	td()
}

func TestLoadHTMLGlobUsingTLS(t *testing.T) {
	td := setupHTMLGlob(t, DebugMode, true)
	// Use InsecureSkipVerify for avoiding `x509: certificate signed by unknown authority` error
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}
	client := &http.Client{Transport: tr}
	res, err := client.Get("https://127.0.0.1:9999/test")
	if err != nil {
		fmt.Println(err)
	}

	resp, _ := ioutil.ReadAll(res.Body)
	assert.Equal(t, "<h1>Hello world</h1>", string(resp[:]))

	td()
}

func TestLoadHTMLGlobFromFuncMap(t *testing.T) {
	time.Now()
	td := setupHTMLGlob(t, DebugMode, false)
	res, err := http.Get("http://127.0.0.1:8888/raw")
	if err != nil {
		fmt.Println(err)
	}

	resp, _ := ioutil.ReadAll(res.Body)
	assert.Equal(t, "Date: 2017/07/01\n", string(resp[:]))

	td()
}

func init() {
	SetMode(TestMode)
}

func TestCreateEngine(t *testing.T) {
	router := New()
	assert.Equal(t, "/", router.basePath)
	assert.Equal(t, router.engine, router)
	assert.Empty(t, router.Handlers)
}

// func TestLoadHTMLDebugMode(t *testing.T) {
// 	router := New()
// 	SetMode(DebugMode)
// 	router.LoadHTMLGlob("*.testtmpl")
// 	r := router.HTMLRender.(render.HTMLDebug)
// 	assert.Empty(t, r.Files)
// 	assert.Equal(t, "*.testtmpl", r.Glob)
//
// 	router.LoadHTMLFiles("index.html.testtmpl", "login.html.testtmpl")
// 	r = router.HTMLRender.(render.HTMLDebug)
// 	assert.Empty(t, r.Glob)
// 	assert.Equal(t, []string{"index.html", "login.html"}, r.Files)
// 	SetMode(TestMode)
// }

func TestLoadHTMLFiles(t *testing.T) {
	td := setupHTMLFiles(t, TestMode, false)
	res, err := http.Get("http://127.0.0.1:8888/test")
	if err != nil {
		fmt.Println(err)
	}

	resp, _ := ioutil.ReadAll(res.Body)
	assert.Equal(t, "<h1>Hello world</h1>", string(resp[:]))
	td()
}

func TestLoadHTMLFiles2(t *testing.T) {
	td := setupHTMLFiles(t, DebugMode, false)
	res, err := http.Get("http://127.0.0.1:8888/test")
	if err != nil {
		fmt.Println(err)
	}

	resp, _ := ioutil.ReadAll(res.Body)
	assert.Equal(t, "<h1>Hello world</h1>", string(resp[:]))
	td()
}

func TestLoadHTMLFiles3(t *testing.T) {
	td := setupHTMLFiles(t, ReleaseMode, false)
	res, err := http.Get("http://127.0.0.1:8888/test")
	if err != nil {
		fmt.Println(err)
	}

	resp, _ := ioutil.ReadAll(res.Body)
	assert.Equal(t, "<h1>Hello world</h1>", string(resp[:]))
	td()
}

func TestLoadHTMLFilesUsingTLS(t *testing.T) {
	td := setupHTMLFiles(t, TestMode, true)
	// Use InsecureSkipVerify for avoiding `x509: certificate signed by unknown authority` error
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}
	client := &http.Client{Transport: tr}
	res, err := client.Get("https://127.0.0.1:9999/test")
	if err != nil {
		fmt.Println(err)
	}

	resp, _ := ioutil.ReadAll(res.Body)
	assert.Equal(t, "<h1>Hello world</h1>", string(resp[:]))
	td()
}

func TestLoadHTMLFilesFuncMap(t *testing.T) {
	time.Now()
	td := setupHTMLFiles(t, TestMode, false)
	res, err := http.Get("http://127.0.0.1:8888/raw")
	if err != nil {
		fmt.Println(err)
	}

	resp, _ := ioutil.ReadAll(res.Body)
	assert.Equal(t, "Date: 2017/07/01\n", string(resp[:]))

	td()
}

func TestAddRoute(t *testing.T) {
	router := New()
	router.addRoute("GET", "/", HandlersChain{func(_ *Context) {}})

	assert.Len(t, router.trees, 1)
	assert.NotNil(t, router.trees.get("GET"))
	assert.Nil(t, router.trees.get("POST"))

	router.addRoute("POST", "/", HandlersChain{func(_ *Context) {}})

	assert.Len(t, router.trees, 2)
	assert.NotNil(t, router.trees.get("GET"))
	assert.NotNil(t, router.trees.get("POST"))

	router.addRoute("POST", "/post", HandlersChain{func(_ *Context) {}})
	assert.Len(t, router.trees, 2)
}

func TestAddRouteFails(t *testing.T) {
	router := New()
	assert.Panics(t, func() { router.addRoute("", "/", HandlersChain{func(_ *Context) {}}) })
	assert.Panics(t, func() { router.addRoute("GET", "a", HandlersChain{func(_ *Context) {}}) })
	assert.Panics(t, func() { router.addRoute("GET", "/", HandlersChain{}) })

	router.addRoute("POST", "/post", HandlersChain{func(_ *Context) {}})
	assert.Panics(t, func() {
		router.addRoute("POST", "/post", HandlersChain{func(_ *Context) {}})
	})
}

func TestCreateDefaultRouter(t *testing.T) {
	router := Default()
	assert.Len(t, router.Handlers, 2)
}

func TestNoRouteWithoutGlobalHandlers(t *testing.T) {
	var middleware0 HandlerFunc = func(c *Context) {}
	var middleware1 HandlerFunc = func(c *Context) {}

	router := New()

	router.NoRoute(middleware0)
	assert.Nil(t, router.Handlers)
	assert.Len(t, router.noRoute, 1)
	assert.Len(t, router.allNoRoute, 1)
	compareFunc(t, router.noRoute[0], middleware0)
	compareFunc(t, router.allNoRoute[0], middleware0)

	router.NoRoute(middleware1, middleware0)
	assert.Len(t, router.noRoute, 2)
	assert.Len(t, router.allNoRoute, 2)
	compareFunc(t, router.noRoute[0], middleware1)
	compareFunc(t, router.allNoRoute[0], middleware1)
	compareFunc(t, router.noRoute[1], middleware0)
	compareFunc(t, router.allNoRoute[1], middleware0)
}

func TestNoRouteWithGlobalHandlers(t *testing.T) {
	var middleware0 HandlerFunc = func(c *Context) {}
	var middleware1 HandlerFunc = func(c *Context) {}
	var middleware2 HandlerFunc = func(c *Context) {}

	router := New()
	router.Use(middleware2)

	router.NoRoute(middleware0)
	assert.Len(t, router.allNoRoute, 2)
	assert.Len(t, router.Handlers, 1)
	assert.Len(t, router.noRoute, 1)

	compareFunc(t, router.Handlers[0], middleware2)
	compareFunc(t, router.noRoute[0], middleware0)
	compareFunc(t, router.allNoRoute[0], middleware2)
	compareFunc(t, router.allNoRoute[1], middleware0)

	router.Use(middleware1)
	assert.Len(t, router.allNoRoute, 3)
	assert.Len(t, router.Handlers, 2)
	assert.Len(t, router.noRoute, 1)

	compareFunc(t, router.Handlers[0], middleware2)
	compareFunc(t, router.Handlers[1], middleware1)
	compareFunc(t, router.noRoute[0], middleware0)
	compareFunc(t, router.allNoRoute[0], middleware2)
	compareFunc(t, router.allNoRoute[1], middleware1)
	compareFunc(t, router.allNoRoute[2], middleware0)
}

func TestNoMethodWithoutGlobalHandlers(t *testing.T) {
	var middleware0 HandlerFunc = func(c *Context) {}
	var middleware1 HandlerFunc = func(c *Context) {}

	router := New()

	router.NoMethod(middleware0)
	assert.Empty(t, router.Handlers)
	assert.Len(t, router.noMethod, 1)
	assert.Len(t, router.allNoMethod, 1)
	compareFunc(t, router.noMethod[0], middleware0)
	compareFunc(t, router.allNoMethod[0], middleware0)

	router.NoMethod(middleware1, middleware0)
	assert.Len(t, router.noMethod, 2)
	assert.Len(t, router.allNoMethod, 2)
	compareFunc(t, router.noMethod[0], middleware1)
	compareFunc(t, router.allNoMethod[0], middleware1)
	compareFunc(t, router.noMethod[1], middleware0)
	compareFunc(t, router.allNoMethod[1], middleware0)
}

func TestRebuild404Handlers(t *testing.T) {

}

func TestNoMethodWithGlobalHandlers(t *testing.T) {
	var middleware0 HandlerFunc = func(c *Context) {}
	var middleware1 HandlerFunc = func(c *Context) {}
	var middleware2 HandlerFunc = func(c *Context) {}

	router := New()
	router.Use(middleware2)

	router.NoMethod(middleware0)
	assert.Len(t, router.allNoMethod, 2)
	assert.Len(t, router.Handlers, 1)
	assert.Len(t, router.noMethod, 1)

	compareFunc(t, router.Handlers[0], middleware2)
	compareFunc(t, router.noMethod[0], middleware0)
	compareFunc(t, router.allNoMethod[0], middleware2)
	compareFunc(t, router.allNoMethod[1], middleware0)

	router.Use(middleware1)
	assert.Len(t, router.allNoMethod, 3)
	assert.Len(t, router.Handlers, 2)
	assert.Len(t, router.noMethod, 1)

	compareFunc(t, router.Handlers[0], middleware2)
	compareFunc(t, router.Handlers[1], middleware1)
	compareFunc(t, router.noMethod[0], middleware0)
	compareFunc(t, router.allNoMethod[0], middleware2)
	compareFunc(t, router.allNoMethod[1], middleware1)
	compareFunc(t, router.allNoMethod[2], middleware0)
}

func compareFunc(t *testing.T, a, b interface{}) {
	sf1 := reflect.ValueOf(a)
	sf2 := reflect.ValueOf(b)
	if sf1.Pointer() != sf2.Pointer() {
		t.Error("different functions")
	}
}

func TestListOfRoutes(t *testing.T) {
	router := New()
	router.GET("/favicon.ico", handlerTest1)
	router.GET("/", handlerTest1)
	group := router.Group("/users")
	{
		group.GET("/", handlerTest2)
		group.GET("/:id", handlerTest1)
		group.POST("/:id", handlerTest2)
	}
	router.Static("/static", ".")

	list := router.Routes()

	assert.Len(t, list, 7)
	assertRoutePresent(t, list, RouteInfo{
		Method:  "GET",
		Path:    "/favicon.ico",
		Handler: "^(.*/vendor/)?github.com/gin-gonic/gin.handlerTest1$",
	})
	assertRoutePresent(t, list, RouteInfo{
		Method:  "GET",
		Path:    "/",
		Handler: "^(.*/vendor/)?github.com/gin-gonic/gin.handlerTest1$",
	})
	assertRoutePresent(t, list, RouteInfo{
		Method:  "GET",
		Path:    "/users/",
		Handler: "^(.*/vendor/)?github.com/gin-gonic/gin.handlerTest2$",
	})
	assertRoutePresent(t, list, RouteInfo{
		Method:  "GET",
		Path:    "/users/:id",
		Handler: "^(.*/vendor/)?github.com/gin-gonic/gin.handlerTest1$",
	})
	assertRoutePresent(t, list, RouteInfo{
		Method:  "POST",
		Path:    "/users/:id",
		Handler: "^(.*/vendor/)?github.com/gin-gonic/gin.handlerTest2$",
	})
}

func assertRoutePresent(t *testing.T, gotRoutes RoutesInfo, wantRoute RouteInfo) {
	for _, gotRoute := range gotRoutes {
		if gotRoute.Path == wantRoute.Path && gotRoute.Method == wantRoute.Method {
			assert.Regexp(t, wantRoute.Handler, gotRoute.Handler)
			return
		}
	}
	t.Errorf("route not found: %v", wantRoute)
}

func handlerTest1(c *Context) {}
func handlerTest2(c *Context) {}
