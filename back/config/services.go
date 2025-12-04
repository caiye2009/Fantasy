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
	// 直接传 db，不再创建 repo
	vendorService := vendor.NewVendorService(db, esSync)
	clientService := client.NewClientService(db, esSync)
	userService := user.NewUserService(db)
	authService := internalAuth.NewAuthService(userService, jwtWang)
	materialService := material.NewMaterialService(db, esSync)
	processService := process.NewProcessService(db, esSync)

	// Price 服务也只传 db
	materialPriceService := materialPrice.NewMaterialPriceService(db, vendorService, rdb)
	processPriceService := processPrice.NewProcessPriceService(db, vendorService, rdb)

	// product 服务传 db
	productService := product.NewProductService(
		db,
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