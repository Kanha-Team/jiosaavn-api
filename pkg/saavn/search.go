// Copyright (c) 2026 Kanha-Team
// All rights reserved.

package saavn

import (
	"context"
	"encoding/json"
	"fmt"
	"html"
	"math"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/Kanha-Team/jiosaavn-api/pkg/crypto"
)

var (
	movieRegex     = regexp.MustCompile(`(?i)[(\[]from\s+[^)\]]+[)\]]`)
	nonAlphaRegex  = regexp.MustCompile(`[^a-zA-Z0-9\s]`)
	noiseWordRegex = regexp.MustCompile(`(?i)\b(song|songs|mp3|audio|video|lyrics|full\s+song|official)\b`)
	unwantedMods   = []string{"lofi", "lo-fi", "slowed", "reverb", "remix", "cover", "reprise", "trending version", "instrumental", "acoustic", "flip"}
)

// CuratedKnownHits maps common search queries to official JioSaavn song IDs.
// This ensures popular hit songs (e.g. Warner Music India tracks suppressed by JioSaavn's search index)
// are always available and accurately returned.
var CuratedKnownHits = map[string][]string{
	"janat ve":             {"Dz6zivf5", "4cJtmy5J"},
	"jannat ve":            {"Dz6zivf5", "4cJtmy5J"},
	"jannat ve darshan":    {"Dz6zivf5"},
	"darshan raval jannat": {"Dz6zivf5"},
	"darshan jannat ve":    {"Dz6zivf5"},
	"jannat":               {"4cJtmy5J", "DjmwRNFg", "Dz6zivf5"},
	"kesariya":             {"7L5vYwR6"},
	"lut gaye":             {"W93hZ63U"},
	"senorita":             {"4-mFwVp7", "_Bf2lP1F"},
	"raataan lambiyan":     {"_mK6eD4B"},
	"ratan lambiyan":       {"_mK6eD4B"},
}

func cleanNoiseKeywords(q string) string {
	q = noiseWordRegex.ReplaceAllString(q, " ")
	return strings.Join(strings.Fields(q), " ")
}

func normalizeQueryText(s string) string {
	s = html.UnescapeString(s)
	s = movieRegex.ReplaceAllString(s, " ")
	s = nonAlphaRegex.ReplaceAllString(s, " ")
	return strings.Join(strings.Fields(strings.ToLower(s)), " ")
}

// phoneticKey normalizes Hinglish/Hindi phonetics by collapsing consecutive duplicate letters
// and standardizing common vowel representations (e.g., jannat -> janat, raataan -> ratan)
func phoneticKey(s string) string {
	s = normalizeQueryText(s)
	if s == "" {
		return ""
	}
	var b strings.Builder
	var prev rune
	for _, r := range s {
		if r == prev && (r >= 'a' && r <= 'z') {
			continue
		}
		b.WriteRune(r)
		prev = r
	}
	res := b.String()
	res = strings.ReplaceAll(res, "ee", "i")
	res = strings.ReplaceAll(res, "oo", "u")
	return res
}

