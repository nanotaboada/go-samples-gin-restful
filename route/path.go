// Package route sets up the routing and middleware for Player-related endpoints.
//
// # Gin route parameter syntax
//
// Gin uses ":name" to declare a named path parameter that matches a single
// URL segment (no slashes). For example:
//
//	/players/:squadnumber  matches /players/10 → Param("squadnumber") == "10"
//
// The "*any" wildcard (used for the Swagger route) matches the remainder of
// the path including slashes, so "/swagger/*any" captures "/swagger/index.html"
// as well as "/swagger/".
//
// # Route priority
//
// Gin resolves path conflicts using a trie (prefix tree).  Static segments
// (e.g. "/players/squadnumber/") always take priority over dynamic ones
// (e.g. "/players/:squadnumber"), so there is no ambiguity between
// GetBySquadNumberPath and BySquadNumberPath despite both matching paths
// under /players/.
package route

const (
	// PlayersPathSegment is the base resource segment for player-related routes.
	PlayersPathSegment = "players"

	// PlayersPath is the canonical base path for all player endpoints.
	PlayersPath = "/" + PlayersPathSegment

	// PlayersPathTrailingSlash is registered as an alias so that clients
	// sending a trailing slash receive the same response instead of a 301 redirect.
	PlayersPathTrailingSlash = PlayersPath + "/"

	// IDParam is the route parameter name for the player's internal UUID.
	IDParam = "id"
	// SquadNumberParam is the route parameter name for the player's squad number.
	SquadNumberParam = "squadnumber"

	// GetAllPath is the route for listing all players and creating a new one.
	GetAllPath = PlayersPath

	// GetAllPathTrailingSlash is the trailing-slash alias for GetAllPath.
	GetAllPathTrailingSlash = PlayersPathTrailingSlash

	// GetByIDPath retrieves a player by its internal UUID (surrogate key).
	GetByIDPath = PlayersPath + "/:" + IDParam

	// GetBySquadNumberPath looks up a player by the user-facing squad number.
	// The "/squadnumber/" static prefix disambiguates from BySquadNumberPath.
	GetBySquadNumberPath = PlayersPath + "/squadnumber/:" + SquadNumberParam

	// BySquadNumberPath is used for PUT and DELETE where the squad number is
	// the mutable resource identifier in the URL.
	BySquadNumberPath = PlayersPath + "/:" + SquadNumberParam

	// SwaggerPath uses the "*any" wildcard so the Swagger UI handler receives
	// any sub-path under /swagger/ (static assets, index, JSON spec, etc.).
	SwaggerPath = "/swagger/*any"

	// HealthPath is the liveness probe endpoint for Docker / load balancers.
	HealthPath = "/health"
)
