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

func CreateRole(db *gorm.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		sqlStorage := storage.NewSQLStore(db)
		roleStorage := storage.NewRoleStore(sqlStorage)
		roleBusiness := business.NewRoleBusiness(roleStorage)

		var reqRole entity.RoleCreatable
		if err := ctx.ShouldBind(&reqRole); err != nil {
			fmt.Println("Error while parse role request to struct: " + err.Error())
			ctx.JSON(http.StatusBadRequest, entity.NewStandardResponse(nil, http.StatusBadRequest, constants.StatusBadRequest, err.Error(), constants.InvalidRoleRequestFormat))
		} else if roleUUID, err := roleBusiness.CreateRole(ctx, &reqRole); err != nil {
			fmt.Println("Error while create new role in role transport: " + err.Error())
			ctx.JSON(http.StatusInternalServerError, entity.NewStandardResponse(nil, http.StatusInternalServerError, constants.StatusInternalServerError, err.Error(), constants.ErrCreateNewRole))
		} else {
			ctx.JSON(http.StatusOK, entity.NewStandardResponse(roleUUID, http.StatusOK, constants.StatusOK, "", constants.CreateNewRoleSuccess))
		}
	}
}
