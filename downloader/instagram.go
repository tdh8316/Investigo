package downloader

import (
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/tidwall/gjson"
)

func downloadInstagram(url string, logger *log.Logger) {
	s := strings.Split(url, "/")
	username := s[len(s)-1]
	r, _ := http.Get(url + "?__a=1")

	bdB, _ := ioutil.ReadAll(r.Body)

	edges := gjson.GetBytes(bdB, "graphql.user.edge_owner_to_timeline_media.edges")

	for i, edge := range edges.Array() {
		uri := edge.Get("node").Get("display_url").String()
		r, e := http.Get(uri)
		if e != nil {
			log.Fatal(e)
		}
		defer r.Body.Close()
		os.Mkdir(username, os.ModePerm)
		file, err := os.Create(username + "/instagram_" + strconv.Itoa(i) + ".jpg")
		if err != nil {
			log.Fatal(err)
			defer file.Close()
		}

		_, err = io.Copy(file, r.Body)
		if err != nil {
			log.Fatal(err)
		}
	}
}
