// Copyright (c) 2026 Kanha-Team
// All rights reserved.

package saavn

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
)

// GetLyrics fetches lyrics for a given song ID
func (c *Client) GetLyrics(ctx context.Context, id string) (*Lyrics, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return nil, fmt.Errorf("song id is required")
	}

	data, err := c.Request(ctx, "lyrics.getLyrics", map[string]string{
		"lyrics_id": id,
	})
	if err != nil {
		return nil, err
	}

	var raw struct {
		Lyrics          string `json:"lyrics"`
		Snippet         string `json:"snippet"`
		LyricsCopyright string `json:"lyrics_copyright"`
		Status          string `json:"status"`
		Error           any    `json:"error"`
	}

	if err := json.Unmarshal(CleanJSONResponse(data), &raw); err != nil {
		return nil, fmt.Errorf("failed to parse lyrics: %w", err)
	}

	if raw.Status == "failure" || raw.Lyrics == "" {
		return nil, fmt.Errorf("lyrics not found for this song")
	}

	return &Lyrics{
		Lyrics:    cleanText(raw.Lyrics),
		Snippet:   cleanText(raw.Snippet),
		Copyright: cleanText(raw.LyricsCopyright),
	}, nil
}
