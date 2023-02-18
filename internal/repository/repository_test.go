package repository

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"

	"makves-task/internal/entity"
)

func TestGormDB_GetItems(t *testing.T) {
	t.Run("Positive", func(t *testing.T) {
		mock := NewMock()
		result := mock.GetItems([]string{"0", "1"})
		require.NotNil(t, result)
	})
	t.Run("Negative", func(t *testing.T) {
		mock := NewMock()
		result := mock.GetItems([]string{})
		require.Equal(t, result, []*entity.Node{})
	})
}

type Mock struct {
	data map[string]string
}

func NewMock() Databaser {
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
