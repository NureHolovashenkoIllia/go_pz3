package routes

import (
	"go_pz3/internal/handlers"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// SetupRouter налаштовує всі маршрути для додатка
func SetupRouter(db *gorm.DB) *gin.Engine {
	router := gin.Default()

	// Створюємо екземпляри обробників
	actorHandler := handlers.NewActorHandler(db)
	performanceHandler := handlers.NewPerformanceHandler(db)
	roleHandler := handlers.NewRoleHandler(db)
	rehearsalHandler := handlers.NewRehearsalHandler(db)
	analyticsHandler := handlers.NewAnalyticsHandler(db)

	// Групуємо ендпоїнти, наприклад, під /api/v1
	api := router.Group("/api/v1")
	{
		// Маршрути для акторів
		actors := api.Group("/actors")
		{
			actors.POST("/", actorHandler.CreateActor)
			actors.GET("/", actorHandler.GetActors)
			actors.GET("/:id", actorHandler.GetActor)
			actors.PUT("/:id", actorHandler.UpdateActor)
			actors.DELETE("/:id", actorHandler.DeleteActor)
		}

		// Маршрути для вистав
		performances := api.Group("/performances")
		{
			performances.POST("/", performanceHandler.CreatePerformance)
			performances.GET("/", performanceHandler.GetPerformances)
			performances.GET("/:id", performanceHandler.GetPerformance)
			performances.PUT("/:id", performanceHandler.UpdatePerformance)
			performances.DELETE("/:id", performanceHandler.DeletePerformance)
		}

		// Маршрути для ролей
		roles := api.Group("/roles")
		{
			roles.POST("/", roleHandler.CreateRole)
			roles.GET("/", roleHandler.GetRoles)
			roles.GET("/:id", roleHandler.GetRole)
			roles.PUT("/:id", roleHandler.UpdateRole)
			roles.DELETE("/:id", roleHandler.DeleteRole)
		}

		// Маршрути для репетицій
		rehearsals := api.Group("/rehearsals")
		{
			rehearsals.POST("/", rehearsalHandler.CreateRehearsal)
			rehearsals.GET("/", rehearsalHandler.GetRehearsals)
			rehearsals.GET("/:id", rehearsalHandler.GetRehearsal)
			rehearsals.PUT("/:id", rehearsalHandler.UpdateRehearsal)
			rehearsals.DELETE("/:id", rehearsalHandler.DeleteRehearsal)
		}

		// Маршрути для аналітики
		analytics := api.Group("/analytics")
		{
			analytics.GET("/most-active-actor", analyticsHandler.GetMostActiveActor)
			analytics.GET("/least-active-actor", analyticsHandler.GetLeastActiveActor)
			analytics.GET("/performance-most-actors", analyticsHandler.GetPerformanceWithMostActors)
			analytics.GET("/most-frequent-performance", analyticsHandler.GetMostFrequentPerformance)
		}
	}

	return router
}
