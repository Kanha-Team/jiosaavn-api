// Copyright (c) 2026 Kanha-Team
// All rights reserved.

package saavn

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
)

// GetPlaylistByID fetches playlist details by list ID
func (c *Client) GetPlaylistByID(ctx context.Context, id string) (*Playlist, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return nil, fmt.Errorf("playlist id is required")
	}

	data, err := c.Request(ctx, "playlist.getDetails", map[string]string{
		"listid": id,
	})
	if err != nil {
		return nil, err
	}

	var raw map[string]any
	if err := json.Unmarshal(CleanJSONResponse(data), &raw); err != nil {
		return nil, fmt.Errorf("failed to parse playlist details: %w", err)
	}

	if _, ok := raw["id"]; !ok {
		return nil, fmt.Errorf("playlist not found")
	}

	playlist := ParseRawPlaylist(raw)
	return &playlist, nil
}

// GetPlaylistByLink fetches playlist details by URL or token
func (c *Client) GetPlaylistByLink(ctx context.Context, linkOrToken string) (*Playlist, error) {
	token := ExtractTokenFromURL(linkOrToken)
	if token == "" {
		return nil, fmt.Errorf("playlist token is required")
	}

	data, err := c.Request(ctx, "webapi.get", map[string]string{
		"token": token,
		"type":  "playlist",
	})
	if err != nil {
		return nil, err
	}

	var raw map[string]any
	if err := json.Unmarshal(CleanJSONResponse(data), &raw); err != nil {
		return nil, fmt.Errorf("failed to parse playlist by token: %w", err)
	}

	if _, ok := raw["id"]; !ok {
		return nil, fmt.Errorf("playlist not found")
	}

	playlist := ParseRawPlaylist(raw)
	return &playlist, nil
}
