package services

import (
	"context"
	"errors"
	"north-post/service/internal/repository"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockPromptRepository struct {
	mock.Mock
}

func (m *mockPromptRepository) GetSystemPrompt(
	ctx context.Context,
	opts repository.GetSystemPromptOptions) (string, error) {
	args := m.Called(ctx, opts)
	if args.Get(0) == nil {
		return "", args.Error(1)
	}
	return args.String(0), args.Error(1)
}

func (m *mockPromptRepository) GetSystemAddressGenerationPrompt(
	ctx context.Context,
	opts repository.GetSystemAddressGenerationPromptOptions,
) (string, error) {
	args := m.Called(ctx, opts)
	if args.Get(0) == nil {
		return "", args.Error(1)
	}
	return args.String(0), args.Error(1)
}

func setupPromptService() (*PromptService, *mockPromptRepository) {
	repo := new(mockPromptRepository)
	service := NewPromptService(repo)
	return service, repo
}

// Tests
func TestPromptService_GetSystemAddressGenerationPrompt(t *testing.T) {
	t.Parallel()
	service, repo := setupPromptService()
	ctx := context.Background()
	input := GetSystemAddressGenerationPromptInput{
		Language: "en",
	}
	expectedPrompt := "system prompt"
	repo.On(
		"GetSystemAddressGenerationPrompt",
		mock.Anything,
		mock.Anything,
	).Return(expectedPrompt, nil).Once()
	output, err := service.GetSystemAddressGenerationPrompt(ctx, input)
	repo.AssertExpectations(t)
	assert.Nil(t, err)
	assert.Equal(t, output.Prompt, expectedPrompt)
}

func TestPromptService_GetSystemAddressGenerationPrompt_Error(t *testing.T) {
	t.Parallel()
	service, repo := setupPromptService()
	ctx := context.Background()
	input := GetSystemAddressGenerationPromptInput{
		Language: "en",
	}
	repo.On(
		"GetSystemAddressGenerationPrompt",
		mock.Anything,
		mock.Anything,
	).Return("", errors.New("error")).Once()
	output, err := service.GetSystemAddressGenerationPrompt(ctx, input)
	repo.AssertExpectations(t)
	assert.Error(t, err)
	assert.Nil(t, output)
}
