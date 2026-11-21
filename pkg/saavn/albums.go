// Copyright (c) 2026 Kanha-Team
// All rights reserved.

package saavn

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
)

// GetAlbumByID fetches album details by album ID
func (c *Client) GetAlbumByID(ctx context.Context, id string) (*Album, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return nil, fmt.Errorf("album id is required")
	}

	data, err := c.Request(ctx, "content.getAlbumDetails", map[string]string{
		"albumid": id,
	})
	if err != nil {
		return nil, err
	}

	var raw map[string]any
	if err := json.Unmarshal(CleanJSONResponse(data), &raw); err != nil {
		return nil, fmt.Errorf("failed to parse album details: %w", err)
	}

	if _, ok := raw["id"]; !ok {
		return nil, fmt.Errorf("album not found")
	}

	album := ParseRawAlbum(raw)
	return &album, nil
}

// GetAlbumByLink fetches album details by its URL or token
func (c *Client) GetAlbumByLink(ctx context.Context, linkOrToken string) (*Album, error) {
	token := ExtractTokenFromURL(linkOrToken)
	if token == "" {
		return nil, fmt.Errorf("album token is required")
	}

	data, err := c.Request(ctx, "webapi.get", map[string]string{
		"token": token,
		"type":  "album",
	})
	if err != nil {
		return nil, err
	}

	var raw map[string]any
	if err := json.Unmarshal(CleanJSONResponse(data), &raw); err != nil {
		return nil, fmt.Errorf("failed to parse album by token: %w", err)
	}

	if _, ok := raw["id"]; !ok {
		return nil, fmt.Errorf("album not found")
	}

	album := ParseRawAlbum(raw)
	return &album, nil
}