func calculateSongScore(song Song, query string, topQueryID string, autoIDMap map[string]bool) float64 {
	qNorm := normalizeQueryText(query)
	tNorm := normalizeQueryText(song.Name)
	if qNorm == "" || tNorm == "" {
		return 0
	}

	qPhonetic := phoneticKey(qNorm)
	tPhonetic := phoneticKey(tNorm)

	var score float64

	// 0. Autocomplete / Curated TopQuery origin bonus
	if song.ID != "" && song.ID == topQueryID {
		score += 160.0
	} else if autoIDMap != nil && autoIDMap[song.ID] {
		score += 80.0
	}

	// 1. Title match (exact string or phonetic Hinglish match)
	hasTitleMatch := false
	if tNorm == qNorm || (qPhonetic != "" && tPhonetic == qPhonetic) {
		score += 220.0
		hasTitleMatch = true
	} else if strings.HasPrefix(tNorm, qNorm) || (qPhonetic != "" && strings.HasPrefix(tPhonetic, qPhonetic)) {
		score += 130.0
		hasTitleMatch = true
	} else if strings.Contains(tNorm, qNorm) || (qPhonetic != "" && strings.Contains(tPhonetic, qPhonetic)) {
		score += 70.0
		hasTitleMatch = true
	}

	// 2. Word overlap (both literal and phonetic)
	qWords := strings.Fields(qNorm)
	tWords := strings.Fields(tNorm)
	tWordMap := make(map[string]bool, len(tWords))
	tPhoneticMap := make(map[string]bool, len(tWords))
	for _, w := range tWords {
		tWordMap[w] = true
		tPhoneticMap[phoneticKey(w)] = true
	}
	matches := 0
	for _, w := range qWords {
		if tWordMap[w] || tPhoneticMap[phoneticKey(w)] {
			matches++
		}
	}
	if len(qWords) > 0 {
		score += (float64(matches) / float64(len(qWords))) * 60.0
	}

	// 3. Artist match if query mentions artist
	artistNorm := normalizeQueryText(song.PrimaryArtists + " " + song.Singers)
	artistPhonetic := phoneticKey(artistNorm)
	hasArtistMatch := false
	for _, qw := range qWords {
		if len(qw) >= 4 && (strings.Contains(artistNorm, qw) || strings.Contains(artistPhonetic, phoneticKey(qw))) {
			score += 50.0
			hasArtistMatch = true
		}
	}

	// GUARD: If a song has ZERO title match, ZERO word matches, and ZERO artist match,
	// it is completely irrelevant to the user query and must NOT receive popularity boosts!
	if matches == 0 && !hasTitleMatch && !hasArtistMatch {
		return -1000.0
	}

	// 4. Penalize unwanted modifiers if not in user query
	rawLower := strings.ToLower(song.Name)
	for _, mod := range unwantedMods {
		if strings.Contains(rawLower, mod) && !strings.Contains(qNorm, mod) {
			score -= 150.0
		}
	}

	// 5. Play count boost (log scale, max ~280 points)
	playCount, _ := strconv.ParseFloat(song.PlayCount, 64)
	if playCount > 0 {
		score += math.Log10(playCount+1) * 35.0
	}

	return score
}

// SortSongsByRelevance re-ranks songs using popularity, autocomplete signal, and title similarity
func SortSongsByRelevance(songs []Song, query string, topQueryID string, autoIDMap map[string]bool) []Song {
	if len(songs) <= 1 || query == "" {
		return songs
	}

	type scoredSong struct {
		song  Song
		score float64
		index int
	}

	scored := make([]scoredSong, len(songs))
	for i, s := range songs {
		scored[i] = scoredSong{
			song:  s,
			score: calculateSongScore(s, query, topQueryID, autoIDMap),
			index: i,
		}
	}

	sort.SliceStable(scored, func(i, j int) bool {
		if scored[i].score != scored[j].score {
			return scored[i].score > scored[j].score
		}
		return scored[i].index < scored[j].index
	})

	sorted := make([]Song, len(songs))
	for i, ss := range scored {
		sorted[i] = ss.song
	}
	return sorted
}

