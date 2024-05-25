package handler

import (
	"errors"
	"net/http"
	"simple-order-go/common"
	"simple-order-go/internal/entity"
	"simple-order-go/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx"
)

type OrderHandler struct {
	orderService service.IOrderService
}

func NewOrderHandler(orderService service.IOrderService) *OrderHandler {
	return &OrderHandler{orderService: orderService}
}

type orderRequest struct {
	CustomerName string `json:"customerName" binding:"required"`
	OrderedAt    string `json:"orderedAt" binding:"required"`
	Items        []struct {
		Name     string `json:"name" binding:"required"`
		Desc     string `json:"description" binding:"required"`
		Quantity int32  `json:"quantity" binding:"required,gt=0"`
	} `json:"items" binding:"required,gt=0,dive"`
}

func (h *OrderHandler) CreateOrder(ctx *gin.Context) {
	var req orderRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	t, err := common.ParseStringToTime(req.OrderedAt)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	items := make(entity.ItemViewModels, len(req.Items))
	for i, item := range req.Items {
		items[i] = entity.ItemViewModel{
			Name:        item.Name,
			Description: item.Desc,
			Quantity:    item.Quantity,
		}
	}

	arg := entity.OrderViewModel{
		CustomerName: req.CustomerName,
		OrderedAt:    t,
		Items:        items,
	}

	order, err := h.orderService.CreateOrder(arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, order)
}

type orderByIDRequest struct {
	ID int64 `uri:"id" binding:"required,gt=0"`
}

func (h *OrderHandler) GetOrderByID(ctx *gin.Context) {
	var req orderByIDRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	order, err := h.orderService.GetOrder(req.ID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, order)
}

func (h *OrderHandler) GetAllOrders(ctx *gin.Context) {
	orders, err := h.orderService.GetAllOrders()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, orders)
}

func (h *OrderHandler) UpdateOrder(ctx *gin.Context) {
	var idReq orderByIDRequest
	if err := ctx.ShouldBindUri(&idReq); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	var req orderRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	t, err := common.ParseStringToTime(req.OrderedAt)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	items := make(entity.ItemViewModels, len(req.Items))
	for i, item := range req.Items {
		items[i] = entity.ItemViewModel{
			Name:        item.Name,
			Description: item.Desc,
			Quantity:    item.Quantity,
		}
	}

	arg := entity.OrderViewModel{
		ID:           idReq.ID,
		CustomerName: req.CustomerName,
		OrderedAt:    t,
		Items:        items,
	}

	err = h.orderService.UpdateOrder(arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, successResponse())
}

func (h *OrderHandler) DeleteOrder(ctx *gin.Context) {
	var req orderByIDRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	err := h.orderService.DeleteOrder(req.ID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, successResponse())
}

func successResponse() gin.H {
	return gin.H{"result": "Success"}
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
