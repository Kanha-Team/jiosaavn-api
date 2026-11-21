// Copyright (c) 2026 Kanha-Team
// All rights reserved.

package saavn

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
)

// GetSongByID fetches song details by its JioSaavn ID (e.g. rjkrTnma)
func (c *Client) GetSongByID(ctx context.Context, id string) (*Song, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return nil, fmt.Errorf("song id is required")
	}

	data, err := c.Request(ctx, "song.getDetails", map[string]string{
		"pids": id,
	})
	if err != nil {
		return nil, err
	}

	var raw struct {
		Songs []map[string]any `json:"songs"`
	}
	if err := json.Unmarshal(CleanJSONResponse(data), &raw); err != nil {
		return nil, fmt.Errorf("failed to parse song details: %w", err)
	}

	if len(raw.Songs) == 0 {
		return nil, fmt.Errorf("song not found")
	}

	song := ParseRawSong(raw.Songs[0])
	return &song, nil
}

// GetSongByLink fetches song details using a JioSaavn URL or token
func (c *Client) GetSongByLink(ctx context.Context, linkOrToken string) (*Song, error) {
	token := ExtractTokenFromURL(linkOrToken)
	if token == "" {
		return nil, fmt.Errorf("song token is required")
	}

	data, err := c.Request(ctx, "webapi.get", map[string]string{
		"token": token,
		"type":  "song",
	})
	if err != nil {
		return nil, err
	}

	var raw struct {
		Songs []map[string]any `json:"songs"`
	}
	if err := json.Unmarshal(CleanJSONResponse(data), &raw); err != nil {
		return nil, fmt.Errorf("failed to parse song by token: %w", err)
	}

	if len(raw.Songs) == 0 {
		return nil, fmt.Errorf("song not found")
	}

	song := ParseRawSong(raw.Songs[0])
	return &song, nil
}

// GetSongsByIDs fetches multiple songs by comma-separated IDs
func (c *Client) GetSongsByIDs(ctx context.Context, ids []string) ([]Song, error) {
	if len(ids) == 0 {
		return nil, nil
	}

	validIDs := make([]string, 0, len(ids))
	for _, id := range ids {
		trimmed := strings.TrimSpace(id)
		if trimmed != "" {
			validIDs = append(validIDs, trimmed)
		}
	}
	if len(validIDs) == 0 {
		return nil, nil
	}

	data, err := c.Request(ctx, "song.getDetails", map[string]string{
		"pids": strings.Join(validIDs, ","),
	})
	if err != nil {
		return nil, err
	}

	var raw struct {
		Songs []map[string]any `json:"songs"`
	}
	if err := json.Unmarshal(CleanJSONResponse(data), &raw); err != nil {
		return nil, fmt.Errorf("failed to parse multiple songs details: %w", err)
	}

	var songs []Song
	for _, item := range raw.Songs {
		songs = append(songs, ParseRawSong(item))
	}
	return songs, nil
}
