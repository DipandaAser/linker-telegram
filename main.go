package main

import (
	"context"
	"github.com/DipandaAser/linker"
	"github.com/DipandaAser/linker-telegram/app"
	"github.com/DipandaAser/linker-telegram/bot"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"net/http"
	"os"
	"time"
)

var db *mongo.Database
var ctx context.Context = context.TODO()

func main() {

	app.Init()
	linker.MongoCtx = &ctx

	// ─── MONGO ──────────────────────────────────────────────────────────────────────
	err := MongoConnect()
	if err != nil {
		log.Fatal("Can't setup mongodb")
	}

	// ─── WE REFRESH THE MONGO CONNECTION EACH 10MINS ──────────────────────────────────────
	ticker := time.NewTicker(time.Minute * 10)
	defer ticker.Stop()
	go func() {
		for range ticker.C {
			go MongoReconnectCheck()
		}
	}()

	err = bot.InitBot()
	if err != nil {
		log.Fatalf("Can't init bot: %s", err)
	}

	// We prepare the reception of incoming call by the
	bot.GetNoBotRouter().POST("/linker", func(c *gin.Context) {
		c.String(http.StatusOK, app.Config.ServiceName)
		return
	})

	serviceUrl := bot.GetNoBotEndpointPath() + "/linker"
	// we put te service online to allow other service to link
	_, err = linker.SetService(app.Config.ServiceName, serviceUrl, app.Config.AuthKey, linker.StatusOnline)
	if err != nil {
		log.Fatal("can't set service")
	}

	// We put the service offline if the program stop
	defer linker.SetService(app.Config.ServiceName, serviceUrl, app.Config.AuthKey, linker.StatusOffline)

	bot.Start()
}

// MongoConnect connects to mongoDB
func MongoConnect() error {

	clientOptions := options.Client().ApplyURI(os.Getenv("MONGO_URI"))

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return err
	}

	// We make sure we have been connected
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return err
	}

	db = client.Database(os.Getenv("DB_NAME"))
	linker.DB = db

	return nil
}

// MongoReconnectCheck reconnects to MongoDB
func MongoReconnectCheck() {

	// We make sure we are still connected
	err := db.Client().Ping(ctx, readpref.Primary())
	if err != nil {
		// We reconnect
		_ = MongoConnect()
	}
}
