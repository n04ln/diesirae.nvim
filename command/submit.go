package command

import "github.com/neovim/go-client/nvim"

func (a *AOJ) Submit(v *nvim.Nvim, args []string) error {
	// NOTE: args[0] -> problem id
	//       args[1] -> language
	//       args[2] -> ...
	v.Command("echo '" + args[0] + " " + args[1] + "'")
	return nil
}
