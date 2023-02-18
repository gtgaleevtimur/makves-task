package repository

import "makves-task/internal/entity"

type Databaser interface {
	GetItems(slice []string) []*entity.Node
	Init()
}
