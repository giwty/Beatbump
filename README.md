
# Beatbump

This project is a continuation of [Beatbump by @snuffyDev](https://github.com/snuffyDev/Beatbump).

An alternative frontend for YouTube Music created using Svelte/SvelteKit, and Golang.

<p align="center">
	  <a href="https://www.gnu.org/licenses/agpl-3.0.en.html">
    <img alt="License: AGPLv3" src="https://shields.io/badge/License-AGPL%20v3-blue.svg">
  </a>
  <a href="https://github.com/humanetech-community/awesome-humane-tech">
    <img alt="Awesome Humane Tech" src="https://raw.githubusercontent.com/humanetech-community/awesome-humane-tech/main/humane-tech-badge.svg?sanitize=true">
  </a>
</p>

## Features

- Automix for continued listening
- No ads
- Background play on mobile devices
- Search for artists, playlists, songs, and albums
  - Note that all playback is audio only (for now)
- Local playlist management
  - Stored in-browser with IndexedDB
  - Can save songs individually under 'Favorites'
  - Peer-to-Peer data synchronization (using WebRTC)
- Group Sessions
  - Achieved with WebRTC in a [mesh](https://en.wikipedia.org/wiki/Mesh_networking)
- Uses a custom wrapper around the YouTube Music API
- Download songs for offline listening \
...and so much more!


## Privacy

All data is stored locally on your device. Data synchronization is done using PeerJS, which uses WebRTC for a
Peer-to-Peer connection between browsers.

## Contributing

Contributions are welcomed

## Running Beatbump

The recommended way to run Beatbump is via `docker-compose`, as Beatbump depends on the `invidious-companion` service to generate valid YouTube sessions and bypass blocking.

In the `docker-compose.yaml` file, you can choose to use the pre built docker image of Beatbump or build it from source.

### Prerequisites
- [Docker](https://docs.docker.com/get-docker/)
- [Docker Compose](https://docs.docker.com/compose/install/)

### Setup via Docker Compose

1. Clone the repository:
   ```bash
   git clone https://github.com/giwty/Beatbump.git
   cd Beatbump
   ```

2. Start the services:
   ```bash
   docker-compose up -d
   ```
   This will start both the Beatbump container and the `invidious-companion` container.

3. Access Beatbump at `http://localhost:8080` or `http://app.localhost:8080`.

## Configuration

To avoid the constant breakage of Beatbump due to YouTube API changes, Beatbump now uses `invidious-companion` to seamlessly generate sessions that pass YouTube's checks, allowing you to listen to music ad-free without encountering playback errors.

The provided `docker-compose.yaml` is already pre-configured to handle the connection between Beatbump and the companion service out of the box.

If you choose to run Beatbump and `invidious-companion` separately, you must provide the following environment variables to the Beatbump container:
- `COMPANION_URL`: The URL of the `invidious-companion` service (e.g., `http://companion:8282`).
- `COMPANION_SECRET_KEY`: The secret key matching the `SERVER_SECRET_KEY` set in `invidious-companion`.

## Downloads (New capability)

Note - the download capability was developed with AI.

Beatbump now has a built in downloader that can download songs, playlists, and albums.

There are different download options for you to try - 

- Download playlists - new download option in the playlists page.
- Download single song - new download option in the songs page.
- Download song mix - choose a "seed" song and specify how many songs you want to download in addition.
- Ongoing listening download - Beatbump will automatically download songs as you listen to them. 

Beatbump uses a local SQLite database (path defined by the `BEATBUMP_DB_PATH` environment variable) to store the download tasks. 
Currently, the database stores:
- **Group Tasks:** Represents high-level tasks like downloading a playlist, an album, or processing an ongoing listening session.
- **Song Tasks:** Represents individual song downloads attached to a group task, including metadata like title, artist, album, and thumbnail.
- **Settings:** Stores application configurations.


**Settings:**
downloads will be enabled once you defined the download path in the settings page.  Optionally you can enable ongoing downloads capability.
- **Download Path:** The directory where music gets saved (mapped to `/downloads` in docker-compose).
- **Ongoing Downloads:** Toggle to automatically save songs to your library as you listen.



## Project Inspirations

- [Invidious](https://github.com/iv-org/invidious) - a privacy focused alternative YouTube front end.
- [yt-dlp](https://github.com/yt-dlp/yt-dlp) -  A feature-rich command-line audio/video downloader 
