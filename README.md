# Song Library API

[![GoDoc](https://godoc.org/github.com/par1ram/song-library?status.svg)](https://godoc.org/github.com/par1ram/song-library)
[![Swagger](https://img.shields.io/badge/OpenAPI-Swagger-blue)](http://localhost:8000/swagger/index.html)

This API provides a way to manage a song library, allowing users to add, update, delete, and retrieve song information. It integrates with an external API to fetch additional song details.

## Table of Contents

* [Introduction](#introduction)
* [Getting Started](#getting-started)
    * [Prerequisites](#prerequisites)
    * [Installation](#installation)
* [Usage](#usage)
    * [Endpoints](#endpoints)
* [Contributing](#contributing)

## Introduction

This project is a RESTful API built with Go and uses PostgreSQL for data storage.  It leverages the Chi router for routing, `sqlc` for type-safe SQL queries, and Swagger for API documentation.  The API interacts with an external service to retrieve details like release date, lyrics, and a related link.

## Getting Started

### Prerequisites

* Go
* PostgreSQL
* An external API providing song details (as configured in `.env`)


### Installation

1. Clone the repository: `git clone https://github.com/par1ram/song-library.git` *(Replace with your repo URL)*
2. Navigate to the project directory: `cd song-library`
3. Create and configure your `.env` file with the necessary environment variables (see `.env.example` for a template). This includes:
    * `PORT` -  the server port
    * `DB_URL` - PostgreSQL connection string
    * `EXTERNAL_API_URL` - the external API endpoint
4. Run database migrations: `goose -dir sql/schema up`
5. Install dependencies: `go mod tidy`
6. Build and run the server: `go run main.go`


## Usage

### Endpoints

The API documentation is available via Swagger at `http://localhost:8000/swagger/index.html` after starting the server.  Use this documentation to explore the available endpoints and their parameters.



The core functionality includes:

* **Retrieving songs:** Retrieve all songs, filtered by group name and song name, with pagination options. Use `/songs/getAll-with-filters` (POST). More advanced filtering with release date is available using `/songs/get-with-filters-and-pagination` (POST).
* **Adding songs:** Add new songs, automatically fetching additional data from an external API. Use  `/songs/add` (POST) with `group_name` and `song_name` in the request body.
* **Updating songs:** Modify existing song information (group name, song name, text, release date, link). Use `/songs/update` (PUT) with the song `id` as a query parameter and the updated fields in the body.
* **Deleting songs:** Remove songs from the library by ID. Use `/songs/delete` (DELETE) with the song  `id` as a query parameter.
* **Getting song verses:**  Retrieve song verses by ID with pagination.  Use `/songs/get-song-verses` (GET) with the song `id`, `limit`, and `offset` as query parameters.


## Contributing

Contributions are welcome! Please feel free to submit issues and pull requests.

