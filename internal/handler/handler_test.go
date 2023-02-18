package handler

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"

	"makves-task/internal/entity"
	"makves-task/internal/repository"
)

func TestNewRouter(t *testing.T) {
	t.Run("NewRouter", func(t *testing.T) {
		mock := NewMock()
		router := NewRouter(mock)
		require.NotNil(t, router)
	})
}

func TestController_GetItems(t *testing.T) {
	t.Run("Positive", func(t *testing.T) {
		mock := NewMock()
		r := NewRouter(mock)
		ts := httptest.NewServer(r)
		defer ts.Close()
		req, err := http.NewRequest(http.MethodGet, ts.URL+"/get-items?id=0&id=1", bytes.NewBuffer([]byte("")))
		require.NoError(t, err)
		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)
		defer resp.Body.Close()
		require.Equal(t, http.StatusOK, resp.StatusCode)
	})
}

type Mock struct {
	data map[string]string
}

func NewMock() repository.Databaser {
	r := &Mock{
		data: make(map[string]string),
	}
	r.Init()
	return r
}

func (m *Mock) Init() {
	for i := 0; i < 2; i++ {
		id := strconv.Itoa(i)
		m.data[id] = id
	}
}

func (m *Mock) GetItems(slice []string) []*entity.Node {
	result := make([]*entity.Node, 0)
	for _, v := range slice {
		if val, ok := m.data[v]; ok {
			node := &entity.Node{
				ID: val,
			}
			result = append(result, node)
		}
	}
	return result
}
