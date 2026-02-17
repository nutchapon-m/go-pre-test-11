package http

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/nutchapon-m/go-pre-test-11/internal/core/domain"
	"github.com/nutchapon-m/go-pre-test-11/internal/core/port"
)

type ProductHandler struct {
	service port.ProductService
}

func NewProductHandler(service port.ProductService) *ProductHandler {
	return &ProductHandler{service: service}
}

type Response struct {
	Successful bool        `json:"successful"`
	ErrorCode  string      `json:"error_code,omitempty"`
	Data       interface{} `json:"data,omitempty"`
}

// CreateProduct godoc
// @Summary Create a new product
// @Description Create a new product with the input payload
// @Tags product
// @Accept json
// @Produce json
// @Param product body port.CreateProductRequest true "Product to create"
// @Success 200 {object} Response
// @Failure 400 {object} Response
// @Failure 500 {object} Response
// @Router /product [post]
func (h *ProductHandler) CreateProduct(c echo.Context) error {
	var req port.CreateProductRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, Response{
			Successful: false,
			ErrorCode:  "INVALID_REQUEST",
		})
	}
	// Simple validation if not using a validator library yet (though `validate:"required"` tag exists in struct)
	// For now rely on service or add manual check if needed, but struct tags suggest validation library usage.
	if req.Name == "" || req.Price <= 0 {
		return c.JSON(http.StatusBadRequest, Response{
			Successful: false,
			ErrorCode:  "INVALID_PAYLOAD",
		})
	}

	product, err := h.service.CreateProduct(c.Request().Context(), &req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, Response{
			Successful: false,
			ErrorCode:  "INTERNAL_SERVER_ERROR",
		})
	}

	return c.JSON(http.StatusOK, Response{
		Successful: true,
		Data:       product,
	})
}

// PatchProduct godoc
// @Summary Patch a product
// @Description Update specific fields of a product
// @Tags product
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Param product body port.PatchProductRequest true "Fields to update"
// @Success 200 {object} Response
// @Failure 400 {object} Response
// @Failure 404 {object} Response
// @Failure 500 {object} Response
// @Router /product/{id} [patch]
func (h *ProductHandler) PatchProduct(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, Response{
			Successful: false,
			ErrorCode:  "INVALID_ID",
		})
	}

	var req port.PatchProductRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, Response{
			Successful: false,
			ErrorCode:  "INVALID_REQUEST",
		})
	}

	err = h.service.PatchProduct(c.Request().Context(), id, &req)
	if err != nil {
		if err == domain.ErrProductNotFound {
			return c.JSON(http.StatusNotFound, Response{
				Successful: false,
				ErrorCode:  "PRODUCT_NOT_FOUND",
			})
		}
		return c.JSON(http.StatusInternalServerError, Response{
			Successful: false,
			ErrorCode:  "INTERNAL_SERVER_ERROR",
		})
	}

	return c.JSON(http.StatusOK, Response{
		Successful: true,
	})
}
