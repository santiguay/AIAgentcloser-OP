package main

import (
	"log"
	"net/http"
	"time"
	"os"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"github.com/joho/godotenv"
)

// Database Models
type Categoria struct {
	ID          uint   `gorm:"primaryKey;column:id" json:"id"`
	Nombre      string `gorm:"unique;column:nombre" json:"nombre"`
	Descripcion string `gorm:"column:descripcion" json:"descripcion,omitempty"`
	CreadaEn    string `gorm:"column:creada_en" json:"creada_en"`
}

type Producto struct {
	ID          uint      `gorm:"primaryKey;column:id" json:"id"`
	Nombre      string    `gorm:"column:nombre" json:"nombre"`
	Precio      float64   `gorm:"column:precio" json:"precio"`
	Stock       int       `gorm:"column:stock" json:"stock"`
	Descripcion string    `gorm:"column:descripcion" json:"descripcion,omitempty"`
	CreadoEn    string    `gorm:"column:creado_en" json:"creado_en"`
	CategoriaID uint      `gorm:"column:categoria_id" json:"categoria_id"`
	Categoria   Categoria `gorm:"foreignKey:CategoriaID" json:"categoria"`
}

type Orden struct {
	ID            uint    `gorm:"primaryKey;column:id" json:"id"`
	NombreCliente string  `gorm:"column:nombre_cliente" json:"nombre_cliente"`
	Domicilio     string  `gorm:"column:domicilio" json:"domicilio"`
	Cedula        string  `gorm:"column:cedula" json:"cedula,omitempty"`
	Telefono      string  `gorm:"column:telefono" json:"telefono,omitempty"`
	Total         float64 `gorm:"column:total" json:"total"`
	Completa      bool    `gorm:"column:completa" json:"completa"`
	CreadaEn      string  `gorm:"column:creada_en" json:"creada_en"`
}

type DetalleVenta struct {
	ID         uint      `gorm:"primaryKey;column:id" json:"id"`
	Cantidad   int       `gorm:"column:cantidad" json:"cantidad"`
	Subtotal   float64   `gorm:"column:subtotal" json:"subtotal"`
	OrdenID    uint      `gorm:"column:orden_id" json:"orden_id"`
	ProductoID uint      `gorm:"column:producto_id" json:"producto_id"`
	Producto   Producto  `gorm:"foreignKey:ProductoID" json:"producto"`
}

// Request para creación de orden
type OrderCreationRequest struct {
	Order        Orden           `json:"order"`
	DetalleVenta []DetalleVenta  `json:"detalle_venta"`
}

// TableName overrides to use actual table names
func (Categoria) TableName() string {
	return "adminApp_categoria"
}

func (Producto) TableName() string {
	return "adminApp_producto"
}

func (Orden) TableName() string {
	return "adminApp_orden"
}

func (DetalleVenta) TableName() string {
	return "adminApp_detalleventa"
}

// Global DB instance
var db *gorm.DB

func initDatabase() {
	var err error
	dsn := os.Getenv("DATABASE_DSN")
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	log.Println("Database connected!")

	// Auto Migrate models (only if you want to create missing tables)
	db.AutoMigrate(&Categoria{}, &Producto{}, &Orden{}, &DetalleVenta{})
	log.Println("Database migrated!")
}

// Handlers for Categoria
func getCategorias(c *gin.Context) {
	var categorias []Categoria
	if result := db.Find(&categorias); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	c.JSON(http.StatusOK, categorias)
}

// Handlers for Producto
func getProductos(c *gin.Context) {
	var productos []Producto
	if result := db.Preload("Categoria").Find(&productos); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	c.JSON(http.StatusOK, productos)
}

// Handlers for Orden
func getOrdenes(c *gin.Context) {
	var ordenes []Orden
	if result := db.Find(&ordenes); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	c.JSON(http.StatusOK, ordenes)
}

func createOrdenWithDetails(c *gin.Context) {
	var orderRequest OrderCreationRequest
	
	// Bind JSON request
	if err := c.ShouldBindJSON(&orderRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Begin a database transaction
	tx := db.Begin()
	if tx.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not begin transaction"})
		return
	}

	// Set creation time if not provided
	if orderRequest.Order.CreadaEn == "" {
		orderRequest.Order.CreadaEn = time.Now().Format(time.RFC3339)
	}

	// Create the order first
	if result := tx.Create(&orderRequest.Order); result.Error != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	// Calculate total for the order
	var total float64 = 0

	// Create sale details and calculate total
	for i := range orderRequest.DetalleVenta {
		// Link the sale detail to the created order
		orderRequest.DetalleVenta[i].OrdenID = orderRequest.Order.ID

		// Fetch product to get price
		var producto Producto
		if result := tx.First(&producto, orderRequest.DetalleVenta[i].ProductoID); result.Error != nil {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Product ID"})
			return
		}

		// Calculate subtotal
		orderRequest.DetalleVenta[i].Subtotal = float64(orderRequest.DetalleVenta[i].Cantidad) * producto.Precio
		total += orderRequest.DetalleVenta[i].Subtotal

		// Create sale detail
		if result := tx.Create(&orderRequest.DetalleVenta[i]); result.Error != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
			return
		}

		// Update product stock
		if result := tx.Model(&producto).Update("stock", producto.Stock - orderRequest.DetalleVenta[i].Cantidad); result.Error != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not update product stock"})
			return
		}
	}

	// Update order total
	if result := tx.Model(&orderRequest.Order).Update("total", total); result.Error != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not update order total"})
		return
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not commit transaction"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"order": orderRequest.Order,
		"detalle_venta": orderRequest.DetalleVenta,
	})
}

func updateOrden(c *gin.Context) {
	id := c.Param("id")
	var orden Orden
	if result := db.First(&orden, id); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Orden not found"})
		return
	}

	var input Orden
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db.Model(&orden).Updates(input)
	c.JSON(http.StatusOK, orden)
}

func deleteOrden(c *gin.Context) {
	id := c.Param("id")
	if result := db.Delete(&Orden{}, id); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Orden deleted"})
}

func main() {
	// Initialize database
	initDatabase()

	// Create Gin router
	r := gin.Default()

	// Routes
	r.GET("/categorias", getCategorias)
	r.GET("/productos", getProductos)
	r.GET("/ordenes", getOrdenes)
	r.POST("/ordenes", createOrdenWithDetails)
	r.PUT("/ordenes/:id", updateOrden)
	r.DELETE("/ordenes/:id", deleteOrden)

	// Run server
	r.Run(":8082")
}