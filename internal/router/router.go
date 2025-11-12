package router

import (
	"net/http"

	"github.com/emmrys-jay/gigmile/internal/handler"
	"github.com/emmrys-jay/gigmile/internal/middleware"
	"github.com/emmrys-jay/gigmile/internal/service"
	"github.com/gorilla/mux"
)

func NewRouter(
	customerService service.CustomerService,
	paymentService service.PaymentService,
	deploymentService service.DeploymentService,
	transactionService service.TransactionService,
	accountService service.AccountService,
) *mux.Router {
	router := mux.NewRouter()

	customerHandler := handler.NewCustomerHandler(customerService)
	paymentHandler := handler.NewPaymentHandler(paymentService)
	deploymentHandler := handler.NewDeploymentHandler(deploymentService)
	transactionHandler := handler.NewTransactionHandler(transactionService)
	accountHandler := handler.NewAccountHandler(accountService)

	// Apply logging middleware
	router.Use(middleware.LoggingMiddleware)

	// Customer routes
	api := router.PathPrefix("/api/v1").Subrouter()
	api.HandleFunc("/customers", customerHandler.CreateCustomer).Methods("POST")
	api.HandleFunc("/customers", customerHandler.GetAllCustomers).Methods("GET")
	api.HandleFunc("/customers/{id}", customerHandler.GetCustomerByID).Methods("GET")
	api.HandleFunc("/customers/{id}", customerHandler.UpdateCustomer).Methods("PUT")
	api.HandleFunc("/customers/{id}", customerHandler.DeleteCustomer).Methods("DELETE")

	// Payment routes
	api.HandleFunc("/payments/notify", paymentHandler.ProcessPaymentNotification).Methods("POST")

	// Deployment routes
	api.HandleFunc("/deployments", deploymentHandler.RecordDeployment).Methods("POST")

	// Transaction routes
	api.HandleFunc("/customers/{id}/transactions", transactionHandler.GetTransactionsByCustomer).Methods("GET")

	// Account routes
	api.HandleFunc("/customers/{id}/account", accountHandler.GetAccountByCustomer).Methods("GET")

	// Health check
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}).Methods("GET")

	return router
}
