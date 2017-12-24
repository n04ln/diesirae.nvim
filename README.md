# diesirae.nvim

![](https://travis-ci.org/NoahOrberg/diesirae.nvim.svg?branch=master)

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

diesirae.nvim provides 5 commands 

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

- Submit source code in current buffer
  - Please enter the ProblemID or URL
- The result(status) is written in a scratch buffer

``` vim
:AojSubmit
```

or use `<C-d>s`

### AojStatus

- Get status of current buffer

``` vim
:AojStatus
```

### AojStatusList

- Get status list of current buffer

``` vim
:AojStatusList
```
