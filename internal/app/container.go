package app

import (
	"context"
	"os"
	"path/filepath"
	"time"

	"music/internal/common/middleware"
	"music/internal/modules/ai"
	"music/internal/modules/analytics"
	"music/internal/modules/auth"
	"music/internal/modules/catalog"
	"music/internal/modules/catalog/album"
	artist "music/internal/modules/catalog/artist"
	"music/internal/modules/catalog/genre"
	"music/internal/modules/catalog/track"
	"music/internal/modules/contribution"
	"music/internal/modules/covers"
	"music/internal/modules/creator"
	"music/internal/modules/dashboard"
	"music/internal/modules/follow"
	"music/internal/modules/gamification"
	"music/internal/modules/health"
	"music/internal/modules/history"
	"music/internal/modules/importcmd"
	"music/internal/modules/importcmd/acquisition"
	importsearch "music/internal/modules/importcmd/search"
	importworker "music/internal/modules/importcmd/worker"
	"music/internal/modules/ingestion"
	"music/internal/modules/ingestion/enrichment"
	"music/internal/modules/ingestion/finalization"
	"music/internal/modules/library"
	"music/internal/modules/lyrics"
	"music/internal/modules/media"
	"music/internal/modules/moderation"
	"music/internal/modules/notification"
	"music/internal/modules/player"
	"music/internal/modules/playlist"
	"music/internal/modules/queue"
	"music/internal/modules/reactions"
	"music/internal/modules/recommendation"
	"music/internal/modules/search"
	"music/internal/modules/social"
	"music/internal/modules/subscription"
	"music/internal/modules/video"
	"music/internal/platform/events"
	platformstorage "music/internal/platform/storage"
	"music/internal/platform/ws"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	opensearch "github.com/opensearch-project/opensearch-go"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type Container struct {
	SQLX  *sqlx.DB
	Bus   *events.Bus
	RDB   *redis.Client
	WSHub *ws.Hub

	// Lifecycle context for background goroutines managed by the container.
	// These are wired from the App's context and WaitGroup.
	enrichCtx  context.Context
	enrichDone context.CancelFunc

	Storage platformstorage.Storage

	AuthMW         gin.HandlerFunc
	OptionalAuthMW gin.HandlerFunc

	HealthHandler *health.Handler

	AuthService *auth.Service
	AuthHandler *auth.Handler

	MediaService *media.Service
	MediaHandler *media.Handler

	ArtistService *artist.Service
	ArtistHandler *artist.Handler

	AlbumService *album.Service
	AlbumHandler *album.Handler

	GenreService *genre.Service
	GenreHandler *genre.Handler

	TrackService *track.Service
	TrackHandler *track.Handler

	LyricsService *lyrics.Service
	LyricsHandler *lyrics.Handler

	PlaylistService *playlist.Service
	PlaylistHandler *playlist.Handler

	LibraryService *library.Service
	LibraryHandler *library.Handler

	QueueService *queue.Service
	QueueHandler *queue.Handler

	FollowService *follow.Service
	FollowHandler *follow.Handler

	RecommendationService recommendation.Service
	RecommendationHandler *recommendation.Handler

	TasteProfileService       *recommendation.TasteProfileService
	ProfileEventHandler       *recommendation.ProfileEventHandler
	OnboardingHandler         *recommendation.OnboardingHandler
	HomeFeedService           *recommendation.HomeFeedService
	HomeRecommendationService *recommendation.HomeRecommendationService
	DiscoverWeeklyService     *recommendation.DiscoverWeeklyService
	ListeningStatsService     *recommendation.ListeningStatsService

	HistoryService *history.Service
	HistoryHandler *history.Handler

	AnalyticsService *analytics.Service
	AnalyticsHandler *analytics.Handler

	NotificationService *notification.Service
	NotificationHandler *notification.Handler

	ModerationService *moderation.Service
	ModerationHandler *moderation.Handler

	SubscriptionService *subscription.Service
	SubscriptionHandler *subscription.Handler

	PlayerService *player.Service
	PlayerHandler *player.Handler

	IngestionService      *ingestion.Service
	IngestionHandler      *ingestion.Handler
	IngestionFinalization *finalization.Service

	SearchService *search.Service
	SearchHandler *search.Handler

	SocialService      *social.Service
	SocialHandler      *social.Handler
	SocialQueueEngine  *social.QueueEngine
	SocialStageManager *social.StageManager

	ReactionsService *reactions.Service
	ReactionsHandler *reactions.Handler

	GamificationService      *gamification.Service
	GamificationHandler      *gamification.Handler
	GamificationEventHandler *gamification.EventHandler

	CreatorService *creator.Service
	CreatorHandler *creator.Handler

	DashboardHandler *dashboard.Handler

	ContributionHandler *contribution.Handler

	ImportService      *importcmd.Service
	ImportHandler      *importcmd.Handler
	ImportSvc          *importcmd.ImportService
	ImportSearch       *importsearch.Service
	ImportAcquisition  *acquisition.Service
	ImportWorker       *importworker.ImportWorker
	ImportArtistSearch *importcmd.ArtistSearcher

	VideoService      *video.Service
	VideoHandler      *video.Handler
	VideoEventHandler *video.EventHandler

	AIHandler *ai.Handler

	MLServiceClient  *lyrics.MLServiceClient
	MLServiceHMACMW  gin.HandlerFunc
	OpenRouterClient *lyrics.OpenRouterClient

	CoverHandler *covers.Handler

	Enricher      *enrichment.Enricher
	LastFMClient  enrichment.LastFMClient
	SpotifyClient enrichment.SpotifyClient
	EnrichHandler *catalog.EnrichHandler
}

func NewContainer(a *App) *Container {
	sqlDB := stdlib.OpenDBFromPool(a.DB)
	sqlxDB := sqlx.NewDb(sqlDB, "pgx")

	logger := zap.L()
	bus := events.NewBus(logger)

	c := &Container{
		SQLX: sqlxDB,
		Bus:  bus,
		RDB:  a.Redis,

		// Lifecycle: derive from the app's cancellable context
		enrichCtx:  a.AddBackground(),
		enrichDone: a.BackgroundDone,
	}

	c.buildHealth(a)
	c.buildStorage(a)
	c.buildAuth(a)
	c.buildMedia(a)
	c.buildCatalog()
	c.buildLyrics(a)
	c.buildCovers(a)
	c.buildPlaylist()
	c.buildLibrary()
	c.buildQueue()
	c.buildFollow()
	c.buildRecommendation()
	c.buildRadio(a)
	c.buildHistory()
	c.buildTasteProfile(a)
	c.buildHomeFeed(a)
	c.buildAnalytics()
	c.buildNotification()
	c.buildSubscription()
	c.buildPlayer()
	c.buildIngestion(a)
	c.buildEnrich()
	c.buildSearch(a)
	c.buildSocial(a)
	c.buildModeration()
	c.buildReactions()
	c.buildGamification()
	c.buildCreator()
	c.buildDashboard()
	c.buildContribution()
	c.buildImport()
	c.buildAI(a)
	c.buildVideo()
	c.subscribeEvents()

	// Wire the transactional outbox store into the event bus for durable,
	// at-least-once event delivery. The relay runs every 5 seconds.
	outboxStore := events.NewSQLOutboxStore(sqlxDB)
	bus.SetOutboxStore(outboxStore, 5*time.Second)

	return c
}

func (c *Container) buildHealth(a *App) {
	c.HealthHandler = health.NewHandler(a.DB, a.Redis)
}

func (c *Container) buildStorage(a *App) {
	storageClient, err := platformstorage.New(a.ctx, platformstorage.Config{
		Driver: a.Config.Storage.Driver,
		Local: platformstorage.LocalConfig{
			BaseDir: a.Config.Storage.Local.BaseDir,
			BaseURL: a.Config.Storage.Local.BaseURL,
		},
		S3: platformstorage.S3Config{
			Bucket:          a.Config.Storage.S3.Bucket,
			Region:          a.Config.Storage.S3.Region,
			Endpoint:        a.Config.Storage.S3.Endpoint,
			AccessKeyID:     a.Config.Storage.S3.AccessKeyID,
			SecretAccessKey: a.Config.Storage.S3.SecretAccessKey,
			PublicBaseURL:   a.Config.Storage.S3.PublicBaseURL,
			UsePathStyle:    a.Config.Storage.S3.UsePathStyle,
			PresignURLs:     a.Config.Storage.S3.PresignURLs,
			PresignTTL:      a.Config.Storage.S3.PresignTTL,
		},
	})
	if err != nil {
		panic(err)
	}

	c.Storage = storageClient
}

func (c *Container) buildAuth(a *App) {
	tokenManager := auth.NewTokenManager(
		a.Config.Auth.JWTAccessSecret,
		a.Config.Auth.JWTRefreshSecret,
		a.Config.Auth.AccessTTL,
		a.Config.Auth.RefreshTTL,
	)

	authRepo := auth.NewRepository(a.DB)
	c.AuthService = auth.NewService(authRepo, tokenManager)
	c.AuthHandler = auth.NewHandler(c.AuthService, a.Validator)
	c.AuthMW = auth.AuthMiddleware(tokenManager)
	c.OptionalAuthMW = auth.OptionalAuthMiddleware(tokenManager)
}

func (c *Container) buildMedia(a *App) {
	mediaRepo := media.NewRepository(c.SQLX)

	c.MediaService = media.NewService(c.Storage, mediaRepo, media.Config{
		MaxFileSizeBytes:  200 << 20, // 200 MB general limit
		MaxImageSizeBytes: 10 << 20,
		MaxAudioSizeBytes: 60 << 20,
		MaxVideoSizeBytes: 200 << 20,
		AllowedImageMime: []string{
			"image/jpeg",
			"image/png",
			"image/webp",
		},
		AllowedAudioMime: []string{
			"audio/mpeg",
			"audio/ogg",
			"audio/flac",
			"audio/wav",
			"audio/x-wav",
			"audio/wave",
			"audio/mp4",
			"audio/aac",
			"application/octet-stream",
		},
		AllowedVideoMime: []string{
			"video/mp4",
			"video/webm",
			"video/quicktime",
			"video/x-msvideo",
			"video/x-matroska",
		},
		StorageProviderName: a.Config.Storage.Driver,
	})

	c.MediaHandler = media.NewHandler(c.MediaService)
}
func (c *Container) buildCatalog() {
	artistRepo := artist.NewRepository(c.SQLX)
	c.ArtistService = artist.NewService(artistRepo)
	c.ArtistHandler = artist.NewHandler(c.ArtistService)

	albumRepo := album.NewRepository(c.SQLX)
	c.AlbumService = album.NewService(albumRepo)
	c.AlbumHandler = album.NewHandler(c.AlbumService)

	genreRepo := genre.NewRepository(c.SQLX)
	c.GenreService = genre.NewService(genreRepo)
	c.GenreHandler = genre.NewHandler(c.GenreService)

	trackRepo := track.NewRepository(c.SQLX)
	c.TrackService = track.NewService(trackRepo)
	c.TrackHandler = track.NewHandler(c.TrackService)
}

func (c *Container) buildEnrich() {
	deezerClient := enrichment.NewDeezerClient()
	coverArtClient := enrichment.NewCoverArtClient()

	c.EnrichHandler = catalog.NewEnrichHandler(
		c.Enricher,
		c.TrackService,
		c.ArtistService,
		c.AlbumService,
		c.LyricsService,
		deezerClient,
		coverArtClient,
		c.LastFMClient,
		c.SpotifyClient,
		zap.L(),
	)
}

func (c *Container) buildLyrics(a *App) {
	lyricsRepo := lyrics.NewRepository(c.SQLX)
	c.LyricsService = lyrics.NewService(lyricsRepo)
	c.LyricsHandler = lyrics.NewHandler(c.LyricsService)

	c.MLServiceClient = lyrics.NewMLServiceClient(a.Config.MLService)
	c.MLServiceHMACMW = middleware.VerifyMLServiceWebhook(a.Config.MLService.WebhookHMACSecret)

	c.LyricsHandler.SetMLClient(c.MLServiceClient)

	c.OpenRouterClient = lyrics.NewOpenRouterClient(a.Config.OpenRouter)
	c.LyricsHandler.SetOpenRouterClient(c.OpenRouterClient)
}

func (c *Container) buildCovers(a *App) {
	sqlDB := stdlib.OpenDBFromPool(a.DB)
	c.CoverHandler = covers.NewHandler(sqlDB)
}

func (c *Container) buildPlaylist() {
	playlistRepo := playlist.NewRepository(c.SQLX)
	c.PlaylistService = playlist.NewService(playlistRepo)
	c.PlaylistHandler = playlist.NewHandler(c.PlaylistService)
}

func (c *Container) buildLibrary() {
	libraryRepo := library.NewRepository(c.SQLX)
	c.LibraryService = library.NewServiceWithPublisher(libraryRepo, c.Bus)
	c.LibraryHandler = library.NewHandler(c.LibraryService)
}

func (c *Container) buildQueue() {
	queueRepo := queue.NewRepository(c.SQLX)
	c.QueueService = queue.NewService(queueRepo)
	c.QueueHandler = queue.NewHandler(c.QueueService)
}

func (c *Container) buildFollow() {
	followRepo := follow.NewRepository(c.SQLX)
	c.FollowService = follow.NewService(followRepo)
	c.FollowHandler = follow.NewHandler(c.FollowService)
}

func (c *Container) buildRecommendation() {
	recommendationRepo := recommendation.NewRepository(c.SQLX)
	c.RecommendationService = recommendation.NewService(recommendationRepo)
	c.RecommendationHandler = recommendation.NewHandler(c.RecommendationService)
}

func (c *Container) buildRadio(a *App) {
	radioRepo := recommendation.NewRadioRepository(c.SQLX)
	mlClient := recommendation.NewSimilarityMLClient(a.Config.MLService)
	radioService := recommendation.NewRadioService(radioRepo, mlClient)
	c.RecommendationHandler.SetRadioService(radioService)
}

func (c *Container) buildHistory() {
	historyRepo := history.NewRepository(c.SQLX)
	c.HistoryService = history.NewServiceWithPublisher(historyRepo, c.Bus)
	c.HistoryHandler = history.NewHandler(c.HistoryService)
}

func (c *Container) buildTasteProfile(a *App) {
	tasteRepo := recommendation.NewTasteProfileRepository(c.SQLX)
	mlClient := recommendation.NewSimilarityMLClient(a.Config.MLService)
	c.TasteProfileService = recommendation.NewTasteProfileService(tasteRepo, c.HistoryService.GetRepo(), mlClient)
	c.ProfileEventHandler = recommendation.NewProfileEventHandler(c.TasteProfileService, c.HistoryService)
	c.OnboardingHandler = recommendation.NewOnboardingHandler(c.TasteProfileService)
}

func (c *Container) buildHomeFeed(a *App) {
	mlClient := recommendation.NewSimilarityMLClient(a.Config.MLService)
	c.HomeFeedService = recommendation.NewHomeFeedService(
		recommendation.NewRepository(c.SQLX),
		c.TasteProfileService,
		mlClient,
		c.RecommendationService,
	)
	c.RecommendationHandler.SetHomeFeedService(c.HomeFeedService)

	c.HomeRecommendationService = recommendation.NewHomeRecommendationService(
		c.TasteProfileService,
		mlClient,
		recommendation.NewRepository(c.SQLX),
	)
	c.RecommendationHandler.SetHomeRecommendationService(c.HomeRecommendationService)

	dwRepo := recommendation.NewDiscoverWeeklyRepository(c.SQLX)
	c.DiscoverWeeklyService = recommendation.NewDiscoverWeeklyService(
		c.TasteProfileService,
		mlClient,
		recommendation.NewRepository(c.SQLX),
		dwRepo,
	)
	c.RecommendationHandler.SetDiscoverWeeklyService(c.DiscoverWeeklyService)

	c.ListeningStatsService = recommendation.NewListeningStatsService(c.SQLX)
	c.RecommendationHandler.SetListeningStatsService(c.ListeningStatsService)
}

func (c *Container) buildAnalytics() {
	analyticsRepo := analytics.NewRepository(c.SQLX)
	c.AnalyticsService = analytics.NewService(analyticsRepo)
	c.AnalyticsHandler = analytics.NewHandler(c.AnalyticsService)
}

func (c *Container) buildNotification() {
	notificationRepo := notification.NewRepository(c.SQLX)
	c.NotificationService = notification.NewService(notificationRepo)
	c.NotificationHandler = notification.NewHandler(c.NotificationService)
}

func (c *Container) buildSubscription() {
	subscriptionRepo := subscription.NewRepository(c.SQLX)
	c.SubscriptionService = subscription.NewService(subscriptionRepo, c.Bus)
	c.SubscriptionHandler = subscription.NewHandler(c.SubscriptionService)
}
func (c *Container) buildPlayer() {
	playerRepo := player.NewRepository(c.SQLX)
	c.PlayerService = player.NewService(playerRepo, c.Storage, c.Bus)
	c.PlayerHandler = player.NewHandler(c.PlayerService)
}

func (c *Container) buildIngestion(a *App) {
	ingestionRepo := ingestion.NewRepository(c.SQLX)

	enrichCfg := enrichment.DefaultConfig()
	enrichCfg.LastFM.APIKey = a.Config.Enrichment.LastFmAPIKey
	enrichCfg.Spotify.ClientID = a.Config.Enrichment.SpotifyClientID
	enrichCfg.Spotify.ClientSecret = a.Config.Enrichment.SpotifyClientSecret

	mbClient := enrichment.NewMusicBrainzClient(enrichCfg.MusicBrainz)
	lfmClient := enrichment.NewLastFMClient(enrichCfg.LastFM)
	c.LastFMClient = lfmClient
	spotClient := enrichment.NewSpotifyClient(enrichCfg.Spotify)
	c.SpotifyClient = spotClient
	lrcClient := enrichment.NewLRCLibClient()
	mlClient := enrichment.NewMLEnrichmentClient(a.Config.MLService.BaseURL)

	enricher := enrichment.NewEnricher(mbClient, lfmClient, spotClient, lrcClient, mlClient, zap.L())
	c.Enricher = enricher

	c.IngestionFinalization = finalization.NewService(c.SQLX, c.Storage, zap.L())
	c.IngestionService = ingestion.NewService(c.Storage, ingestionRepo, enricher)
	c.IngestionHandler = ingestion.NewHandler(c.IngestionService)
	c.IngestionHandler.SetFinalizer(c.IngestionFinalization)

	cleanupSvc := ingestion.NewCleanupService(ingestionRepo, zap.L())
	c.IngestionHandler.SetCleanup(cleanupSvc)
	c.IngestionHandler.SetMLClient(c.MLServiceClient)
	c.IngestionHandler.SetStorage(c.Storage, a.Config.Storage.Driver)
	cleanupCtx := a.AddBackground()
	go func() {
		defer a.BackgroundDone()
		cleanupSvc.StartPeriodicCleanup(cleanupCtx, 1*time.Hour, 24*time.Hour)
	}()
}

func (c *Container) buildSearch(a *App) {
	var osClient *opensearch.Client

	if a.Config.OpenSearch.URL != "" {
		cfg := opensearch.Config{
			Addresses: []string{a.Config.OpenSearch.URL},
			Username:  a.Config.OpenSearch.Username,
			Password:  a.Config.OpenSearch.Password,
		}

		var err error
		osClient, err = opensearch.NewClient(cfg)
		if err != nil {
			zap.L().Warn("failed to initialize opensearch client, falling back to SQL search", zap.Error(err))
		}
	}

	searchRepo := search.NewRepository(c.SQLX)
	c.SearchService = search.NewService(osClient, searchRepo)
	c.SearchHandler = search.NewHandler(c.SearchService)
}

func (c *Container) buildModeration() {
	moderationRepo := moderation.NewRepository(c.SQLX)
	c.ModerationService = moderation.NewService(moderationRepo)
	c.ModerationHandler = moderation.NewHandler(c.ModerationService)
}

func (c *Container) buildSocial(a *App) {
	c.WSHub = ws.NewHub(a.Logger, a.Config.CORS.AllowedOrigins)
	wsCtx := a.AddBackground()
	go func() {
		defer a.BackgroundDone()
		if err := c.WSHub.Run(wsCtx); err != nil && err != context.Canceled {
			a.Logger.Error("ws hub exited with error", zap.Error(err))
		}
	}()

	socialRepo := social.NewRepository(c.SQLX)
	partyBroadcaster := social.NewPartyBroadcaster(c.WSHub)
	roomBroadcaster := social.NewRoomBroadcaster(c.WSHub)
	c.SocialService = social.NewServiceWithPlaylist(socialRepo, partyBroadcaster, roomBroadcaster, c.PlaylistService)

	c.SocialQueueEngine = social.NewQueueEngine(socialRepo, roomBroadcaster, a.Logger)
	c.SocialStageManager = social.NewStageManager(socialRepo, roomBroadcaster, c.Bus, a.Logger)
	c.SocialHandler = social.NewHandlerFull(c.SocialService, c.SocialQueueEngine, c.SocialStageManager)
}

func (c *Container) buildReactions() {
	reactionsRepo := reactions.NewRepository(c.SQLX)
	c.ReactionsService = reactions.NewService(reactionsRepo)
	c.ReactionsHandler = reactions.NewHandler(c.ReactionsService)
}

func (c *Container) buildGamification() {
	gamificationRepo := gamification.NewRepository(c.SQLX)
	c.GamificationService = gamification.NewService(gamificationRepo)
	c.GamificationHandler = gamification.NewHandler(c.GamificationService)
	c.GamificationEventHandler = gamification.NewEventHandler(c.GamificationService)
}

func (c *Container) buildCreator() {
	creatorRepo := creator.NewRepository(c.SQLX)
	c.CreatorService = creator.NewService(creatorRepo)
	c.CreatorHandler = creator.NewHandler(c.CreatorService)
}

func (c *Container) buildAI(a *App) {
	aiRepo := ai.NewRepository(c.SQLX)
	aiClient := ai.NewOpenAIClient(
		a.Config.AI.OpenAIEndpoint,
		a.Config.AI.OpenAIKey,
		a.Config.AI.EmbeddingModel,
	)
	c.AIHandler = ai.NewHandler(ai.NewService(aiRepo, aiClient, zap.L(), a.Config.AI.Enabled))
}

func (c *Container) buildVideo() {
	videoRepo := video.NewRepository(c.SQLX)
	c.VideoService = video.NewService(videoRepo, c.FollowService)
	commentRepo := video.NewCommentRepository(c.SQLX)
	c.VideoHandler = video.NewHandler(c.VideoService, commentRepo)
	c.VideoEventHandler = video.NewEventHandler(c.VideoService)
}

func (c *Container) buildDashboard() {
	c.DashboardHandler = dashboard.NewHandler(c.SQLX, c.RDB)
}

func (c *Container) buildContribution() {
	contributionRepo := contribution.NewRepository(c.SQLX)
	contributionSvc := contribution.NewService(contributionRepo, contribution.NewAIVerifier(), nil)
	c.ContributionHandler = contribution.NewHandler(contributionSvc)
}

func (c *Container) buildImport() {
	c.ImportService = importcmd.NewService(zap.L())
	downloadDir := os.Getenv("IMPORT_DOWNLOAD_DIR")
	if downloadDir == "" {
		downloadDir = filepath.Join(os.TempDir(), "music-imports")
	}

	proxy := os.Getenv("IMPORT_PROXY")

	// Build discovery providers
	discoveryProviders := []importsearch.Provider{
		importsearch.NewLocalProvider(c.SQLX),
		importsearch.NewiTunesProvider(),
		importsearch.NewDeezerProvider(),
	}

	lastfmKey := os.Getenv("LASTFM_API_KEY")
	if lastfmKey != "" {
		discoveryProviders = append(discoveryProviders, importsearch.NewLastFMProvider(lastfmKey))
	}

	spotifyID := os.Getenv("SPOTIFY_CLIENT_ID")
	spotifySecret := os.Getenv("SPOTIFY_CLIENT_SECRET")
	if spotifyID != "" && spotifySecret != "" {
		discoveryProviders = append(discoveryProviders, importsearch.NewSpotifyProvider(spotifyID, spotifySecret))
	}

	discoveryProviders = append(discoveryProviders, importsearch.NewMusicBrainzProvider())

	c.ImportSearch = importsearch.NewService(zap.L(), c.RDB, discoveryProviders)

	// Build acquisition layer
	c.ImportAcquisition = acquisition.NewService(zap.L(), proxy)

	// Build import worker
	dlSvc := importworker.NewDownloadService(proxy)
	c.ImportWorker = importworker.NewImportWorker(
		zap.L(),
		c.RDB,
		dlSvc,
		c.IngestionService,
		c.ImportAcquisition,
		downloadDir,
	)

	c.ImportSvc = importcmd.NewImportService(
		zap.L(),
		c.ImportService,
		c.IngestionService,
		c.ImportAcquisition,
		c.ImportSearch,
		c.ImportWorker,
		downloadDir,
	)

	// Build artist searcher (Deezer + MusicBrainz, no keys needed)
	c.ImportArtistSearch = importcmd.NewArtistSearcher(zap.L())

	c.ImportHandler = importcmd.NewHandler(c.ImportSvc, c.ImportArtistSearch)
}

func (c *Container) subscribeEvents() {
	analyticsEvents := analytics.NewEventHandler(c.AnalyticsService)
	notificationEvents := notification.NewEventHandler(c.NotificationService)
	recommendationEvents := recommendation.NewEventHandler(c.RecommendationService)
	historyEvents := history.NewEventHandler(c.HistoryService)

	c.Bus.Subscribe(events.EventTrackPlayed, analyticsEvents.OnTrackPlayed)
	c.Bus.Subscribe(events.EventTrackPlayed, historyEvents.OnTrackPlayed)
	c.Bus.Subscribe(events.EventTrackPlayed, recommendationEvents.OnTrackPlayed)
	c.Bus.Subscribe(events.EventTrackPlayed, c.VideoEventHandler.OnTrackPlayed)
	c.Bus.Subscribe(events.EventTrackPlayed, c.GamificationEventHandler.OnTrackPlayed)

	c.Bus.Subscribe(events.EventPlaybackSignalRecorded, c.ProfileEventHandler.OnPlaybackSignalRecorded)
	c.Bus.Subscribe(events.EventTrackLiked, c.ProfileEventHandler.OnTrackLiked)
	c.Bus.Subscribe(events.EventTrackLiked, c.GamificationEventHandler.OnTrackLiked)

	c.Bus.Subscribe(events.EventPlaylistCreated, analyticsEvents.OnPlaylistCreated)
	c.Bus.Subscribe(events.EventPlaylistCreated, notificationEvents.OnPlaylistCreated)

	c.Bus.Subscribe(events.EventSubscriptionPurchased, analyticsEvents.OnSubscriptionPurchased)
	c.Bus.Subscribe(events.EventSubscriptionPurchased, notificationEvents.OnSubscriptionPurchased)

	c.Bus.Subscribe(events.EventUserRegistered, analyticsEvents.OnUserRegistered)
	c.Bus.Subscribe(events.EventUserRegistered, notificationEvents.OnUserRegistered)
}
