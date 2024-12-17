package routes

import (
	"api-gateway/internal/clients"
	"api-gateway/internal/handlers"
	"api-gateway/internal/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

type Routes struct {
	router chi.Router
}

func NewRouter(r chi.Router) *Routes {
	r.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))
	return &Routes{router: r}
}

func (r *Routes) SetupAuthRoutes(authHandlers *handlers.AuthHandler) {
	r.router.Post("/register", authHandlers.Register)
	r.router.Post("/login", authHandlers.Login)
}

func (r *Routes) SetupInventoryRoutes(authClient *clients.AuthClient, inventoryHandlers *handlers.InventoryHandler) {
	r.router.Post("/inventory", middleware.AuthMiddleware(authClient)(inventoryHandlers.AddItem))
	r.router.Put("/inventory", middleware.AuthMiddleware(authClient)(inventoryHandlers.UpdateItem))
	r.router.Get("/inventory/{userID}/{itemID}", middleware.AuthMiddleware(authClient)(inventoryHandlers.GetItem))
	r.router.Get("/inventory/{userID}", middleware.AuthMiddleware(authClient)(inventoryHandlers.GetInventory))
	r.router.Delete("/inventory/{userID}/{itemID}", middleware.AuthMiddleware(authClient)(inventoryHandlers.DeleteItem))
}

func (r *Routes) SetupStockRoutes(stockHandlers *handlers.StockHandlers) {
	r.router.Get("/ws", stockHandlers.NotificationHandler)
}
