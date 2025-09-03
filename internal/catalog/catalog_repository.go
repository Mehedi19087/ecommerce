package catalog

import (
	"gorm.io/gorm"
)

type ProductRepository interface {

	//product methods
	Create(product *Product) error
	FindByID(id uint) (*Product, error)
	FindAll() ([]*Product, error)
	Update(product *Product) error
	Delete(id uint) error
	FindBySearchTerm(searchTerm string) ([]Product, error)

	//category methods
	CreateCategory(category *Category) error

	//products by category method
	FindProductsByCategory(CategoryID uint, limit, offset int) ([]Product, int64, error)
}

type productRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepository{
		db: db,
	}
}

func (r *productRepository) Create(product *Product) error {
	return r.db.Create(product).Error
}

func (r *productRepository) FindByID(id uint) (*Product, error) {
	var product Product
	if err := r.db.First(&product, id).Error; err != nil {
		return nil, err
	}
	return &product, nil

}

func (r *productRepository) FindAll() ([]*Product, error) {
	var products []*Product

	if err := r.db.Find(&products).Error; err != nil {
		return nil, err
	}
	return products, nil
}

func (r *productRepository) Update(product *Product) error {
	return r.db.Save(product).Error
}

func (r *productRepository) Delete(id uint) error {
	return r.db.Delete(&Product{}, id).Error
}

func (r *productRepository) CreateCategory(category *Category) error {
	return r.db.Create(category).Error
}

func (r *productRepository) FindProductsByCategory(categoryID uint, limit, offset int) ([]Product, int64, error) {
	var products []Product
	var total int64

	r.db.Model(&Product{}).Where("category_id = ?", categoryID).Count(&total)

	err := r.db.Where("category_id = ?", categoryID).
		Limit(limit).
		Offset(offset).
		Find(&products).Error

	return products, total, err
}

func (r *productRepository) FindBySearchTerm(searchTerm string) ([]Product, error) {
	var products []Product
	searchPattern := "%" + searchTerm + "%"
	err := r.db.Where("name ILIKE ? OR description ILIKE ?", searchPattern, searchPattern).
		Find(&products).Error

	return products, err
}
