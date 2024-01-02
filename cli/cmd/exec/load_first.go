package exec

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"pmon3/cli/proxy"
)

func loadFirst(execPath string, flags string) ([]string, error) {
	data, err := proxy.RunProcess([]string{"start", execPath, flags})
	if err != nil {
		return nil, err
	}

	var tbData []string
	_ = json.Unmarshal(data, &tbData)

	return tbData, nil
}

func getExecFile(args []string) (string, error) {
	execFile := args[0]
	_, err := os.Stat(execFile)
	if os.IsNotExist(err) {
		return "", fmt.Errorf("%s does not exist", execFile)
	}

	if path.IsAbs(execFile) {
		return execFile, nil
	}

	absPath, err := filepath.Abs(execFile)
	if err != nil {
		return "", fmt.Errorf("get file path error: %v", err.Error())
	}

	return absPath, nil
}
