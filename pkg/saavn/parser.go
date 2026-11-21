// Copyright (c) 2026 Kanha-Team
// All rights reserved.

package saavn

import (
	"html"
	"strings"

	"github.com/Kanha-Team/jiosaavn-api/pkg/crypto"
)

// cleanText unescapes HTML entities and trims whitespace
func cleanText(s string) string {
	return strings.TrimSpace(html.UnescapeString(s))
}

// ParseRawSong converts a raw song map into a normalized Song struct
func ParseRawSong(raw map[string]any) Song {
	id, _ := raw["id"].(string)
	name, _ := raw["title"].(string)
	if name == "" {
		name, _ = raw["song"].(string)
	}

	songType, _ := raw["type"].(string)
	year, _ := raw["year"].(string)
	permaURL, _ := raw["perma_url"].(string)
	playCount, _ := raw["play_count"].(string)
	language, _ := raw["language"].(string)
	explicit, _ := raw["explicit_content"].(string)
	imageURL, _ := raw["image"].(string)

	var (
		duration       string
		releaseDate    string
		label          string
		copyright      string
		hasLyrics      string
		lyricsSnippet  string
		encryptedMedia string
		albumName      string
		albumID        string
		albumURL       string
		primaryArtists string
		singers        string
		artists        []ArtistMini
	)

	if moreInfo, ok := raw["more_info"].(map[string]any); ok {
		duration, _ = moreInfo["duration"].(string)
		releaseDate, _ = moreInfo["release_date"].(string)
		label, _ = moreInfo["label"].(string)
		copyright, _ = moreInfo["copyright_text"].(string)
		hasLyrics, _ = moreInfo["has_lyrics"].(string)
		lyricsSnippet, _ = moreInfo["lyrics_snippet"].(string)
		encryptedMedia, _ = moreInfo["encrypted_media_url"].(string)
		albumName, _ = moreInfo["album"].(string)
		albumID, _ = moreInfo["album_id"].(string)
		albumURL, _ = moreInfo["album_url"].(string)

		if pa, ok := moreInfo["artistMap"].(map[string]any); ok {
			if paList, ok := pa["primary_artists"].([]any); ok {
				for _, item := range paList {
					if aMap, ok := item.(map[string]any); ok {
						aID, _ := aMap["id"].(string)
						aName, _ := aMap["name"].(string)
						aRole, _ := aMap["role"].(string)
						aType, _ := aMap["type"].(string)
						aImage, _ := aMap["image"].(string)
						aURL, _ := aMap["perma_url"].(string)
						artists = append(artists, ArtistMini{
							ID:    aID,
							Name:  cleanText(aName),
							Role:  aRole,
							Type:  aType,
							Image: crypto.FormatImageLinks(aImage),
							URL:   aURL,
						})
					}
				}
			}
		}

		if s, ok := moreInfo["singers"].(string); ok {
			singers = s
		}
		if s, ok := moreInfo["music"].(string); ok && singers == "" {
			singers = s
		}
	}

	if primaryArtists == "" {
		if pa, ok := raw["primary_artists"].(string); ok && pa != "" {
			primaryArtists = pa
		}
	}
	if primaryArtists == "" && len(artists) > 0 {
		var names []string
		for _, a := range artists {
			if a.Name != "" {
				names = append(names, a.Name)
			}
		}
		primaryArtists = strings.Join(names, ", ")
	}
	if primaryArtists == "" {
		if sub, ok := raw["subtitle"].(string); ok {
			primaryArtists = sub
		}
	}

	if albumName == "" {
		if a, ok := raw["album"].(string); ok {
			albumName = a
		}
	}

	if albumName != "" {
		primaryArtists = strings.TrimSuffix(primaryArtists, " - "+albumName)
		primaryArtists = strings.TrimSuffix(primaryArtists, " · "+albumName)
	}

	if albumID == "" {
		if a, ok := raw["albumid"].(string); ok {
			albumID = a
		}
	}
	if albumURL == "" {
		if a, ok := raw["album_url"].(string); ok {
			albumURL = a
		}
	}
	if duration == "" {
		if d, ok := raw["duration"].(string); ok {
			duration = d
		}
	}
	if label == "" {
		if l, ok := raw["label"].(string); ok {
			label = l
		}
	}
	if copyright == "" {
		if c, ok := raw["copyright_text"].(string); ok {
			copyright = c
		}
	}
	if hasLyrics == "" {
		if hl, ok := raw["has_lyrics"].(string); ok {
			hasLyrics = hl
		}
	}
	if encryptedMedia == "" {
		if em, ok := raw["encrypted_media_url"].(string); ok {
			encryptedMedia = em
		}
	}
	if singers == "" {
		if s, ok := raw["singers"].(string); ok {
			singers = s
		}
	}

	var downloadLinks []crypto.QualityLink
	if encryptedMedia != "" {
		decrypted, err := crypto.DecryptURL(encryptedMedia)
		if err == nil && decrypted != "" {
			downloadLinks = crypto.GetAudioDownloadLinks(decrypted)
		}
	}

	return Song{
		ID:              id,
		Name:            cleanText(name),
		Type:            songType,
		Year:            year,
		ReleaseDate:     releaseDate,
		Duration:        duration,
		Label:           cleanText(label),
		ExplicitContent: explicit,
		PlayCount:       playCount,
		Language:        language,
		HasLyrics:       hasLyrics,
		LyricsID:        lyricsSnippet,
		URL:             permaURL,
		Copyright:       cleanText(copyright),
		Album: AlbumMini{
			ID:    albumID,
			Name:  cleanText(albumName),
			URL:   albumURL,
			Image: crypto.FormatImageLinks(imageURL),
		},
		Artists:        artists,
		PrimaryArtists: cleanText(primaryArtists),
		Singers:        cleanText(singers),
		Image:          crypto.FormatImageLinks(imageURL),
		DownloadURL:    downloadLinks,
	}
}

