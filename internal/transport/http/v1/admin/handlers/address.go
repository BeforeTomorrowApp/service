package handlers

import (
	"log/slog"
	"net/http"
	"north-post/service/internal/domain/v1/models"
	"north-post/service/internal/services"
	"north-post/service/internal/transport/http/v1/dto"
	"north-post/service/internal/transport/http/v1/utils"

	"github.com/gin-gonic/gin"
)

type AddressHandler struct {
	service *services.AddressService
	logger  *slog.Logger
}

func NewAddressHandler(service *services.AddressService, logger *slog.Logger) *AddressHandler {
	return &AddressHandler{
		service: service,
		logger:  logger,
	}
}

func (h *AddressHandler) GetAddresses(c *gin.Context) {
	var req dto.GetAllAddressesRequest
	if !utils.BindJSON(c, &req, h.logger) {
		return
	}
	if !utils.ValidateLanguage(c, req.Language, h.logger) {
		return
	}
	input := services.GetAddressesInput{
		Language: req.Language,
		Tags:     req.Tags,
		Limit:    req.Limit,
	}
	output, err := h.service.GetAddresses(c.Request.Context(), input)
	if err != nil {
		h.logger.Error("failed to get addresses", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch address"})
		return
	}
	response := dto.GetAllAddressResponse{
		Data:  dto.ToAddressDTOs(output.Addresses),
		Count: output.Count,
	}
	c.JSON(http.StatusOK, response)
}

func (h *AddressHandler) GetAddressById(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		h.logger.Warn("missing address id parameter")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Address ID is required"})
		return
	}
	languageStr := c.Query("language")
	if languageStr == "" {
		h.logger.Warn("missing language parameter")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Language is required"})
		return
	}
	language := models.Language(languageStr)
	if !utils.ValidateLanguage(c, language, h.logger) {
		return
	}
	input := services.GetAddressByIdInput{
		Language: language,
		ID:       id,
	}
	output, err := h.service.GetAddressById(c.Request.Context(), input)
	if err != nil {
		h.logger.Error("failed to get address", "addressId", input.ID, "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch address"})
		return
	}
	response := dto.GetAddressByIdResponse{
		Data: dto.ToAddressDTO(output.Address),
	}
	c.JSON(http.StatusOK, response)
}
