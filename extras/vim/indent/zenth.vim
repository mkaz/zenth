" Vim indent file
" Language: Zenth
" Filenames: *.zn

if exists("b:did_indent")
  finish
endif
let b:did_indent = 1

setlocal indentexpr=ZenthIndent()
setlocal indentkeys=0{,0},!^F,o,O

function! ZenthIndent()
  let prevlnum = prevnonblank(v:lnum - 1)
  if prevlnum == 0
    return 0
  endif

  let prevline = getline(prevlnum)
  let curline = getline(v:lnum)
  let ind = indent(prevlnum)

  " Increase indent after opening brace
  if prevline =~ '{\s*$'
    let ind += shiftwidth()
  endif

  " Decrease indent on closing brace
  if curline =~ '^\s*}'
    let ind -= shiftwidth()
  endif

  return ind
endfunction
