package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/smartwalle/alipay/v3"

	"github.com/sheet-platform/backend/internal/config"
	"github.com/sheet-platform/backend/internal/dto/response"
	"github.com/sheet-platform/backend/internal/pkg"
	"github.com/sheet-platform/backend/internal/service"
)

type OrderHandler struct {
	orderService *service.OrderService
	cfg          *config.Config
}

func NewOrderHandler(orderService *service.OrderService, cfg *config.Config) *OrderHandler {
	return &OrderHandler{orderService: orderService, cfg: cfg}
}

func (h *OrderHandler) RegisterRoutes(public, auth, admin *gin.RouterGroup) {
	public.POST("/alipay/notify", h.Notify)
	public.GET("/orders/:order_no/pay", h.Pay)
	auth.GET("/alipay/order/status", h.OrderStatus)
	auth.POST("/orders", h.Create)
	auth.GET("/orders", h.List)
}

func (h *OrderHandler) Create(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var req struct {
		SheetMusicID uint `json:"sheet_music_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(400, "参数错误"))
		return
	}

	createReq := &service.CreateOrderReq{
		UserID:       userID.(uint),
		SheetMusicID: req.SheetMusicID,
	}

	order, payURL, err := h.orderService.CreateOrder(createReq)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error(400, err.Error()))
		return
	}

	c.JSON(http.StatusOK, response.Success(gin.H{
		"id":       order.ID,
		"order_no": order.OrderNo,
		"amount":   order.Amount,
		"status":   order.Status,
		"pay_url":  payURL,
	}))
}

func (h *OrderHandler) List(c *gin.Context) {
	userID, _ := c.Get("user_id")

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	orders, total, err := h.orderService.ListOrders(userID.(uint), page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(500, "查询订单列表失败"))
		return
	}

	c.JSON(http.StatusOK, response.Success(response.PageResponse{
		List:     orders,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}))
}

func (h *OrderHandler) Pay(c *gin.Context) {
	orderNo := c.Param("order_no")
	if orderNo == "" {
		c.JSON(http.StatusBadRequest, response.Error(400, "缺少订单号"))
		return
	}

	token := c.Query("token")
	if token == "" {
		c.JSON(http.StatusUnauthorized, response.Error(401, "未登录"))
		return
	}

	claims, err := pkg.ParseJWT(h.cfg.JWT.Secret, token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, response.Error(401, "登录已过期，请重新登录"))
		return
	}

	payURL, err := h.orderService.GetPayURL(orderNo, claims.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error(400, err.Error()))
		return
	}

	c.Redirect(http.StatusFound, payURL)
}

func (h *OrderHandler) OrderStatus(c *gin.Context) {
	userID, _ := c.Get("user_id")
	orderNo := c.Query("order_no")
	if orderNo == "" {
		c.JSON(http.StatusBadRequest, response.Error(400, "缺少订单号"))
		return
	}
	status, err := h.orderService.GetOrderStatus(orderNo, userID.(uint))
	if err != nil {
		c.JSON(http.StatusBadRequest, response.Error(400, err.Error()))
		return
	}
	c.JSON(http.StatusOK, response.Success(status))
}

func (h *OrderHandler) Notify(c *gin.Context) {
	if err := c.Request.ParseForm(); err != nil {
		c.String(http.StatusBadRequest, "fail")
		return
	}

	if c.Request.PostForm.Get("out_trade_no") == "" || c.Request.PostForm.Get("trade_status") == "" {
		c.String(http.StatusBadRequest, "fail")
		return
	}

	if err := h.orderService.HandleAlipayNotify(c.Request.PostForm); err != nil {
		c.String(http.StatusOK, "fail")
		return
	}

	alipay.ACKNotification(c.Writer)
}
