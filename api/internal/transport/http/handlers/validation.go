package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

// Register custom struct field names
func init() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
			if name == "-" || name == "" {
				return fld.Name
			}
			return name
		})
	}

}

func PayloadErrorJSON(err error) map[string]any {
	var (
		result = map[string]any{
			"validation_error": false,
		}
		verr   validator.ValidationErrors
		unmErr = new(json.UnmarshalTypeError)
	)

	if errors.As(err, &verr) {
		errorMsgs := GetValidationErrorMessages(verr)
		result["validation_error"] = true
		result["errors"] = errorMsgs

	} else if errors.As(err, &unmErr) {
		result["validation_error"] = true
		result["errors"] = GetUnmarshalErrorMessages(*unmErr)

	} else {
		result["detail"] = err.Error()
	}
	return result
}

func GetValidationErrorMessages(errs validator.ValidationErrors) map[string]string {
	var errorMessages = make(map[string]string)
	for _, fieldErr := range errs {
		var msg string
		k := fieldErr.Field()

		switch fieldErr.Kind() {
		case reflect.Struct:
			if k == "" {
				k = fieldErr.Tag()
				msg = "field is required"
			} else {
				switch fieldErr.Tag() {
				case "required":
					msg = "required"
				case "min":
					msg = fmt.Sprintf("%s should be more than %s", k, fieldErr.Param())
				}
			}
		case reflect.Slice:
			k = fieldErr.StructField()
			k = strings.ToLower(k) // stub

			switch fieldErr.Tag() {
			case "required":
				msg = "field is required"
			case "min":
				msg = fmt.Sprintf("%s should be more than %s", k, fieldErr.Param())
			default:
				msg = "unknown"
			}
		case reflect.String:
			if ok, _ := regexp.MatchString(`\[*\]`, k); ok { // checks for dive
				tmp := strings.Split(k, "[")[0]
				k = strings.ToLower(tmp)
			}
			msg = "field is required"
		default:
			msg = "unknown error"
		}

		errorMessages[k] = msg
	}
	return errorMessages
}

func GetUnmarshalErrorMessages(err json.UnmarshalTypeError) map[string]string {
	var msg string
	errors := map[string]string{}

	field := err.Field

	switch err.Type.Kind() {
	case reflect.Slice:
		msg = "field should be array of " + err.Type.Elem().String()
	case reflect.Struct:
		if field == "" {
			field = err.Type.Field(0).Tag.Get("json")
		}
		msg = "field is required" //TODO: stub
	}
	errors[field] = msg
	return errors

}
