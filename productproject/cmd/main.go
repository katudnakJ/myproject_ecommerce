// main.go

package main

import (
	"context"
	"log"
	"productproject/internal/config"
	"productproject/internal/handlers"

	product "productproject/internal/product"

	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func TimeoutMiddleware(timeout time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), timeout)
		defer cancel()

		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	db, err := product.NewPostgresDatabase(cfg.GetConnectionString())
	if err != nil {
		log.Printf("Failed to connect to database: %v", err)
	}
	if db != nil {
		defer db.Close()
	}

	store := product.NewStore(db)
	h := handlers.NewProductHandlers(store)
	customerHandlers := handlers.NewCustomerHandlers(store)

	go func() {
		for {
			time.Sleep(10 * time.Second)
			if err := db.Ping(); err != nil {
				log.Printf("Database connection lost: %v", err)
				// พยายามเชื่อมต่อใหม่
				if reconnErr := db.Reconnect(cfg.GetConnectionString()); reconnErr != nil {
					log.Printf("Failed to reconnect: %v", reconnErr)
				} else {
					log.Printf("Successfully reconnected to the database")
				}
			}
		}
	}()

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	// กำหนดค่า CORS
	configCors := cors.Config{
		AllowOrigins:     []string{"http://localhost:8080", "http://localhost:3000", "http://localhost:4000", "http://localhost:8085"}, // เพิ่ม URL ของแอป React
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}
	r.Use(cors.New(configCors))

	r.Use(TimeoutMiddleware(5 * time.Second))

	r.GET("/health", h.HealthCheck)

	// API v1
	v1 := r.Group("/api/v1")
	{
		products := v1.Group("/products")
		{
			products.GET("", h.GetProducts)
			// products.GET("/brand", h.GetBrand)
			products.GET("/new", h.GetNewProdShop)
			products.GET("/rec", h.GetRecProducts)

			//Mook
			products.GET("/:Prod_ID", h.GetProduct)

			//man
			products.GET("GetAllProducts", h.GetAllProducts)
			products.GET("GetProductByCategoryRoom/:room_name", h.GetProductsByCategoryRoom)
			products.GET("GetProductByCategoryFurniture/:fur_name", h.GetProductsByCategoryFurniture)
			products.GET("GetProductByCategoryRoom&Fur/:room_name/:fur_name", h.GetProductsByCategoryRoomAndFurniture)

			//krit
			products.POST("/register", h.RegisterCustomer)
			products.GET("/latest", h.GetLatestProduct)
			products.GET("/all", h.GetAllProducts)
			products.GET("/shop/:shop_id", h.GetProductsByShopID)              // New route
			products.GET("/shop/:shop_id/latest", h.GetLatestProductsByShopID) // New route
		}
		shops := v1.Group("/shops")
		{
			shops.GET("/info/:shop_id", h.GetShopInfo)
		}
		customers := v1.Group("/customers")
		{
			customers.GET("/", customerHandlers.GetCustomers)
			customers.GET("/:id", customerHandlers.GetCustomerByID)
		}
		auth := v1.Group("/auth")
		{
			auth.POST("/login", h.Login)
		}

		//Q
		Cart := v1.Group("/Cart")
		{
			Cart.GET("/:cust_id", h.GetCartHandler)
			Cart.POST("/:cust_id/:prod_id", h.AddToCartHandler)
		}
		CartDELETE := v1.Group("/CartDELETE")
		{
			CartDELETE.DELETE("/:cust_id/:cart_item_id", h.DeleteFromCartHandler)
		}
	}

	if err := r.Run(":" + cfg.AppPort); err != nil {
		log.Printf("Failed to run server: %v", err)
	}
}
