package link

import (
	"encoding/json"
	"strings"

	msg "github.com/aziontech/azion-cli/messages/init"
	"github.com/aziontech/azion-cli/pkg/contracts"
	"github.com/aziontech/azion-cli/utils"
)

func (cmd *LinkCmd) createTemplateAzion(info *LinkInfo) error {

	err := cmd.Mkdir(info.PathWorkingDir+"/azion", 0755) // 0755 is the permission mode for the new directories
	if err != nil {
		return msg.ErrorFailedCreatingAzionDirectory
	}

	azionJson := &contracts.AzionApplicationOptions{
		Name:     info.Name,
		Env:      "production",
		Template: strings.ToLower(info.Preset),
		Mode:     strings.ToLower(info.Mode),
		Prefix:   "",
	}
	azionJson.Function.Name = "__DEFAULT__"
	azionJson.Function.File = "./out/worker.js"
	azionJson.Function.Args = "./azion/args.json"
	azionJson.Domain.Name = "__DEFAULT__"
	azionJson.Application.Name = "__DEFAULT__"
	azionJson.Origin.Name = "__DEFAULT__"
	azionJson.RtPurge.PurgeOnPublish = true

	return cmd.createJsonFile(azionJson, info)

}

func (cmd *LinkCmd) createJsonFile(options *contracts.AzionApplicationOptions, info *LinkInfo) error {
	data, err := json.MarshalIndent(options, "", "  ")
	if err != nil {
		return msg.ErrorUnmarshalAzionFile
	}

	err = cmd.WriteFile(info.PathWorkingDir+"/azion/azion.json", data, 0644)
	if err != nil {
		return utils.ErrorInternalServerError
	}
	return nil
}
