if exists('g:loaded_diesirae')
  finish
endif
let g:loaded_diesirae = 1

function! s:RequireDiesIrae(host) abort
  return jobstart(['diesirae.nvim'], { 'rpc': v:true })
endfunction


call remote#host#Register('diesirae.nvim', '0', function('s:RequireDiesIrae'))
call remote#host#RegisterPlugin('diesirae.nvim', '0', [
\ {'type': 'command', 'name': 'AojSelf', 'sync': 1, 'opts': {}},
\ {'type': 'command', 'name': 'AojSession', 'sync': 1, 'opts': {}},
\ {'type': 'command', 'name': 'AojStatus', 'sync': 1, 'opts': {}},
\ {'type': 'command', 'name': 'AojStatusList', 'sync': 1, 'opts': {}},
\ {'type': 'command', 'name': 'AojSubmit', 'sync': 0, 'opts': {}},
\ ])
