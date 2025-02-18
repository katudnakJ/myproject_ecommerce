// product_handlers.go
package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	product "productproject/internal/product"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type ProductHandlers struct {
	store *product.Store
}

func NewProductHandlers(store *product.Store) *ProductHandlers {
	return &ProductHandlers{store: store}
}

// func (h *ProductHandlers) GetBrand(c *gin.Context) {
// 	brand, err := h.store.GetBrand(c.Request.Context())
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		return
// 	}
// 	c.JSON(http.StatusOK, brand)
// }

// Bam
func (h *ProductHandlers) GetNewProdShop(c *gin.Context) {
	np, err := h.store.GetNewProdShop(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"err": err.Error()})
		return
	}
	c.JSON(http.StatusOK, np)
}

func (h *ProductHandlers) GetProducts(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "20")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit value"})
		return
	}

	params := product.ProductQueryParams{
		Limit:  limit,
		Order:  c.Query("order"),
		Search: c.Query("search"),
	}

	response, err := h.store.GetProducts(c.Request.Context(), params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *ProductHandlers) GetRecProducts(c *gin.Context) {
	np, err := h.store.GetRecProducts(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, np)

}

// Mook
func convertTimesToUserTimezone(products *product.Products, loc *time.Location) {
	products.Created_At = products.Created_At.In(loc)
	products.Updated_At = products.Updated_At.In(loc)
}

func (h *ProductHandlers) GetProduct(c *gin.Context) {
	Prod_ID := c.Param("Prod_ID")

	// ใช้ Prod_ID โดยตรงในฟังก์ชัน GetProduct
	products, err := h.store.GetProduct(c.Request.Context(), Prod_ID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	userTimezone := "Asia/Bangkok"
	loc, err := time.LoadLocation(userTimezone)
	if err != nil {
		log.Fatal("ไม่สามารถโหลด timezone ได้:", err)
	}

	convertTimesToUserTimezone(products, loc)

	c.JSON(http.StatusOK, products)
}

// Man
func (h *ProductHandlers) GetAllProducts(c *gin.Context) {
	products, err := h.store.GetAllProducts(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, products)
}

func (h *ProductHandlers) GetProductsByCategoryRoom(c *gin.Context) {
	Room_Name := c.Param("room_name")

	products, err := h.store.GetProductsByCategoryRoom(c.Request.Context(), Room_Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, products)
}

func (h *ProductHandlers) GetProductsByCategoryFurniture(c *gin.Context) {
	Fur_Name := c.Param("fur_name")

	products, err := h.store.GetProductsByCategoryFurniture(c.Request.Context(), Fur_Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, products)
}

func (h *ProductHandlers) GetProductsByCategoryRoomAndFurniture(c *gin.Context) {
	Room_Name := c.Param("room_name")
	Fur_Name := c.Param("fur_name")

	products, err := h.store.GetProductsByCategoryRoomAndFurniture(c.Request.Context(), Room_Name, Fur_Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, products)
}

// Krit
type RegisterCustomerRequest struct {
	FirstName string `json:"cust_fname"`
	LastName  string `json:"cust_lname"`
	Email     string `json:"cust_email"`
	Phone     string `json:"cust_phonenumber"`
	Username  string `json:"cust_username"`
	Password  string `json:"cust_password"`
}

func (h *ProductHandlers) RegisterCustomer(c *gin.Context) {
	var req RegisterCustomerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	customer := product.Customers{
		FName:    req.FirstName,
		LName:    req.LastName,
		CusEM:    req.Email,
		CusPhone: req.Phone,
		Username: req.Username,
		CusPW:    req.Password,
	}

	if err := h.store.CreateCustomer(c.Request.Context(), &customer); err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "customer not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusCreated, customer)
}

func (h *ProductHandlers) GetLatestProduct(c *gin.Context) {
	products, err := h.store.GetLatestProduct(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(products) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "no products found"})
		return
	}

	c.JSON(http.StatusOK, products)
}

func (h *ProductHandlers) GetAllsProducts(c *gin.Context) {
	products, err := h.store.GetAllsProducts(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(products) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "no products found"})
		return
	}

	c.JSON(http.StatusOK, products)
}

func (h *ProductHandlers) GetShopInfo(c *gin.Context) {
	shop_id := c.Param("shop_id")
	shops, err := h.store.GetShopInfo(c.Request.Context(), shop_id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(shops) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "no shops found"})
		return
	}

	c.JSON(http.StatusOK, shops)
}

type CustomerHandlers struct {
	db product.EcommerceDatabase
}

func NewCustomerHandlers(db product.EcommerceDatabase) *CustomerHandlers {
	return &CustomerHandlers{
		db: db,
	}
}

func (h *CustomerHandlers) GetCustomers(c *gin.Context) {
	customers, err := h.db.GetCustomers(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(customers) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "no customers found"})
		return
	}

	c.JSON(http.StatusOK, customers)
}

func (h *CustomerHandlers) GetCustomerByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid customer ID"})
		return
	}

	customer, err := h.db.GetCustomerByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if customer == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "customer not found"})
		return
	}

	c.JSON(http.StatusOK, customer)
}

type LoginRequest struct {
	Username string `json:"cust_username"`
	Password string `json:"cust_password"`
}

func (h *ProductHandlers) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	customer, err := h.store.GetCustomerByUsernameAndPassword(c.Request.Context(), req.Username, req.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid username or password"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, customer)
}

