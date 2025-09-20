package admin

import (
	"fmt"
	"github.com/JIIL07/jcloud/internal/server/cookies"
	j "github.com/JIIL07/jcloud/pkg/json"
	"net/http"
	"os/exec"
	"runtime"
)

func HandleCmdExec(w http.ResponseWriter, r *http.Request) {
	store := cookies.GetSession(r, "admin")

	if store.IsNew {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("Unauthorized")) // nolint:errcheck
		return
	}

	var req j.Request
	req.Command = r.URL.Query().Get("command")

	output, err := ExecuteCommand(req.Command)
	if err != nil {
		j.RespondWithError(w, err)
		return
	}

	j.RespondWithJSON(w, output)

}

func ExecuteCommand(command string) (string, error) {
	var cmd *exec.Cmd

	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/C", command)
	} else {
		cmd = exec.Command("sh", "-c", command)
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("command execution failed: %w", err)
	}

	return string(output), nil
}
