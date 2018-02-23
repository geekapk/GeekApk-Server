package models

import (
	"log"
	"time"
	"modelmap"
	"encoding/base64"
	"github.com/jinzhu/gorm"
	"github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserModel struct {
	db *gorm.DB
}

func NewUserModel(db *gorm.DB) *UserModel {
	db.AutoMigrate(&User{})

	return &UserModel {
		db: db,
	}
}

type User struct {
	Id string `gorm:"type:varchar(36);primary_key"`
	Email string `gorm:"unique_index"`
	Password string `gorm:"type:varchar(128)"`
	PendingAuth bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

type UserCreateInfo struct {
	Email string `json:"email"`
	Name string `json:"name"`
	DisplayName string `json:"display_name"`
	RawPassword string `json:"password"`
}

type UserUpdateInfo struct {
	Id string `json:"id"`
	Email string `json:"email"`
	DisplayName string `json:"display_name"`
	IconId string `json:"icon_id"`
}

func (m *UserModel) GetName() string {
	return "User"
}

func (m *UserModel) Create(rc *modelmap.RequestContext, loadCreateInfo modelmap.Deserializer) interface{} {
	var cinfo UserCreateInfo
	loadCreateInfo(&cinfo)

	log.Println(&cinfo)

	if !checkBasicString(cinfo.Email) ||
		!checkBasicString(cinfo.Name) ||
		!checkBasicString(cinfo.DisplayName) ||
		!checkBasicString(cinfo.RawPassword) {
			return "Illegal input"
		}
	if !checkPassword(cinfo.RawPassword) {
		return "Invalid password"
	}

	encPw, err := bcrypt.GenerateFromPassword(
		[]byte(cinfo.RawPassword),
		bcrypt.DefaultCost,
	)
	if err != nil {
		panic(err)
	}

	u := User {
		Id: uuid.Must(uuid.NewV4()).String(),
		Email: cinfo.Email,
		Password: base64.StdEncoding.EncodeToString(encPw),
		PendingAuth: false, // Disable for now
	}
	acc := Account {
		Id: u.Id,
		Name: cinfo.Name,
		DisplayName: cinfo.DisplayName,
		Disabled: false,
		IconId: "",
		Subtype: "User",
	}
	err = m.db.Create(&acc).Error
	if err != nil {
		return "Account already exists"
	}
	err = m.db.Create(&u).Error
	if err != nil {
		// Rollback
		m.db.Delete(&acc)
		return "Duplicated email address"
	}
	return "User created"
}

func (m *UserModel) Read(rc *modelmap.RequestContext, filter map[string]modelmap.FilterRule) interface{} {
	if f, ok := filter["id"]; ok {
		var u User
		err := m.db.Where("id = ?", f.Value).First(&u).Error
		if err != nil {
			if gorm.IsRecordNotFoundError(err) {
				return "User not found"
			} else {
				panic(err)
			}
		}

		var acc Account
		err = m.db.Where("id = ?", f.Value).First(&acc).Error

		// Indicating that the internal data structure has corrupted.
		if err != nil {
			panic(err)
		}

		return map[string]interface{} {
			"id": acc.Id,
			"name": acc.Name,
			"display_name": acc.DisplayName,
			"disabled": acc.Disabled,
			"icon_id": acc.IconId,
			"email": u.Email,
			"pending_auth": u.PendingAuth,
		};
	} else {
		return "Invalid filter"
	}
}

func (m *UserModel) Update(
	rc *modelmap.RequestContext,
	filter map[string]modelmap.FilterRule,
	loadUpdateInfo modelmap.Deserializer,
) interface{} {
	return "Not implemented"
}

func (m *UserModel) Delete(rc *modelmap.RequestContext, filter map[string]modelmap.FilterRule) interface{} {
	return "Not implemented"
	//return "Account is not allowed to be deleted directly. Delete a User instead."
}

