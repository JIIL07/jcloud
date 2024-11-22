package routes

import (
	"github.com/JIIL07/jcloud/internal/server/admin"
	"github.com/JIIL07/jcloud/internal/server/handlers"
	"github.com/JIIL07/jcloud/internal/server/middleware"
	"github.com/JIIL07/jcloud/internal/server/static"
	"github.com/JIIL07/jcloud/internal/server/storage"
	"github.com/JIIL07/jcloud/internal/server/utils"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func SetupRouter(b *static.Files, s *storage.Storage) *mux.Router {
	router := mux.NewRouter()

	router.Use(middleware.StorageMiddleware(s))
	router.HandleFunc("/", handlers.RootHandler).Methods(http.MethodGet).Name("root")
	router.HandleFunc("/static/{filename}", b.BinaryHandler).Methods(http.MethodGet)
	router.HandleFunc("/check", CheckerHandler).Methods(http.MethodGet)

	api := router.PathPrefix("/api/v1").Name("api").Subrouter()
	api.Use(middleware.StorageMiddleware(s))
	api.HandleFunc("/login", handlers.LoginHandler).Methods(http.MethodPost).Name("login")
	api.HandleFunc("/logout", handlers.LogoutHandler).Methods(http.MethodGet).Name("logout")
	api.HandleFunc("/healthcheck", handlers.HealthCheckHandler).Methods(http.MethodGet).Name("healthcheck")

	user := api.PathPrefix("/user/{user}").Name("user").Subrouter()
	user.Use(middleware.UserMiddleware)
	user.HandleFunc("", handlers.CurrentUserHandler).Methods(http.MethodGet).Name("current-user")
	user.HandleFunc("/profile", handlers.ProfileHandler).Methods(http.MethodGet).Name("profile")

	f := user.PathPrefix("/files").Name("files").Subrouter()
	f.Use(middleware.UserMiddleware)
	f.HandleFunc("/list", handlers.ListFilesHandler).Methods(http.MethodGet).Name("list-files")
	f.HandleFunc("/upload", handlers.AddFileHandler).Methods(http.MethodPost).Name("add-file")
	f.HandleFunc("/download", handlers.DownloadFileHandler).Methods(http.MethodGet).Name("download-file")
	f.HandleFunc("/delete", handlers.DeleteFileHandler).Methods(http.MethodDelete).Name("delete-file")

	private := router.PathPrefix("/admin").Subrouter()
	private.HandleFunc("/admin", admin.AuthHandler).Methods(http.MethodGet)
	private.HandleFunc("/admin/checkadmin", admin.CheckHandler).Methods(http.MethodGet)
	private.HandleFunc("/sql", admin.HandleSQLQuery).Methods(http.MethodGet)
	private.HandleFunc("/cmd", admin.HandleCmdExec).Methods(http.MethodGet)

	return router
}

func CheckerHandler(w http.ResponseWriter, r *http.Request) {
	s := utils.ProvideStorage(r, w)

	updates := map[string]interface{}{
		"email":        "aida@example.com",
		"password":     "newhashedpassword",
		"hashprotocol": "md5",
	}

	err := s.UpdateUser("aida", updates)
	if err != nil {
		log.Fatalf("Failed to update user: %v", err)
	}

}
