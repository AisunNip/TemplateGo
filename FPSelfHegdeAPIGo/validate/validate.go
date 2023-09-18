package validate

import (
	"FPSelfHegdeAPIGo/common"
	"FPSelfHegdeAPIGo/constant"
	"FPSelfHegdeAPIGo/logging"
	"errors"
	"fmt"

	"gopkg.in/go-playground/validator.v9"
	validate "gopkg.in/go-playground/validator.v9"
)

type Validate struct {
	transID string
}

func (this *Validate) SetTransID(transID string) {
	this.transID = transID
}

func (this *Validate) GetTransID() string {
	return this.transID
}

var logger = logging.InitScheduleLogger("CRMOPrepaidGo", logging.FpOutbound)

func (this *Validate) CheckRequest(bean interface{}) (response common.ResponseBean, err error) {

	logger.Info(this.transID, "Start process function CheckRequest")

	valid := validate.New()
	err = valid.Struct(bean)

	if err != nil {
		for _, e := range err.(validate.ValidationErrors) {
			resp, err := this.mapErrorResponse(e.Tag(), e.Field(), e.Param())
			if err != nil {
				logger.Error(this.transID, "Process mapErrorResponse Exception : %s", err.Error())
				return response, err
			}
			response = resp
			return response, err
		}
	}
	return
}

func (this *Validate) mapErrorResponse(tag string, field string, param string) (response common.ResponseBean, err error) {
	field = common.GetLowerFirstVariable(field)
	var statusCode string
	var message string
	switch tag {
	case constant.Required:
		statusCode = constant.RequiredField
		message = "Error : " + field + " is required."
	case constant.Max:
		statusCode = constant.OverScope
		message = "Error : " + field + " is not over " + param
	case constant.Min:
		statusCode = constant.UnderScope
		message = "Error : " + field + " must be greater than " + param
	case constant.Length:
		statusCode = constant.OverLengthCode
		message = "Error : " + field + " length not more than " + param
	case constant.Numeric:
		statusCode = constant.NoMatchTypeNumeric
		message = "Error : " + field + " is not numeric."
	default:
		statusCode = constant.ConditionNotFound
		message = "Other error."
	}
	response.Code = statusCode
	response.Msg = message
	response.TransID = this.transID
	return
}

func ValidateStruct(dataStruct interface{}) error {
	validate := validate.New()
	err := validate.Struct(dataStruct)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			return errors.New(fmt.Sprintf("%s: %s", err.StructField(), err.Tag()))
		}
	} else {
		return nil
	}
	return err
}
