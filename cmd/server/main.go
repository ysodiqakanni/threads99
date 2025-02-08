package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/ysodiqakanni/threads99/internal/comment"
	"github.com/ysodiqakanni/threads99/internal/community"
	"github.com/ysodiqakanni/threads99/internal/config"
	"github.com/ysodiqakanni/threads99/internal/post"
	"github.com/ysodiqakanni/threads99/internal/user"
	"github.com/ysodiqakanni/threads99/pkg/dbcontext"
	"github.com/ysodiqakanni/threads99/pkg/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"net/http"
	"net/url"
	"os"
	"time"
)

// Version indicates the current version of the application.
var Version = "1.0.0"

var flagConfig = flag.String("config", "./config/dev.yml", "path to the config file")

func main() {
	fmt.Println("threading..")
	flag.Parse()
	// create root logger tagged with server version
	logger := log.New().With(nil, "version", Version)

	// load application configurations
	cfg, err := config.Load(*flagConfig, logger)
	if err != nil {
		logger.Errorf("failed to load application configuration: %s", err)
		os.Exit(-1)
	}
	// connect to the mongo database
	escapedPassword := url.QueryEscape(cfg.DbPassword)
	connStr := fmt.Sprintf(cfg.DbConnectionString, escapedPassword)
	db, err := SetupMongoDB(connStr, cfg.DbName)

	if err != nil {
		logger.Error(err)
		os.Exit(-1)
	}

	defer db.Client().Disconnect(context.Background())

	// build HTTP server
	address := fmt.Sprintf(":%v", cfg.ServerPort)
	hs := &http.Server{
		Addr:    address,
		Handler: buildHandler(logger, dbcontext.New(db), cfg),
	}

	// start the HTTP server with graceful shutdown
	// go routing.GracefulShutdown(hs, 10*time.Second, logger.Infof)
	logger.Infof("server %v is running at %v", Version, address)
	if err := hs.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Error(err)
		os.Exit(-1)
	}
}

func SetupMongoDB(connStr, dbName string) (*mongo.Database, error) {
	fmt.Println("attempting to connect to the db.." + dbName + connStr)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(connStr))
	if err != nil {
		return nil, err
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	fmt.Println("Ping successful!")
	// mongoClient = client

	/*
		// Now create collections and set rules
		//collection := client.Database(dbName).Collection("communities")
		// Define validation rule
		communitiesValidation := bson.M{
			"$jsonSchema": bson.M{
				"bsonType": "object",
				"required": []string{"created_by_user_id"},
				"properties": bson.M{
					"email": bson.M{
						"bsonType":    "string",
						"description": "must be a ObjectId and is required",
					},
					// Add other properties as needed
				},
			},
		}

		// Create collection with validation
		opts := options.CreateCollection().SetValidator(communitiesValidation)
		if err := client.Database(dbName).CreateCollection(context.Background(), "communities", opts); err != nil {
			fmt.Println("Validation rule setup failed!")
			return nil, err
		}*/

	return client.Database(dbName), nil
}

func buildHandler(logger log.Logger, db *dbcontext.DB, cfg *config.Config) http.Handler {
	r := mux.NewRouter()
	r.HandleFunc("/api/healthcheck", HealthCheckHandler).Methods("GET")

	post.RegisterHandlers(r,
		post.NewService(post.NewRepository(db, logger), community.NewRepository(db, logger), logger),
		logger,
		cfg.JWTSigningKey)

	comment.RegisterHandlers(r,
		comment.NewService(comment.NewRepository(db, logger), post.NewRepository(db, logger), logger),
		logger,
		cfg.JWTSigningKey)

	community.RegisterHandlers(r,
		community.NewService(community.NewRepository(db, logger), logger),
		logger,
		cfg.JWTSigningKey)

	user.RegisterHandlers(r,
		user.NewService(user.NewRepository(db, logger), logger, cfg.JWTSigningKey, cfg.JWTExpiration),
		logger)

	//auth.RegisterHandlers(r,
	//	auth.NewService(cfg.JWTSigningKey, cfg.JWTExpiration, logger, user.NewRepository(db, logger)),
	//	logger)

	return r
}
func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Everything is dope from this side :)")
}
