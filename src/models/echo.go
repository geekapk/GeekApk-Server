package models

import (
	"modelmap"
)

type EchoModel struct {

}

type EchoModelCreateOrUpdateInfo struct {
	Name string `json:"name"`
	Description string `json:"description"`
}

func (m *EchoModel) GetName() string {
	return "Echo"
}

func (m *EchoModel) Create(createInfo modelmap.Deserializer) interface{} {
	var info EchoModelCreateOrUpdateInfo
	createInfo(&info)
	return &info
}

func (m *EchoModel) Read(filter map[string]modelmap.FilterRule) interface{} {
	return filter
}

func (m *EchoModel) Update(
	filter map[string]modelmap.FilterRule,
	updateInfo modelmap.Deserializer,
) interface{} {
	var info EchoModelCreateOrUpdateInfo
	updateInfo(&info)

	return map[string]interface{} {
		"filter": filter,
		"info": &info,
	}
}

func (m *EchoModel) Delete(filter map[string]modelmap.FilterRule) interface{} {
	return filter
}
