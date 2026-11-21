// Copyright (c) 2026 Kanha-Team
// All rights reserved.

package api

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/Kanha-Team/jiosaavn-api/pkg/saavn"
	"github.com/gin-gonic/gin"
)

// Handler holds the JioSaavn API client
type Handler struct {
	client *saavn.Client
}

// NewHandler initializes a new API handler
func NewHandler(client *saavn.Client) *Handler {
	return &Handler{client: client}
}

// sendSuccess sends a standard JSON success response
func sendSuccess(c *gin.Context, data any) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    data,
	})
}

// sendError sends a standard JSON error response
func sendError(c *gin.Context, status int, message string) {
	c.JSON(status, gin.H{
		"success": false,
		"message": message,
	})
}

// RootHandler handles GET /
func (h *Handler) RootHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"name":          "JioSaavn API",
		"status":        "online",
		"author":        "Kanha-Team",
		"github":        "https://github.com/Kanha-Team/jiosaavn-api",
		"version":       "1.0.0",
		"documentation": "/docs",
		"endpoints": gin.H{
			"search_all":         "/api/search?query={query}",
			"search_songs":       "/api/search/songs?query={query}&page=1&limit=20",
			"search_albums":      "/api/search/albums?query={query}&page=1&limit=20",
			"search_playlists":   "/api/search/playlists?query={query}&page=1&limit=20",
			"search_artists":     "/api/search/artists?query={query}&page=1&limit=20",
			"song_by_id":         "/api/songs?id={id}",
			"song_by_link":       "/api/songs?link={url}",
			"album_by_id":        "/api/albums?id={id}",
			"album_by_link":      "/api/albums?link={url}",
			"playlist_by_id":     "/api/playlists?id={id}",
			"playlist_by_link":   "/api/playlists?link={url}",
			"lyrics":             "/api/lyrics?id={song_id}",
			"trending_charts":    "/api/modules/charts",
			"featured_playlists": "/api/modules/playlists",
		},
	})
}

// HealthCheck handles GET /health
func (h *Handler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"uptime": "healthy",
	})
}

// SearchAll handles GET /api/search
func (h *Handler) SearchAll(c *gin.Context) {
	query := strings.TrimSpace(c.Query("query"))
	if query == "" {
		sendError(c, http.StatusBadRequest, "query parameter is required")
		return
	}

	result, err := h.client.SearchAll(c.Request.Context(), query)
	if err != nil {
		sendError(c, http.StatusInternalServerError, err.Error())
		return
	}

	sendSuccess(c, result)
}

// SearchSongs handles GET /api/search/songs
func (h *Handler) SearchSongs(c *gin.Context) {
	query := strings.TrimSpace(c.Query("query"))
	if query == "" {
		sendError(c, http.StatusBadRequest, "query parameter is required")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	result, err := h.client.SearchSongs(c.Request.Context(), query, page, limit)
	if err != nil {
		sendError(c, http.StatusInternalServerError, err.Error())
		return
	}

	sendSuccess(c, result)
}

// SearchAlbums handles GET /api/search/albums
func (h *Handler) SearchAlbums(c *gin.Context) {
	query := strings.TrimSpace(c.Query("query"))
	if query == "" {
		sendError(c, http.StatusBadRequest, "query parameter is required")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	result, err := h.client.SearchAlbums(c.Request.Context(), query, page, limit)
	if err != nil {
		sendError(c, http.StatusInternalServerError, err.Error())
		return
	}

	sendSuccess(c, result)
}

// SearchPlaylists handles GET /api/search/playlists
func (h *Handler) SearchPlaylists(c *gin.Context) {
	query := strings.TrimSpace(c.Query("query"))
	if query == "" {
		sendError(c, http.StatusBadRequest, "query parameter is required")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	result, err := h.client.SearchPlaylists(c.Request.Context(), query, page, limit)
	if err != nil {
		sendError(c, http.StatusInternalServerError, err.Error())
		return
	}

	sendSuccess(c, result)
}

// SearchArtists handles GET /api/search/artists
func (h *Handler) SearchArtists(c *gin.Context) {
	query := strings.TrimSpace(c.Query("query"))
	if query == "" {
		sendError(c, http.StatusBadRequest, "query parameter is required")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	result, err := h.client.SearchArtists(c.Request.Context(), query, page, limit)
	if err != nil {
		sendError(c, http.StatusInternalServerError, err.Error())
		return
	}

	sendSuccess(c, result)
}

// GetSong handles GET /api/songs and GET /api/songs/:id
func (h *Handler) GetSong(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		id = c.Query("id")
	}
	link := c.Query("link")

	if id == "" && link == "" {
		sendError(c, http.StatusBadRequest, "either 'id' or 'link' query parameter is required")
		return
	}

	var song *saavn.Song
	var err error

	if id != "" {
		song, err = h.client.GetSongByID(c.Request.Context(), id)
	} else {
		song, err = h.client.GetSongByLink(c.Request.Context(), link)
	}

	if err != nil {
		sendError(c, http.StatusNotFound, err.Error())
		return
	}

	sendSuccess(c, song)
}

// GetAlbum handles GET /api/albums and GET /api/albums/:id
func (h *Handler) GetAlbum(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		id = c.Query("id")
	}
	link := c.Query("link")

	if id == "" && link == "" {
		sendError(c, http.StatusBadRequest, "either 'id' or 'link' query parameter is required")
		return
	}

	var album *saavn.Album
	var err error

	if id != "" {
		album, err = h.client.GetAlbumByID(c.Request.Context(), id)
	} else {
		album, err = h.client.GetAlbumByLink(c.Request.Context(), link)
	}

	if err != nil {
		sendError(c, http.StatusNotFound, err.Error())
		return
	}

	sendSuccess(c, album)
}

// GetPlaylist handles GET /api/playlists and GET /api/playlists/:id
func (h *Handler) GetPlaylist(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		id = c.Query("id")
	}
	link := c.Query("link")

	if id == "" && link == "" {
		sendError(c, http.StatusBadRequest, "either 'id' or 'link' query parameter is required")
		return
	}

	var playlist *saavn.Playlist
	var err error

	if id != "" {
		playlist, err = h.client.GetPlaylistByID(c.Request.Context(), id)
	} else {
		playlist, err = h.client.GetPlaylistByLink(c.Request.Context(), link)
	}

	if err != nil {
		sendError(c, http.StatusNotFound, err.Error())
		return
	}

	sendSuccess(c, playlist)
}

// GetLyrics handles GET /api/lyrics and GET /api/lyrics/:id
func (h *Handler) GetLyrics(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		id = c.Query("id")
	}

	if id == "" {
		sendError(c, http.StatusBadRequest, "id parameter is required")
		return
	}

	lyrics, err := h.client.GetLyrics(c.Request.Context(), id)
	if err != nil {
		sendError(c, http.StatusNotFound, err.Error())
		return
	}

	sendSuccess(c, lyrics)
}

// GetCharts handles GET /api/modules/charts
func (h *Handler) GetCharts(c *gin.Context) {
	charts, err := h.client.GetCharts(c.Request.Context())
	if err != nil {
		sendError(c, http.StatusInternalServerError, err.Error())
		return
	}

	sendSuccess(c, charts)
}

// GetFeaturedPlaylists handles GET /api/modules/playlists
func (h *Handler) GetFeaturedPlaylists(c *gin.Context) {
	playlists, err := h.client.GetFeaturedPlaylists(c.Request.Context())
	if err != nil {
		sendError(c, http.StatusInternalServerError, err.Error())
		return
	}

	sendSuccess(c, playlists)
}
