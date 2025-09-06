package catalog

import (
	"net/http"
	"strconv"
    "os"
    "context"

	"github.com/gin-gonic/gin"
    "github.com/cloudinary/cloudinary-go/v2"
    "github.com/cloudinary/cloudinary-go/v2/api/uploader"

	// "bytes"
	// "io"
	// "log"
)

type ProductController struct {
	productService ProductService
}

func NewProductController(productService ProductService) *ProductController {
	return &ProductController{productService: productService}
}


func (ct *ProductController) UploadImage(c *gin.Context) {
    // Get uploaded file
    file, err := c.FormFile("image") // assuming frontend sends "image" field
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "No file uploaded"})
        return
    }

    // Open the file
    src, err := file.Open()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not open file"})
        return
    }
    defer src.Close()

    // Initialize Cloudinary
    cld, err := cloudinary.NewFromParams(
        os.Getenv("CLOUDINARY_CLOUD_NAME"),
        os.Getenv("CLOUDINARY_API_KEY"),
        os.Getenv("CLOUDINARY_API_SECRET"),
    )
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Cloudinary config failed"})
        return
    }

    // Upload to Cloudinary
    result, err := cld.Upload.Upload(
        context.Background(),
        src,
        uploader.UploadParams{
            ResourceType: "image", // for images
            Folder:       "uploads", // optional: organize in folder
        },
    )

    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Upload failed: " + err.Error()})
        return
    }

    // Return the secure URL
    c.JSON(http.StatusOK, gin.H{
        "url": result.SecureURL,
        "public_id": result.PublicID,
    })
}

// request structs
type createProductRequest struct {
	Name        string   `json:"name" binding:"required"`
	Image       []string `json:"images" binding:"required,min=1"`
	Description string   `json:"description"`
	SKU         string   `json:"sku" binding:"required"`
	Price       float64  `json:"price" binding:"required,min=0"`
	Stock       int      `json:"stock" binding:"min=0"`
	CategoryID  uint     `json:"category_id"`
	SubCategoryID    *uint `json:"sub_category_id,omitempty"`          // Optional
    SubSubCategoryID *uint `json:"sub_sub_category_id,omitempty"` 

}

func (c *ProductController) CreateProduct(ctx *gin.Context) {

	var req createProductRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	product, err := c.productService.CreateProduct(
		req.Name,
		req.Image,
		req.Description,
		req.SKU,
		req.Price,
		req.Stock,
		req.CategoryID,
		req.SubCategoryID,
		req.SubSubCategoryID,
	)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"message": "Product created successfully",
		"product": product,
	})
}

func (c *ProductController) GetProductByID(ctx *gin.Context) {
    idStr := ctx.Param("id")
    id, err := strconv.ParseUint(idStr, 10, 64)
    if err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{
            "error": "invalid product id Format",
        })
        return
    }
    product, err := c.productService.GetProductByID(uint(id))
    if err != nil {
        ctx.JSON(http.StatusNotFound, gin.H{
            "error": err.Error(),
        })
        return
    }

    // ✅ Convert filenames to full URLs


    ctx.JSON(http.StatusOK, gin.H{
        "product": map[string]interface{}{
            "id":          product.ID,
            "name":        product.Name,
            "images":      []string(product.Image), // ✅ Full URLs instead of filenames
            "description": product.Description,
            "sku":         product.SKU,
            "price":       product.Price,
            "stock":       product.Stock,
            "category_id": product.CategoryID,
        },
    })
}
func (c *ProductController) ListProducts(ctx *gin.Context) {
	products, err := c.productService.ListProducts()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve products",
		})
		return
	}

	var productsResponse []map[string]interface{}
    for _, product := range products {
        productsResponse = append(productsResponse, map[string]interface{}{
            "id":          product.ID,
            "name":        product.Name,
            "images":      []string(product.Image), // ✅ Return original URLs
            "description": product.Description,
            "sku":         product.SKU,
            "price":       product.Price,
            "stock":       product.Stock,
            "category_id": product.CategoryID,
            "sub_category_id": product.SubCategoryID,
            "sub_sub_category_id": product.SubSubCategoryID,
            "created_at":  product.CreatedAt,
            "updated_at":  product.UpdatedAt,
        })
    }

    ctx.JSON(http.StatusOK, gin.H{
        "products": productsResponse,
    })
}

