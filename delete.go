package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/shkh/lastfm-go/lastfm"
	"hawx.me/code/xdg"
)

var (
	auth  = flag.String("auth", "", "")
	after = flag.String("after", "730h", "")
	save  = flag.String("save", "", "")
	help  = flag.Bool("help", false, "")
)

const HELP = `Usage: la-delete [options]

  Deletes old scrobbles. Note: If --save is not given data is not saved!

    --auth PATH         # Path to file with auth details
    --after DUR         # Duration to delete after (default: '730h')
    --save DIR          # Directory to save scrobbles to
    --help              # Display this help message
`

type Saver interface {
	Save(id string, track interface{}) error
}

type emptySaver struct{}

func (_ *emptySaver) Save(_ string, _ interface{}) error { return nil }

type fileSaver struct {
	loc string
}

func (s *fileSaver) Save(id string, track interface{}) error {
	trackLoc := filepath.Join(s.loc, id+".json")

	data, err := json.Marshal(track)
	if err != nil {
		return err
	}

	log.Println("writing:", trackLoc)
	err = ioutil.WriteFile(trackLoc, data, 0644)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	flag.Parse()

	if *help {
		fmt.Println(HELP)
		return
	}

	authPath := xdg.Config("la-delete/auth")
	if *auth != "" {
		authPath = *auth
	}

	var conf struct {
		ApiKey, ApiSecret, Username, Password string
	}

	if _, err := toml.DecodeFile(authPath, &conf); err != nil {
		log.Fatal(err)
	}

	api := lastfm.New(conf.ApiKey, conf.ApiSecret)
	if err := api.Login(conf.Username, conf.Password); err != nil {
		log.Fatal(err)
	}

	dur, err := time.ParseDuration(*after)
	if err != nil {
		log.Fatal(err)
	}

	var saver Saver = &emptySaver{}
	if *save != "" {
		saver = &fileSaver{loc: *save}
	}

	for {
		result, _ := api.User.GetRecentTracks(lastfm.P{
			"user":  conf.Username,
			"to":    time.Now().UTC().Add(-dur).Unix(),
			"limit": 200,
		})

		if len(result.Tracks) == 0 {
			break
		}

		for _, t := range result.Tracks {
			if err := saver.Save(t.Date.Uts, t); err != nil {
				log.Println(err)
				break
			}

			err := api.Library.RemoveScrobble(lastfm.P{
				"artist":    t.Artist.Name,
				"track":     t.Name,
				"timestamp": t.Date.Uts,
			})

			if err != nil {
				log.Println(err)
				break
			}

			log.Println(t.Name, "by", t.Artist.Name)
		}
	}
}
