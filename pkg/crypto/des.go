// Copyright (c) 2026 Kanha-Team
// All rights reserved.

package crypto

import (
	"crypto/des"
	"encoding/base64"
	"fmt"
	"regexp"
	"strings"
)

var (
	// JioSaavn DES encryption key
	desKey = []byte("38346591")

	// Regex to match audio quality suffixes in decrypted CDN URLs
	qualityRegex = regexp.MustCompile(`_(12|48|96|160|320)\.(mp4|m4a|mp3)`)
)

// QualityLink represents an audio file or image URL with its quality level
type QualityLink struct {
	Quality string `json:"quality"`
	URL     string `json:"url"`
}

// DecryptURL decrypts JioSaavn's encrypted_media_url using DES-ECB
func DecryptURL(encryptedBase64 string) (string, error) {
	if strings.TrimSpace(encryptedBase64) == "" {
		return "", fmt.Errorf("empty encrypted url")
	}

	ciphertext, err := base64.StdEncoding.DecodeString(encryptedBase64)
	if err != nil {
		return "", fmt.Errorf("failed to decode base64: %w", err)
	}

	block, err := des.NewCipher(desKey)
	if err != nil {
		return "", fmt.Errorf("failed to create DES cipher: %w", err)
	}

	bs := block.BlockSize()
	if len(ciphertext)%bs != 0 {
		return "", fmt.Errorf("ciphertext length %d is not a multiple of block size %d", len(ciphertext), bs)
	}

	plaintext := make([]byte, len(ciphertext))
	for i := 0; i < len(ciphertext); i += bs {
		block.Decrypt(plaintext[i:i+bs], ciphertext[i:i+bs])
	}

	// Remove PKCS5 padding
	plaintext = pkcs5Unpad(plaintext, bs)
	return string(plaintext), nil
}

// pkcs5Unpad removes PKCS5/PKCS7 padding bytes
func pkcs5Unpad(src []byte, blockSize int) []byte {
	length := len(src)
	if length == 0 {
		return src
	}
	pad := int(src[length-1])
	if pad <= 0 || pad > blockSize || pad > length {
		return src
	}
	// Verify padding
	for i := length - pad; i < length; i++ {
		if src[i] != byte(pad) {
			return src
		}
	}
	return src[:length-pad]
}

// GetAudioDownloadLinks generates download links for all supported bitrates
func GetAudioDownloadLinks(decryptedURL string) []QualityLink {
	if decryptedURL == "" {
		return []QualityLink{}
	}

	qualities := []string{"12kbps", "48kbps", "96kbps", "160kbps", "320kbps"}
	bitrates := []string{"12", "48", "96", "160", "320"}

	var links []QualityLink
	for i, q := range qualities {
		bitrate := bitrates[i]
		url := qualityRegex.ReplaceAllString(decryptedURL, fmt.Sprintf("_%s.$2", bitrate))
		links = append(links, QualityLink{
			Quality: q,
			URL:     url,
		})
	}
	return links
}

// FormatImageLinks returns thumbnail images in 50x50, 150x150, and 500x500 qualities
func FormatImageLinks(rawURL string) []QualityLink {
	if rawURL == "" {
		return []QualityLink{}
	}

	// Saavn CDN image dimensions
	dims := []struct {
		Quality string
		Dim     string
	}{
		{"50x50", "50x50.jpg"},
		{"150x150", "150x150.jpg"},
		{"500x500", "500x500.jpg"},
	}

	re := regexp.MustCompile(`(50x50|150x150|500x500)\.(jpg|png|jpeg)`)
	var links []QualityLink
	for _, d := range dims {
		imgURL := re.ReplaceAllString(rawURL, d.Dim)
		links = append(links, QualityLink{
			Quality: d.Quality,
			URL:     imgURL,
		})
	}
	return links
}
