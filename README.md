# JioSaavn REST API (Golang)

An unofficial, high-performance REST API wrapper for JioSaavn built in Go.

## Features

- **High Performance:** Built with the Go Gin framework for ultra-low latency and low memory usage.
- **Audio Stream Resolution:** Resolves audio URLs across multiple quality levels (12kbps, 48kbps, 96kbps, 160kbps, 320kbps).
- **Search:** Endpoints for songs, albums, playlists, and artists.
- **Albums & Playlists:** Full album and playlist retrieval with tracks and artwork.
- **Lyrics:** Song lyrics fetching endpoint.
- **Trending & Charts:** Top charts and featured playlists.
- **Docker Support:** Multi-stage Dockerfile and docker-compose included.

## Quick Start

### Local Run

```bash
go mod download
go run ./cmd/api
```

The server runs on `http://localhost:8080`.

### Docker Run

```bash
docker build -t jiosaavn-api .
docker run -p 8080:8080 jiosaavn-api
```

## API Endpoints

### Search

- `GET /api/search?query={query}` - Autocomplete search across all categories
- `GET /api/search/songs?query={query}&page=1&limit=20` - Search songs
- `GET /api/search/albums?query={query}&page=1&limit=20` - Search albums
- `GET /api/search/playlists?query={query}&page=1&limit=20` - Search playlists
- `GET /api/search/artists?query={query}&page=1&limit=20` - Search artists

### Songs

- `GET /api/songs?id={id}` - Get song details by ID
- `GET /api/songs?link={url}` - Get song details by URL
- `GET /api/songs/:id` - Get song details by path ID

### Albums

- `GET /api/albums?id={id}` - Get album details by ID
- `GET /api/albums?link={url}` - Get album details by URL
- `GET /api/albums/:id` - Get album details by path ID

### Playlists

- `GET /api/playlists?id={id}` - Get playlist details by ID
- `GET /api/playlists?link={url}` - Get playlist details by URL
- `GET /api/playlists/:id` - Get playlist details by path ID

### Lyrics & Modules

- `GET /api/lyrics?id={song_id}` - Get lyrics
- `GET /api/lyrics/:id` - Get lyrics by path ID
- `GET /api/modules/charts` - Get top charts
- `GET /api/modules/playlists` - Get featured playlists

## Community & Credits

- **Developer:** [Kanha](https://github.com/Kanha-Team)
- **Contact Owner:** [@xoknha](https://t.me/xoknha)
- **Official Channel:** [@meowcloner](https://t.me/meowcloner)
- **Official Group:** [@meowclone](https://t.me/meowclone)
- **GitHub:** [Kanha-Team](https://github.com/Kanha-Team)

## License

MIT License. Maintained by [Kanha-Team](https://github.com/Kanha-Team).
