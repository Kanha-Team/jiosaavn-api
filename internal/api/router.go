// Copyright (c) 2026 Kanha-Team
// All rights reserved.

package api

import (
	"github.com/Kanha-Team/jiosaavn-api/pkg/saavn"
	"github.com/gin-gonic/gin"
)

// SetupRouter sets up Gin engine routes and middleware
func SetupRouter(client *saavn.Client) *gin.Engine {
	r := gin.New()

	// Global Middlewares
	r.Use(gin.Recovery())
	r.Use(RequestLogger())
	r.Use(CORSMiddleware())

	h := NewHandler(client)

	// General routes
	r.GET("/", h.RootHandler)
	r.GET("/api/", h.RootHandler)
	r.GET("/health", h.HealthCheck)

	// API Group
	apiGroup := r.Group("/api")
	{
		apiGroup.GET("", h.RootHandler)

		// Search
		search := apiGroup.Group("/search")
		{
			search.GET("", h.SearchAll)
			search.GET("/songs", h.SearchSongs)
			search.GET("/albums", h.SearchAlbums)
			search.GET("/playlists", h.SearchPlaylists)
			search.GET("/artists", h.SearchArtists)
		}

		// Songs
		apiGroup.GET("/songs", h.GetSong)
		apiGroup.GET("/songs/:id", h.GetSong)

		// Albums
		apiGroup.GET("/albums", h.GetAlbum)
		apiGroup.GET("/albums/:id", h.GetAlbum)

		// Playlists
		apiGroup.GET("/playlists", h.GetPlaylist)
		apiGroup.GET("/playlists/:id", h.GetPlaylist)

		// Lyrics
		apiGroup.GET("/lyrics", h.GetLyrics)
		apiGroup.GET("/lyrics/:id", h.GetLyrics)

		// Modules / Trending
		modules := apiGroup.Group("/modules")
		{
			modules.GET("/charts", h.GetCharts)
			modules.GET("/playlists", h.GetFeaturedPlaylists)
		}
	}

	return r
}
