package main

import (
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv" // Adicione esta linha
	"github.com/tiago-prog/novels-api/internal/controller"
	"github.com/tiago-prog/novels-api/internal/db"
	"github.com/tiago-prog/novels-api/internal/repository"
	"github.com/tiago-prog/novels-api/internal/usecase"
)

func main() {
	// Carregar variáveis do .env
	err := godotenv.Load()
	if err != nil {
		log.Println("Arquivo .env não encontrado ou não pôde ser carregado")
	}

	// Configuração do banco
	pgHost := os.Getenv("PG_HOST")
	if pgHost == "" {
		pgHost = "localhost"
	}
	// ... (similar para outras variáveis)

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET não definido")
	}

	pgUser := os.Getenv("PG_USER")
	if pgUser == "" {
		pgUser = "postgres"
	}
	pgPassword := os.Getenv("PG_PASSWORD")
	if pgPassword == "" {
		pgPassword = "1234"
	}
	pgDatabase := os.Getenv("PG_DATABASE")
	if pgDatabase == "" {
		pgDatabase = "postgres"
	}

	database, err := db.Connect(pgHost, "5432", pgUser, pgPassword, pgDatabase)
	if err != nil {
		log.Fatalf("Erro ao conectar ao banco: %v", err)
	}
	defer database.Close()

	userRepo := repository.NewUserRepository(database)
	if err := userRepo.CreateAdminIfNotExists(); err != nil {
		log.Printf("Erro ao criar admin: %v", err)
	} else {
		log.Println("Admin criado ou já existente")
	}

	router := gin.Default()

	// Configurar CORS
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"}, // Altere para os domínios permitidos em produção!
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	//Camada de repository
	UserRepository := repository.NewUserRepository(database)
	//Camada usecase
	UserUseCase := usecase.NewUserUsecase(UserRepository)
	//Camada de controllers
	UserController := controller.NewUserController(UserUseCase)

	// Rotas públicas
	public := router.Group("")
	{
		public.GET("/health", func(c *gin.Context) { c.JSON(200, gin.H{"status": "ok"}) })
		public.POST("/register", UserController.Register)
		public.POST("/login", UserController.Login)
	}

	// Rotas protegidas
	protected := router.Group("")
	protected.Use(controller.Auth(jwtSecret))
	{
		protected.GET("/user/:email", UserController.GetUserByEmail)
		protected.POST("/suspend/:executor_id/:target_id", UserController.SuspendUser)
		protected.GET("/users", UserController.GetAllUsers)
	}

	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Erro ao iniciar servidor: %v", err)
	}
}
