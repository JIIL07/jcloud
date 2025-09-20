package routes

import (
	"net/http"

	"github.com/JIIL07/jcloud/internal/server/admin"
	"github.com/JIIL07/jcloud/internal/server/handlers"
	"github.com/JIIL07/jcloud/internal/server/middleware"
	"github.com/JIIL07/jcloud/internal/server/static"
	"github.com/JIIL07/jcloud/internal/server/storage"
	"github.com/JIIL07/jcloud/internal/server/types"
	"github.com/gorilla/mux"
)

type RouterConfig struct {
	StorageService    types.StorageService
	FileService       types.FileService
	ResponseService   types.ResponseService
	ValidationService types.ValidationService
	StaticFiles       *static.Files
}

func SetupRouter(config *RouterConfig) *mux.Router {
	router := mux.NewRouter()
	fileHandler := handlers.NewFileHandler(config.FileService, config.ResponseService)
	fileVersionHandler := handlers.NewFileVersionHandler(config.StorageService, config.ResponseService)

	storageConcrete := config.StorageService.(*handlers.StorageAdapter).Storage
	router.Use(middleware.StorageMiddleware(storageConcrete))

	router.HandleFunc("/", handlers.RootHandler).Methods(http.MethodGet).Name("root")
	router.HandleFunc("/static/{filename}", config.StaticFiles.BinaryHandler).Methods(http.MethodGet)
	router.HandleFunc("/check", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("OK"))
	}).Methods(http.MethodPost, http.MethodGet)

	api := router.PathPrefix("/api/v1").Name("api").Subrouter()
	api.Use(middleware.StorageMiddleware(storageConcrete))
	api.HandleFunc("/login", handlers.LoginHandler).Methods(http.MethodPost).Name("login")
	api.HandleFunc("/logout", handlers.LogoutHandler).Methods(http.MethodGet).Name("logout")
	api.HandleFunc("/healthcheck", handlers.HealthCheckHandler).Methods(http.MethodGet).Name("healthcheck")

	user := api.PathPrefix("/user/{user}").Name("user").Subrouter()
	user.Use(middleware.UserMiddleware)
	user.HandleFunc("", handlers.CurrentUserHandler).Methods(http.MethodGet).Name("current-user")
	user.HandleFunc("/profile", handlers.ProfileHandler).Methods(http.MethodGet).Name("user-profile")
	user.HandleFunc("/update", handlers.UpdateUserHandler).Methods(http.MethodPost).Name("update-user")
	user.HandleFunc("/delete", handlers.DeleteUserHandler).Methods(http.MethodDelete).Name("delete-user")

	files := user.PathPrefix("/files").Name("files").Subrouter()
	files.Use(middleware.UserMiddleware)

	files.HandleFunc("/upload", fileHandler.AddFileHandler).Methods(http.MethodPost).Name("add-file")
	files.HandleFunc("/list", fileHandler.ListFilesHandler).Methods(http.MethodGet).Name("list-files")
	files.HandleFunc("/images", fileHandler.ImageGalleryHandler).Methods(http.MethodGet).Name("list-images")

	currentFile := files.PathPrefix("/{filename}").Name("current-file").Subrouter()
	currentFile.Use(middleware.UserMiddleware)

	currentFile.HandleFunc("/versions", fileVersionHandler.AddFileVersionHandler).Methods(http.MethodPost)
	currentFile.HandleFunc("/versions", fileVersionHandler.GetFileVersionsHandler).Methods(http.MethodGet)
	currentFile.HandleFunc("/versions/{version}", fileVersionHandler.GetFileVersionHandler).Methods(http.MethodGet)
	currentFile.HandleFunc("/versions/last", fileVersionHandler.GetLastFileVersionHandler).Methods(http.MethodGet)
	currentFile.HandleFunc("/versions/{version}", fileVersionHandler.DeleteFileVersionHandler).Methods(http.MethodDelete)
	currentFile.HandleFunc("/versions", fileVersionHandler.DeleteFileVersionsHandler).Methods(http.MethodDelete)
	currentFile.HandleFunc("/restore", fileVersionHandler.RestoreFileToVersionHandler).Methods(http.MethodGet)

	currentFile.HandleFunc("/metadata", fileHandler.UpdateMetadataHandler).Methods(http.MethodPatch).Name("update-metadata")
	currentFile.HandleFunc("/info", fileHandler.FileInfoHandler).Methods(http.MethodGet).Name("file-info")
	currentFile.HandleFunc("/hash-sum", fileHandler.HashSumHandler).Methods(http.MethodGet).Name("hash-sum")
	currentFile.HandleFunc("/download", fileHandler.DownloadFileHandler).Methods(http.MethodGet).Name("download-file")
	currentFile.HandleFunc("/delete", fileHandler.DeleteFileHandler).Methods(http.MethodDelete).Name("delete-file")

	private := router.PathPrefix("/admin").Subrouter()
	private.HandleFunc("/admin", admin.AuthHandler).Methods(http.MethodGet)
	private.HandleFunc("/admin/checkadmin", admin.CheckHandler).Methods(http.MethodGet)
	private.HandleFunc("/all-users", admin.AllUsersHandler).Methods(http.MethodGet)
	private.HandleFunc("/sql", admin.HandleSQLQuery).Methods(http.MethodGet)
	private.HandleFunc("/cmd", admin.HandleCmdExec).Methods(http.MethodGet)

	return router
}

func LegacySetupRouter(b *static.Files, s *storage.Storage) *mux.Router {
	responseService := handlers.NewResponseService()
	validationService := handlers.NewValidationService()
	storageAdapter := handlers.NewStorageAdapter(s)
	fileService := handlers.NewFileService(storageAdapter, validationService, responseService)

	config := &RouterConfig{
		StorageService:    storageAdapter,
		FileService:       fileService,
		ResponseService:   responseService,
		ValidationService: validationService,
		StaticFiles:       b,
	}

	return SetupRouter(config)
}