// SearchSongs performs Hybrid Search combining autocomplete suggestions and catalog results
func (c *Client) SearchSongs(ctx context.Context, query string, page, limit int) (*SearchSongsResult, error) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 20
	}

	cleanQuery := cleanNoiseKeywords(query)
	if cleanQuery == "" {
		cleanQuery = query
	}

	var (
		topQueryID     string
		autoIDs        []string
		detectedArtist string
	)

	// Pass 1: Autocomplete suggestions (provides curated top hits & CTR-ranked tracks)
	autoData, err := c.Request(ctx, "autocomplete.get", map[string]string{
		"query": cleanQuery,
	})
	if err == nil {
		var autoRaw struct {
			TopQuery struct {
				Data []map[string]any `json:"data"`
			} `json:"topquery"`
			Songs struct {
				Data []map[string]any `json:"data"`
			} `json:"songs"`
			Albums struct {
				Data []map[string]any `json:"data"`
			} `json:"albums"`
		}
		if jsonErr := json.Unmarshal(CleanJSONResponse(autoData), &autoRaw); jsonErr == nil {
			if len(autoRaw.TopQuery.Data) > 0 {
				tq := autoRaw.TopQuery.Data[0]
				tqType, _ := tq["type"].(string)
				tqID, _ := tq["id"].(string)
				if tqType == "song" && tqID != "" {
					topQueryID = tqID
					autoIDs = append(autoIDs, tqID)
				} else if tqType == "artist" {
					detectedArtist, _ = tq["title"].(string)
				}
			}

			for _, s := range autoRaw.Songs.Data {
				sid, _ := s["id"].(string)
				if sid != "" {
					autoIDs = append(autoIDs, sid)
				}
			}

			for _, a := range autoRaw.Albums.Data {
				if mi, ok := a["more_info"].(map[string]any); ok {
					if pidsStr, ok := mi["song_pids"].(string); ok && pidsStr != "" {
						for _, pid := range strings.Split(pidsStr, ",") {
							pid = strings.TrimSpace(pid)
							if pid != "" {
								autoIDs = append(autoIDs, pid)
								break
							}
						}
					}
				}
			}
		}
	}

	// Pass 2: Main catalog search
	searchData, err := c.Request(ctx, "search.getResults", map[string]string{
		"q": cleanQuery,
		"p": strconv.Itoa(page),
		"n": strconv.Itoa(limit),
	})
	if err != nil && len(autoIDs) == 0 {
		return nil, err
	}

	var raw struct {
		Total   int              `json:"total"`
		Start   int              `json:"start"`
		Results []map[string]any `json:"results"`
	}
	if err == nil {
		_ = json.Unmarshal(CleanJSONResponse(searchData), &raw)
	}

	// Map candidates to deduplicate
	candidates := make(map[string]Song)
	for _, item := range raw.Results {
		song := ParseRawSong(item)
		if song.ID != "" {
			candidates[song.ID] = song
		}
	}

	// Fallback if 0 results and artist detected: try searching without artist name
	if len(candidates) == 0 && detectedArtist != "" {
		cleanedWithoutArtist := strings.TrimSpace(strings.ReplaceAll(strings.ToLower(cleanQuery), strings.ToLower(detectedArtist), ""))
		if cleanedWithoutArtist != "" {
			if retryData, retryErr := c.Request(ctx, "search.getResults", map[string]string{
				"q": cleanedWithoutArtist,
				"p": "1",
				"n": strconv.Itoa(limit),
			}); retryErr == nil {
				var retryRaw struct {
					Total   int              `json:"total"`
					Start   int              `json:"start"`
					Results []map[string]any `json:"results"`
				}
				if json.Unmarshal(CleanJSONResponse(retryData), &retryRaw) == nil {
					for _, item := range retryRaw.Results {
						song := ParseRawSong(item)
						if song.ID != "" {
							candidates[song.ID] = song
						}
					}
					raw.Total = retryRaw.Total
				}
			}
		}
	}

	// Inject curated hits for known tracks suppressed by JioSaavn catalog search
	pCleanQuery := phoneticKey(cleanQuery)
	for key, ids := range CuratedKnownHits {
		pKey := phoneticKey(key)
		if pCleanQuery == pKey || strings.Contains(pCleanQuery, pKey) || strings.Contains(pKey, pCleanQuery) {
			for _, cid := range ids {
				if _, exists := candidates[cid]; !exists {
					autoIDs = append(autoIDs, cid)
					if topQueryID == "" {
						topQueryID = cid
					}
				}
			}
		}
	}

	// Pass 3: Hydrate missing autocomplete IDs that have no details yet
	autoIDMap := make(map[string]bool)
	var missingIDs []string
	for _, id := range autoIDs {
		autoIDMap[id] = true
		if _, exists := candidates[id]; !exists {
			missingIDs = append(missingIDs, id)
			if len(missingIDs) >= 8 {
				break
			}
		}
	}

	if len(missingIDs) > 0 {
		if hydrated, hErr := c.GetSongsByIDs(ctx, missingIDs); hErr == nil {
			for _, song := range hydrated {
				if song.ID != "" {
					candidates[song.ID] = song
				}
			}
		}
	}

	// Convert map to slice
	allSongs := make([]Song, 0, len(candidates))
	for _, s := range candidates {
		allSongs = append(allSongs, s)
	}

	// Pass 4: Smart Re-ranking
	rankedSongs := SortSongsByRelevance(allSongs, query, topQueryID, autoIDMap)

	totalCount := raw.Total
	if len(rankedSongs) > totalCount {
		totalCount = len(rankedSongs)
	}

	if limit > 0 && len(rankedSongs) > limit {
		rankedSongs = rankedSongs[:limit]
	}

	return &SearchSongsResult{
		Total:   totalCount,
		Start:   raw.Start,
		Results: rankedSongs,
	}, nil
}

// SearchAlbums searches for albums matching query
func (c *Client) SearchAlbums(ctx context.Context, query string, page, limit int) (*SearchAlbumsResult, error) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 20
	}

	data, err := c.Request(ctx, "search.getAlbumResults", map[string]string{
		"q": query,
		"p": strconv.Itoa(page),
		"n": strconv.Itoa(limit),
	})
	if err != nil {
		return nil, err
	}

	var raw struct {
		Total   int              `json:"total"`
		Start   int              `json:"start"`
		Results []map[string]any `json:"results"`
	}
	if err := json.Unmarshal(CleanJSONResponse(data), &raw); err != nil {
		return nil, fmt.Errorf("failed to parse search albums response: %w", err)
	}

	var albums []Album
	for _, item := range raw.Results {
		albums = append(albums, ParseRawAlbum(item))
	}

	return &SearchAlbumsResult{
		Total:   raw.Total,
		Start:   raw.Start,
		Results: albums,
	}, nil
}

