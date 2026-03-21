" Vim syntax file
" Language: Zenth
" Filenames: *.zn

if exists("b:current_syntax")
  finish
endif

" Keywords
syn keyword zenthKeyword      fn let var const return if else for in match range
syn keyword zenthKeyword      obj enum interface import from as break continue type

" Constants
syn keyword zenthConstant     true false nil INT_MAX INT_MIN
syn keyword zenthIdentifier   self

" Built-in types
syn keyword zenthType         Int Float Bool Str Byte Date
syn keyword zenthType         Array Tuple Hashmap Set Fn

" Built-in functions
syn keyword zenthBuiltin      Print Println Len Str Hashmap File Args Env Exit Range Rangei Int Float Abs Min Max Clamp Round Floor Ceil Pow Sqrt Assert AssertEq Zip Tuple Date

" Operators
syn match zenthOperator       /[+\-*/%=<>!&|^~]/
syn match zenthOperator       /=>/

" Numbers
syn match zenthNumber         /\<\d\+\>/
syn match zenthFloat          /\<\d\+\.\d*\>/

" Strings
syn region zenthString        start=/"""/ end=/"""/
syn region zenthString        start=/"/ skip=/\\"/ end=/"/
syn region zenthString        start=/'/ skip=/\\'/ end=/'/

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
