package main

import (
	"fmt"
	"log"
	"net/http"

	flatbuffers "github.com/google/flatbuffers/go"
	"github.com/julienschmidt/httprouter"

	"github.com/ko/feedparser"
	"github.com/ko/feedserver/schema/feeds"
)

func AuthCheck(h httprouter.Handle, requiredUser, requiredPassword string) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		if 1 == 1 {
			h(w, r, ps)
		} else {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		}
	}
}

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	feed := feedparser.JsonToItunesFeed(feedparser.XmlToJson(feedparser.GetFeedLocal()))
	episodes := feedparser.GetEpisodes(feed)
	episode1 := feedparser.GetEpisode(episodes, 0)
	shownotes := feedparser.GetEpisodeNotes(episode1)
	fmt.Fprint(w, shownotes)
}

func SearchPodcasts(w http.ResponseWriter, r *http.Request, params httprouter.Params) {

	query := params.ByName("query")
	podcasts := feedparser.Search(query)
	barray, err := feedparser.SearchResultsItemsToJson(podcasts)
	if err != nil {
		log.Fatal(err)
	}
	str := string(barray)
	fmt.Fprint(w, str)
}

func SecretRoute(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Secret Route!\n")
}

func MakeFeed(b *flatbuffers.Builder) []byte {
	b.Reset()

	title := []byte("this is the title")
	title_position := b.CreateByteString(title)

	feeds.ItunesFeedStart(b)
	feeds.ItunesFeedAddTitle(b, title_position)
	feed_position := feeds.ItunesFeedEnd(b)

	b.Finish(feed_position)

	return b.Bytes[b.Head():]
}

func ReadFeed(buf []byte) (feed *feeds.ItunesFeed, title string) {
	feed = feeds.GetRootAsItunesFeed(buf, 0)

	title = string(feed.Title())

	return
}

func TestFeedRead(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	builder := flatbuffers.NewBuilder(0)
	buf := MakeFeed(builder)
	_, title := ReadFeed(buf)
	fmt.Fprint(w, title)
}

func main() {
	user := "username"
	pass := "password"
	router := httprouter.New()
	router.GET("/", Index)
	router.GET("/search/:query", SearchPodcasts)
	router.GET("/a/", AuthCheck(SecretRoute, user, pass))
	router.GET("/test/", TestFeedRead)

	log.Fatal(http.ListenAndServe(":8080", router))
}
