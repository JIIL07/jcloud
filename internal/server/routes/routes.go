package routes

import (
	"encoding/base64"
	"fmt"
	"github.com/JIIL07/jcloud/internal/server/admin"
	"github.com/JIIL07/jcloud/internal/server/handlers"
	"github.com/JIIL07/jcloud/internal/server/middleware"
	"github.com/JIIL07/jcloud/internal/server/static"
	"github.com/JIIL07/jcloud/internal/server/storage"
	"github.com/JIIL07/jcloud/internal/server/utils"
	"github.com/gorilla/mux"
	"net/http"
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

	currentFile.HandleFunc("/metadata", handlers.UpdateMetadataHandler).Methods(http.MethodPatch).Name("update-metadata")
	currentFile.HandleFunc("/history", handlers.FileHistoryHandler).Methods(http.MethodGet).Name("file-history")
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

	files, err := s.GetImageFiles(1)
	if err != nil {
		http.Error(w, "Failed to retrieve images", http.StatusInternalServerError)
		return
	}

	html := "<html><body><h1>Image Gallery</h1><div style='display: flex; flex-wrap: wrap;'>"
	for _, file := range files {
		imageDataURL := fmt.Sprintf("data:image/%s;base64,%s", file.Metadata.Extension, base64.StdEncoding.EncodeToString(file.Data))
		html += fmt.Sprintf(
			"<div style='margin: 10px;'><img src='%s' alt='%s' style='width: 200px; height: auto;'></div>",
			imageDataURL,
			file.Metadata.Name,
		)
	}
	html += "</div></body></html>"

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(html)) // nolint:errcheck
}
