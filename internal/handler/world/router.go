package world

import "github.com/gin-gonic/gin"

type WorldModuleHandlers struct {
	StoryHandler     *StoryHandler
	CategoryHandler  *CategoryHandler
	LocationHandler  *LocationHandler
	FactionHandler   *FactionHandler
	CharacterHandler *CharacterHandler
	PlotHandler      *PlotHandler
	ItemHandler      *ItemHandler
}

func RegisterWorldRoutes(rg *gin.RouterGroup, authMiddleware gin.HandlerFunc, h *WorldModuleHandlers) {
	stories := rg.Group("/stories", authMiddleware)
	{
		// 1. Stories (Tác phẩm & Chuyển sang Chính thức)
		stories.POST("", h.StoryHandler.CreateStory)
		stories.GET("", h.StoryHandler.GetStories)
		stories.GET("/:id", h.StoryHandler.GetStoryDetail)
		stories.PUT("/:id", h.StoryHandler.UpdateStory)
		stories.DELETE("/:id", h.StoryHandler.DeleteStory)
		stories.POST("/:id/publish", h.StoryHandler.PublishStory)

		// 2. Categories (Từ điển hằng số & Smart Recheck)
		stories.POST("/:id/categories", h.CategoryHandler.CreateCategory)
		stories.GET("/:id/categories", h.CategoryHandler.GetCategories)
		stories.POST("/:id/categories/:category_id/items", h.CategoryHandler.CreateCategoryItem)
		stories.POST("/:id/recheck", h.CategoryHandler.RecheckValue)

		// 3. Locations (Bản đồ / Địa danh phân cấp cây)
		stories.POST("/:id/locations", h.LocationHandler.CreateLocation)
		stories.GET("/:id/locations", h.LocationHandler.GetLocations)
		stories.GET("/:id/locations/:location_id", h.LocationHandler.GetLocationByID)
		stories.PUT("/:id/locations/:location_id", h.LocationHandler.UpdateLocation)
		stories.DELETE("/:id/locations/:location_id", h.LocationHandler.DeleteLocation)

		// 4. Factions (Thế lực / Tông môn / Gia tộc)
		stories.POST("/:id/factions", h.FactionHandler.CreateFaction)
		stories.GET("/:id/factions", h.FactionHandler.GetFactions)
		stories.GET("/:id/factions/:faction_id", h.FactionHandler.GetFactionByID)
		stories.PUT("/:id/factions/:faction_id", h.FactionHandler.UpdateFaction)
		stories.DELETE("/:id/factions/:faction_id", h.FactionHandler.DeleteFaction)

		// 5. Characters (Nhân vật & Phát hiện từ mới)
		stories.POST("/:id/characters", h.CharacterHandler.CreateCharacter)
		stories.GET("/:id/characters", h.CharacterHandler.GetCharacters)
		stories.GET("/:id/characters/:character_id", h.CharacterHandler.GetCharacterByID)
		stories.PUT("/:id/characters/:character_id", h.CharacterHandler.UpdateCharacter)
		stories.DELETE("/:id/characters/:character_id", h.CharacterHandler.DeleteCharacter)

		// 6. Plots (Tình tiết & Timeline)
		stories.POST("/:id/plots", h.PlotHandler.CreatePlot)
		stories.GET("/:id/plots", h.PlotHandler.GetPlots)
		stories.GET("/:id/plots/:plot_id", h.PlotHandler.GetPlotByID)
		stories.PUT("/:id/plots/:plot_id", h.PlotHandler.UpdatePlot)
		stories.DELETE("/:id/plots/:plot_id", h.PlotHandler.DeletePlot)

		// 7. Items (Vật phẩm & Phẩm giai)
		stories.POST("/:id/items", h.ItemHandler.CreateItem)
		stories.GET("/:id/items", h.ItemHandler.GetItems)
		stories.GET("/:id/items/:item_id", h.ItemHandler.GetItemByID)
		stories.PUT("/:id/items/:item_id", h.ItemHandler.UpdateItem)
		stories.DELETE("/:id/items/:item_id", h.ItemHandler.DeleteItem)
	}
}
