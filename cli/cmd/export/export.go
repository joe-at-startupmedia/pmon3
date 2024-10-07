package export

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
	"pmon3/cli"
	"pmon3/cli/cmd/base"
	"pmon3/pmond/model"
)

type flags struct {
	format  string
	orderBy string
}

var flag flags

var Cmd = &cobra.Command{
	Use:   "export",
	Short: "Export Process Configuration",
	Run: func(cmd *cobra.Command, args []string) {
		base.OpenSender()
		defer base.CloseSender()
		Export(flag)
	},
}

func init() {
	Cmd.Flags().StringVarP(&flag.format, "format", "f", "json", "the format to export")
	Cmd.Flags().StringVarP(&flag.format, "order", "o", "json", "the field by which to order by")
}

func jsonPrettyPrint(in string) string {
	var out bytes.Buffer
	err := json.Indent(&out, []byte(in), "", "  ")
	if err != nil {
		return in
	}
	return out.String()
}

func tomlPrettyPrint(ac *model.ProcessConfig) string {
	out := new(bytes.Buffer)
	encoder := toml.NewEncoder(out)
	if err := encoder.Encode(ac); err != nil {
		cli.Log.Fatal(err)
	}
	return out.String()
}

func yamlPrettyPrint(ac *model.ProcessConfig) string {
	out := new(bytes.Buffer)
	encoder := yaml.NewEncoder(out)
	if err := encoder.Encode(ac); err != nil {
		cli.Log.Fatal(err)
	}
	return out.String()
}

func Export(f flags) {
	exportString, err := GetExportString(f.format, f.orderBy)
	if err != nil {
		base.OutputError(err.Error())
	} else {
		fmt.Println(exportString)
	}
}

func GetExportString(format string, orderBy string) (string, error) {
	sent := base.SendCmd("export", orderBy)
	newCmdResp := base.GetResponse(sent)
	if len(newCmdResp.GetError()) > 0 {
		return "", errors.New(newCmdResp.GetError())
	}
	jsonOutput := newCmdResp.GetValueStr()

	var ac model.ProcessConfig
	if format == "toml" || format == "yaml" {
		err := json.Unmarshal([]byte(jsonOutput), &ac)
		if err != nil {
			return "", err
		}
	}

	var exportString string
	switch format {
	case "toml":
		exportString = tomlPrettyPrint(&ac)
	case "yaml":
		exportString = yamlPrettyPrint(&ac)
	case "json":
		exportString = jsonPrettyPrint(jsonOutput)
	default:

		return "", errors.New("accepted formats: json, toml or yaml")
	}
	return exportString, nil
}
