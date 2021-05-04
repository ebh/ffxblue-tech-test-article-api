package router

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPostArticleValidation(t *testing.T) {
	tcs := []struct {
		name string
		body string
	}{
		{
			name: "Missing id",
			body: `{
				"title": "latest science shows that potato chips are better for you than sugar",
				"date" : "2016-09-22T00:00:00Z",
				"body" : "some text, potentially containing simple markup about how potato chips are great",
				"tags" : ["health", "fitness", "science"]
			}`,
		},
		{
			name: "Missing title",
			body: `{
				"id": 1,
				"date" : "2016-09-22T00:00:00Z",
				"body" : "some text, potentially containing simple markup about how potato chips are great",
				"tags" : ["health", "fitness", "science"]
			}`,
		},
		{
			name: "Title longer than 100 characters",
			body: `{
				"id": 1,
				"title": "latest science shows that potato chips are better for you than sugar - sooooooooooooooooooo very long",
				"date" : "2016-09-22T00:00:00Z",
				"body" : "some text, potentially containing simple markup about how potato chips are great",
				"tags" : ["health", "fitness", "science"]
			}`,
		},
		{
			name: "Missing date",
			body: `{
				"id": 1,
				"title": "latest science shows that potato chips are better for you than sugar",
				"body" : "some text, potentially containing simple markup about how potato chips are great",
				"tags" : ["health", "fitness", "science"]
			}`,
		},
		{
			name: "Missing body",
			body: `{
				"id": 1,
				"title": "latest science shows that potato chips are better for you than sugar",
				"date" : "2016-09-22T00:00:00Z",
				"tags" : ["health", "fitness", "science"]
			}`,
		},
		{
			name: "Missing tags",
			body: `{
				"id": 1,
				"title": "latest science shows that potato chips are better for you than sugar",
				"date" : "2016-09-22T00:00:00Z",
				"body" : "some text, potentially containing simple markup about how potato chips are great",
			}`,
		},
		{
			name: "Empty set of tags title",
			body: `{
				"id": 1,
				"title": "latest science shows that potato chips are better for you than sugar",
				"date" : "2016-09-22T00:00:00Z",
				"body" : "some text, potentially containing simple markup about how potato chips are great",
				"tags" : []
			}`,
		},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			r := InitRouter()

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("POST", "/v1/articles", bytes.NewBuffer([]byte(tc.body)))
			r.ServeHTTP(w, req)

			assert.Equal(t, http.StatusBadRequest, w.Code)
			assert.JSONEq(t, "{\"msg\":\"invalid parameters\",\"data\":null}", w.Body.String())
		})
	}
}

func TestGetTagValidation(t *testing.T) {
	tcs := []struct {
		path         string
		expectedBody string
	}{
		{
			path:         "/v1/tags/health/2020010",
			expectedBody: "{\"msg\":\"invalid date\",\"data\":null}",
		},
		{
			path:         "/v1/tags//20200101",
			expectedBody: "{\"msg\":\"invalid tag name\",\"data\":null}",
		},
	}

	for _, tc := range tcs {
		t.Run(tc.path, func(t *testing.T) {
			r := InitRouter()

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", tc.path, nil)
			r.ServeHTTP(w, req)

			assert.Equal(t, http.StatusBadRequest, w.Code)
			b := w.Body.String()
			assert.JSONEq(t, tc.expectedBody, b)
		})
	}
}

func TestGetArticleValidation(t *testing.T) {
	r := InitRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/v1/articles/0", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.JSONEq(t, "{\"msg\":\"invalid id\",\"data\":null}", w.Body.String())
}

