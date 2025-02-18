// ecommerce.go
package ecommerce

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	_ "github.com/lib/pq"
)

type Customers struct {
	CusID    int    `db:"cust_id" json:"cust_id"`
	FName    string `db:"cust_fname" json:"cust_fname"`
	LName    string `db:"cust_lname" json:"cust_lname"`
	CusEM    string `db:"cust_email" json:"cust_email"`
	CusPhone string `db:"cust_phonenumber" json:"cust_phonenumber"`
	Username string `db:"cust_username" json:"cust_username"`
	CusPW    string `db:"cust_password" json:"cust_password"`
}

type Products struct {
	Prod_ID      string    `json:"Prod_ID"`
	Prod_Name    string    `json:"Prod_Name"`
	Prod_Price   float64   `json:"Prod_Price"`
	Prod_Details string    `json:"Prod_Details"`
	Prod_Image   string    `json:"Prod_Image"`
	Stock        int       `json:"Stock"`
	Sales_Amount int       `json:"Sales_Amount"`
	Status_Rec   bool      `json:"Status_Rec"`
	Created_At   time.Time `json:"Created_At"`
	Updated_At   time.Time `json:"Updated_At"`
	Room_type
	Fur_type
	Brand
	Shops
}

type Shops struct {
	Shop_ID          string `json:"Shop_ID"`
	Shop_Name        string `json:"Shop_Name"`
	Shop_Des         string `json:"Shop_Des"`
	Shop_Image       string `json:"Shop_Image"`
	Shop_Address     string `json:"Shop_Address"`
	Shop_PhoneNumber string `json:"Shop_PhoneNumber"`
	Shop_Email       string `json:"Shop_Email"`
}

type Brand struct {
	Brand_ID    string `json:"Brand_ID"`
	Brand_Name  string `json:"Brand_Name"`
	Brand_Image string `json:"Brand_Image"`
}

type Room_type struct {
	Room_ID   string `json:"Room_ID"`
	Room_Name string `json:"Room_Name"`
}

type Fur_type struct {
	Fur_ID   string `json:"Fur_ID"`
	Fur_Name string `json:"Fur_Name"`
}

// Query path
type ProductQueryParams struct {
	Limit  int    `json:"limit"`
	Order  string `json:"order"`
	Search string `json:"search"`
}

// Query path response
type ProductResponse struct {
	Items []Products `json:"items"`
	Limit int        `json:"limit"`
}

type CartItem struct {
	Prod_ID      string  `json:"Prod_ID"`
	Prod_Image   string  `json:"Prod_Image"`
	Prod_Name    string  `json:"Prod_Name"`
	Brand_Name   string  `json:"brand_name"`
	Prod_Price   float64 `json:"Prod_Price"`
	TotalPrice   float64 `json:"TotalPrice"`
	Quantity     int     `json:"Quantity"`
	Cart_Item_ID int     `json:"Cart_Item_ID"`
}

type Cart struct {
	Cart_ID    string     `json:"Cart_ID"`    // รหัสตะกร้า
	Cust_ID    string     `json:"Cust_ID"`    // รหัสลูกค้า (FK เชื่อมกับ Customers)
	Cart_Items []CartItem `json:"Cart_Items"` // รายการสินค้าในตะกร้า
	TotalPrice float64    `json:"TotalPrice"` // ราคารวมทั้งหมดของสินค้าทั้งหมดในตะกร้า
	Created_At time.Time  `json:"Created_At"` // วันที่สร้างตะกร้า
	Updated_At time.Time  `json:"Updated_At"` // วันที่อัพเดตตะกร้า
}

type EcommerceDatabase interface {
	// GetBrand(ctx context.Context) ([]Brand, error)
	GetNewProdShop(ctx context.Context) ([]Products, error)
	GetRecProducts(ctx context.Context) ([]Products, error)
	GetProducts(ctx context.Context, params ProductQueryParams) (*ProductResponse, error)

	//Mook
	GetProduct(ctx context.Context, Prod_ID string) (*Products, error)

	//Man
	GetAllProducts(ctx context.Context) ([]Products, error)
	GetProductsByCategoryRoom(ctx context.Context, Room_Name string) ([]Products, error)
	GetProductsByCategoryFurniture(ctx context.Context, Fur_Name string) ([]Products, error)
	GetProductsByCategoryRoomAndFurniture(ctx context.Context, Room_Name string, Fur_Name string) ([]Products, error)

	//krit
	CreateCustomer(ctx context.Context, customer *Customers) error
	GetCustomers(ctx context.Context) ([]Customers, error)
	GetCustomerByID(ctx context.Context, id int) (*Customers, error)
	GetLatestProduct(ctx context.Context) ([]Products, error)
	GetAllsProducts(ctx context.Context) ([]Products, error)
	GetShopInfo(ctx context.Context, id string) ([]Shops, error)
	GetCustomerByUsernameAndPassword(ctx context.Context, username, password string) (*Customers, error)
	GetProductsByShopID(ctx context.Context, shopID int) ([]Products, error)
	GetLatestProductsByShopID(ctx context.Context, shopID int) ([]Products, error)

	//Q
	GetCart(ctx context.Context, Cust_ID string) (Cart, error)
	AddToCart(ctx context.Context, Cust_ID string, Prod_ID string, Quantity int) error
	DeleteFromCart(ctx context.Context, Cust_ID string, Cart_Item_ID int, Quantity int) error

	Close() error
	Ping() error
	Reconnect(connStr string) error
}

