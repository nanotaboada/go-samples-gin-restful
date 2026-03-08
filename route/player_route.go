// Package route sets up the routing and middleware for Player-related endpoints.
package route

import (
	"fmt"
	"time"

	"github.com/gin-contrib/cache"
	"github.com/gin-contrib/cache/persistence"
	"github.com/gin-gonic/gin"
	"github.com/nanotaboada/go-samples-gin-restful/controller"
)

// RegisterPlayerRoutes wires all player endpoints to the router.
//
// In Gin the first argument to router.GET/POST/etc. is the path pattern and
// the remaining arguments are handler functions chained left to right.  When
// multiple handlers are listed (e.g. a middleware followed by the real handler)
// Gin calls them in order; any handler can stop the chain early by calling
// context.Abort().
//
// # Caching strategy
//
// Read endpoints (GET) are wrapped with cache.CachePage, which:
//  1. Computes a cache key from the full request URL.
//  2. On first hit, calls the real handler and stores the response in the
//     InMemoryStore with the given TTL (1 hour here).
//  3. On subsequent hits within the TTL, replays the cached response without
//     calling the handler — no DB round-trip.
//
// Write endpoints (POST, PUT, DELETE) are wrapped with ClearCache, which
// deletes the affected cache keys before delegating to the real handler, so
// the next GET always fetches fresh data.
func RegisterPlayerRoutes(router *gin.Engine, controller *controller.PlayerController, store *persistence.InMemoryStore) {
	// Register routes for /players (without trailing slash)
	router.GET(GetAllPath, cache.CachePage(store, time.Hour, controller.GetAll))
	router.POST(GetAllPath, ClearCache(store, controller.Post))

	// Register alias routes for /players/ (with trailing slash).
	// Gin does not automatically redirect trailing-slash variants; registering
	// them explicitly avoids 301 redirects that some clients don't follow.
	router.GET(GetAllPathTrailingSlash, cache.CachePage(store, time.Hour, controller.GetAll))
	router.POST(GetAllPathTrailingSlash, ClearCache(store, controller.Post))

	// GET by squad number (user-facing identifier)
	router.GET(GetBySquadNumberPath, cache.CachePage(store, time.Hour, controller.GetBySquadNumber))

	// GET by internal UUID (surrogate key)
	router.GET(GetByIDPath, cache.CachePage(store, time.Hour, controller.GetByID))

	// PUT and DELETE use squad number as the mutable resource identifier
	router.PUT(BySquadNumberPath, ClearCache(store, controller.Put))
	router.DELETE(BySquadNumberPath, ClearCache(store, controller.Delete))
}

// ClearCache is a middleware factory that invalidates cached responses before
// a mutating handler (POST, PUT, DELETE) runs.
//
// It returns a gin.HandlerFunc — Gin's standard handler type, which is just
// func(*gin.Context).  The closure captures `store` and `handler`, so each
// call to ClearCache produces an independent middleware instance.
//
// Cache keys are derived from URL paths using cache.CreateKey (the same
// function cache.CachePage uses internally), ensuring the keys match exactly.
// The squad-number-specific key is only added when the route has a
// :squadnumber parameter (PUT / DELETE), not for collection-level mutations.
func ClearCache(store persistence.CacheStore, handler gin.HandlerFunc) gin.HandlerFunc {
	return func(context *gin.Context) {
		squadNumber := context.Param(SquadNumberParam)

		// Always bust the collection-level cache (the GET /players response).
		keys := []string{
			cache.CreateKey(PlayersPath),
			cache.CreateKey(PlayersPathTrailingSlash),
		}
		// Also bust the individual resource cache when operating on a specific player.
		if squadNumber != "" {
			keys = append(keys, cache.CreateKey(fmt.Sprintf("%s/squadnumber/%s", PlayersPath, squadNumber)))
		}
		for _, key := range keys {
			// Ignore delete errors: a cache-miss on delete is harmless.
			_ = store.Delete(key)
		}
		handler(context)
	}
}
