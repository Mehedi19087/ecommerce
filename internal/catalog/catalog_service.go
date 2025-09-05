package catalog

import (
	"errors"
	"log"
)

type ProductService interface {

	//products method
	CreateProduct(name, image, description, sku string, price float64, stock int, categoryId uint, subCategoryID, subSubCategoryID *uint) (*Product, error)
	GetProductByID(id uint) (*Product, error)
	ListProducts() ([]*Product, error)
	UpdateProduct(id uint, name, description, sku string, price float64, stock int, categoryId uint,subCategoryID, subSubCategoryID *uint) (*Product, error)
	DeleteProduct(id uint) error
	SearchProducts(searchTerm string) ([]Product, error)

	//category methods

	CreateCategory(name string) (*Category, error)

	GetProductsByCategory(categoryID uint, page, pageSize int) ([]Product, int64, error)


	CreateSubCategory(name string, categoryID uint) (*SubCategory, error)
    CreateSubSubCategory(name string, subCategoryID uint) (*SubSubCategory, error)
    GetCategoryHierarchy() ([]Category, error)
    GetSubCategoriesByCategoryID(categoryID uint) ([]SubCategory, error)
    GetSubSubCategoriesBySubCategoryID(subCategoryID uint) ([]SubSubCategory, error)

	ListCategories() ([]Category, error)
    GetCategoryByID(id uint) (*Category, error)
    GetSubCategoryByID(id uint) (*SubCategory, error)
    GetSubSubCategoryByID(id uint) (*SubSubCategory, error)

	 GetProductsBySubCategoryID(subCategoryID uint) ([]Product, error)
    GetProductsBySubSubCategoryID(subSubCategoryID uint) ([]Product, error)

	 DeleteCategory(id uint) error
    DeleteSubCategory(id uint) error
    DeleteSubSubCategory(id uint) error
}
type productService struct {
	repo ProductRepository
}

func NewProductService(repo ProductRepository) ProductService {
	return &productService{repo: repo}
}

func (s *productService) CreateProduct(name, image, description, sku string, price float64, stock int, categoryId uint , subCategoryID, subSubCategoryID *uint) (*Product, error) {
	if name == "" {
		return nil, errors.New("product name is required")
	}
	if image == "" {
		return nil, errors.New("product image is required")
	}
	if price < 0 {
		return nil, errors.New("price cannot be negative")
	}
	if stock < 0 {
		return nil, errors.New("stock cannot be negative")
	}
	product := &Product{
		Name:        name,
		Image:       image,
		Description: description,
		SKU:         sku,
		Price:       price,
		Stock:       stock,
		CategoryID:  categoryId,
		SubCategoryID: subCategoryID,
		SubSubCategoryID: subSubCategoryID,
	}
	// Save to database
	if err := s.repo.Create(product); err != nil {
		return nil, err
	}
	return product, nil
}

func (s *productService) GetProductByID(id uint) (*Product, error) {
	if id == 0 {
		return nil, errors.New("product id is required")
	}
	product, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	return product, nil
}

func (s *productService) ListProducts() ([]*Product, error) {
	products, err := s.repo.FindAll()
	if err != nil {
		return nil, err
	}
	return products, nil
}

func (s *productService) UpdateProduct(id uint, name, description, sku string, price float64, stock int, categoryId uint, subCategoryID, subSubCategoryID *uint) (*Product, error) {
	if id == 0 {
		return nil, errors.New("product id is required")
	}
	product, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("product not found")
	}

	// Update fields if provided
	if name != "" {
		product.Name = name
	}
	if description != "" {
		product.Description = description
	}

	if sku != "" {
		product.SKU = sku
	}

	if price > 0 {
		product.Price = price
	}

	if stock >= 0 {
		product.Stock = stock
	}

	if categoryId != 0 {
		product.CategoryID = categoryId
	}
	 if subCategoryID != nil {
        product.SubCategoryID = subCategoryID
    }
    if subSubCategoryID != nil {
        product.SubSubCategoryID = subSubCategoryID
    }

	// Save updated product
	if err := s.repo.Update(product); err != nil {
		return nil, err
	}
	return product, nil
}

func (s *productService) DeleteProduct(id uint) error {
	// Verify product exists
	product, err := s.repo.FindByID(id)
	if err != nil {
		return errors.New("product not found")
	}

	// Delete the product
	return s.repo.Delete(product.ID)
}

func (s *productService) CreateCategory(name string) (*Category, error) {
	if name == "" {
		return nil, errors.New("category name is required")
	}
	category := &Category{Name: name}
	err := s.repo.CreateCategory(category)
	return category, err

}

func (s *productService) GetProductsByCategory(categoryID uint, page, pageSize int) ([]Product, int64, error) {
	offset := (page - 1) * pageSize
	return s.repo.FindProductsByCategory(categoryID, pageSize, offset)
}

func (s *productService) SearchProducts(searchTerm string) ([]Product, error) {
	if searchTerm == "" {
		return nil, errors.New("search term is required")
	}
	if len(searchTerm) < 2 {
		return nil, errors.New("search term must be at least 2 characters")
	}
	products, err := s.repo.FindBySearchTerm(searchTerm)
	if err != nil {
		return nil, err
	}
	return products, nil
}


