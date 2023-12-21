package globals

const (
	EntityNotFoundMsg    = "Entidade não encontrada"
	InvoiceNotFoundMsg   = "NFA não encontrada"
	CancelingNotFoundMsg = "Cancelamento de NFA não encontrado"
	MetricsNotFoundMsg   = "Registro de métricas não encontrado"
	PrintingNotFoundMsg  = "Impressão de NFA não encontrada"
	UserNotFoundMsg      = "Usuário não encontrado"
)

const (
	InternalServerErrMsg = "Ocorreu um erro inesperado no nosso servidor. Por favor tente novamente daqui a pouco."
)

const (
	MandatoryFieldMsg      = "Campo obrigatório"
	ValueTooLongMsg        = "Valor maior do que o suportado"
	InvalidFormatMsg       = "Formato inválido"
	IlogicalDatesMsg       = "Data inicial deve ser anterior à final"
	TimeRangeTooLongMsg    = "Período não pode ser maior que 365 dias"
	UnacceptableValueMsg   = "Valor inaceitável"
	MustHaveItemsMsg       = "A NF deve ter pelo menos 1 produto"
	InvalidItemsMsg        = "Dados dos produtos inválidos"
	MustHaveIeOrAddressMsg = "Ie OU endereço completo obrigatórios"
)

const (
	ReqCardNotVisibleMsg = "A operação foi inciada com sucesso, porém devido aos filtros não aparecerá imediatamente na lista de requerimentos."
)

const (
	DefaultFiltersDaysRange = 10
)
