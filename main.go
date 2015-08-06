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
	HomeKey     = flag.String("homeKey", "config_example.json", "S3 Key")
	localConfig = flag.String("localConfig", "", "Override S3 path for local testing")
	HomeEmail   = flag.String("homeEmail", "tim@magnetic.io", "Email address")
	Port        = flag.Int("port", 1980, "Tatanka's API port")
	Noop        = flag.Bool("noop", false, "Start in noop mode")
	Store       *store.Store
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

	// Call home and get state
	if len(*HomeBucket) > 0 || len(*HomeKey) > 0 || len(*localConfig) > 0 {

		Store = store.New(*HomeRegion, *HomeBucket, *HomeKey, *localConfig)

	} else {
		log.Fatal("Please provide a Home location")
		os.Exit(1)
	}

	state, err := Store.GetState()
	if err != nil {
		log.Fatal("Error getting state: " + err.Error())
	}

	/*
		######################
		Initialize sub systems
		######################
	*/

	log.Println("Initializing mailer...")

	mailer := mailer.New(*HomeEmail)

	// if err := mailer.Send("hi there! I'm going to try and run on a " + state.InstanceType); err != nil {
	// 	log.Println("Error sending email: " + err.Error())
	// }

	controller := &control.Controller{}

	if !*Noop {
		controller.Start(state, *mailer)
	}

	if *Noop {

		evacChan := make(chan (bool))
		evacNotice := false
		controller.StartTerminationChecker(state.TerminationUrl, state.LifeTime, evacChan)

		go func() {
			for evacNotice == false {
				select {
				case evacNotice := <-evacChan:
					if evacNotice {
						log.Println("Got evac notice")
					} else {
						log.Println("Not evacing just yet")
					}
				}
			}
		}()
	}

	log.Println("Initializing api...")

	if api, err := api.New(Version); err != nil {
		panic("failed to create REST Api")
	} else {
		api.Run("0.0.0.0:" + strconv.Itoa(*Port))
	}

}
