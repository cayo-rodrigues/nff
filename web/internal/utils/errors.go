package utils

type EntityNotFoundError struct{}
type InternalServerError struct{}

func (*EntityNotFoundError) Error() string {
	return "Entidade não encontrada"
}
func (*InternalServerError) Error() string {
	return "Ocorreu um erro inesperado. Por favor tente novamente daqui a pouco."
}
