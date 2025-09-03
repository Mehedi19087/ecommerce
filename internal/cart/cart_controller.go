package cart

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type CartController struct {
	cartService CartService
}

func NewCartController(cartService CartService) *CartController {
	return &CartController{cartService: cartService}
}

func (c *CartController) GetCart(ctx *gin.Context) {
	userID, exists := ctx.Get("userID")
    if !exists {
    ctx.JSON(http.StatusUnauthorized, gin.H{
        "error": "User not authenticated",
    })
    return
    }

// Convert to uint (since middleware sets it as uint)
    userIDUint := userID.(uint)

	cart, err := c.cartService.GetCartByUserID(userIDUint)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve cart",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"cart":  cart,
		"total": cart.CalculateTotal(),
	})
}

type addItemRequest struct {
	ProductID uint `json:"product_id" binding:"required"`
	Quantity  int  `json:"quantity" binding:"required,min=1"`
}

func (c *CartController) AddItemToCart(ctx *gin.Context) {
	userID, exists := ctx.Get("userID")
    if !exists {
    ctx.JSON(http.StatusUnauthorized, gin.H{
        "error": "User not authenticated",
    })
    return
    }

// Convert to uint (since middleware sets it as uint)
    userIDUint := userID.(uint)

	var req addItemRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	cart, err := c.cartService.AddItemToCart(userIDUint, req.ProductID, req.Quantity)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Item added to cart",
		"cart":    cart,
		"total":   cart.CalculateTotal(),
	})
}

type updateItemRequest struct {
	Quantity int `json:"quantity" binding:"required,min=1"`
}

func (c *CartController) UpdateCartItem(ctx *gin.Context) {
	userID, exists := ctx.Get("userID")
    if !exists {
    ctx.JSON(http.StatusUnauthorized, gin.H{
        "error": "User not authenticated",
    })
    return
    }

// Convert to uint (since middleware sets it as uint)
    userIDUint := userID.(uint)

	itemID, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid item ID",
		})
		return
	}

	var req updateItemRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	cart, err := c.cartService.UpdateCartItem(userIDUint, uint(itemID), req.Quantity)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Cart item updated",
		"cart":    cart,
		"total":   cart.CalculateTotal(),
	})
}

func (c *CartController) RemoveCartItem(ctx *gin.Context) {
	userID, exists := ctx.Get("userID")
    if !exists {
    ctx.JSON(http.StatusUnauthorized, gin.H{
        "error": "User not authenticated",
    })
    return
    }

// Convert to uint (since middleware sets it as uint)
    userIDUint := userID.(uint)
	itemID, err := strconv.ParseUint(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid item ID",
		})
		return
	}

	cart, err := c.cartService.RemoveCartItem(userIDUint, uint(itemID))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Item removed from cart",
		"cart":    cart,
		"total":   cart.CalculateTotal(),
	})
}

func (c *CartController) ClearCart(ctx *gin.Context) {
	userID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"error": "You must be logged in to clear your cart",
		})
		return
	}

	err := c.cartService.ClearCart(userID.(uint))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to clear cart",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Cart cleared successfully",
	})
}