type updateProductRequest struct {
	Name        string  `json:"name"`
	Image       []string `json:"images"` 
	Description string  `json:"description"`
	SKU         string  `json:"sku"`
	Price       float64 `json:"price" binding:"min=0"`
	Stock       int     `json:"stock" binding:"min=0"`
	CategoryID  uint    `json:"category_id"`
	SubCategoryID    *uint   `json:"sub_category_id,omitempty"`
    SubSubCategoryID *uint   `json:"sub_sub_category_id,omitempty"`
}

func (c *ProductController) UpdateProduct(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid product ID format",
		})
		return
	}
	var req updateProductRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	product, err := c.productService.UpdateProduct(
		uint(id),
		req.Name,
        req.Image,
		req.Description,
		req.SKU,
		req.Price,
		req.Stock,
		req.CategoryID,
		req.SubCategoryID,
        req.SubSubCategoryID,
	)

	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Product updated successfully",
		"product": product,
	})
}

func (c *ProductController) DeleteProduct(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid product ID format",
		})
		return
	}

	err = c.productService.DeleteProduct(uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Product deleted successfully",
	})
}

// category controllers
type CreateCategoryRequest struct {
	Name string `json:"name" binding:"required"`
}

func (c *ProductController) CreateCategory(ctx *gin.Context) {
	var req CreateCategoryRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error()})
		return
	}

	category, err := c.productService.CreateCategory(req.Name)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{"category": category})
}

func (c *ProductController) GetProductsByCategory(ctx *gin.Context) {
	categoryIDStr := ctx.Param("id")
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "20"))

	var products []Product
	var total int64
	var err error

	if categoryIDStr != "" {
		categoryID, parseErr := strconv.ParseUint(categoryIDStr, 10, 32)
		if parseErr != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid category_id format",
			})
			return
		}

		// ✅ Correct: Get products by specific category
		products, total, err = c.productService.GetProductsByCategory(uint(categoryID), page, pageSize)
	}

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve products",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"products":  products,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

func (c *ProductController) SearchProducts(ctx *gin.Context) {
	searchTerm := ctx.Query("q")
	if searchTerm == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "Search Term is Required",
		})
		return
	}
	products, err := c.productService.SearchProducts(searchTerm)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"products": products,
	})
}


type CreateSubCategoryRequest struct {
    Name       string `json:"name" binding:"required"`
    CategoryID uint   `json:"category_id" binding:"required"`
}

type CreateSubSubCategoryRequest struct {
    Name          string `json:"name" binding:"required"`
    SubCategoryID uint   `json:"sub_category_id" binding:"required"`
}

func (c *ProductController) CreateSubCategory(ctx *gin.Context) {
    var req CreateSubCategoryRequest
    if err := ctx.ShouldBindJSON(&req); err != nil {
        ctx.JSON(400, gin.H{"error": err.Error()})
        return
    }

    subCategory, err := c.productService.CreateSubCategory(req.Name, req.CategoryID)
    if err != nil {
        ctx.JSON(500, gin.H{"error": err.Error()})
        return
    }

    ctx.JSON(201, gin.H{"subcategory": subCategory})
}

func (c *ProductController) CreateSubSubCategory(ctx *gin.Context) {
    var req CreateSubSubCategoryRequest
    if err := ctx.ShouldBindJSON(&req); err != nil {
        ctx.JSON(400, gin.H{"error": err.Error()})
        return
    }

    subSubCategory, err := c.productService.CreateSubSubCategory(req.Name, req.SubCategoryID)
    if err != nil {
        ctx.JSON(500, gin.H{"error": err.Error()})
        return
    }

    ctx.JSON(201, gin.H{"sub_subcategory": subSubCategory})
}

