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

func (m *EchoModel) Create(rc *modelmap.RequestContext, createInfo modelmap.Deserializer) interface{} {
	var info EchoModelCreateOrUpdateInfo
	createInfo(&info)
	return &info
}

func (m *EchoModel) Read(rc *modelmap.RequestContext, filter map[string]modelmap.FilterRule) interface{} {
	return filter
}

func (m *EchoModel) Update(
	rc *modelmap.RequestContext,
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

func (m *EchoModel) Delete(rc *modelmap.RequestContext, filter map[string]modelmap.FilterRule) interface{} {
	return filter
}