type PostgresDatabase struct {
	db *sql.DB
}

func NewPostgresDatabase(connStr string) (*PostgresDatabase, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %v", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(5 * time.Minute)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %v", err)
	}

	return &PostgresDatabase{db: db}, nil
}

// func (pdb *PostgresDatabase) GetBrand(ctx context.Context) ([]Brand, error) {
// 	rows, err := pdb.db.QueryContext(ctx, "select Brand_ID,Brand_Name from brand")
// 	if err != nil {
// 		return nil, fmt.Errorf("Fail to Query Brand")
// 	}
// 	var response []Brand
// 	for rows.Next() {
// 		var p Brand
// 		if err := rows.Scan(&p.Brand_ID, &p.Brand_Name); err != nil {
// 			return nil, fmt.Errorf("fail to Scan Brands")
// 		}
// 		response = append(response, p)
// 	}
// 	defer rows.Close()

//		return response, nil
//	}

// Bam
func (pdb *PostgresDatabase) GetRecProducts(ctx context.Context) ([]Products, error) {
	rows, err := pdb.db.QueryContext(ctx, `select Prod_ID,Prod_Name,Prod_Price,Prod_Image,Created_At,Updated_At,Brand_Name,Status_Rec,Sales_Amount
	from products p left join brand b on p.Brand_ID = b.Brand_ID
	left join shops s on p.Shop_ID = s.shop_ID where Status_Rec = true LIMIT 3`)
	if err != nil {
		return nil, fmt.Errorf("fail to get Products")
	}
	var response []Products
	for rows.Next() {
		var p Products
		err := rows.Scan(&p.Prod_ID, &p.Prod_Name, &p.Prod_Price, &p.Prod_Image, &p.Created_At, &p.Updated_At, &p.Brand_Name, &p.Status_Rec, &p.Sales_Amount)
		if err != nil {
			return nil, fmt.Errorf("fail to Scan products")
		}

		response = append(response, p)
	}
	defer rows.Close()
	return response, nil

}
func (pdb *PostgresDatabase) GetNewProdShop(ctx context.Context) ([]Products, error) {
	rows, err := pdb.db.QueryContext(ctx, `SELECT DISTINCT ON (s.shop_ID) 
    p.Prod_ID, 
    p.Prod_Name, 
    p.Prod_Price, 
    p.Prod_Image, 
    p.Sales_Amount, 
    p.Created_At, 
    p.Updated_At, 
    b.Brand_Name, 
    s.shop_ID
FROM 
    products p
LEFT JOIN 
    brand b ON p.Brand_ID = b.Brand_ID
LEFT JOIN 
    shops s ON p.Shop_ID = s.shop_ID
ORDER BY 
    s.shop_ID, p.Created_At DESC;
`)
	if err != nil {
		return nil, fmt.Errorf("Fail to get product")
	}

	var response []Products
	for rows.Next() {
		var p Products
		err := rows.Scan(&p.Prod_ID, &p.Prod_Name, &p.Prod_Price, &p.Prod_Image, &p.Sales_Amount, &p.Created_At, &p.Updated_At, &p.Brand_Name, &p.Shop_ID)
		if err != nil {
			return nil, fmt.Errorf("Fail to scan products")
		}
		response = append(response, p)
	}
	defer rows.Close()
	return response, nil
}

// Search
func (pdb *PostgresDatabase) GetProducts(ctx context.Context, params ProductQueryParams) (*ProductResponse, error) {
	query := `
        SELECT p.Prod_ID, p.Prod_Name, p.Prod_Details, p.Prod_Price, p.Sales_Amount, b.Brand_Name, p.Prod_Image
        FROM products p 
		LEFT JOIN brand b on p.Brand_ID = b.Brand_ID
		LEFT JOIN shops s on p.Shop_ID = s.Shop_ID
		LEFT JOIN Room_Type rt on p.Room_ID = rt.Room_ID
		LEFT JOIN Fur_type ft on p.Fur_ID = ft.Fur_ID 
        WHERE 1=1`

	args := []interface{}{}
	placeholderCount := 1
	// การจัดการพารามิเตอร์ search
	if params.Search != "" {
		query += fmt.Sprintf(" AND (p.Prod_Name ILIKE $%d OR s.Shop_Name ILIKE $%d OR rt.Room_Name ILIKE $%d OR ft.Fur_Name ILIKE $%d)", placeholderCount, placeholderCount+1, placeholderCount+2, placeholderCount+3)
		args = append(args, "%"+params.Search+"%", "%"+params.Search+"%", "%"+params.Search+"%", "%"+params.Search+"%")
		placeholderCount += 4
	}

	// การจัดการ ORDER BY ด้วย sort และ order
	orderDirections := map[string]string{
		"asc":  "ASC",
		"desc": "DESC",
	}

	orderDirection, ok := orderDirections[strings.ToLower(params.Order)]
	if !ok {
		query += fmt.Sprintf(" ORDER BY p.Prod_Name %s", orderDirection)
	} else {
		query += fmt.Sprintf(" ORDER BY p.Prod_Price %s", orderDirection)
	}
	// การกำหนดค่า limit
	limit := 20
	if params.Limit > 0 && params.Limit <= 100 {
		limit = params.Limit
	}
	query += fmt.Sprintf(" LIMIT $%d", placeholderCount)
	args = append(args, limit+1)
	placeholderCount++

	// ดำเนินการ query และประมวลผลผลลัพธ์
	rows, err := pdb.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get products: %v", err)
	}
	defer rows.Close()

	var products []Products
	for rows.Next() {
		var product Products
		if err := rows.Scan(
			&product.Prod_ID, &product.Prod_Name, &product.Prod_Details, &product.Prod_Price, &product.Sales_Amount, &product.Brand_Name, &product.Prod_Image); err != nil {
			return nil, fmt.Errorf("failed to scan product: %v", err)
		}

		products = append(products, product)
		if len(products) == limit+1 {
			break
		}
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed to iterate over products: %v", err)
	}

	response := &ProductResponse{
		Items: products[:min(len(products), limit)],
		Limit: limit,
	}

	return response, nil
}