// ParseRawAlbum converts a raw album map into a normalized Album struct
func ParseRawAlbum(raw map[string]any) Album {
	id, _ := raw["id"].(string)
	name, _ := raw["title"].(string)
	year, _ := raw["year"].(string)
	albumType, _ := raw["type"].(string)
	playCount, _ := raw["play_count"].(string)
	language, _ := raw["language"].(string)
	explicit, _ := raw["explicit_content"].(string)
	listCount, _ := raw["list_count"].(string)
	permaURL, _ := raw["perma_url"].(string)
	imageURL, _ := raw["image"].(string)

	primaryArtists, _ := raw["subtitle"].(string)
	if primaryArtists == "" {
		if mi, ok := raw["more_info"].(map[string]any); ok {
			primaryArtists, _ = mi["music"].(string)
		}
	}

	var songs []Song
	if list, ok := raw["list"].([]any); ok {
		for _, item := range list {
			if sMap, ok := item.(map[string]any); ok {
				songs = append(songs, ParseRawSong(sMap))
			}
		}
	}

	return Album{
		ID:              id,
		Name:            cleanText(name),
		Year:            year,
		Type:            albumType,
		PlayCount:       playCount,
		Language:        language,
		ExplicitContent: explicit,
		SongCount:       listCount,
		URL:             permaURL,
		PrimaryArtists:  cleanText(primaryArtists),
		Image:           crypto.FormatImageLinks(imageURL),
		Songs:           songs,
	}
}

// ParseRawPlaylist converts a raw playlist map into a normalized Playlist struct
func ParseRawPlaylist(raw map[string]any) Playlist {
	id, _ := raw["id"].(string)
	name, _ := raw["title"].(string)
	if name == "" {
		name, _ = raw["listname"].(string)
	}
	if name == "" {
		name, _ = raw["name"].(string)
	}
	playlistType, _ := raw["type"].(string)
	playCount, _ := raw["play_count"].(string)
	language, _ := raw["language"].(string)
	explicit, _ := raw["explicit_content"].(string)
	listCount, _ := raw["list_count"].(string)
	permaURL, _ := raw["perma_url"].(string)
	headerDesc, _ := raw["header_desc"].(string)
	imageURL, _ := raw["image"].(string)

	var songs []Song
	if list, ok := raw["list"].([]any); ok {
		for _, item := range list {
			if sMap, ok := item.(map[string]any); ok {
				songs = append(songs, ParseRawSong(sMap))
			}
		}
	}

	return Playlist{
		ID:              id,
		Name:            cleanText(name),
		Type:            playlistType,
		PlayCount:       playCount,
		Language:        language,
		ExplicitContent: explicit,
		SongCount:       listCount,
		URL:             permaURL,
		HeaderDesc:      cleanText(headerDesc),
		Image:           crypto.FormatImageLinks(imageURL),
		Songs:           songs,
	}
}

// CleanJSONResponse strips leading junk or BOM if returned by Saavn
func CleanJSONResponse(data []byte) []byte {
	s := strings.TrimSpace(string(data))
	// Strip UTF-8 BOM if present
	s = strings.TrimPrefix(s, "\ufeff")
	return []byte(s)
}

// ExtractTokenFromURL extracts perma_url token
func ExtractTokenFromURL(rawURL string) string {
	parts := strings.Split(strings.Trim(rawURL, "/"), "/")
	if len(parts) > 0 {
		return parts[len(parts)-1]
	}
	return rawURL
}
