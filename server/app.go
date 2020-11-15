package server

import (
	"context"
	"github.com/abylq/folder/folder"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/abylq/folder/auth"
	"github.com/abylq/folder/bookmark"

	authhttp "github.com/abylq/folder/auth/delivery/http"
	authmongo "github.com/abylq/folder/auth/repository/mongo"
	authusecase "github.com/abylq/folder/auth/usecase"
	bmhttp "github.com/abylq/folder/bookmark/delivery/http"
	bmmongo "github.com/abylq/folder/bookmark/repository/mongo"
	bmusecase "github.com/abylq/folder/bookmark/usecase"
	fmhttp "github.com/abylq/folder/folder/delivery/http"
	fmmongo "github.com/abylq/folder/folder/repository/mongo"
	fmusecase "github.com/abylq/folder/folder/usecase"
)

type App struct {
	httpServer *http.Server
	bookmarkUC bookmark.UseCase
	authUC     auth.UseCase
	folderUC   folder.UseCase

}

func NewApp() *App {
	db := initDB()

	userRepo := authmongo.NewUserRepository(db, viper.GetString("mongo.user_collection"))
	bookmarkRepo := bmmongo.NewBookmarkRepository(db, viper.GetString("mongo.bookmark_collection"))
	folderRepo := fmmongo.NewFolderRepository(db,viper.GetString("mongo.folder_collection"))
	return &App{
		bookmarkUC: bmusecase.NewBookmarkUseCase(bookmarkRepo),
		authUC: authusecase.NewAuthUseCase(
			userRepo,
			viper.GetString("auth.hash_salt"),
			[]byte(viper.GetString("auth.signing_key")),
			viper.GetDuration("auth.token_ttl"),
		),
		folderUC: fmusecase.NewFolderUseCase(folderRepo),
	}
}

func (a *App) Run(port string) error {
	// Init gin handler
	router := gin.Default()
	router.Use(
		gin.Recovery(),
		gin.Logger(),
	)

	// Set up http handlers
	// SignUp/SignIn endpoints
	authhttp.RegisterHTTPEndpoints(router, a.authUC)

	// API endpoints
	authMiddleware := authhttp.NewAuthMiddleware(a.authUC)
	api := router.Group("/api", authMiddleware)

	bmhttp.RegisterHTTPEndpoints(api, a.bookmarkUC)
	fmhttp.RegisterHTTPEndpoints(api, a.folderUC)

	// HTTP Server
	a.httpServer = &http.Server{
		Addr:           ":" + port,
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	go func() {
		if err := a.httpServer.ListenAndServe(); err != nil {
			log.Fatalf("Failed to listen and serve: %+v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, os.Interrupt)

	<-quit

	ctx, shutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdown()

	return a.httpServer.Shutdown(ctx)
}

func initDB() *mongo.Database {
	client, err := mongo.NewClient(options.Client().ApplyURI(viper.GetString("mongo.uri")))
	if err != nil {
		log.Fatalf("Error occured while establishing connection to mongoDB")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}

	return client.Database(viper.GetString("mongo.name"))
}