func (c *ProductController) GetCategoryHierarchy(ctx *gin.Context) {
    categories, err := c.productService.GetCategoryHierarchy()
    if err != nil {
        ctx.JSON(500, gin.H{"error": err.Error()})
        return
    }

    ctx.JSON(200, gin.H{"categories": categories})
}

// Add these methods to your existing ProductController

// Get all categories
func (c *ProductController) ListCategories(ctx *gin.Context) {
    categories, err := c.productService.ListCategories()
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{
            "error": "Failed to retrieve categories",
        })
        return
    }

    ctx.JSON(http.StatusOK, gin.H{
        "categories": categories,
    })
}

// Get category by ID
func (c *ProductController) GetCategoryByID(ctx *gin.Context) {
    idStr := ctx.Param("id")
    id, err := strconv.ParseUint(idStr, 10, 64)
    if err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{
            "error": "Invalid category ID format",
        })
        return
    }

    category, err := c.productService.GetCategoryByID(uint(id))
    if err != nil {
        ctx.JSON(http.StatusNotFound, gin.H{
            "error": "Category not found",
        })
        return
    }

    ctx.JSON(http.StatusOK, gin.H{
        "category": category,
    })
}

// Get subcategories by category ID
func (c *ProductController) GetSubCategoriesByCategory(ctx *gin.Context) {
    categoryIDStr := ctx.Param("id")
    categoryID, err := strconv.ParseUint(categoryIDStr, 10, 64)
    if err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{
            "error": "Invalid category ID format",
        })
        return
    }

    subCategories, err := c.productService.GetSubCategoriesByCategoryID(uint(categoryID))
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{
            "error": "Failed to retrieve subcategories",
        })
        return
    }

    ctx.JSON(http.StatusOK, gin.H{
        "subcategories": subCategories,
    })
}

// Get specific subcategory by ID
func (c *ProductController) GetSubCategoryByID(ctx *gin.Context) {
    idStr := ctx.Param("id")
    id, err := strconv.ParseUint(idStr, 10, 64)
    if err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{
            "error": "Invalid subcategory ID format",
        })
        return
    }

    subCategory, err := c.productService.GetSubCategoryByID(uint(id))
    if err != nil {
        ctx.JSON(http.StatusNotFound, gin.H{
            "error": "Subcategory not found",
        })
        return
    }

    ctx.JSON(http.StatusOK, gin.H{
        "subcategory": subCategory,
    })
}

// Get sub-subcategories by subcategory ID
func (c *ProductController) GetSubSubCategoriesBySubCategory(ctx *gin.Context) {
    subCategoryIDStr := ctx.Param("id")
    subCategoryID, err := strconv.ParseUint(subCategoryIDStr, 10, 64)
    if err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{
            "error": "Invalid subcategory ID format",
        })
        return
    }

    subSubCategories, err := c.productService.GetSubSubCategoriesBySubCategoryID(uint(subCategoryID))
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{
            "error": "Failed to retrieve sub-subcategories",
        })
        return
    }

    ctx.JSON(http.StatusOK, gin.H{
        "sub_subcategories": subSubCategories,
    })
}

// Get specific sub-subcategory by ID
func (c *ProductController) GetSubSubCategoryByID(ctx *gin.Context) {
    idStr := ctx.Param("id")
    id, err := strconv.ParseUint(idStr, 10, 64)
    if err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{
            "error": "Invalid sub-subcategory ID format",
        })
        return
    }

    subSubCategory, err := c.productService.GetSubSubCategoryByID(uint(id))
    if err != nil {
        ctx.JSON(http.StatusNotFound, gin.H{
            "error": "Sub-subcategory not found",
        })
        return
    }

    ctx.JSON(http.StatusOK, gin.H{
        "sub_subcategory": subSubCategory,
    })
}