// Mook Func
func (pdb *PostgresDatabase) GetProduct(ctx context.Context, Prod_ID string) (*Products, error) {
	var products Products

	query := `
	SELECT 
		p.Prod_ID, p.Prod_Name, p.Prod_Price, p.Prod_Details, p.Prod_Image, 
		p.Stock, p.Sales_Amount, p.Status_Rec, p.Created_At, p.Updated_At,
		rt.Room_ID, rt.Room_Name, 
		ft.Fur_ID, ft.Fur_Name, 
		b.Brand_ID, b.Brand_Name, b.Brand_Image,
		s.Shop_ID, s.Shop_Name, s.Shop_Des, s.Shop_Image, s.Shop_Address, s.Shop_PhoneNumber, s.Shop_Email
	FROM 
		products p
	LEFT JOIN room_type rt ON p.room_id = rt.Room_ID
	LEFT JOIN fur_type ft ON p.fur_id = ft.Fur_ID
	LEFT JOIN brand b ON p.brand_id = b.Brand_ID
	LEFT JOIN shops s ON p.shop_id = s.Shop_ID
	WHERE 
		p.Prod_ID = $1`

	err := pdb.db.QueryRowContext(ctx, query, Prod_ID).Scan(
		&products.Prod_ID, &products.Prod_Name, &products.Prod_Price, &products.Prod_Details, &products.Prod_Image,
		&products.Stock, &products.Sales_Amount, &products.Status_Rec, &products.Created_At, &products.Updated_At,
		&products.Room_type.Room_ID, &products.Room_type.Room_Name,
		&products.Fur_type.Fur_ID, &products.Fur_type.Fur_Name,
		&products.Brand.Brand_ID, &products.Brand.Brand_Name, &products.Brand.Brand_Image,
		&products.Shops.Shop_ID, &products.Shops.Shop_Name, &products.Shops.Shop_Des, &products.Shops.Shop_Image, &products.Shops.Shop_Address, &products.Shops.Shop_PhoneNumber, &products.Shops.Shop_Email,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("product not found")
		}
		return nil, fmt.Errorf("failed to get product: %v", err)
	}

	return &products, nil
}

// Man func
func (pdb *PostgresDatabase) GetAllProducts(ctx context.Context) ([]Products, error) {
	rows, err := pdb.db.QueryContext(ctx, `select p.prod_id, p.prod_name, p.prod_price, p.prod_details, p.prod_image, p.stock,
										p.sales_amount, p.status_rec, p.created_at, p.updated_at, p.brand_id, b.brand_name,
										p.room_id, r.room_name, p.fur_id, f.fur_name, p.shop_id, s.shop_name, s.shop_des,
										s.shop_image, s.shop_address, s.shop_phonenumber, s.shop_email from products p
										join brand b on p.brand_id = b.brand_id
										join room_type r on p.room_id = r.room_id
										join fur_type f on p.fur_id = f.fur_id
										join shops s on p.shop_id = s.shop_id`)
	if err != nil {
		return nil, fmt.Errorf("failed to get all products: %v", err)
	}
	defer rows.Close()

	var products []Products
	for rows.Next() {
		var prod Products
		if err := rows.Scan(&prod.Prod_ID, &prod.Prod_Name, &prod.Prod_Price, &prod.Prod_Details,
			&prod.Prod_Image, &prod.Stock, &prod.Sales_Amount, &prod.Status_Rec, &prod.Created_At,
			&prod.Updated_At, &prod.Brand_ID, &prod.Brand_Name, &prod.Room_ID, &prod.Room_Name,
			&prod.Fur_ID, &prod.Fur_Name, &prod.Shop_ID, &prod.Shop_Name, &prod.Shop_Des, &prod.Shop_Image,
			&prod.Shop_Address, &prod.Shop_PhoneNumber, &prod.Shop_Email); err != nil {
			return nil, fmt.Errorf("failed to scan product: %v", err)
		}
		products = append(products, prod)
	}
	return products, nil
}

