package api

import (
	"fmt"
	"net/http"

	"github.com/astaxie/beego"
	"github.com/cloudsonic/sonic-server/api/responses"
	"github.com/cloudsonic/sonic-server/engine"
)

type SearchingController struct {
	search       engine.Search
	query        string
	artistCount  int
	artistOffset int
	albumCount   int
	albumOffset  int
	songCount    int
	songOffset   int
}

func NewSearchingController(search engine.Search) *SearchingController {
	return &SearchingController{search: search}
}

func (c *SearchingController) getParams(r *http.Request) error {
	var err error
	c.query, err = RequiredParamString(r, "query", "Parameter query required")
	if err != nil {
		return err
	}
	c.artistCount = ParamInt(r, "artistCount", 20)
	c.artistOffset = ParamInt(r, "artistOffset", 0)
	c.albumCount = ParamInt(r, "albumCount", 20)
	c.albumOffset = ParamInt(r, "albumOffset", 0)
	c.songCount = ParamInt(r, "songCount", 20)
	c.songOffset = ParamInt(r, "songOffset", 0)
	return nil
}

func (c *SearchingController) searchAll() (engine.Entries, engine.Entries, engine.Entries) {
	as, err := c.search.SearchArtist(c.query, c.artistOffset, c.artistCount)
	if err != nil {
		beego.Error("Error searching for Artists:", err)
	}
	als, err := c.search.SearchAlbum(c.query, c.albumOffset, c.albumCount)
	if err != nil {
		beego.Error("Error searching for Albums:", err)
	}
	mfs, err := c.search.SearchSong(c.query, c.songOffset, c.songCount)
	if err != nil {
		beego.Error("Error searching for MediaFiles:", err)
	}

	beego.Debug(fmt.Sprintf("Searching for [%s] resulted in %d songs, %d albums and %d artists", c.query, len(mfs), len(als), len(as)))
	return mfs, als, as
}

func (c *SearchingController) Search2(w http.ResponseWriter, r *http.Request) (*responses.Subsonic, error) {
	err := c.getParams(r)
	if err != nil {
		return nil, err
	}
	mfs, als, as := c.searchAll()

	response := NewEmpty()
	searchResult2 := &responses.SearchResult2{}
	searchResult2.Artist = make([]responses.Artist, len(as))
	for i, e := range as {
		searchResult2.Artist[i] = responses.Artist{Id: e.Id, Name: e.Title}
	}
	searchResult2.Album = ToChildren(als)
	searchResult2.Song = ToChildren(mfs)
	response.SearchResult2 = searchResult2
	return response, nil
}

func (c *SearchingController) Search3(w http.ResponseWriter, r *http.Request) (*responses.Subsonic, error) {
	err := c.getParams(r)
	if err != nil {
		return nil, err
	}
	mfs, als, as := c.searchAll()

	response := NewEmpty()
	searchResult3 := &responses.SearchResult3{}
	searchResult3.Artist = make([]responses.ArtistID3, len(as))
	for i, e := range as {
		searchResult3.Artist[i] = responses.ArtistID3{
			Id:         e.Id,
			Name:       e.Title,
			CoverArt:   e.CoverArt,
			AlbumCount: e.AlbumCount,
		}
	}
	searchResult3.Album = ToAlbums(als)
	searchResult3.Song = ToChildren(mfs)
	response.SearchResult3 = searchResult3
	return response, nil
}
