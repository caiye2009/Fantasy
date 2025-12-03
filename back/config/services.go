package config

import (
	"gorm.io/gorm"
	"github.com/redis/go-redis/v9"
	"github.com/elastic/go-elasticsearch/v8"

	"back/pkg/auth"
	"back/pkg/es"
	"back/internal/vendor"
	"back/internal/client"
	"back/internal/user"
	internalAuth "back/internal/auth"
	"back/internal/material"
	materialPrice "back/internal/material/price"
	"back/internal/process"
	processPrice "back/internal/process/price"
	"back/internal/product"
	"back/internal/search"
)

type Services struct {
	Vendor        *vendor.VendorService
	Client        *client.ClientService
	User          *user.UserService
	Auth          *internalAuth.AuthService
	Material      *material.MaterialService
	MaterialPrice *materialPrice.MaterialPriceService
	Process       *process.ProcessService
	ProcessPrice  *processPrice.ProcessPriceService
	Product       *product.ProductService
	Search        *search.SearchService
}

func InitServices(db *gorm.DB, rdb *redis.Client, esClient *elasticsearch.Client, jwtWang *auth.JWTWang, esSync *es.ESSync) *Services {
	vendorRepo := vendor.NewVendorRepo(db)
	vendorService := vendor.NewVendorService(vendorRepo, esSync)

	clientRepo := client.NewClientRepo(db)
	clientService := client.NewClientService(clientRepo, esSync)

	userRepo := user.NewUserRepo(db)
	userService := user.NewUserService(userRepo)

	authService := internalAuth.NewAuthService(userService, jwtWang)

	materialRepo := material.NewMaterialRepo(db)
	materialService := material.NewMaterialService(materialRepo, esSync)

	materialPriceRepo := materialPrice.NewMaterialPriceRepo(db)
	materialPriceService := materialPrice.NewMaterialPriceService(materialPriceRepo, vendorService, rdb)

	processRepo := process.NewProcessRepo(db)
	processService := process.NewProcessService(processRepo, esSync)

	processPriceRepo := processPrice.NewProcessPriceRepo(db)
	processPriceService := processPrice.NewProcessPriceService(processPriceRepo, vendorService, rdb)

	productRepo := product.NewProductRepo(db)
	productService := product.NewProductService(
		productRepo,
		materialPriceService,
		processPriceService,
		materialService,
		processService,
		esSync,
	)

	searchService := search.NewSearchService(esClient)

	return &Services{
		Vendor:        vendorService,
		Client:        clientService,
		User:          userService,
		Auth:          authService,
		Material:      materialService,
		MaterialPrice: materialPriceService,
		Process:       processService,
		ProcessPrice:  processPriceService,
		Product:       productService,
		Search:        searchService,
	}
}