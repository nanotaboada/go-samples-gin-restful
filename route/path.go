package route

const (
	// PlayersPathSegment is the base resource segment for player-related routes
	PlayersPathSegment = "players"

	// PlayersPath is the base path for all player endpoints
	PlayersPath = "/" + PlayersPathSegment

	// Common path params
	IDParam          = "id"
	SquadNumberParam = "squadnumber"

	// GetAllPath is the route path for retrieving all players
	GetAllPath = PlayersPath

	// GetByIDPath is the route path for retrieving a player by ID
	GetByIDPath = PlayersPath + "/:" + IDParam

	// GetBySquadNumberPath is the route path for retrieving a player by squad number
	GetBySquadNumberPath = PlayersPath + "/squadnumber/:" + SquadNumberParam

	// SwaggerPath is the route path for Swagger UI documentation
	SwaggerPath = "/swagger/*any"

	// HealthPath is the route path for health check endpoint
	HealthPath = "/health"
)
