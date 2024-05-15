package models

import (
	"fmt"
	"strings"
)

type Filters struct {
	query  *strings.Builder
	values []any
}

func (f *Filters) Where(condition string) *Filters {
	f.query.WriteString(" WHERE " + condition)
	return f
}

func (f *Filters) And(condition string) *Filters {
	f.query.WriteString(" AND " + condition)
	return f
}

func (f *Filters) Or() *Filters {
	f.query.WriteString(" OR ")
	return f
}

func (f *Filters) ILike(condition ...string) *Filters {
	conditionsCount := len(condition)
	if conditionsCount == 0 {
		f.query.WriteString(" ILIKE ")
		return f
	}

	if conditionsCount == 1 {
		f.query.WriteString(" ILIKE '%' || " + condition[0] + " || '%'")
		return f
	}

	cond := new(strings.Builder)
	for _, c := range condition {
		cond.WriteString(c)
	}
	f.query.WriteString(" ILIKE '%' || " + cond.String() + " || '%'")
	return f
}

func (f *Filters) OrderBy(column string) *Filters {
	f.query.WriteString(" ORDER BY " + column)
	return f
}

func (f *Filters) Asc() *Filters {
	f.query.WriteString(" ASC ")
	return f
}

func (f *Filters) Desc() *Filters {
	f.query.WriteString(" DESC ")
	return f
}

func (f *Filters) String() string {
	return f.query.String()
}

func (f *Filters) AppendValue(v any) {
	f.values = append(f.values, v)
}

func (f *Filters) Values() []any {
	return f.values
}

func (f *Filters) ValuesCount() int {
	return len(f.values)
}

func (f *Filters) Placeholder(value any) *Filters {
	position := f.ValuesCount() + 1
	f.AppendValue(value)

	f.query.WriteString(fmt.Sprintf("$%d", position))

	return f
}

func (f *Filters) WildPlaceholder(value any) *Filters {
	position := f.ValuesCount() + 1
	f.AppendValue(value)

	f.query.WriteString(fmt.Sprintf("%s || $%d || %s", "'%'", position, "'%'"))

	return f
}

func NewFilters() *Filters {
	return &Filters{
		query: new(strings.Builder),
	}
}
