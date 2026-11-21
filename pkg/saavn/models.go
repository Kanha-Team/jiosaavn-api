// Copyright (c) 2026 Kanha-Team
// All rights reserved.

package saavn

import "github.com/Kanha-Team/jiosaavn-api/pkg/crypto"

// APIResponse represents standard API output
type APIResponse[T any] struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	Data    T      `json:"data"`
}

// ArtistMini represents brief artist info
type ArtistMini struct {
	ID    string               `json:"id"`
	Name  string               `json:"name"`
	Role  string               `json:"role,omitempty"`
	Type  string               `json:"type,omitempty"`
	Image []crypto.QualityLink `json:"image,omitempty"`
	URL   string               `json:"url,omitempty"`
}

// AlbumMini represents minimal album information inside a song
type AlbumMini struct {
	ID    string               `json:"id"`
	Name  string               `json:"name"`
	URL   string               `json:"url"`
	Image []crypto.QualityLink `json:"image,omitempty"`
}

// Song represents normalized song details
type Song struct {
	ID              string               `json:"id"`
	Name            string               `json:"name"`
	Type            string               `json:"type"`
	Year            string               `json:"year"`
	ReleaseDate     string               `json:"release_date"`
	Duration        string               `json:"duration"`
	Label           string               `json:"label"`
	ExplicitContent string               `json:"explicit_content"`
	PlayCount       string               `json:"play_count"`
	Language        string               `json:"language"`
	HasLyrics       string               `json:"has_lyrics"`
	LyricsID        string               `json:"lyrics_id,omitempty"`
	URL             string               `json:"url"`
	Copyright       string               `json:"copyright"`
	Album           AlbumMini            `json:"album"`
	Artists         []ArtistMini         `json:"artists"`
	PrimaryArtists  string               `json:"primary_artists"`
	Singers         string               `json:"singers"`
	Image           []crypto.QualityLink `json:"image"`
	DownloadURL     []crypto.QualityLink `json:"download_url"`
}

// Album represents normalized album details
type Album struct {
	ID              string               `json:"id"`
	Name            string               `json:"name"`
	Year            string               `json:"year"`
	Type            string               `json:"type"`
	PlayCount       string               `json:"play_count"`
	Language        string               `json:"language"`
	ExplicitContent string               `json:"explicit_content"`
	SongCount       string               `json:"song_count"`
	URL             string               `json:"url"`
	PrimaryArtists  string               `json:"primary_artists"`
	Artists         []ArtistMini         `json:"artists"`
	Image           []crypto.QualityLink `json:"image"`
	Songs           []Song               `json:"songs"`
}

// Playlist represents normalized playlist details
type Playlist struct {
	ID              string               `json:"id"`
	Name            string               `json:"name"`
	Type            string               `json:"type"`
	PlayCount       string               `json:"play_count"`
	Language        string               `json:"language"`
	ExplicitContent string               `json:"explicit_content"`
	SongCount       string               `json:"song_count"`
	URL             string               `json:"url"`
	HeaderDesc      string               `json:"header_desc"`
	Image           []crypto.QualityLink `json:"image"`
	Songs           []Song               `json:"songs"`
}

// Lyrics represents song lyrics
type Lyrics struct {
	Lyrics    string `json:"lyrics"`
	Snippet   string `json:"snippet"`
	Copyright string `json:"copyright"`
}

// Chart represents a top chart or featured playlist item
type Chart struct {
	ID    string               `json:"id"`
	Title string               `json:"title"`
	Type  string               `json:"type"`
	Count int                  `json:"count"`
	URL   string               `json:"url"`
	Image []crypto.QualityLink `json:"image"`
}

// SearchSongsResult represents paginated song search results
type SearchSongsResult struct {
	Total   int    `json:"total"`
	Start   int    `json:"start"`
	Results []Song `json:"results"`
}

// SearchAlbumsResult represents paginated album search results
type SearchAlbumsResult struct {
	Total   int     `json:"total"`
	Start   int     `json:"start"`
	Results []Album `json:"results"`
}

// SearchPlaylistsResult represents paginated playlist search results
type SearchPlaylistsResult struct {
	Total   int        `json:"total"`
	Start   int        `json:"start"`
	Results []Playlist `json:"results"`
}

// SearchArtistsResult represents paginated artist search results
type SearchArtistsResult struct {
	Total   int          `json:"total"`
	Start   int          `json:"start"`
	Results []ArtistMini `json:"results"`
}

// SearchAllResult represents all-in-one autocomplete search
type SearchAllResult struct {
	Songs     []Song       `json:"songs"`
	Albums    []Album      `json:"albums"`
	Playlists []Playlist   `json:"playlists"`
	Artists   []ArtistMini `json:"artists"`
	TopQuery  []any        `json:"top_query,omitempty"`
}
