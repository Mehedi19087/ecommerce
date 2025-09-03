package order

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type OrderController struct {
	orderService OrderService
}

func NewOrderController(orderService OrderService) *OrderController {
	return &OrderController{orderService: orderService}
}

func (c *OrderController) CreateOrder(ctx *gin.Context) {

	userID, exists := ctx.Get("userID")
    if !exists {
    ctx.JSON(http.StatusUnauthorized, gin.H{
        "error": "User not authenticated",
    })
    return
}
      userIDUint := userID.(uint)
	var req CreateOrderRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	order, err := c.orderService.CreateOrderFromCart(userIDUint, req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "Order created successfully",
		"order":   order,
	})
}

func (c *OrderController) GetUserOrders(ctx *gin.Context) {
	userID, exists := ctx.Get("userID")
if !exists {
    ctx.JSON(http.StatusUnauthorized, gin.H{
        "error": "User not authenticated",
    })
    return
}
userIDUint := userID.(uint)

	orders, err := c.orderService.GetUserOrders(userIDUint)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get orders",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"orders": orders,
		"count":  len(orders),
	})
}

func (c *OrderController) GetOrderByID(ctx *gin.Context) {
	userID, exists := ctx.Get("userID")
if !exists {
    ctx.JSON(http.StatusUnauthorized, gin.H{
        "error": "User not authenticated",
    })
    return
}
userIDUint := userID.(uint)

	orderID, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid order ID",
		})
		return
	}

	order, err := c.orderService.GetOrderByID(uint(orderID), userIDUint)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"order": order,
	})
}

func (c *OrderController) CancelOrder(ctx *gin.Context) {
	userID, exists := ctx.Get("userID")
if !exists {
    ctx.JSON(http.StatusUnauthorized, gin.H{
        "error": "User not authenticated",
    })
    return
}
userIDUint := userID.(uint)

	orderID, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid order ID",
		})
		return
	}

	err = c.orderService.CancelOrder(uint(orderID), userIDUint)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Order cancelled successfully",
	})
}

func (c *OrderController) SubmitPaymentProof(ctx *gin.Context) {
	userID, exists := ctx.Get("userID")
if !exists {
    ctx.JSON(http.StatusUnauthorized, gin.H{
        "error": "User not authenticated",
    })
    return
}
userIDUint := userID.(uint)

	orderID, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid order ID",
		})
		return
	}

	var req SubmitPaymentProofRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	proof, err := c.orderService.SubmitPaymentProof(uint(orderID), userIDUint, req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message":       "Payment proof submitted successfully",
		"payment_proof": proof,
	})
}

func (c *OrderController) GetPaymentProof(ctx *gin.Context) {
	userID, exists := ctx.Get("userID")
if !exists {
    ctx.JSON(http.StatusUnauthorized, gin.H{
        "error": "User not authenticated",
    })
    return
}
userIDUint := userID.(uint)

	orderID, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid order ID",
		})
		return
	}

	proof, err := c.orderService.GetPaymentProof(uint(orderID), userIDUint)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"payment_proof": proof,
	})
}
func (c *OrderController) UpdatePaymentProof(ctx *gin.Context) {
	userID, exists := ctx.Get("userID")
if !exists {
    ctx.JSON(http.StatusUnauthorized, gin.H{
        "error": "User not authenticated",
    })
    return
}
userIDUint := userID.(uint)

	orderID, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid order ID",
		})
		return
	}

	var req SubmitPaymentProofRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	proof, err := c.orderService.UpdatePaymentProof(uint(orderID), userIDUint, req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"message":       "Payment proof updated successfully",
		"payment_proof": proof,
	})
}

func (c *OrderController) GetAllOrdersAdmin(ctx *gin.Context) {
	orders, err := c.orderService.GetAllOrdersAdmin()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get orders",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"orders": orders,
		"count":  len(orders),
	})
}

func (c *OrderController) UpdateOrderStatusAdmin(ctx *gin.Context) {
	orderID, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid order ID",
		})
		return
	}

	var req struct {
		Status string `json:"status" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	err = c.orderService.UpdateOrderStatusAdmin(uint(orderID), req.Status)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Order status updated successfully",
	})
}

func (c *OrderController) ReviewPaymentProofAdmin(ctx *gin.Context) {
	userID, exists := ctx.Get("userID")
if !exists {
    ctx.JSON(http.StatusUnauthorized, gin.H{
        "error": "User not authenticated",
    })
    return
}
userIDUint := userID.(uint)

	proofID, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid payment proof ID",
		})
		return
	}

	var req struct {
		Status     string `json:"status" binding:"required"`
		AdminNotes string `json:"admin_notes"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Validate status
	if req.Status != "approved" && req.Status != "rejected" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Status must be 'approved' or 'rejected'",
		})
		return
	}

	err = c.orderService.ReviewPaymentProofAdmin(uint(proofID), req.Status, req.AdminNotes, userIDUint)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Payment proof reviewed successfully",
	})
}
