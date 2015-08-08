package main

import (
	"flag"
	"github.com/tnolet/tatanka/api"
	"github.com/tnolet/tatanka/control"
	"github.com/tnolet/tatanka/helpers"
	"github.com/tnolet/tatanka/mailer"
	"github.com/tnolet/tatanka/store"
	"log"
	"os"
	"strconv"
)

var (
	HomeRegion  = flag.String("homeRegion", "eu-west-1", "Home region")
	HomeBucket  = flag.String("homeBucket", "tnolet-tatanka", "S3 Bucket")
	HomeKey     = flag.String("homeKey", "conf/config_example.json", "S3 Key")
	localConfig = flag.String("localConfig", "", "Override S3 path for local testing")
	HomeEmail   = flag.String("homeEmail", "tim@magnetic.io", "Email address")
	Port        = flag.Int("port", 1980, "Tatanka's API port")
	Noop        = flag.Bool("noop", false, "Start in noop mode")
)

func init() {
	log.SetFlags(log.LstdFlags)
}

func main() {
	flag.Parse()
	helpers.SetValueFromEnv(&HomeRegion, "TATANKA_HOME_REGION")
	helpers.SetValueFromEnv(&HomeBucket, "TATANKA_HOME_BUCKET")
	helpers.SetValueFromEnv(&HomeKey, "TATANKA_HOME_KEY")
	helpers.SetValueFromEnv(&localConfig, "TATANKA_LOCAL_CONFIG")
	helpers.SetValueFromEnv(&Port, "TATANKA_API_PORT")
	helpers.SetValueFromEnv(&Noop, "TATANKA_NOOP")

	log.Println("Starting Tatanka...")

	// Create store with channel
	stateChan := make(chan (store.State))

	if len(*HomeBucket) > 0 || len(*HomeKey) > 0 || len(*localConfig) > 0 {

		store := store.New(*HomeRegion, *HomeBucket, *HomeKey, *localConfig, stateChan)
		store.Start()

	} else {
		log.Fatal("Please provide a Home location")
		os.Exit(1)
	}

	mailer := mailer.New(*HomeEmail, *HomeRegion)

	controlChan := make(chan (control.Message))
	controller := control.New(*mailer, controlChan, stateChan)

	controller.Start()

	controlChan <- &control.Init{}
	controlChan <- &control.StartDeathWatch{}

	if api, err := api.New(Version); err != nil {
		panic("failed to create REST Api")
	} else {
		api.Run("0.0.0.0:" + strconv.Itoa(*Port))
	}

}