func (h *ProductHandlers) GetProductsByShopID(c *gin.Context) {
	shopID, err := strconv.Atoi(c.Param("shop_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid shop ID"})
		return
	}

	products, err := h.store.GetProductsByShopID(c.Request.Context(), shopID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(products) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "no products found"})
		return
	}

	c.JSON(http.StatusOK, products)
}

func (h *ProductHandlers) GetLatestProductsByShopID(c *gin.Context) {
	shopID, err := strconv.Atoi(c.Param("shop_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid shop ID"})
		return
	}

	products, err := h.store.GetLatestProductsByShopID(c.Request.Context(), shopID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if len(products) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "no products found"})
		return
	}

	c.JSON(http.StatusOK, products)
}

// Q
func (h *ProductHandlers) GetCartHandler(c *gin.Context) {
	// รับค่า Cust_ID จาก URL parameter
	Cust_ID := c.Param("cust_id")

	if Cust_ID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Customer ID is required"})
		return
	}

	// เรียกใช้ฟังก์ชัน GetCart เพื่อดึงข้อมูลตะกร้าของลูกค้า
	cart, err := h.store.GetCart(c.Request.Context(), Cust_ID)
	if err != nil {
		// ถ้าตะกร้าไม่พบ, ส่งกลับ 404
		if err.Error() == fmt.Sprintf("no cart found for customer ID: %s", Cust_ID) {
			c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("Cart not found for customer ID %s", Cust_ID)})
		} else {
			// ถ้ามีข้อผิดพลาดอื่น ๆ ส่งกลับ 500 Internal Server Error
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Error fetching cart: %v", err)})
		}
		return
	}

	// คำนวณ total_price โดยรวมราคาสินค้าทุกชิ้นในตะกร้า
	totalPrice := 0.0
	for _, item := range cart.Cart_Items {
		totalPrice += item.TotalPrice // แน่ใจว่า TotalPrice เป็น float64
	}

	// อัพเดตค่า total_price ของตะกร้า
	cart.TotalPrice = totalPrice

	// ใช้ json.MarshalIndent() เพื่อจัดรูปแบบข้อมูลให้มีการเว้นบรรทัด
	response, err := json.MarshalIndent(gin.H{
		"cart_id":     cart.Cart_ID,
		"cust_id":     cart.Cust_ID,
		"total_price": cart.TotalPrice,
		"cart_items":  cart.Cart_Items,
		"created_at":  cart.Created_At,
		"updated_at":  cart.Updated_At,
	}, "", "  ")

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to format JSON"})
		return
	}

	// ส่งผลลัพธ์ในรูปแบบ JSON ที่จัดรูปแบบให้อ่านง่าย
	c.Data(http.StatusOK, "application/json", response)
}

func (h *ProductHandlers) AddToCartHandler(c *gin.Context) {
	// รับค่า Cust_ID และ Prod_ID จาก URL parameter
	Cust_ID := c.Param("cust_id")
	Prod_ID := c.Param("prod_id")

	// รับค่า Quantity จาก Body ของ request
	var request struct {
		Quantity int `json:"quantity"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	Quantity := request.Quantity

	if Cust_ID == "" || Prod_ID == "" || Quantity <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Cust_ID, Prod_ID, or Quantity"})
		return
	}

	// เรียกใช้ฟังก์ชัน AddToCart จาก Store เพื่อเพิ่มสินค้าในตะกร้า
	err := h.store.AddToCart(c.Request.Context(), Cust_ID, Prod_ID, Quantity)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to add product to cart: %v", err)})
		return
	}

	// ส่งผลลัพธ์การเพิ่มสินค้าในตะกร้า
	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("Product %s has been added to cart successfully", Prod_ID),
	})
}

func (h *ProductHandlers) DeleteFromCartHandler(c *gin.Context) {
	// รับค่า Cust_ID และ Cart_Item_ID จาก URL parameter
	Cust_ID := c.Param("cust_id")
	Cart_Item_ID := c.Param("cart_item_id")

	// ตรวจสอบว่า Cust_ID หรือ Cart_Item_ID ว่างเปล่าหรือไม่
	if Cust_ID == "" || Cart_Item_ID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Customer ID and Cart Item ID are required"})
		return
	}

	// รับค่า Quantity จาก query parameter (ถ้ามี)
	Quantity := c.DefaultQuery("quantity", "1") // ใช้ค่าเริ่มต้นเป็น 1

	// แปลง Cart_Item_ID เป็น integer
	itemID, err := strconv.Atoi(Cart_Item_ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid cart item ID"})
		return
	}

	// แปลง Quantity เป็น integer
	quantityInt, err := strconv.Atoi(Quantity)
	if err != nil || quantityInt <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid quantity"})
		return
	}

	// เรียกใช้ฟังก์ชัน DeleteFromCart เพื่อดำเนินการลบสินค้าจากตะกร้า
	err = h.store.DeleteFromCart(c.Request.Context(), Cust_ID, itemID, quantityInt)
	if err != nil {
		// กรณีเกิดข้อผิดพลาดในฟังก์ชัน DeleteFromCart
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to delete item from cart: %v", err)})
		return
	}

	// ส่งข้อความว่าลบสินค้าสำเร็จ
	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("%d items deleted successfully from the cart", quantityInt)})
}

func (h *ProductHandlers) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "Healthy"})
}
