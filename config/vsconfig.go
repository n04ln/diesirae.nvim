package config

import (
	"encoding/json"

	"github.com/NoahOrberg/diesirae.nvim/util"
	"github.com/neovim/go-client/nvim"
)

/* e.g.
 * g:diesirae_config = {
 *   "build": {
 *     ".go":{
 *       "buildcommands": ["go", "build", "-o", "*bin*", "*source*"],
 *       "runcommands": ["*bin*"]
 *     },
 *     ".cpp":{
 *       "buildcommands": ["g++", "-o", "*bin*", "*source*"],
 *       "runcommands": ["*bin*"]
 *     },
 *     ".c":{
 *       "buildcommands": ["gcc", "-o", "*bin*", "*source*"],
 *       "runcommands": ["*bin*"]
 *     },
 *     ".py":{
 *       "buildcommands": [],
 *       "runcommands": ["python3", "*bin*"]
 *     }
 *   }
 * }
 */

type VSConfig struct {
	Commands map[string]struct {
		Language     string   `json:"language"`
		BuildCommand []string `json:"build_command"`
		ExecCommand  []string `json:"exec_command"`
	} `json:"commands"`
}

func GetVSConfig(v *nvim.Nvim) (*VSConfig, error) {
	vsConfig := &VSConfig{}
	// 雑に取得
	x := new(interface{})
	if err := v.Var(util.ConfigVarName, x); err != nil {
		return vsConfig, nil
	}
	// 雑にマーシャル
	y, err := json.Marshal(x)
	if err != nil {
		return nil, err
	}
	// 丁寧にアンマーシャル
	err = json.Unmarshal(y, &vsConfig)

	return vsConfig, err
}
