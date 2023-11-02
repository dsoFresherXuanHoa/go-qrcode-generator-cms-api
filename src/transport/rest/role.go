package rest

import (
	"fmt"
	"go-qrcode-generator-cms-api/src/business"
	"go-qrcode-generator-cms-api/src/constants"
	"go-qrcode-generator-cms-api/src/entity"
	"go-qrcode-generator-cms-api/src/storage"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var (
	InvalidCreateRoleRequest         = "Invalid Role Incoming Request: Check Swagger For More Information."
	ValidateCreateRoleRequestFailure = "Invalid Role Incoming Request: Check Swagger For More Information."
	CreateRoleFailure                = "Error Create Role: Make Sure You Had Permission To Do This."
	CreateRoleSuccess                = "Create Role Success: Congrats."
)

func CreateRole(db *gorm.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		sqlStorage := storage.NewSQLStore(db)
		roleStorage := storage.NewRoleStore(sqlStorage)
		roleBusiness := business.NewRoleBusiness(roleStorage)

		var reqRole entity.RoleCreatable
		if err := ctx.ShouldBind(&reqRole); err != nil {
			fmt.Println("Error while parse role request: " + err.Error())
			ctx.JSON(http.StatusBadRequest, entity.NewStandardResponse(nil, http.StatusBadRequest, constants.StatusBadRequest, err.Error(), InvalidCreateRoleRequest))
		} else if err := reqRole.Validate(); err != nil {
			fmt.Println("Error while validate role request: " + err.Error())
			ctx.JSON(http.StatusBadRequest, entity.NewStandardResponse(nil, http.StatusBadRequest, constants.StatusBadRequest, err.Error(), ValidateCreateRoleRequestFailure))
		} else if roleUUID, err := roleBusiness.CreateRole(ctx, &reqRole); err != nil {
			fmt.Println("Error while create role: " + err.Error())
			ctx.JSON(http.StatusInternalServerError, entity.NewStandardResponse(nil, http.StatusInternalServerError, constants.StatusInternalServerError, err.Error(), CreateRoleFailure))
		} else {
			ctx.JSON(http.StatusOK, entity.NewStandardResponse(roleUUID, http.StatusOK, constants.StatusOK, "", CreateRoleSuccess))
		}
	}
}
