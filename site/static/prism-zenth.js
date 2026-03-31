Prism.languages.zenth = {
  'comment': [
    { pattern: /\/\*[\s\S]*?\*\//, greedy: true },
    { pattern: /\/\/.*/, greedy: true }
  ],
  'string': [
    { pattern: /"""[\s\S]*?"""/, greedy: true },
    { pattern: /"(?:[^"\\]|\\.)*"/, greedy: true },
    { pattern: /'(?:[^'\\]|\\.)*'/, greedy: true }
  ],
  'keyword': /\b(?:fn|let|var|const|return|if|else|for|in|match|range|obj|enum|interface|import|from|as|break|continue|type)\b/,
  'builtin': /\b(?:Print|Println|Len|Str|Hashmap|File|Args|Env|Exit|Range|Rangei|Int|Float|Ord|Chr|Abs|Min|Max|Clamp|Round|Floor|Ceil|Pow|Sqrt|Assert|AssertEq|Zip|Tuple|Date|Input|Set)\b/,
  'type': /\b(?:Int|Float|Bool|Str|Byte|Date|Array|Tuple|Hashmap|Set|Fn)\b/,
  'constant': /\b(?:true|false|nil|INT_MAX|INT_MIN)\b/,
  'self': /\bself\b/,
  'number': [
    { pattern: /\b\d+\.\d*\b/ },
    { pattern: /\b\d+\b/ }
  ],
  'operator': /=>|[+\-*/%=<>!&|^~]=?|&&|\|\||\+\+|--/,
  'punctuation': /[{}()\[\];,.:]/
};

Prism.highlightAll();
