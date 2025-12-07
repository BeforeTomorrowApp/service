package services

import (
	"context"
	"fmt"

	"north-post/service/internal/domain/v1/models"
	"north-post/service/internal/repository"
)

const defaultPageSize = 100

type AddressService struct {
	repo *repository.AddressRepository
}

func NewAddressService(repo *repository.AddressRepository) *AddressService {
	return &AddressService{
		repo: repo,
	}
}

type GetAddressesInput struct {
	Language models.Language
	Tags     []string
	Limit    int
}

type GetAddressesOutput struct {
	Addresses []models.AddressItem
	Count     int
}

type GetAddressByIdInput struct {
	Language models.Language
	ID       string
}

type GetAddressByIdOutput struct {
	Address models.AddressItem
}

// GetAddresses godoc
// @Summary Get all addresses
// @Description Get all addresses by language and filtered optional tags
// @Param request body dto.GetAddressesRequest true "Request body"
// @Success 200 {object} dto.GetAddressesResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /api/v1/admin/addresses [post]
func (s *AddressService) GetAddresses(ctx context.Context, input GetAddressesInput) (*GetAddressesOutput, error) {
	limit := input.Limit
	if limit <= 0 || limit > defaultPageSize {
		limit = defaultPageSize
	}
	opts := repository.GetAllAddressesOptions{
		Language: input.Language,
		Tags:     input.Tags,
		Limit:    limit,
	}
	addresses, err := s.repo.GetAllAddresses(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to get addresses: %w", err)
	}
	return &GetAddressesOutput{Addresses: addresses, Count: len(addresses)}, nil
}

func (s *AddressService) GetAddressById(ctx context.Context, input GetAddressByIdInput) (*GetAddressByIdOutput, error) {
	opts := repository.GetAddressByIdOption{
		Language:  input.Language,
		AddressID: input.ID,
	}
	address, err := s.repo.GetAddressById(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to get address: %w", err)
	}
	return &GetAddressByIdOutput{Address: *address}, nil
}
