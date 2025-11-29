package payment

import (
	"net/http"
	"strconv"

	"github.com/example/ai-avatar-studio/internal/pkg/middleware"
	"github.com/example/ai-avatar-studio/internal/pkg/response"
	paymentsvc "github.com/example/ai-avatar-studio/internal/service/payment"
	"github.com/gin-gonic/gin"
)

// Handler wires recharge/payment endpoints.
type Handler struct {
	service *paymentsvc.Service
	secret  string
}

func NewHandler(service *paymentsvc.Service, secret string) *Handler {
	return &Handler{service: service, secret: secret}
}

func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	auth := middleware.Authenticator(h.secret)
	rg.POST("/store/payments", auth, h.create)
	rg.GET("/store/payments", auth, h.listUser)
	rg.GET("/store/payments/:out_trade_no", auth, h.query)
	// Notify/return are public callbacks from gateway.
	rg.GET("/store/payments/notify", h.notify)
	rg.POST("/store/payments/notify", h.notify)
	rg.GET("/store/payments/return", h.returnPage)

	// Admin audit
	admin := middleware.AdminOnly(h.secret)
	rg.GET("/admin/payments", admin, h.listAdmin)
}

func (h *Handler) create(c *gin.Context) {
	var req struct {
		Amount  float64 `json:"amount"`
		PayType string  `json:"pay_type"`
	}
	if err := c.ShouldBindJSON(&req); err != nil || req.Amount <= 0 {
		response.Error(c, http.StatusBadRequest, "invalid amount")
		return
	}
	userID := middleware.CurrentUserID(c)
	res, err := h.service.CreateOrder(c.Request.Context(), userID, req.Amount, req.PayType)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	response.Success(c, gin.H{
		"out_trade_no": res.OutTradeNo,
		"pay_url":      res.PayURL,
		"coins":        res.Order.Coins,
		"amount":       req.Amount,
	})
}

func (h *Handler) query(c *gin.Context) {
	outTradeNo := c.Param("out_trade_no")
	order, err := h.service.Query(c.Request.Context(), outTradeNo)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	if order == nil {
		response.Error(c, http.StatusNotFound, "order not found")
		return
	}
	// Only owner can query
	if middleware.CurrentUserID(c) != order.UserID {
		response.Error(c, http.StatusForbidden, "forbidden")
		return
	}
	response.Success(c, order)
}

func (h *Handler) listUser(c *gin.Context) {
	userID := middleware.CurrentUserID(c)
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	status := c.DefaultQuery("status", "paid")
	orders, err := h.service.ListUserOrders(c.Request.Context(), userID, limit, offset, status)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, orders)
}

func (h *Handler) listAdmin(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "100"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	orders, err := h.service.ListAll(c.Request.Context(), limit, offset)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.Success(c, orders)
}

func (h *Handler) notify(c *gin.Context) {
	payload := readPayload(c)
	order, err := h.service.HandleNotify(c.Request.Context(), payload)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	// gateway expects plain success
	c.String(http.StatusOK, "success")
	_ = order
}

func (h *Handler) returnPage(c *gin.Context) {
	payload := readPayload(c)
	order, err := h.service.HandleNotify(c.Request.Context(), payload)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":       "ok",
		"out_trade_no": order.OutTradeNo,
		"amount":       payload["money"],
		"coins":        order.Coins,
	})
}

// readPayload collects query/form params into map.
func readPayload(c *gin.Context) map[string]string {
	payload := map[string]string{}
	for k, v := range c.Request.URL.Query() {
		if len(v) > 0 {
			payload[k] = v[0]
		}
	}
	_ = c.Request.ParseForm()
	for k, v := range c.Request.PostForm {
		if len(v) > 0 {
			payload[k] = v[0]
		}
	}
	return payload
}
