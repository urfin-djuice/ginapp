package ginapp

import (
	"context"
	"flag"
	"log"
	"net/http"
	"oko/docs"
	"oko/pkg/cfg"
	"oko/pkg/cli"
	"oko/pkg/e"
	"oko/pkg/env"
	"oko/pkg/ginapp/controller"
	"oko/pkg/valid"
	"os"
	"os/signal"
	"time"

	"github.com/gin-contrib/cache/persistence"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/gomodule/redigo/redis"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	"gopkg.in/go-playground/validator.v9"
)

type App struct {
	Engine          *gin.Engine
	cachePool       *redis.Pool
	CacheStore      *persistence.RedisStore
	Srv             *http.Server
	RootRote        string
	RootHandlers    controller.HandlerList
	Head            *gin.RouterGroup
	Ctrls           []controller.Ctrl
	Validators      []valid.Item
	ValidatorEngine *validator.Validate
	ValidatorMsgs   map[string]string
}

func (a *App) Init() {
	cfg.Load()
	a.initCache()
	a.initEngine()
	a.initValidator()

	flushCache := flag.Bool("flush-cache", false, "Flush redis page cache")
	cli.Init()
	if *flushCache {
		log.Println("Flush redis page cache")
		err := a.CacheStore.Flush()
		if err != nil {
			log.Fatal(err)
		}
		log.Println("Redis page cache flushed")
	}
}

func (a *App) Do() {
	defer a.close()
	a.Srv = &http.Server{
		Addr:    cfg.App.APIListen,
		Handler: a.Engine,
	}
	go func() {
		if err := a.Srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("listen: ", err)
		}
	}()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Shutdown Server ...")
}

func (a *App) close() {
	if err := a.cachePool.Close(); err != nil {
		log.Printf("Fail to close cache redis pool: [%+v]", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := a.Srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	log.Println("Server exiting")
}

func (a *App) initEngine() {
	a.Engine = gin.New()
	a.Engine.Use(gin.Logger())
	a.addCors()
	a.Engine.GET("/", pingPage)

	docs.SwaggerInfo.Host = cfg.App.APIHost
	config := &ginSwagger.Config{
		URL: cfg.App.APISchema + "://" + cfg.App.APIHost + "/swagger/doc.json", // The url pointing to API definition
	}
	a.Engine.GET("/swagger/*any", ginSwagger.CustomWrapHandler(config, swaggerFiles.Handler))
	a.Head = a.Engine.Group(a.RootRote, a.RootHandlers...)
	if len(a.Ctrls) > 0 {
		for _, c := range a.Ctrls {
			a.addCtrl(c)
		}
	}
}

func (a *App) addCors() {
	a.Engine.Use(cors.New(cors.Config{
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "HEAD"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", cfg.App.AuthTokenKey},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return true
		},
		MaxAge: 12 * time.Hour,
	}))
}

func (a *App) addCtrl(c controller.Ctrl) {
	if len(c.Acts) > 0 {
		r := a.Head.Group("/"+c.Name+"/", c.Handlers...)
		for _, a := range c.Acts {
			switch a.Method {
			case "GET":
				r.GET(a.Route, a.Handlers...)
			case "POST":
				r.POST(a.Route, a.Handlers...)
			case "DELETE":
				r.DELETE(a.Route, a.Handlers...)
			case "PUT":
				r.PUT(a.Route, a.Handlers...)
			case "HEAD":
				r.HEAD(a.Route, a.Handlers...)
			default:
				log.Panicln("Unsupported method " + a.Method)
			}
		}
	}
}

func (a *App) initCache() {
	a.cachePool = &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.DialURL(env.GetEnvOrPanic("REDIS_API_CACHE_URL"))
			if err != nil {
				return nil, err
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
	a.CacheStore = persistence.NewRedisCacheWithPool(a.cachePool, time.Second)
}

func (a *App) initValidator() {
	var ok bool
	if a.ValidatorEngine, ok = binding.Validator.Engine().(*validator.Validate); !ok {
		log.Println(a.ValidatorEngine)
		log.Panicln("Fail to binding validator engine")
	}
	if len(a.Validators) > 0 {
		for _, item := range a.Validators {
			if item.Handler != nil {
				if err := a.ValidatorEngine.RegisterValidation(item.Key, item.Handler); err != nil {
					log.Panicln("Fail to register validator "+item.Key, err)
				}
			}
			e.ValidatorMessages[item.Key] = item.Message
		}
	}
}

func pingPage(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "OK"})
}
