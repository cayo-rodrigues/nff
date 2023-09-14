package utils

import (
	"errors"
	"fmt"
	"html/template"
	"net/http"
)

var EntityNotFoundErr = errors.New("Entidade inexistente")
var InternalServerErr = errors.New("Ocorreu um erro inesperado no nosso servidor. Por favor tente novamente daqui a pouco.")

func ErrorResponse(w http.ResponseWriter, err error, event string, srcTmpl *template.Template, targetTmpl string, tmplData interface{}) {
	// w.WriteHeader(statusCode)
	eventMsg := fmt.Sprintf("{\"%s\": \"%s\"}", event, err.Error())
	w.Header().Add("HX-Trigger-After-Settle", eventMsg)
	srcTmpl.ExecuteTemplate(w, targetTmpl, tmplData)
}

func GeneralErrorResponse(w http.ResponseWriter, err error, srcTmpl *template.Template, targetTmpl string, tmplData interface{}) {
	ErrorResponse(w, err, "general-error", srcTmpl, targetTmpl, tmplData)
}
