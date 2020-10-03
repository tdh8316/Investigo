package downloader

import (
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/tidwall/gjson"
)

func downloadInstagram(url string, logger *log.Logger) {
	s := strings.Split(url, "/")
	username := s[len(s)-1]
	os.Mkdir(username, os.ModePerm)

	r, _ := http.Get(url + "?__a=1")
	bdB, _ := ioutil.ReadAll(r.Body)
	r.Body.Close()

	profilePicURLHd := gjson.GetBytes(bdB, "graphql.user.profile_pic_url_hd").String()
	r, e := http.Get(profilePicURLHd)
	if e != nil {
		log.Fatal(e)
	}
	defer r.Body.Close()
	file, err := os.Create(username + "/instagram_profile_pic_hd.jpg")
	if err != nil {
		log.Fatal(err)
	}
	_, err = io.Copy(file, r.Body)
	if err != nil {
		log.Fatal(err)
	}

	edges := gjson.GetBytes(bdB, "graphql.user.edge_owner_to_timeline_media.edges")
	var wg sync.WaitGroup
	for i, edge := range edges.Array() {
		wg.Add(1)
		go func(edge gjson.Result, i int) {
			defer wg.Done()
			uri := edge.Get("node").Get("display_url").String()
			r, e := http.Get(uri)
			if e != nil {
				log.Fatal(e)
			}

			file, err := os.Create(username + "/instagram_" + strconv.Itoa(i) + ".jpg")
			if err != nil {
				log.Fatal(err)
			}

			_, err = io.Copy(file, r.Body)
			if err != nil {
				log.Fatal(err)
			}

			r.Body.Close()
			file.Close()
		}(edge, i)
	}
	wg.Wait()
}