func (pdb *PostgresDatabase) GetProductsByCategoryRoom(ctx context.Context, Room_Name string) ([]Products, error) {
	rows, err := pdb.db.QueryContext(ctx, `select p.prod_id, p.prod_name, p.prod_price, p.prod_details, p.prod_image, p.stock,
									p.sales_amount, p.status_rec, p.created_at, p.updated_at, p.brand_id, b.brand_name,
									p.room_id, r.room_name, p.fur_id, f.fur_name, p.shop_id, s.shop_name, s.shop_des,
									s.shop_image, s.shop_address, s.shop_phonenumber, s.shop_email from products p
									join brand b on p.brand_id = b.brand_id
									join room_type r on p.room_id = r.room_id
									join fur_type f on p.fur_id = f.fur_id
									join shops s on p.shop_id = s.shop_id
									where r.room_name = $1`, Room_Name)
	if err != nil {
		return nil, fmt.Errorf("failed to get products by category room: %v", err)
	}
	defer rows.Close()

	var products []Products
	for rows.Next() {
		var prod Products
		if err := rows.Scan(&prod.Prod_ID, &prod.Prod_Name, &prod.Prod_Price, &prod.Prod_Details,
			&prod.Prod_Image, &prod.Stock, &prod.Sales_Amount, &prod.Status_Rec, &prod.Created_At,
			&prod.Updated_At, &prod.Brand_ID, &prod.Brand_Name, &prod.Room_ID, &prod.Room_Name,
			&prod.Fur_ID, &prod.Fur_Name, &prod.Shop_ID, &prod.Shop_Name, &prod.Shop_Des, &prod.Shop_Image,
			&prod.Shop_Address, &prod.Shop_PhoneNumber, &prod.Shop_Email); err != nil {
			return nil, fmt.Errorf("failed to scan product: %v", err)
		}
		products = append(products, prod)
	}
	return products, nil
}

func (pdb *PostgresDatabase) GetProductsByCategoryFurniture(ctx context.Context, Fur_Name string) ([]Products, error) {
	rows, err := pdb.db.QueryContext(ctx, `select p.prod_id, p.prod_name, p.prod_price, p.prod_details, p.prod_image, p.stock,
									p.sales_amount, p.status_rec, p.created_at, p.updated_at, p.brand_id, b.brand_name,
									p.room_id, r.room_name, p.fur_id, f.fur_name, p.shop_id, s.shop_name, s.shop_des,
									s.shop_image, s.shop_address, s.shop_phonenumber, s.shop_email from products p
									join brand b on p.brand_id = b.brand_id
									join room_type r on p.room_id = r.room_id
									join fur_type f on p.fur_id = f.fur_id
									join shops s on p.shop_id = s.shop_id
									where f.fur_name = $1`, Fur_Name)
	if err != nil {
		return nil, fmt.Errorf("failed to get products by category furniture: %v", err)
	}
	defer rows.Close()

	var products []Products
	for rows.Next() {
		var prod Products
		if err := rows.Scan(&prod.Prod_ID, &prod.Prod_Name, &prod.Prod_Price, &prod.Prod_Details,
			&prod.Prod_Image, &prod.Stock, &prod.Sales_Amount, &prod.Status_Rec, &prod.Created_At,
			&prod.Updated_At, &prod.Brand_ID, &prod.Brand_Name, &prod.Room_ID, &prod.Room_Name,
			&prod.Fur_ID, &prod.Fur_Name, &prod.Shop_ID, &prod.Shop_Name, &prod.Shop_Des, &prod.Shop_Image,
			&prod.Shop_Address, &prod.Shop_PhoneNumber, &prod.Shop_Email); err != nil {
			return nil, fmt.Errorf("failed to scan product: %v", err)
		}
		products = append(products, prod)
	}
	return products, nil
}

func (pdb *PostgresDatabase) GetProductsByCategoryRoomAndFurniture(ctx context.Context, Room_Name string, Fur_Name string) ([]Products, error) {
	rows, err := pdb.db.QueryContext(ctx, `select p.prod_id, p.prod_name, p.prod_price, p.prod_details, p.prod_image, p.stock,
									p.sales_amount, p.status_rec, p.created_at, p.updated_at, p.brand_id, b.brand_name,
									p.room_id, r.room_name, p.fur_id, f.fur_name, p.shop_id, s.shop_name, s.shop_des,
									s.shop_image, s.shop_address, s.shop_phonenumber, s.shop_email from products p
									left join brand b on p.brand_id = b.brand_id
									left join room_type r on p.room_id = r.room_id
									left join fur_type f on p.fur_id = f.fur_id
									left join shops s on p.shop_id = s.shop_id
									WHERE r.room_name = $1 AND f.fur_name = $2`, Room_Name, Fur_Name)
	if err != nil {
		return nil, fmt.Errorf("failed to get products by category room and furniture: %v", err)
	}
	defer rows.Close()

	var products []Products
	for rows.Next() {
		var prod Products
		if err := rows.Scan(&prod.Prod_ID, &prod.Prod_Name, &prod.Prod_Price, &prod.Prod_Details,
			&prod.Prod_Image, &prod.Stock, &prod.Sales_Amount, &prod.Status_Rec, &prod.Created_At,
			&prod.Updated_At, &prod.Brand_ID, &prod.Brand_Name, &prod.Room_ID, &prod.Room_Name,
			&prod.Fur_ID, &prod.Fur_Name, &prod.Shop_ID, &prod.Shop_Name, &prod.Shop_Des, &prod.Shop_Image,
			&prod.Shop_Address, &prod.Shop_PhoneNumber, &prod.Shop_Email); err != nil {
			return nil, fmt.Errorf("failed to scan product: %v", err)
		}
		products = append(products, prod)
	}
	return products, nil
}

