package seeds

import (
	"github.com/google/uuid"
	"github.com/junaid9001/lattrix-backend/internal/domain/models"
	"gorm.io/gorm"
)

var permissions = []models.Permission{
	{
		ID:          uuid.New(),
		Code:        "api:create",
		Description: "Create API",
	},
	{
		ID:          uuid.New(),
		Code:        "api:update",
		Description: "Update API",
	},
	{
		ID:          uuid.New(),
		Code:        "api:delete",
		Description: "Delete API",
	},
	{
		ID:          uuid.New(),
		Code:        "dashboard:view",
		Description: "View dashboard",
	}, {
		ID:          uuid.New(),
		Code:        "api-group:create",
		Description: "Create api group",
	},
	{
		ID:          uuid.New(),
		Code:        "api-group:delete",
		Description: "delete api group",
	}, {
		ID:          uuid.New(),
		Code:        "api-group:update",
		Description: "Update api group",
	}, {
		ID:          uuid.New(),
		Code:        "role:superadmin",
		Description: "Full permission",
	},
}

func SeedPermissions(db *gorm.DB) error {

	for _, val := range permissions {
		var existing models.Permission
		err := db.Where("code = ?", val.Code).First(&existing).Error

		if err == gorm.ErrRecordNotFound {
			if err := db.Create(&val).Error; err != nil {
				return err
			}
		}
	}
	return nil
}