// 👈 NEW: Add these service implementations
func (s *productService) CreateSubCategory(name string, categoryID uint) (*SubCategory, error) {
    if name == "" {
		log.Println("name needed")
        return nil, errors.New("subcategory name is required")
    }

    // Check if parent category exists
    _, err := s.repo.FindCategoryByID(categoryID)
    if err != nil {
		log.Println("parent category not found")
        return nil, errors.New("parent category not found")
    }

    subCategory := &SubCategory{
        Name:       name,
        CategoryID: categoryID,
    }

    err = s.repo.CreateSubCategory(subCategory)
    if err != nil {
		log.Println("subcategory is not created , db error")
        return nil, err
    }

    return subCategory, nil
}

func (s *productService) CreateSubSubCategory(name string, subCategoryID uint) (*SubSubCategory, error) {
    if name == "" {
        return nil, errors.New("sub-subcategory name is required")
    }

    subSubCategory := &SubSubCategory{
        Name:          name,
        SubCategoryID: subCategoryID,
        ProductCount:  0,
    }

    err := s.repo.CreateSubSubCategory(subSubCategory)
    if err != nil {
        return nil, errors.New("failed to create sub-subcategory")
    }

    return subSubCategory, nil
}


func (s *productService) GetCategoryHierarchy() ([]Category, error) {
    return s.repo.FindCategoriesWithHierarchy()
}

func (s *productService) GetSubCategoriesByCategoryID(categoryID uint) ([]SubCategory, error) {
    return s.repo.FindSubCategoriesByCategoryID(categoryID)
}

func (s *productService) GetSubSubCategoriesBySubCategoryID(subCategoryID uint) ([]SubSubCategory, error) {
    return s.repo.FindSubSubCategoriesBySubCategoryID(subCategoryID)
}

// Add these methods to your productService struct (after your existing methods)

// Get all categories
func (s *productService) ListCategories() ([]Category, error) {
    return s.repo.FindAllCategories()
}

// Get specific category by ID
func (s *productService) GetCategoryByID(id uint) (*Category, error) {
    if id == 0 {
        return nil, errors.New("category id is required")
    }
    return s.repo.FindCategoryByID(id)
}

// Get specific subcategory by ID
func (s *productService) GetSubCategoryByID(id uint) (*SubCategory, error) {
    if id == 0 {
        return nil, errors.New("subcategory id is required")
    }
    return s.repo.FindSubCategoryByID(id)
}

// Get specific sub-subcategory by ID
func (s *productService) GetSubSubCategoryByID(id uint) (*SubSubCategory, error) {
    if id == 0 {
        return nil, errors.New("sub-subcategory id is required")
    }
    return s.repo.FindSubSubCategoryByID(id)
}

func (s *productService) GetProductsBySubCategoryID(subCategoryID uint) ([]Product, error) {
    if subCategoryID == 0 {
        return nil, errors.New("subcategory ID is required")
    }
    return s.repo.FindProductsBySubCategoryID(subCategoryID)
}

func (s *productService) GetProductsBySubSubCategoryID(subSubCategoryID uint) ([]Product, error) {
    if subSubCategoryID == 0 {
        return nil, errors.New("sub-subcategory ID is required")
    }
    return s.repo.FindProductsBySubSubCategoryID(subSubCategoryID)
}

func (s *productService) DeleteCategory(id uint) error {
    if id == 0 {
        return errors.New("category ID is required")
    }

    // Check if category has subcategories
    subCategories, err := s.repo.FindSubCategoriesByCategoryID(id)
    if err == nil && len(subCategories) > 0 {
        return errors.New("cannot delete category: it has subcategories")
    }

    // Check if category has products
    products, _, err := s.repo.FindProductsByCategory(id, 1, 0)
    if err == nil && len(products) > 0 {
        return errors.New("cannot delete category: it has products")
    }

    return s.repo.DeleteCategory(id)
}

func (s *productService) DeleteSubCategory(id uint) error {
    if id == 0 {
        return errors.New("subcategory ID is required")
    }

    // Check if subcategory has sub-subcategories
    subSubCategories, err := s.repo.FindSubSubCategoriesBySubCategoryID(id)
    if err == nil && len(subSubCategories) > 0 {
        return errors.New("cannot delete subcategory: it has sub-subcategories")
    }

    // Check if subcategory has products
    products, err := s.repo.FindProductsBySubCategoryID(id)
    if err == nil && len(products) > 0 {
        return errors.New("cannot delete subcategory: it has products")
    }

    return s.repo.DeleteSubCategory(id)
}

func (s *productService) DeleteSubSubCategory(id uint) error {
    if id == 0 {
        return errors.New("sub-subcategory ID is required")
    }

    // Check if sub-subcategory has products
    products, err := s.repo.FindProductsBySubSubCategoryID(id)
    if err == nil && len(products) > 0 {
        return errors.New("cannot delete sub-subcategory: it has products")
    }

    return s.repo.DeleteSubSubCategory(id)
}