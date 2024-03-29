package services

type scanner interface {
	Scan(args ...any) error
}

type scannableModel interface {
	Values() []any
}

func scan(row scanner, models ...scannableModel) error {
	values := []any{}
	for _, model := range models {
		values = append(values, model.Values()...)
	}
	return row.Scan(values...)
}
