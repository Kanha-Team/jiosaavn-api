// Copyright (c) 2026 Kanha-Team
// All rights reserved.

package saavn

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Kanha-Team/jiosaavn-api/pkg/crypto"
)

// GetCharts fetches top charts from JioSaavn
func (c *Client) GetCharts(ctx context.Context) ([]Chart, error) {
	data, err := c.Request(ctx, "content.getCharts", nil)
	if err != nil {
		return nil, err
	}

	var raw []map[string]any
	if err := json.Unmarshal(CleanJSONResponse(data), &raw); err != nil {
		return nil, fmt.Errorf("failed to parse charts response: %w", err)
	}

	var charts []Chart
	for _, item := range raw {
		id, _ := item["id"].(string)
		title, _ := item["title"].(string)
		itemType, _ := item["type"].(string)
		perma, _ := item["perma_url"].(string)
		image, _ := item["image"].(string)
		countVal := 0
		if cnt, ok := item["count"].(float64); ok {
			countVal = int(cnt)
		}

		charts = append(charts, Chart{
			ID:    id,
			Title: cleanText(title),
			Type:  itemType,
			Count: countVal,
			URL:   perma,
			Image: crypto.FormatImageLinks(image),
		})
	}

	return charts, nil
}

// GetFeaturedPlaylists fetches featured playlists from JioSaavn
func (c *Client) GetFeaturedPlaylists(ctx context.Context) ([]Playlist, error) {
	data, err := c.Request(ctx, "content.getFeaturedPlaylists", map[string]string{
		"p": "1",
		"n": "20",
	})
	if err != nil {
		return nil, err
	}

	var raw struct {
		Data []map[string]any `json:"data"`
	}
	if err := json.Unmarshal(CleanJSONResponse(data), &raw); err != nil {
		return nil, fmt.Errorf("failed to parse featured playlists: %w", err)
	}

	var playlists []Playlist
	for _, item := range raw.Data {
		playlists = append(playlists, ParseRawPlaylist(item))
	}

	return playlists, nil
}
