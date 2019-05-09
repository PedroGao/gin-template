package middleware

import (
	"fmt"
	"github.com/PedroGao/jerry/libs/erro"
	lv "github.com/PedroGao/jerry/libs/validator"
	"github.com/gin-gonic/gin"
	"gopkg.in/go-playground/validator.v9"
	"net/http"
	"unicode"
)

func ErrorHandler(c *gin.Context) {
	c.Next()
	// 取最后一个Error为返回的Error
	length := len(c.Errors)
	if length > 0 {
		e := c.Errors[length-1]
		switch e1 := e.Err.(type) {
		case *erro.HttpErr:
			writeHttpError(c, *e1)
		case erro.HttpErr:
			writeHttpError(c, e1)
		case validator.ValidationErrors:
			//log.Infoln(e1.Translate(lv.Trans))
			mapErrors := make(map[string]string)
			for _, err := range e1 {
				//fieldName, mg := validationErrorToText(err)
				mapErrors[err.Field()] = err.Translate(lv.Trans)
			}
			writeParamError(c, mapErrors)
		case *validator.ValidationErrors:
			mapErrors := make(map[string]string)
			for _, err := range *e1 {
				//log.Infoln(err.Translate(lv.Trans))
				//fieldName, mg := validationErrorToText(err)
				//mapErrors[fieldName] = mg
				mapErrors[err.Field()] = err.Translate(lv.Trans)
			}
			writeParamError(c, mapErrors)
		default:
			writeError(c, e.Err.Error())
		}
	}
}

func validationErrorToText(e validator.FieldError) (string, string) {
	runes := []rune(e.Field())
	runes[0] = unicode.ToLower(runes[0])
	fieldName := string(runes)
	switch e.Tag() {
	case "required":
		return fieldName, fmt.Sprintf("%s is required", fieldName)
	case "max":
		return fieldName, fmt.Sprintf("%s must be less or equal to %s", fieldName, e.Param)
	case "min":
		return fieldName, fmt.Sprintf("%s must be more or equal to %s", fieldName, e.Param)
	}
	return fieldName, fmt.Sprintf("%s: is not valid", fieldName)
}

func writeError(ctx *gin.Context, errString interface{}) {
	status := http.StatusBadRequest
	if ctx.Writer.Status() != http.StatusOK {
		status = ctx.Writer.Status()
	}
	s := erro.UnKnown.SetMsg(errString).SetUrl(ctx.Request.URL.String())
	ctx.JSON(status, s)
}

func writeParamError(ctx *gin.Context, errString interface{}) {
	status := http.StatusBadRequest
	if ctx.Writer.Status() != http.StatusOK {
		status = ctx.Writer.Status()
	}
	s := erro.ParamsErr.SetMsg(errString).SetUrl(ctx.Request.URL.String())
	ctx.JSON(status, s)
}

func writeHttpError(ctx *gin.Context, httpErr erro.HttpErr) {
	s := httpErr.SetUrl(ctx.Request.URL.String())
	ctx.JSON(httpErr.HttpCode, s)
}
