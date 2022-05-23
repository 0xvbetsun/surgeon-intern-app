package adminService

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	qlmodel "github.com/vbetsun/surgeon-intern-app/internal/adminService/graph/model"
	"go.uber.org/zap"
)

type (
	RestApi struct {
		service *Service
	}
)

func NewRestApi(service *Service) *RestApi {
	return &RestApi{service: service}
}

func (r *RestApi) InitializeRouter(rg *gin.RouterGroup) {
	rg.POST("/activateUser", r.ActivateUser)
	rg.GET("/verifyActivationParameters", r.VerifyActivationParameters)
}

func (r *RestApi) ActivateUser(c *gin.Context) {
	var input qlmodel.ActivateUserInput
	zap.S().Error("Activate user")
	err := c.ShouldBindJSON(&input)
	if err != nil {
		zap.S().Error(err.Error())
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	retUser, err := r.service.adminUserService.ActivateUser(context.TODO(), input)
	if err != nil {
		zap.S().Error(err.Error())
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	c.IndentedJSON(http.StatusCreated, retUser)
}

func (r *RestApi) VerifyActivationParameters(c *gin.Context) {
	var input qlmodel.ActivationVerificationInput
	err := c.BindQuery(&input)
	if err != nil {
		zap.S().Error(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, qlmodel.ActivationVerification{Status: qlmodel.ActivationVerificationStatusUnknownError})
		return
	}
	validationStatus, err := r.service.adminUserService.ValidateActivation(context.TODO(), &input)
	if err != nil {
		zap.S().Error(err.Error())
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	if validationStatus != qlmodel.ActivationVerificationStatusActive {
		c.AbortWithStatusJSON(http.StatusBadRequest, qlmodel.ActivationVerification{Status: validationStatus})
		return
	}
	c.IndentedJSON(http.StatusOK, qlmodel.ActivationVerification{Status: validationStatus})
}
