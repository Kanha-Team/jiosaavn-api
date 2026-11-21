// Copyright (c) 2026 Kanha-Team
// All rights reserved.

package crypto

import (
	"testing"
)

func TestDecryptURL(t *testing.T) {
	// Sample encrypted string from JioSaavn
	sampleEnc := "ID2ieOjCrwfgWvL5sXl4B1ImC5QfbsDy803Qf6rBxBcKiDi7MTsK3bjCem46t5reztaAkfUtgBtTWZuU/IFhrRw7tS9a8Gtq"
	url, err := DecryptURL(sampleEnc)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if url == "" {
		t.Fatal("expected non-empty decrypted url")
	}

	links := GetAudioDownloadLinks(url)
	if len(links) != 5 {
		t.Fatalf("expected 5 quality links, got %d", len(links))
	}

	found320 := false
	for _, l := range links {
		if l.Quality == "320kbps" {
			found320 = true
			break
		}
	}
	if !found320 {
		t.Fatal("320kbps link missing")
	}
}
