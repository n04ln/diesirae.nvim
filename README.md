# diesirae.nvim

![](https://travis-ci.org/NoahOrberg/diesirae.nvim.svg?branch=master)

AOJ client for NeoVim  

## Demo

[![asciicast](https://asciinema.org/a/188775.png)](https://asciinema.org/a/188775)

## How to install
### Require

- NeoVim (>=0.2.1)
- go (>=1.8)
- glide

### Install 

1. get
``` sh
$ go get github.com/NoahOrberg/diesirae.nvim
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

diesirae.nvim provides 6 commands 

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
- This is Asynchronous
  - And this is function. not command
- Submit source code in current buffer
  - Please enter the ProblemID or URL
- The result(status) is written in a scratch buffer

``` vim
:call AojSubmit("<PROBLEM_ID_OR_URL>")
```

or use `<C-d>s` in normal mode. 
This is interactive(`<C-d>s` is `:call AojSubmit(input("problem id: "))`) 

### AojRunSample

- Try Sample Input/Output
  - It is the same as using `AojSubmit`

- This command requires setting of the `g:diesirae_config` variable

``` vim
let g:diesirae_config = {
\  'commands': {
\    'py': {
\      'build_command': [], 
\      'exec_command': ['python3', '*source*']
\    },
\    'go': {
\      'build_command': ['go', 'build', '-o', '*bin*', '*source*'], 
\      'exec_command': ['*bin*']
\    }
\  }
\}
```

- But It is not recorded on the AOJ server

``` vim
:call AojRunSample("<PROBLEM_ID_OR_URL>")
```

or use `<C-d>t` in normal mode.
This is interactive(`<C-d>t` is `:call AojRunSample(input("problem id: "))`) 

### AojDescription

- Get Description

``` vim
:call AojDescription("<PROBLEM_ID_OR_URL>")
```

or use `<C-d>d` in normal mode.
This is interactive(`<C-d>d` is `:call AojDescription(input("problem id: "))`) 

### AojStatus

- Get status of current buffer

``` vim
:AojStatus
```

### AojStatusList

- Get status list of all buffers

``` vim
:AojStatusList
```

