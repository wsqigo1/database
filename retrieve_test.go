package redis

import "testing"

func TestGetByPk(t *testing.T) {
	GetByPk()
}

func TestGetToMap(t *testing.T) {
	GetToMap()
}

func TestGetPluck(t *testing.T) {
	GetPluck()
}

func TestGetPluckExp(t *testing.T) {
	GetPluckExp()
}

func TestGetSelect(t *testing.T) {
	GetSelect()
}

func TestGetDistinct(t *testing.T) {
	GetDistinct()
}

func TestWhereMethod(t *testing.T) {
	WhereMethod()
}

func TestWhereType(t *testing.T) {
	WhereType()
}

func TestPlaceHolder(t *testing.T) {
	PlaceHolder()
}

func TestOrderBy(t *testing.T) {
	OrderBy()
}

func TestPagination(t *testing.T) {
	Pagination(Pager{3, 15})
}

func TestPaginationScope(t *testing.T) {
	PaginationScope(Pager{3, 15})
}

func TestGroupHaving(t *testing.T) {
	GroupHaving()
}

func TestCountPage(t *testing.T) {
	CountPage(Pager{3, 15})
}
