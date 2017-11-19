# diesirae.nvim
AOJ client by NeoVim

## How to install
### Require
- go (>=1.8)
- glide
### Install 
1. build
``` sh
$ go get github.com/NoahOrberg/diesirae.nvim
$ cd $GOPATH/src/github.com/NoahOrberg/diesirae.nvim/
$ glide install
$ go install
```
2. write in your `init.vim` (or `.vimrc`)
``` vim
" using dein.nvim
call dein#add("NoahOrberg/diesirae.nvim")
```

## Usage
you need to set the following two environment variables
``` sh
$ export AOJ_ID=xxxxxxxx
$ export AOJ_RAWPASSWORD=yyyyyyyy
```
diesirae.nvim provide 4 commands 
### AojSelf
- Check session
``` vim
:AojSelf
```
### AojSession
- Get cookie when session does not exist
``` vim
:AojSession
```
### AojSubmit
- Submit contents in current buffer
- Need a argument
  - Problem ID
- The result(status) is written in a scratch buffer
``` vim
:AojSubmit ITP1_1_A
```
### AojStatus
- Get status of current buffer
``` vim
:AojStatus
```
