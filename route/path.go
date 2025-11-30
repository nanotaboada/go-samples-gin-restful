package route

const (
	// PlayersPathSegment is the base resource segment for player-related routes
	PlayersPathSegment = "players"

	// PlayersPath is the base path for all player endpoints
	PlayersPath = "/" + PlayersPathSegment

	// PlayersPathTrailingSlash is the base path with trailing slash for alias routes
	PlayersPathTrailingSlash = PlayersPath + "/"

	// Common path params
	IDParam          = "id"
	SquadNumberParam = "squadnumber"

	// GetAllPath is the route path for retrieving all players
	GetAllPath = PlayersPath

	// GetAllPathTrailingSlash is the alias route path with trailing slash
	GetAllPathTrailingSlash = PlayersPathTrailingSlash

	// GetByIDPath is the route path for retrieving a player by ID
	GetByIDPath = PlayersPath + "/:" + IDParam

	// GetBySquadNumberPath is the route path for retrieving a player by squad number
	GetBySquadNumberPath = PlayersPath + "/squadnumber/:" + SquadNumberParam

	// SwaggerPath is the route path for Swagger UI documentation
	SwaggerPath = "/swagger/*any"

	// HealthPath is the route path for health check endpoint
	HealthPath = "/health"
)