// Krit func
func (p *PostgresDatabase) CreateCustomer(ctx context.Context, customer *Customers) error {
	query := `
		INSERT INTO customers (Cust_Fname, Cust_Lname, Cust_Email, Cust_PhoneNumber, Cust_Username, Cust_Password)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING Cust_ID
	`

	// Execute the query and get the generated ID
	err := p.db.QueryRowContext(ctx, query,
		customer.FName, customer.LName, customer.CusEM, customer.CusPhone, customer.Username, customer.CusPW).Scan(&customer.CusID)

	if err != nil {
		return fmt.Errorf("failed to create customer: %v", err)
	}

	return nil
}

func (p *PostgresDatabase) GetCustomers(ctx context.Context) ([]Customers, error) {
	query := `
		SELECT cust_id, cust_fname, cust_lname, cust_email, cust_phonenumber, cust_username
		FROM customers
	`

	rows, err := p.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var customers []Customers
	for rows.Next() {
		var customer Customers
		err := rows.Scan(
			&customer.CusID, &customer.FName, &customer.LName, &customer.CusEM, &customer.CusPhone, &customer.Username,
		)
		if err != nil {
			return nil, err
		}
		customers = append(customers, customer)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return customers, nil
}

func (p *PostgresDatabase) GetCustomerByID(ctx context.Context, id int) (*Customers, error) {
	query := `
		SELECT cust_id, cust_fname, cust_lname, cust_email, cust_phonenumber, cust_username
		FROM customers
		WHERE cust_id = $1
	`

	var customer Customers
	err := p.db.QueryRowContext(ctx, query, id).Scan(
		&customer.CusID, &customer.FName, &customer.LName, &customer.CusEM, &customer.CusPhone, &customer.Username,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &customer, nil
}

func (p *PostgresDatabase) GetLatestProduct(ctx context.Context) ([]Products, error) {
	query := `
		SELECT prod_id, prod_name, prod_price, prod_details, prod_image, 
			   stock, sales_amount, status_rec, created_at, updated_at, 
			   room_id, fur_id
		FROM products 
		ORDER BY created_at DESC 
		LIMIT 3
	`

	rows, err := p.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []Products
	for rows.Next() {
		var product Products
		err := rows.Scan(
			&product.Prod_ID, &product.Prod_Name, &product.Prod_Price,
			&product.Prod_Details, &product.Prod_Image, &product.Stock,
			&product.Sales_Amount, &product.Status_Rec, &product.Created_At,
			&product.Updated_At, &product.Room_ID, &product.Fur_ID,
		)
		if err != nil {
			return nil, err
		}
		products = append(products, product)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return products, nil
}

func (p *PostgresDatabase) GetAllsProducts(ctx context.Context) ([]Products, error) {
	query := `
		SELECT prod_id, prod_name, prod_price, prod_details, prod_image, 
			   stock, sales_amount, status_rec, created_at, updated_at, 
			   room_id, fur_id
		FROM products
		ORDER BY created_at ASC
	`

	rows, err := p.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []Products
	for rows.Next() {
		var product Products
		err := rows.Scan(
			&product.Prod_ID, &product.Prod_Name, &product.Prod_Price,
			&product.Prod_Details, &product.Prod_Image, &product.Stock,
			&product.Sales_Amount, &product.Status_Rec, &product.Created_At,
			&product.Updated_At, &product.Room_ID, &product.Fur_ID,
		)
		if err != nil {
			return nil, err
		}
		products = append(products, product)
	}

	return products, nil
}

func (p *PostgresDatabase) GetShopInfo(ctx context.Context, id string) ([]Shops, error) {
	query := `
		SELECT shop_id, shop_name, shop_des, shop_image, shop_address, shop_phonenumber, shop_email
		FROM shops Where shop_id = $1
	`

	rows, err := p.db.QueryContext(ctx, query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var shops []Shops
	for rows.Next() {
		var shop Shops
		err := rows.Scan(
			&shop.Shop_ID, &shop.Shop_Name, &shop.Shop_Des, &shop.Shop_Image,
			&shop.Shop_Address, &shop.Shop_PhoneNumber, &shop.Shop_Email,
		)
		if err != nil {
			return nil, err
		}
		shops = append(shops, shop)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return shops, nil
}

func (p *PostgresDatabase) GetCustomerByUsernameAndPassword(ctx context.Context, username, password string) (*Customers, error) {
	query := `
		SELECT cust_id, cust_fname, cust_lname, cust_email, cust_phonenumber, cust_username
		FROM customers
		WHERE cust_username = $1 AND cust_password = $2
	`

	var customer Customers
	err := p.db.QueryRowContext(ctx, query, username, password).Scan(
		&customer.CusID, &customer.FName, &customer.LName, &customer.CusEM, &customer.CusPhone, &customer.Username,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &customer, nil
}

func (p *PostgresDatabase) GetProductsByShopID(ctx context.Context, shopID int) ([]Products, error) {
	query := `
        SELECT prod_id, prod_name, prod_price, prod_details, prod_image, 
               stock, sales_amount, status_rec, created_at, updated_at, 
               room_id, fur_id
        FROM products
        WHERE shop_id = $1
        ORDER BY created_at ASC
    `

	rows, err := p.db.QueryContext(ctx, query, shopID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []Products
	for rows.Next() {
		var product Products
		err := rows.Scan(
			&product.Prod_ID, &product.Prod_Name, &product.Prod_Price,
			&product.Prod_Details, &product.Prod_Image, &product.Stock,
			&product.Sales_Amount, &product.Status_Rec, &product.Created_At,
			&product.Updated_At, &product.Room_ID, &product.Fur_ID,
		)
		if err != nil {
			return nil, err
		}
		products = append(products, product)
	}

	return products, nil
}

func (p *PostgresDatabase) GetLatestProductsByShopID(ctx context.Context, shopID int) ([]Products, error) {
	query := `
        SELECT prod_id, prod_name, prod_price, prod_details, prod_image, 
               stock, sales_amount, status_rec, created_at, updated_at, 
               room_id, fur_id
        FROM products
        WHERE shop_id = $1
        ORDER BY created_at DESC
        LIMIT 3
    `

	rows, err := p.db.QueryContext(ctx, query, shopID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []Products
	for rows.Next() {
		var product Products
		err := rows.Scan(
			&product.Prod_ID, &product.Prod_Name, &product.Prod_Price,
			&product.Prod_Details, &product.Prod_Image, &product.Stock,
			&product.Sales_Amount, &product.Status_Rec, &product.Created_At,
			&product.Updated_At, &product.Room_ID, &product.Fur_ID,
		)
		if err != nil {
			return nil, err
		}
		products = append(products, product)
	}

	return products, nil
}

// Q func
func (db *PostgresDatabase) GetCart(ctx context.Context, Cust_ID string) (Cart, error) {
	// สร้างคำสั่ง SQL เพื่อดึงข้อมูลตะกร้าของลูกค้า
	query := `
        SELECT c.cart_id, c.cust_id, c.totalprice, c.created_at, c.updated_at
        FROM cart c
        WHERE c.cust_id = $1
    `
	var cart Cart
	err := db.db.QueryRowContext(ctx, query, Cust_ID).Scan(&cart.Cart_ID, &cart.Cust_ID, &cart.TotalPrice, &cart.Created_At, &cart.Updated_At)
	if err != nil {
		if err == sql.ErrNoRows {
			return Cart{}, fmt.Errorf("no cart found for customer ID: %s", Cust_ID)
		}
		return Cart{}, err
	}

	// ดึงรายการสินค้าจากตะกร้า
	itemQuery := `
        SELECT ci.cart_item_id, ci.prod_id, p.prod_image, p.prod_name, ci.brand_name, ci.prod_price, ci.quantity, ci.totalprice
        FROM cart_items ci
        JOIN products p ON ci.prod_id = p.prod_id
        JOIN brand b ON ci.brand_name = b.brand_name
        WHERE ci.cart_id = $1
    `
	rows, err := db.db.QueryContext(ctx, itemQuery, cart.Cart_ID)
	if err != nil {
		return Cart{}, fmt.Errorf("error fetching cart items: %v", err)
	}
	defer rows.Close()

	// เพิ่มรายการสินค้าในตะกร้า
	for rows.Next() {
		var item CartItem
		if err := rows.Scan(
			&item.Cart_Item_ID,
			&item.Prod_ID,
			&item.Prod_Image,
			&item.Prod_Name,  // ชื่อลูกค้าสินค้า
			&item.Brand_Name, // ชื่อแบรนด์
			&item.Prod_Price,
			&item.Quantity,
			&item.TotalPrice); err != nil {
			return Cart{}, fmt.Errorf("error scanning cart item: %v", err)
		}
		cart.Cart_Items = append(cart.Cart_Items, item)
	}

	if err := rows.Err(); err != nil {
		return Cart{}, fmt.Errorf("error reading rows: %v", err)
	}

	// คืนค่าตะกร้าและรายการสินค้า
	return cart, nil
}

func (s *Store) GetCart(ctx context.Context, Cust_ID string) (Cart, error) {
	// เรียกฟังก์ชัน GetProducts จาก PostgresDatabase
	return s.db.GetCart(ctx, Cust_ID)
}

func (db *PostgresDatabase) createNewCart(ctx context.Context, Cust_ID string) (string, error) {
	query := `
		INSERT INTO cart (cust_id, totalprice, created_at, updated_at)
		VALUES ($1, $2, NOW(), NOW()) RETURNING cart_id
	`
	var cartID string
	err := db.db.QueryRowContext(ctx, query, Cust_ID, 0).Scan(&cartID)
	if err != nil {
		return "", err
	}
	return cartID, nil
}

func (db *PostgresDatabase) AddToCart(ctx context.Context, Cust_ID string, Prod_ID string, Quantity int) error {
	// ขั้นตอนที่ 1: ตรวจสอบว่ามีตะกร้าของลูกค้าหรือไม่
	query := `
		SELECT cart_id FROM cart WHERE cust_id = $1
	`
	var cartID string
	err := db.db.QueryRowContext(ctx, query, Cust_ID).Scan(&cartID)
	if err != nil {
		if err == sql.ErrNoRows {
			// ถ้าไม่มีตะกร้า (กรณีนี้จะสร้างตะกร้าหรือเพิ่มตะกร้าใหม่)
			cartID, err = db.createNewCart(ctx, Cust_ID)
			if err != nil {
				return fmt.Errorf("failed to create new cart for customer %s: %v", Cust_ID, err)
			}
		} else {
			return fmt.Errorf("failed to retrieve cart for customer %s: %v", Cust_ID, err)
		}
	}

	// ขั้นตอนที่ 2: ตรวจสอบว่าสินค้าที่จะเพิ่มอยู่ในตะกร้าแล้วหรือยัง
	itemQuery := `
		SELECT cart_item_id, prod_price, quantity FROM cart_items WHERE cart_id = $1 AND prod_id = $2
	`
	var cartItemID int
	var prodPrice float64
	var currentQuantity int
	err = db.db.QueryRowContext(ctx, itemQuery, cartID, Prod_ID).Scan(&cartItemID, &prodPrice, &currentQuantity)

	// ถ้ามีสินค้าตัวนี้ในตะกร้าแล้วให้เพิ่มจำนวนสินค้า
	if err == nil {
		updateQuery := `
			UPDATE cart_items 
			SET quantity = quantity + $1, totalprice = prod_price * (quantity + $1)
			WHERE cart_item_id = $2
		`
		_, err := db.db.ExecContext(ctx, updateQuery, Quantity, cartItemID)
		if err != nil {
			return fmt.Errorf("failed to update quantity for item %s in cart: %v", Prod_ID, err)
		}
		return nil
	}

	// ถ้ายังไม่มีสินค้าตัวนี้ในตะกร้าให้เพิ่มรายการใหม่
	if err == sql.ErrNoRows {
		// ดึงข้อมูลของสินค้า (รวมถึง brand_name)
		productQuery := `
			SELECT prod_image, prod_name, prod_price, b.brand_name 
			FROM products p
			JOIN brand b ON p.brand_id = b.brand_id
			WHERE p.prod_id = $1
		`
		var prodName, brandName, prodImage string
		err = db.db.QueryRowContext(ctx, productQuery, Prod_ID).Scan(&prodImage, &prodName, &prodPrice, &brandName)
		if err != nil {
			return fmt.Errorf("failed to retrieve product details for product %s: %v", Prod_ID, err)
		}

		// คำนวณ totalprice ใน Go ก่อน
		totalPrice := prodPrice * float64(Quantity)

		// เพิ่มสินค้าลงในตะกร้า
		insertItemQuery := `
			INSERT INTO cart_items (cart_id, prod_id, prod_image, prod_name, prod_price, quantity, brand_name, totalprice)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		`
		_, err = db.db.ExecContext(ctx, insertItemQuery, cartID, Prod_ID, prodImage, prodName, prodPrice, Quantity, brandName, totalPrice)
		if err != nil {
			return fmt.Errorf("failed to insert new item %s into cart: %v", Prod_ID, err)
		}
		return nil
	}

	// กรณีเกิดข้อผิดพลาด
	return fmt.Errorf("unexpected error: %v", err)
}

func (bs *Store) AddToCart(ctx context.Context, Cust_ID string, Prod_ID string, Quantity int) error {
	// เรียกใช้ฟังก์ชัน AddToCart จาก PostgresDatabase และส่งพารามิเตอร์ไปด้วย
	return bs.db.AddToCart(ctx, Cust_ID, Prod_ID, Quantity)
}

func (db *PostgresDatabase) DeleteFromCart(ctx context.Context, Cust_ID string, Cart_Item_ID int, Quantity int) error {
	// ตรวจสอบว่ามีตะกร้าของลูกค้าหรือไม่
	query := `SELECT cart_id FROM cart WHERE cust_id = $1`
	var cartID string
	err := db.db.QueryRowContext(ctx, query, Cust_ID).Scan(&cartID)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("Cart not found for customer ID %s", Cust_ID)
		}
		return err
	}

	// ตรวจสอบสินค้าตาม cart_item_id
	itemQuery := `SELECT cart_item_id, quantity, prod_price FROM cart_items WHERE cart_id = $1 AND cart_item_id = $2`
	var existingItemID, existingQuantity int
	var prodPrice float64
	err = db.db.QueryRowContext(ctx, itemQuery, cartID, Cart_Item_ID).Scan(&existingItemID, &existingQuantity, &prodPrice)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("Item not found in the cart for item ID %d", Cart_Item_ID)
		}
		return err
	}

	// ถ้าจำนวนที่ระบุมากกว่าหรือเท่ากับจำนวนที่มีในตะกร้า, ให้ลบสินค้าทั้งหมด
	if Quantity >= existingQuantity {
		deleteQuery := `DELETE FROM cart_items WHERE cart_id = $1 AND cart_item_id = $2`
		_, err = db.db.ExecContext(ctx, deleteQuery, cartID, Cart_Item_ID)
		if err != nil {
			return fmt.Errorf("Failed to delete item from cart: %v", err)
		}
	} else {
		// ลดจำนวนสินค้าลง
		updateQuery := `
		UPDATE cart_items 
		SET quantity = quantity - $1
		WHERE cart_id = $2 
		  AND cart_item_id = $3
		`
		_, err = db.db.ExecContext(ctx, updateQuery, Quantity, cartID, Cart_Item_ID)
		if err != nil {
			return fmt.Errorf("Failed to update quantity in cart: %v", err)
		}
	}

	// คำนวณ total price ใหม่ของตะกร้า
	err = db.updateCartTotalPrice(ctx, cartID)
	if err != nil {
		return fmt.Errorf("Failed to update total price after deleting item: %v", err)
	}

	return nil
}

func (db *PostgresDatabase) updateCartTotalPrice(ctx context.Context, cartID string) error {
	// คำนวณ total price ใหม่หลังจากลบสินค้า
	totalPriceQuery := `SELECT SUM(totalprice) FROM cart_items WHERE cart_id = $1`
	var totalPrice float64
	err := db.db.QueryRowContext(ctx, totalPriceQuery, cartID).Scan(&totalPrice)
	if err != nil {
		return fmt.Errorf("Failed to calculate total price for cart: %v", err)
	}

	// อัพเดต total price ในตะกร้า
	updateQuery := `UPDATE cart SET totalprice = $1 WHERE cart_id = $2`
	_, err = db.db.ExecContext(ctx, updateQuery, totalPrice, cartID)
	if err != nil {
		return fmt.Errorf("Failed to update cart total price: %v", err)
	}

	return nil
}

func (s *Store) DeleteFromCart(ctx context.Context, Cust_ID string, Cart_Item_ID int, Quantity int) error {
	return s.db.DeleteFromCart(ctx, Cust_ID, Cart_Item_ID, Quantity)
}

// Other func
func (pdb *PostgresDatabase) Close() error {
	return pdb.db.Close()
}

func (pdb *PostgresDatabase) Ping() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return pdb.db.PingContext(ctx)
}

func (pdb *PostgresDatabase) Reconnect(connStr string) error {
	if pdb.db != nil {
		pdb.db.Close()
	}

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %v", err)
	}

	// ตั้งค่า connection pool
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(5 * time.Minute)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return fmt.Errorf("failed to ping database: %v", err)
	}

	pdb.db = db
	return nil
}

type Store struct {
	db EcommerceDatabase
}

func NewStore(db EcommerceDatabase) *Store {
	return &Store{db: db}
}

//	func (s *Store) GetBrand(ctx context.Context) ([]Brand, error) {
//		return s.db.GetBrand(ctx)
//	}

func (s *Store) GetProducts(ctx context.Context, params ProductQueryParams) (*ProductResponse, error) {
	return s.db.GetProducts(ctx, params)
}

func (s *Store) GetNewProdShop(ctx context.Context) ([]Products, error) {
	return s.db.GetNewProdShop(ctx)
}

func (s *Store) GetRecProducts(ctx context.Context) ([]Products, error) {
	return s.db.GetRecProducts(ctx)
}

// Mook
func (s *Store) GetProduct(ctx context.Context, Prod_ID string) (*Products, error) {
	return s.db.GetProduct(ctx, Prod_ID)
}

func (s *Store) Close() error {
	return s.db.Close()
}

func (s *Store) Ping() error {
	if s.db == nil {
		return fmt.Errorf("database connection is not initialized")
	}
	return s.db.Ping()
}

func (s *Store) Reconnect(connStr string) error {
	return s.db.Reconnect(connStr)
}

//Man

func (bs *Store) GetAllProducts(ctx context.Context) ([]Products, error) {
	return bs.db.GetAllProducts(ctx)
}

func (bs *Store) GetProductsByCategoryRoom(ctx context.Context, Room_Name string) ([]Products, error) {
	return bs.db.GetProductsByCategoryRoom(ctx, Room_Name)
}

func (bs *Store) GetProductsByCategoryFurniture(ctx context.Context, Fur_Name string) ([]Products, error) {
	return bs.db.GetProductsByCategoryFurniture(ctx, Fur_Name)
}

func (bs *Store) GetProductsByCategoryRoomAndFurniture(ctx context.Context, Room_Name string, Fur_Name string) ([]Products, error) {
	return bs.db.GetProductsByCategoryRoomAndFurniture(ctx, Room_Name, Fur_Name)
}

// Krit
func (s *Store) CreateCustomer(ctx context.Context, customer *Customers) error {
	return s.db.CreateCustomer(ctx, customer)
}

func (s *Store) GetCustomers(ctx context.Context) ([]Customers, error) {
	return s.db.GetCustomers(ctx)
}

func (s *Store) GetCustomerByID(ctx context.Context, id int) (*Customers, error) {
	return s.db.GetCustomerByID(ctx, id)
}

func (s *Store) GetLatestProduct(ctx context.Context) ([]Products, error) {
	return s.db.GetLatestProduct(ctx)
}

func (s *Store) GetAllsProducts(ctx context.Context) ([]Products, error) {
	return s.db.GetAllProducts(ctx)
}

func (s *Store) GetShopInfo(ctx context.Context, id string) ([]Shops, error) {
	return s.db.GetShopInfo(ctx, id)
}

func (s *Store) GetCustomerByUsernameAndPassword(ctx context.Context, username, password string) (*Customers, error) {
	return s.db.GetCustomerByUsernameAndPassword(ctx, username, password)
}

func (s *Store) GetProductsByShopID(ctx context.Context, shopID int) ([]Products, error) {
	return s.db.GetProductsByShopID(ctx, shopID)
}

func (s *Store) GetLatestProductsByShopID(ctx context.Context, shopID int) ([]Products, error) {
	return s.db.GetLatestProductsByShopID(ctx, shopID)
}
