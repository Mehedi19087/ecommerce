package catalog

import (
	"errors"
)

type ProductService interface {

	//products method
	CreateProduct(name, image, description, sku string, price float64, stock int, categoryId uint) (*Product, error)
	GetProductByID(id uint) (*Product, error)
	ListProducts() ([]*Product, error)
	UpdateProduct(id uint, name, description, sku string, price float64, stock int, categoryId uint) (*Product, error)
	DeleteProduct(id uint) error
	SearchProducts(searchTerm string) ([]Product, error)

	//category methods

	CreateCategory(name string) (*Category, error)

	GetProductsByCategory(categoryID uint, page, pageSize int) ([]Product, int64, error)
}
type productService struct {
	repo ProductRepository
}

func NewProductService(repo ProductRepository) ProductService {
	return &productService{repo: repo}
}

func (s *productService) CreateProduct(name, image, description, sku string, price float64, stock int, categoryId uint) (*Product, error) {
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

func (s *productService) UpdateProduct(id uint, name, description, sku string, price float64, stock int, categoryId uint) (*Product, error) {
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

	if categoryId == 0 {
		product.CategoryID = categoryId
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
