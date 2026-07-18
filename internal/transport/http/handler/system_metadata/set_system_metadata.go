package system_metadata

import (
	"context"

	"github.com/danielgtaylor/huma/v2"
	"github.com/vandordev/vkit-tango/internal/contract"
	"github.com/vandordev/vkit-tango/internal/transport/http/method"
	"github.com/vandordev/vkit-tango/internal/usecase"
)

type SetSystemMetadataHandler struct {
	api     huma.API
	command contract.Command[usecase.SetSystemMetadataInput, usecase.SetSystemMetadataResult]
}

var _ contract.HTTPHandler = (*SetSystemMetadataHandler)(nil)

func NewSetSystemMetadataHandler(api huma.API, command contract.Command[usecase.SetSystemMetadataInput, usecase.SetSystemMetadataResult]) *SetSystemMetadataHandler {
	return &SetSystemMetadataHandler{api: api, command: command}
}

type setSystemMetadataInput struct {
	Key  string `path:"key" minLength:"1"`
	Body struct {
		Value map[string]any `json:"value"`
	}
}

type setSystemMetadataOutput struct {
	Body struct {
		Success bool `json:"success"`
		Data    struct {
			Key string `json:"key"`
		} `json:"data"`
	}
}

func (handler *SetSystemMetadataHandler) RegisterRoutes() {
	method.PUT(handler.api, "/api/v1/system-metadata/{key}", "Set system metadata", []string{"system-metadata"}, handler.Handle)
}

func (handler *SetSystemMetadataHandler) Handle(ctx context.Context, input *setSystemMetadataInput) (*setSystemMetadataOutput, error) {
	result, err := handler.command.Execute(ctx, usecase.SetSystemMetadataInput{Key: input.Key, Value: input.Body.Value})
	if err != nil {
		return nil, err
	}
	output := &setSystemMetadataOutput{}
	output.Body.Success = true
	output.Body.Data.Key = result.Key
	return output, nil
}
