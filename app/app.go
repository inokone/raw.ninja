package app

import (
	"fmt"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/inokone/photostorage/common"
	"github.com/inokone/photostorage/docs"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	web "github.com/inokone/photostorage/web"
	"gorm.io/gorm"

	"github.com/inokone/photostorage/image"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
)

var config *common.AppConfig
var DB *gorm.DB
var IS *image.Store

func init() {
	conf, err := common.LoadConfig()
	if err != nil {
		log.Error().Err(err).Msg("Failed to load application configuration.")
		os.Exit(1)
	}
	config = conf
	initLog()
	log.Info().Msg("Photostorage app starting up...")
}

func initLog() {
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	zerolog.TimeFieldFormat = time.RFC3339
	level, err := zerolog.ParseLevel(config.Log.LogLevel)
	if err != nil {
		log.Warn().Err(err).Msg("Failed to parse log level, default is debug.")
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else {
		zerolog.SetGlobalLevel(level)
	}
	if config.Log.PrettyLog {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}
}

func App(port int) {
	var err error
	if err = initDb(config.Database); err != nil {
		log.Error().Msg("Failed to set up connection to database. Application spinning down.")
		os.Exit(1)
	}
	if err = initStore(config.Store); err != nil {
		log.Error().Msg("Failed to set up image store. Application spinning down.")
		os.Exit(1)
	}

	r := gin.New()

	// Setup middleware
	r.Use(gin.Recovery())
	r.Use(web.LoggingMiddleware)
	r.MaxMultipartMemory = 8 << 20

	// Set up Swagger
	docs.SwaggerInfo.BasePath = "/api/v1"
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	// Set up routes
	v1 := r.Group("/api/v1")
	web.Init(v1, DB, *IS, *config)

	r.Run(fmt.Sprintf(":%d", port))
}
