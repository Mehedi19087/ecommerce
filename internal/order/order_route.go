package order

import (
	"github.com/gin-gonic/gin"
)

func SetupOrderRoutes(router *gin.Engine, orderController *OrderController) {
	v1 := router.Group("/api/v1")
	orders := v1.Group("/orders")
	{
		orders.POST("", orderController.CreateOrder)
		orders.GET("", orderController.GetUserOrders)
		orders.GET("/:id", orderController.GetOrderByID)
		orders.PUT("/:id/cancel", orderController.CancelOrder)
		orders.POST("/:id/payment-proof", orderController.SubmitPaymentProof)
		orders.GET("/:id/payment-proof", orderController.GetPaymentProof)
		orders.PUT("/:id/payment-proof", orderController.UpdatePaymentProof)

	}
	admin := v1.Group("/admin")
	{
		admin.GET("/orders", orderController.GetAllOrdersAdmin)
		admin.PUT("/orders/:id/status", orderController.UpdateOrderStatusAdmin)
		admin.PUT("/payment-proofs/:id/review", orderController.ReviewPaymentProofAdmin)
	}
}


// Order Created

// User Email:
// “Your order is pending. Please complete your payment and submit the payment proof.”
// Admin Email:
// “A new order has been placed and is awaiting payment.”
// Payment Proof Submitted

// Admin Email:
// “Payment proof has been submitted for order #XYZ. Please review and approve/reject.”
// Payment Approved by Admin

// User Email:
// “Your payment has been approved and your order is confirmed. Thank you!”