// Get products by subcategory
func (c *ProductController) GetProductsBySubCategory(ctx *gin.Context) {
    subCategoryIDStr := ctx.Param("id")
    subCategoryID, err := strconv.ParseUint(subCategoryIDStr, 10, 64)
    if err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{
            "error": "Invalid subcategory ID format",
        })
        return
    }

    products, err := c.productService.GetProductsBySubCategoryID(uint(subCategoryID))
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{
            "error": "Failed to retrieve products",
        })
        return
    }

    ctx.JSON(http.StatusOK, gin.H{
        "products": products,
    })
}

// Get products by sub-subcategory
func (c *ProductController) GetProductsBySubSubCategory(ctx *gin.Context) {
    subSubCategoryIDStr := ctx.Param("id")
    subSubCategoryID, err := strconv.ParseUint(subSubCategoryIDStr, 10, 64)
    if err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{
            "error": "Invalid sub-subcategory ID format",
        })
        return
    }

    products, err := c.productService.GetProductsBySubSubCategoryID(uint(subSubCategoryID))
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{
            "error": "Failed to retrieve products",
        })
        return
    }

    ctx.JSON(http.StatusOK, gin.H{
        "products": products,
    })
}

// Delete category
func (c *ProductController) DeleteCategory(ctx *gin.Context) {
    idStr := ctx.Param("id")
    id, err := strconv.ParseUint(idStr, 10, 64)
    if err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{
            "error": "Invalid category ID format",
        })
        return
    }

    err = c.productService.DeleteCategory(uint(id))
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{
            "error": err.Error(),
        })
        return
    }

    ctx.JSON(http.StatusOK, gin.H{
        "message": "Category deleted successfully",
    })
}

// Delete subcategory
func (c *ProductController) DeleteSubCategory(ctx *gin.Context) {
    idStr := ctx.Param("id")
    id, err := strconv.ParseUint(idStr, 10, 64)
    if err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{
            "error": "Invalid subcategory ID format",
        })
        return
    }

    err = c.productService.DeleteSubCategory(uint(id))
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{
            "error": err.Error(),
        })
        return
    }

    ctx.JSON(http.StatusOK, gin.H{
        "message": "Subcategory deleted successfully",
    })
}

// Delete sub-subcategory
func (c *ProductController) DeleteSubSubCategory(ctx *gin.Context) {
    idStr := ctx.Param("id")
    id, err := strconv.ParseUint(idStr, 10, 64)
    if err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{
            "error": "Invalid sub-subcategory ID format",
        })
        return
    }

    err = c.productService.DeleteSubSubCategory(uint(id))
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{
            "error": err.Error(),
        })
        return
    }

    ctx.JSON(http.StatusOK, gin.H{
        "message": "Sub-subcategory deleted successfully",
    })
}

type UpdateCategoryRequest struct {
    Name string `json:"name" binding:"required"`
}

func (c *ProductController) UpdateCategory(ctx *gin.Context) {
    idStr := ctx.Param("id")
    id, err := strconv.ParseUint(idStr, 10, 64)
    if err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{
            "error": "Invalid category ID format",
        })
        return
    }

    var req UpdateCategoryRequest
    if err := ctx.ShouldBindJSON(&req); err != nil {
        ctx.JSON(http.StatusBadRequest, gin.H{
            "error": err.Error(),
        })
        return
    }

    category, err := c.productService.UpdateCategory(uint(id), req.Name)
    if err != nil {
        ctx.JSON(http.StatusNotFound, gin.H{
            "error": err.Error(),
        })
        return
    }

    ctx.JSON(http.StatusOK, gin.H{
        "message":  "Category updated successfully",
        "category": category,
    })
}

// ✅ ADD THIS METHOD
func (c *ProductController) ListSubCategories(ctx *gin.Context) {
    subCategories, err := c.productService.ListSubCategories()
    if err != nil {
        ctx.JSON(http.StatusInternalServerError, gin.H{
            "error": "Failed to retrieve subcategories",
        })
        return
    }

    ctx.JSON(http.StatusOK, gin.H{
        "subcategories": subCategories,
    })
}
