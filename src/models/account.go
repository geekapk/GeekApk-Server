package models

import (
	"time"
	"modelmap"
	"github.com/jinzhu/gorm"
)

type AccountModel struct {
	db *gorm.DB
}

func NewAccountModel(db *gorm.DB) *AccountModel {
	db.AutoMigrate(&Account{})

	return &AccountModel {
		db: db,
	}
}

type Account struct {
	Id string `gorm:"type:varchar(36);primary_key"`
	Name string `gorm:"type:varchar(32);unique_index"`
	DisplayName string `gorm:"type:varchar(64)"`
	Disabled bool
	IconId string `gorm:"type:varchar(36)"`
	Subtype string `gorm:"type:varchar(16)"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type AccountUpdateInfo struct {
	Id string `json:"id"`
	DisplayName string `json:"display_name"`
	IconId string `json:"icon_id"`
}

func (m *AccountModel) GetName() string {
	return "Account"
}

func (m *AccountModel) Create(rc *modelmap.RequestContext, createInfo modelmap.Deserializer) interface{} {
	return "Account is not allowed to be created directly. Create a User instead."
}

func (m *AccountModel) Read(rc *modelmap.RequestContext, filter map[string]modelmap.FilterRule) interface{} {
	if f, ok := filter["id"]; ok {
		var acc Account
		err := m.db.Where("id = ?", f.Value).First(&acc).Error
		if err != nil {
			if gorm.IsRecordNotFoundError(err) {
				return "Account not found"
			} else {
				panic(err)
			}
		}
		return map[string]interface{} {
			"id": acc.Id,
			"name": acc.Name,
			"display_name": acc.DisplayName,
			"disabled": acc.Disabled,
			"icon_id": acc.IconId,
			"subtype": acc.Subtype,
		};
	} else {
		return "Invalid filter"
	}
}

func (m *AccountModel) Update(
	rc *modelmap.RequestContext,
	filter map[string]modelmap.FilterRule,
	loadUpdateInfo modelmap.Deserializer,
) interface{} {
	return "Account is not allowed to be updated directly. Update a User instead."
}

func (m *AccountModel) Delete(rc *modelmap.RequestContext, filter map[string]modelmap.FilterRule) interface{} {
	return "Account is not allowed to be deleted directly. Delete a User instead."
}
