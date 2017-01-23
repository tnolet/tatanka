package main

import (
	"flag"
	"github.com/tnolet/tatanka/api"
	"github.com/tnolet/tatanka/control"
	"github.com/tnolet/tatanka/helpers"
	"github.com/tnolet/tatanka/mail"
	"github.com/tnolet/tatanka/store"
	"log"
	"os"
	"runtime"
	"strconv"
)

const Version = "0.1.0"

var (
	HomeRegion  = flag.String("homeRegion", "eu-west-1", "Home region")
	HomeBucket  = flag.String("homeBucket", "tnolet-tatanka", "S3 Bucket")
	HomeKey     = flag.String("homeKey", "conf/config_example.json", "S3 Key")
	localConfig = flag.String("localConfig", "", "Override S3 path for local testing")
	HomeEmail   = flag.String("homeEmail", "tim@unumotors.com", "Email address")
	Port        = flag.Int("port", 1980, "Tatanka's API port")
	Noop        = flag.Bool("noop", false, "Start in noop mode")
	NumCPU      int
)

func init() {
	log.SetFlags(log.LstdFlags)
	NumCPU = runtime.NumCPU()
	runtime.GOMAXPROCS(NumCPU)
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

	mailChan := make(chan string, 100)
	mailer := mail.New(*HomeEmail, *HomeRegion, mailChan)
	mailer.Start()

	stateChan := make(chan (store.State))

	if len(*HomeBucket) > 0 || len(*HomeKey) > 0 || len(*localConfig) > 0 {

		Store := store.New(*HomeRegion, *HomeBucket, *HomeKey, *localConfig, stateChan)
		Store.Start()

	} else {
		log.Fatal("Please provide a Home location")
		mailChan <- control.FatalErrorMail("No or invalid home location")
		os.Exit(1)
	}

	controlChan := make(chan (control.Message), 10)
	controller := control.New(mailChan, controlChan, stateChan, *Noop)

	controller.Start()

	controlChan <- &control.Init{}
	controlChan <- &control.StartDeathWatch{}

	if Api, err := api.New(Version, controller); err != nil {
		panic("failed to create REST Api")
	} else {
		Api.Run("0.0.0.0:" + strconv.Itoa(*Port))
	}
}
