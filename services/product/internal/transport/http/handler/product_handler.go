package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/nashinnov8/go-ecommerce/pkg/apperror"
	"github.com/nashinnov8/go-ecommerce/pkg/response"
	"github.com/nashinnov8/go-ecommerce/services/product/internal/dto"
	"github.com/nashinnov8/go-ecommerce/services/product/internal/service"
)

type ProductHandler struct {
	svc service.ProductService
}

func NewProductHandler(svc service.ProductService) *ProductHandler {
	return &ProductHandler{svc: svc}
}

func (h *ProductHandler) CreateProduct(c *gin.Context) {
	var req dto.CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	res, err := h.svc.CreateProduct(c.Request.Context(), req)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Created(c, res)
}

func (h *ProductHandler) GetProductByID(c *gin.Context) {
	id, err := parseUUIDParam(c, "id")
	if err != nil {
		response.Error(c, err)
		return
	}

	res, err := h.svc.GetProductByID(c.Request.Context(), id)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.OK(c, res)
}

func (h *ProductHandler) ListProducts(c *gin.Context) {
	var query dto.ListProductsQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	items, total, err := h.svc.ListProducts(c.Request.Context(), query)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Paginated(c, items, query.Page, query.Limit, total)
}

func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	id, err := parseUUIDParam(c, "id")
	if err != nil {
		response.Error(c, err)
		return
	}

	var req dto.UpdateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	res, err := h.svc.UpdateProduct(c.Request.Context(), id, req)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.OK(c, res)
}

func (h *ProductHandler) UpdateStock(c *gin.Context) {
	id, err := parseUUIDParam(c, "id")
	if err != nil {
		response.Error(c, err)
		return
	}

	var req dto.UpdateStockRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	res, err := h.svc.UpdateStock(c.Request.Context(), id, req)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.OK(c, res)
}

func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	id, err := parseUUIDParam(c, "id")
	if err != nil {
		response.Error(c, err)
		return
	}

	if err := h.svc.DeleteProduct(c.Request.Context(), id); err != nil {
		response.Error(c, err)
		return
	}

	response.OK(c, gin.H{"deleted": true})
}

func parseUUIDParam(c *gin.Context, param string) (uuid.UUID, error) {
	raw := c.Param(param)
	id, err := uuid.Parse(raw)
	if err != nil {
		return uuid.Nil, apperror.BadRequest("INVALID_UUID", "Invalid product ID format", err)
	}
	return id, nil
}