// SearchPlaylists searches for playlists matching query
func (c *Client) SearchPlaylists(ctx context.Context, query string, page, limit int) (*SearchPlaylistsResult, error) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 20
	}

	data, err := c.Request(ctx, "search.getPlaylistResults", map[string]string{
		"q": query,
		"p": strconv.Itoa(page),
		"n": strconv.Itoa(limit),
	})
	if err != nil {
		return nil, err
	}

	var raw struct {
		Total   int              `json:"total"`
		Start   int              `json:"start"`
		Results []map[string]any `json:"results"`
	}
	if err := json.Unmarshal(CleanJSONResponse(data), &raw); err != nil {
		return nil, fmt.Errorf("failed to parse search playlists response: %w", err)
	}

	var playlists []Playlist
	for _, item := range raw.Results {
		playlists = append(playlists, ParseRawPlaylist(item))
	}

	return &SearchPlaylistsResult{
		Total:   raw.Total,
		Start:   raw.Start,
		Results: playlists,
	}, nil
}

// SearchArtists searches for artists matching query
func (c *Client) SearchArtists(ctx context.Context, query string, page, limit int) (*SearchArtistsResult, error) {
	if page <= 0 {
		page = 1
	}
	if limit <= 0 {
		limit = 20
	}

	data, err := c.Request(ctx, "search.getArtistResults", map[string]string{
		"q": query,
		"p": strconv.Itoa(page),
		"n": strconv.Itoa(limit),
	})
	if err != nil {
		return nil, err
	}

	var raw struct {
		Total   int              `json:"total"`
		Start   int              `json:"start"`
		Results []map[string]any `json:"results"`
	}
	if err := json.Unmarshal(CleanJSONResponse(data), &raw); err != nil {
		return nil, fmt.Errorf("failed to parse search artists response: %w", err)
	}

	var artists []ArtistMini
	for _, item := range raw.Results {
		id, _ := item["id"].(string)
		name, _ := item["name"].(string)
		role, _ := item["role"].(string)
		image, _ := item["image"].(string)
		perma, _ := item["perma_url"].(string)

		artists = append(artists, ArtistMini{
			ID:    id,
			Name:  cleanText(name),
			Role:  role,
			Type:  "artist",
			Image: crypto.FormatImageLinks(image),
			URL:   perma,
		})
	}

	return &SearchArtistsResult{
		Total:   raw.Total,
		Start:   raw.Start,
		Results: artists,
	}, nil
}

// SearchAll performs autocomplete search across all categories
func (c *Client) SearchAll(ctx context.Context, query string) (*SearchAllResult, error) {
	data, err := c.Request(ctx, "autocomplete.get", map[string]string{
		"query": query,
	})
	if err != nil {
		return nil, err
	}

	var raw struct {
		Albums struct {
			Data []map[string]any `json:"data"`
		} `json:"albums"`
		Songs struct {
			Data []map[string]any `json:"data"`
		} `json:"songs"`
		Playlists struct {
			Data []map[string]any `json:"data"`
		} `json:"playlists"`
		Artists struct {
			Data []map[string]any `json:"data"`
		} `json:"artists"`
	}

	if err := json.Unmarshal(CleanJSONResponse(data), &raw); err != nil {
		return nil, fmt.Errorf("failed to parse search all response: %w", err)
	}

	var songs []Song
	for _, item := range raw.Songs.Data {
		songs = append(songs, ParseRawSong(item))
	}

	var albums []Album
	for _, item := range raw.Albums.Data {
		albums = append(albums, ParseRawAlbum(item))
	}

	var playlists []Playlist
	for _, item := range raw.Playlists.Data {
		playlists = append(playlists, ParseRawPlaylist(item))
	}

	var artists []ArtistMini
	for _, item := range raw.Artists.Data {
		id, _ := item["id"].(string)
		name, _ := item["title"].(string)
		image, _ := item["image"].(string)
		perma, _ := item["perma_url"].(string)
		role, _ := item["role"].(string)

		artists = append(artists, ArtistMini{
			ID:    id,
			Name:  cleanText(name),
			Role:  role,
			Type:  "artist",
			Image: crypto.FormatImageLinks(image),
			URL:   perma,
		})
	}

	return &SearchAllResult{
		Songs:     songs,
		Albums:    albums,
		Playlists: playlists,
		Artists:   artists,
	}, nil
}
