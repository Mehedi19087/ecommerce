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

	CreateSubCategory(subCategory *SubCategory) error
    CreateSubSubCategory(subSubCategory *SubSubCategory) error
    FindCategoriesWithHierarchy() ([]Category, error)
    FindSubCategoriesByCategoryID(categoryID uint) ([]SubCategory, error)
    FindSubSubCategoriesBySubCategoryID(subCategoryID uint) ([]SubSubCategory, error)
	FindCategoryByID(id uint) (*Category, error)

	FindAllCategories() ([]Category, error)
    FindSubCategoryByID(id uint) (*SubCategory, error)
    FindSubSubCategoryByID(id uint) (*SubSubCategory, error)

	FindProductsBySubCategoryID(subCategoryID uint) ([]Product, error)
	FindProductsBySubSubCategoryID(subSubCategoryID uint) ([]Product, error)

	DeleteCategory(id uint) error
    DeleteSubCategory(id uint) error
    DeleteSubSubCategory(id uint) error
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

func (r *productRepository) CreateSubCategory(subCategory *SubCategory) error {
    return r.db.Create(subCategory).Error
}

func (r *productRepository) CreateSubSubCategory(subSubCategory *SubSubCategory) error {
    return r.db.Create(subSubCategory).Error
}

func (r *productRepository) FindCategoriesWithHierarchy() ([]Category, error) {
    var categories []Category
    err := r.db.Preload("SubCategories.SubSubCategories").Find(&categories).Error
    return categories, err
}

func (r *productRepository) FindSubCategoriesByCategoryID(categoryID uint) ([]SubCategory, error) {
    var subCategories []SubCategory
    err := r.db.Where("category_id = ?", categoryID).Find(&subCategories).Error
    return subCategories, err
}

func (r *productRepository) FindSubSubCategoriesBySubCategoryID(subCategoryID uint) ([]SubSubCategory, error) {
    var subSubCategories []SubSubCategory
    err := r.db.Where("sub_category_id = ?", subCategoryID).Find(&subSubCategories).Error
    return subSubCategories, err
}

// Add this method to your productRepository struct
func (r *productRepository) FindCategoryByID(id uint) (*Category, error) {
    var category Category
    err := r.db.First(&category, id).Error
    if err != nil {
        return nil, err
    }
    return &category, nil
}

// Add these methods to your productRepository struct (after your existing methods)

// Get all categories
func (r *productRepository) FindAllCategories() ([]Category, error) {
    var categories []Category
    err := r.db.Find(&categories).Error
    return categories, err
}

// Get specific subcategory by ID
func (r *productRepository) FindSubCategoryByID(id uint) (*SubCategory, error) {
    var subCategory SubCategory
    err := r.db.First(&subCategory, id).Error
    if err != nil {
        return nil, err
    }
    return &subCategory, nil
}

// Get specific sub-subcategory by ID
func (r *productRepository) FindSubSubCategoryByID(id uint) (*SubSubCategory, error) {
    var subSubCategory SubSubCategory
    err := r.db.First(&subSubCategory, id).Error
    if err != nil {
        return nil, err
    }
    return &subSubCategory, nil
}

func (r *productRepository) FindProductsBySubCategoryID(subCategoryID uint) ([]Product, error) {
    var products []Product
    err := r.db.Where("sub_category_id = ?", subCategoryID).Find(&products).Error
    return products, err
}

func (r *productRepository) FindProductsBySubSubCategoryID(subSubCategoryID uint) ([]Product, error) {
    var products []Product
    err := r.db.Where("sub_sub_category_id = ?", subSubCategoryID).Find(&products).Error
    return products, err
}

func (r *productRepository) DeleteCategory(id uint) error {
    return r.db.Delete(&Category{}, id).Error
}

func (r *productRepository) DeleteSubCategory(id uint) error {
    return r.db.Delete(&SubCategory{}, id).Error
}

func (r *productRepository) DeleteSubSubCategory(id uint) error {
    return r.db.Delete(&SubSubCategory{}, id).Error
}