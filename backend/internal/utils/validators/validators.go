package validators

import (
	"encoding/json"
	"net/http"

	httpcommon "chi-mysql-boilerplate/internal/domain/http_common"
	"chi-mysql-boilerplate/internal/utils/helpers"

	"github.com/go-playground/validator/v10"
)

type Validator struct {
	validator *validator.Validate
	helpers   *helpers.Message
}

func NewValidator(helpers *helpers.Message) *Validator {
	validator := validator.New()
	// register any custom validator here
	validator.RegisterValidation("pwdminlen", ValidatePassword)

	return &Validator{
		validator: validator,
		helpers:   helpers,
	}
}

// custom validator for password
func ValidatePassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()
	return len(password) >= 6
}

// parse JSON from the request body and validate the struct
func (v *Validator) BindJSONAndValidate(w http.ResponseWriter, r *http.Request, body interface{}) error {
	if err := helpers.ReadJSON(w, r, body); err != nil {
		v.HandleError(w, err)
	}

	if err := v.validator.Struct(body); err != nil {
		v.HandleError(w, err)
		return err
	}
	return nil
}

func (v *Validator) HandleError(w http.ResponseWriter, err error) {
	var httpErr httpcommon.Error

	switch t := err.(type) {
	case *json.UnmarshalTypeError:
		httpErr = httpcommon.Error{
			Message: httpcommon.ErrorMessage.InvalidDataType,
			Code:    httpcommon.ErrorResponseCode.InvalidDataType,
			Field:   t.Field,
		}
	case *json.SyntaxError:
		httpErr = httpcommon.Error{
			Message: err.Error(),
			Code:    httpcommon.ErrorResponseCode.InvalidRequest,
		}
	case validator.ValidationErrors:
		httpErr = httpcommon.Error{
			Message: err.Error(),
			Code:    httpcommon.ErrorResponseCode.InvalidRequest,
		}
	default:
		httpErr = httpcommon.Error{
			Message: err.Error(),
			Code:    httpcommon.ErrorResponseCode.InvalidRequest,
			Field:   "",
		}
	}

	helpers.MessageLogs.ErrorLog.Println(err)
	helpers.WriteJSON(w, http.StatusBadRequest, httpcommon.NewErrorResponse(httpErr))
}
