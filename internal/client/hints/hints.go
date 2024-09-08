package hints

import (
	"fmt"
	"github.com/JIIL07/jcloud/internal/config"
)

type Hint struct {
	Message string   `json:"message"`
	Hint    []string `json:"hint"`
}

type HintCondition string

const (
	UnknownCommand         HintCondition = "unknownCommand"
	FlagNotRecognized      HintCondition = "flagNotRecognized"
	NoArgsProvided         HintCondition = "noArgsProvided"
	MissingFileOrDirectory HintCondition = "missingFileOrDirectory"
	MissingRequiredFlag    HintCondition = "missingRequiredFlag"

	EmptyPath            HintCondition = "emptyPath"
	LoginRequired        HintCondition = "loginRequired"
	NoFilesMatched       HintCondition = "noFilesMatched"
	AllFlagMissing       HintCondition = "allFlagMissing"
	NoModifiedFiles      HintCondition = "noModifiedFiles"
	DryRunWithoutFiles   HintCondition = "dryRunWithoutFiles"
	ExcludeNotMatched    HintCondition = "excludeNotMatched"
	UnmodifiedWithUpdate HintCondition = "unmodifiedWithUpdate"
)

func DisplayHint(command string, condition HintCondition, c *config.ClientConfig) string {
	hintKey := fmt.Sprintf("%s.%s", command, condition)
	if enabled, ok := c.CmdHints.Config[hintKey]; !ok || !enabled {
		return ""
	}

	if cmdHints, exists := c.CmdHints.Hint[hintKey]; exists {
		hintMessage := cmdHints.Message + "\n"
		for _, h := range cmdHints.Hints {
			hintMessage += fmt.Sprintf("hint: %s\n", h)
		}
		return hintMessage
	}

	return ""
}
