package utils

import (
	"errors"
	"html/template"
	"net/http"
)

var EntityNotFoundErr = errors.New("Entidade não encontrada")
var InternalServerErr = errors.New("Ocorreu um erro inesperado no nosso servidor. Por favor tente novamente daqui a pouco.")

func ErrorResponse(w http.ResponseWriter, event string, srcTmpl *template.Template, targetTmpl string, tmplData interface{}) {
	// w.WriteHeader(statusCode)
	w.Header().Add("HX-Trigger-After-Settle", event)
	srcTmpl.ExecuteTemplate(w, targetTmpl, tmplData)
}

func GeneralErrorResponse(w http.ResponseWriter, err error, srcTmpl *template.Template) {
	w.Header().Add("HX-Retarget", "#general-error-msg")
	w.Header().Add("HX-Reswap", "outerHTML")
	ErrorResponse(w, "general-error", srcTmpl, "general-error-msg", err.Error())
}
