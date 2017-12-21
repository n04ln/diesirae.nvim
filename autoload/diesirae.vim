" GetWindowList:
"  return Dictionary Type.
"    key is window number(String), and value is buffername(String).
function! diesirae#getWindowList()
    let res = {}

    for i in range(1, winnr('$'))
        let res[i] = bufname(winbufnr(i))   
    endfor

    return res
endfunction
