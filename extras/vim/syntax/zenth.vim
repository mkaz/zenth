" Vim syntax file
" Language: Zenth
" Filenames: *.zn

if exists("b:current_syntax")
  finish
endif

" Keywords
syn keyword zenthKeyword      fn let var const return if else for in match range
syn keyword zenthKeyword      obj enum interface import from as break continue

" Constants
syn keyword zenthConstant     true false nil
syn keyword zenthIdentifier   self

" Built-in types
syn keyword zenthType         int i8 i16 i32 i64 u8 u16 u32 u64
syn keyword zenthType         f32 f64 bool str byte

" Built-in functions
syn keyword zenthBuiltin      print println len str map append push pop file flag

" Operators
syn match zenthOperator       /[+\-*/%=<>!&|^~]/
syn match zenthOperator       /=>/

" Numbers
syn match zenthNumber         /\<\d\+\>/
syn match zenthFloat          /\<\d\+\.\d*\>/

" Strings
syn region zenthString        start=/"/ skip=/\\"/ end=/"/

" Comments
syn region zenthComment       start="//" end="$"
syn region zenthCommentBlock  start="/\*" end="\*/"

" Highlighting
hi def link zenthKeyword      Statement
hi def link zenthConstant     Constant
hi def link zenthIdentifier   Identifier
hi def link zenthType         Type
hi def link zenthBuiltin      Function
hi def link zenthOperator     Operator
hi def link zenthNumber       Number
hi def link zenthFloat        Float
hi def link zenthString       String
hi def link zenthComment      Comment
hi def link zenthCommentBlock Comment

let b:current_syntax = "zenth"
