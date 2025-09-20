package routes

import (
	"github.com/JIIL07/jcloud/internal/server/admin"
	"github.com/JIIL07/jcloud/internal/server/handlers"
	"github.com/JIIL07/jcloud/internal/server/middleware"
	"github.com/JIIL07/jcloud/internal/server/static"
	"github.com/JIIL07/jcloud/internal/server/storage"
	"github.com/JIIL07/jcloud/internal/server/utils"
	"github.com/gorilla/mux"
	"net/http"
	"time"
)

func SetupRouter(b *static.Files, s *storage.Storage) *mux.Router {
	router := mux.NewRouter()

	router.Use(middleware.StorageMiddleware(s))
	router.HandleFunc("/", handlers.RootHandler).Methods(http.MethodGet).Name("root")
	router.HandleFunc("/static/{filename}", b.BinaryHandler).Methods(http.MethodGet)
	router.HandleFunc("/check", CheckerHandler).Methods(http.MethodPost, http.MethodGet)

	api := router.PathPrefix("/api/v1").Name("api").Subrouter()
	api.Use(middleware.StorageMiddleware(s))
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
	files.HandleFunc("/upload", handlers.AddFileHandler).Methods(http.MethodPost).Name("add-file")

	files.HandleFunc("/list", handlers.ListFilesHandler).Methods(http.MethodGet).Name("list-files")
	files.HandleFunc("/images", handlers.ImageGalleryHandler).Methods(http.MethodGet).Name("list-images")

	currentFile := files.PathPrefix("/{filename}").Name("current-file").Subrouter()
	currentFile.Use(middleware.UserMiddleware)

	currentFile.HandleFunc("/versions", handlers.AddFileVersionHandler).Methods(http.MethodPost)
	currentFile.HandleFunc("/versions", handlers.GetFileVersionsHandler).Methods(http.MethodGet)
	currentFile.HandleFunc("/versions/{version}", handlers.GetFileVersionHandler).Methods(http.MethodGet)
	currentFile.HandleFunc("/versions/last", handlers.GetLastFileVersionHandler).Methods(http.MethodGet)
	currentFile.HandleFunc("/versions/{version}", handlers.DeleteFileVersionHandler).Methods(http.MethodDelete)
	currentFile.HandleFunc("/versions", handlers.DeleteFileVersionsHandler).Methods(http.MethodDelete)
	currentFile.HandleFunc("/restore", handlers.RestoreFileToVersionHandler).Methods(http.MethodGet)
	currentFile.HandleFunc("/history", handlers.GetFileHistoryHandler).Methods(http.MethodGet)

	currentFile.HandleFunc("/metadata", handlers.UpdateMetadataHandler).Methods(http.MethodPatch).Name("update-metadata")
	currentFile.HandleFunc("/share", handlers.ShareFileHandler).Methods(http.MethodPost).Name("share-file")
	currentFile.HandleFunc("/permissions", handlers.FilePermissionsHandler).Methods(http.MethodGet).Name("file-permissions")
	currentFile.HandleFunc("/permissions", handlers.UpdatePermissionsHandler).Methods(http.MethodPatch).Name("update-permissions")
	currentFile.HandleFunc("/info", handlers.FileInfoHandler).Methods(http.MethodGet).Name("file-info")
	currentFile.HandleFunc("/partial-update", handlers.PartialUpdateHandler).Methods(http.MethodPatch).Name("partial-update")
	currentFile.HandleFunc("/hash-sum", handlers.HashSumHandler).Methods(http.MethodGet).Name("hash-sum")
	currentFile.HandleFunc("/download", handlers.DownloadFileHandler).Methods(http.MethodGet).Name("download-file")
	currentFile.HandleFunc("/delete", handlers.DeleteFileHandler).Methods(http.MethodDelete).Name("delete-file")
	currentFile.HandleFunc("/data", handlers.FileDataHandler).Methods(http.MethodPost).Name("file-data")

	private := router.PathPrefix("/admin").Subrouter()
	private.HandleFunc("/admin", admin.AuthHandler).Methods(http.MethodGet)
	private.HandleFunc("/admin/checkadmin", admin.CheckHandler).Methods(http.MethodGet)
	private.HandleFunc("/all-users", admin.AllUsersHandler).Methods(http.MethodGet)
	private.HandleFunc("/sql", admin.HandleSQLQuery).Methods(http.MethodGet)
	private.HandleFunc("/cmd", admin.HandleCmdExec).Methods(http.MethodGet)

	return router
}

func CheckerHandler(w http.ResponseWriter, r *http.Request) {
	s := utils.ProvideStorage(r, w)

	//f := &storage.File{
	//	UserID:        1,
	//	LastVersionID: 0,
	//	Metadata: storage.FileMetadata{
	//		Name:        "test",
	//		Extension:   "txt",
	//		Size:        len("test file"),
	//		HashSum:     jhash.Hash([]byte("test file")),
	//		Description: "test description",
	//	},
	//	Data:       []byte("test file"),
	//	Status:     "upload",
	//	CreatedAt:  time.Now(),
	//	ModifiedAt: time.Now(),
	//}

	v := storage.FileVersion{
		FileID:      1,
		UserID:      1,
		Version:     1,
		FullVersion: true,
		Delta:       []byte("test file"),
		ChangeType:  "upload",
		CreatedAt:   time.Now(),
	}

	err := s.AddFileVersion(v)
	if err != nil {
		http.Error(w, "Failed to add file version", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("File saved successfully")) // nolint:errcheck
}
