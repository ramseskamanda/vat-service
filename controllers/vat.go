package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ramseskamanda/vat-service/config"
	"github.com/ramseskamanda/vat-service/pkg/soap"
)

//XXX: move to /models if they grow
type CheckVatRequest struct {
	VATNumber string `json:"vatNumber"`
}

type CheckVatResponse struct {
	VATNumber string `json:"vatNumber"`
	Valid     bool   `json:"valid"`
	Message   string `json:"message"`
}

type VATController struct{}

func (controller *VATController) CheckVATNumber(ctx *gin.Context) {
	var req CheckVatRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	skip := config.GetBool("services.vat.skipPreflight")
	if !skip {
		if err := controller.preflightCheck(&req); err != nil {
			ctx.JSON(http.StatusOK, &CheckVatResponse{VATNumber: req.VATNumber, Valid: false, Message: err.Error()})
			return
		}
	}

	valid, err := controller.remoteCheck(&req)
	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, &CheckVatResponse{VATNumber: req.VATNumber, Valid: valid})
}

func (controller *VATController) preflightCheck(req *CheckVatRequest) error {
	vat := req.VATNumber

	if len(vat) != 11 {
		return errors.New("VAT number should be `DE` followed by 9 digits")
	}

	countryCode, vatNumber := vat[:2], vat[2:]

	if countryCode != "DE" {
		return errors.New("VAT numbers are restricted to Germany")
	}

	if _, err := strconv.Atoi(vatNumber); err != nil {
		return errors.New("VAT number can only contain `DE` followed by 9 digits")
	}

	return nil
}

//FIXME: returns a 500 (soap:Server:SERVICE_UNAVAILABLE) when in `test` mode and passing a VAT number that's not supported by VIES
func (controller *VATController) remoteCheck(req *CheckVatRequest) (bool, error) {
	client := new(soap.Client)
	url := config.GetString("services.vat.url")
	fmt.Printf("url: %v\n", url)
	response, err := client.Request(url, req)
	if err != nil {
		fmt.Printf("err: %v\n", err)
		return false, err
	}

	return response.SoapBody.Body.Valid == "true", nil
}