func TestRoutes(t *testing.T) {
	t.Run("Post article and then make GET calls", func(t *testing.T) {
		articleBody := `{
				"id": 1,
				"title": "latest science shows that potato chips are better for you than sugar",
				"date" : "2016-09-22T00:00:00Z",
				"body" : "some text, potentially containing simple markup about how potato chips are great",
				"tags" : ["health", "fitness", "science"]
			}`

		r := InitRouter()

		// POST
		w := httptest.NewRecorder()
		postReq, _ := http.NewRequest("POST", "/v1/articles", bytes.NewBuffer([]byte(articleBody)))
		r.ServeHTTP(w, postReq)

		require.Equal(t, http.StatusOK, w.Code)
		require.JSONEq(t, "{\"msg\":\"article saved\",\"data\":null}", w.Body.String())

		// GET - Articles
		w = httptest.NewRecorder()
		getArticlesReq, _ := http.NewRequest("GET", "/v1/articles/1", nil)
		r.ServeHTTP(w, getArticlesReq)

		assert.Equal(t, http.StatusOK, w.Code)
		require.JSONEq(t, articleBody, w.Body.String())

		// GET - Tags
		w = httptest.NewRecorder()
		getTagReq, _ := http.NewRequest("GET", "/v1/tags/health/20160922", nil)
		r.ServeHTTP(w, getTagReq)

		assert.Equal(t, http.StatusOK, w.Code)
		tagRespBody := `{
			"tag": "health",
			"count": 1,
			"articles": [1],
			"related_tags":["fitness","science"]
		}`
		require.JSONEq(t, tagRespBody, w.Body.String())
	})

	t.Run("Get article that does not exist", func(t *testing.T) {
		r := InitRouter()

		w := httptest.NewRecorder()
		getArticlesReq, _ := http.NewRequest("GET", "/v1/articles/123", nil)
		r.ServeHTTP(w, getArticlesReq)

		assert.Equal(t, http.StatusNotFound, w.Code)
		require.JSONEq(t, "{\"msg\":\"no article with that ID found\",\"data\":null}", w.Body.String())
	})

	t.Run("Get for tag that does not exist", func(t *testing.T) {
		r := InitRouter()

		// GET - Tags
		w := httptest.NewRecorder()
		getTagReq, _ := http.NewRequest("GET", "/v1/tags/xyz/20160922", nil)
		r.ServeHTTP(w, getTagReq)

		assert.Equal(t, http.StatusOK, w.Code)
		tagRespBody := `{
			"tag": "xyz",
			"count": 0,
			"articles": null,
			"related_tags": null
		}`
		require.JSONEq(t, tagRespBody, w.Body.String())
	})

	t.Run("Get for tag that exists but date does not exist", func(t *testing.T) {
		r := InitRouter()

		// GET - Tags
		w := httptest.NewRecorder()
		getTagReq, _ := http.NewRequest("GET", "/v1/tags/health/20200922", nil)
		r.ServeHTTP(w, getTagReq)

		assert.Equal(t, http.StatusOK, w.Code)
		tagRespBody := `{
			"tag": "health",
			"count": 0,
			"articles": null,
			"related_tags": null
		}`
		require.JSONEq(t, tagRespBody, w.Body.String())
	})

	t.Run("Post article twice with same ID", func(t *testing.T) {
		body := `{
				"id": 2,
				"title": "latest science shows that potato chips are better for you than sugar",
				"date" : "2016-09-22T00:00:00Z",
				"body" : "some text, potentially containing simple markup about how potato chips are great",
				"tags" : ["health", "fitness", "science"]
			}`

		r := InitRouter()

		w := httptest.NewRecorder()
		postReq1, _ := http.NewRequest("POST", "/v1/articles", bytes.NewBuffer([]byte(body)))
		r.ServeHTTP(w, postReq1)

		require.Equal(t, http.StatusOK, w.Code)
		require.JSONEq(t, "{\"msg\":\"article saved\",\"data\":null}", w.Body.String())

		w = httptest.NewRecorder()
		postReq2, _ := http.NewRequest("POST", "/v1/articles", bytes.NewBuffer([]byte(body)))
		r.ServeHTTP(w, postReq2)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		require.JSONEq(t, "{\"msg\":\"article with that ID already exists\",\"data\":null}", w.Body.String())
	})
}
