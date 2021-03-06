package ui

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/joshuathompson/baton/api"
	"github.com/joshuathompson/baton/utils"
	"github.com/jroimartin/gocui"
)

type PlaylistTrackTable struct {
	data     *api.PlaylistTracksPaged
	playlist *api.SimplePlaylist
}

func NewPlaylistTrackTable(playlistTracksPaged *api.PlaylistTracksPaged, playlist *api.SimplePlaylist) *PlaylistTrackTable {
	return &PlaylistTrackTable{
		data:     playlistTracksPaged,
		playlist: playlist,
	}
}

func (t *PlaylistTrackTable) getColumnWidths(maxX int) map[string]int {
	m := make(map[string]int)
	m["length"] = maxX / 8
	m["artist"] = maxX / 4
	m["album"] = maxX / 5
	m["popularity"] = maxX / 10
	m["name"] = maxX - m["track_number"] - m["length"] - m["artist"] - m["album"] - m["popularity"]

	return m
}

func (t *PlaylistTrackTable) renderHeader(v *gocui.View, maxX int) {
	columnWidths := t.getColumnWidths(maxX)

	namesHeader := utils.LeftPaddedString("NAME", columnWidths["name"], 2)
	artistHeader := utils.LeftPaddedString("ARTIST", columnWidths["artist"], 2)
	albumHeader := utils.LeftPaddedString("ALBUM", columnWidths["album"], 2)
	lengthHeader := utils.LeftPaddedString("LENGTH", columnWidths["length"], 2)
	popularityHeader := utils.LeftPaddedString("POPULARITY", columnWidths["popularity"], 2)

	loadedLength := maxX / 3
	loadedHeader := utils.LeftPaddedString(fmt.Sprintf("Showing %d of %d tracks", len(t.data.Items), t.data.Total), loadedLength, 2)
	titleLength := maxX - loadedLength

	fmt.Fprintf(v, "\u001b[1m%s %s[0m\n\n", utils.LeftPaddedString("TRACKS", titleLength, 2), loadedHeader)
	fmt.Fprintf(v, "\u001b[1m%s %s %s %s %s\u001b[0m\n", namesHeader, artistHeader, albumHeader, lengthHeader, popularityHeader)
}

func (t *PlaylistTrackTable) render(v *gocui.View, maxX int) {
	columnWidths := t.getColumnWidths(maxX)

	for _, item := range t.data.Items {
		name := utils.LeftPaddedString(item.Track.Name, columnWidths["name"], 2)
		var artistNames []string
		for _, artist := range item.Track.Artists {
			artistNames = append(artistNames, artist.Name)
		}
		artists := utils.LeftPaddedString(strings.Join(artistNames, ", "), columnWidths["artist"], 2)
		album := utils.LeftPaddedString(item.Track.Album.Name, columnWidths["album"], 2)
		length := utils.LeftPaddedString(utils.MillisecondsToFormattedTime(item.Track.DurationMs), columnWidths["length"], 2)
		popularity := utils.LeftPaddedString(strconv.Itoa(item.Track.Popularity), columnWidths["popularity"], 2)

		fmt.Fprintf(v, "\n%s %s %s %s %s", name, artists, album, length, popularity)
	}
}

func (t *PlaylistTrackTable) getTableLength() int {
	return len(t.data.Items)
}

func (t *PlaylistTrackTable) loadNextRecords() error {
	if t.data.Next != "" {
		nextTracks, err := api.GetNextTracksForPlaylist(t.data.Next)

		if err != nil {
			return err
		}

		t.data.Href = nextTracks.Href
		t.data.Offset = nextTracks.Offset
		t.data.Next = nextTracks.Next
		t.data.Previous = nextTracks.Previous
		t.data.Items = append(t.data.Items, nextTracks.Items...)
	}
	return nil
}

func (t *PlaylistTrackTable) playSelected(selectedIndex int) error {
	item := t.data.Items[selectedIndex]
	playerOptions := api.PlayerOptions{
		ContextURI: t.playlist.URI,
		Offset: &api.PlayerOffsetOptions{
			URI: item.Track.URI,
		},
	}
	return api.StartPlayback(&playerOptions)
}

func (t *PlaylistTrackTable) newTableFromSelection(selectedIndex int) (Table, error) {
	item := t.data.Items[selectedIndex]
	playerOptions := api.PlayerOptions{
		ContextURI: t.playlist.URI,
		Offset: &api.PlayerOffsetOptions{
			URI: item.Track.URI,
		},
	}

	return nil, api.StartPlayback(&playerOptions)
}
