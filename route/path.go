package route

const (
	// PlayersPathSegment is the base resource segment for player-related routes
	PlayersPathSegment = "players" // base resource segment (no slashes)

	// PlayersPath is the base path for all player endpoints
	PlayersPath = "/" + PlayersPathSegment

	// GetAllPath is the route path for retrieving all players
	GetAllPath = PlayersPath + "/"

	// GetByIDPath is the route path for retrieving a player by ID
	GetByIDPath = PlayersPath + "/:id"

	// GetBySquadNumberPath is the route path for retrieving a player by squad number
	GetBySquadNumberPath = PlayersPath + "/squadnumber/:squadnumber"

	// SwaggerPath is the route path for Swagger UI documentation
	SwaggerPath = "/swagger/*any"

	// HealthPath is the route path for health check endpoint
	HealthPath = "/health"
)
