" diesirae#getWindowList:
"  return Dictionary Type.
"    key is window number(String), and value is buffername(String).
function! diesirae#getWindowList() abort
    let res = {}

    for i in range(1, winnr('$'))
        let res[i] = bufname(winbufnr(i))   
    endfor

    return res
endfunction

" diesirae#aojSubmit:
"  void
"   submission function within `input()`
function! diesirae#aojSubmit() abort
    let pid = input("problem id: ")
    AojSubmit(pid)
endfunction
