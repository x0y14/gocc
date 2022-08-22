[![CI](https://github.com/x0y14/gocc/actions/workflows/ci.yml/badge.svg)](https://github.com/x0y14/gocc/actions/workflows/ci.yml)

[aarch64向け9ccのc実装(未完成)](https://github.com/x0y14/9cc)のgo再実装  

参考:
- [低レイヤを知りたい人のためのCコンパイラ作成入門](https://www.sigbus.info/compilerbook)
- [Jun's Homepage](https://www.mztn.org/dragon/arm6400idx.html#toc)
- [arm developer](https://developer.arm.com/documentation/102374/latest/)
- [Overview of ARM64 ABI conventions](https://docs.microsoft.com/en-us/cpp/build/arm64-windows-abi-conventions?view=msvc-170)
- [modexp](https://modexp.wordpress.com/2018/10/30/arm64-assembly/)
- https://stackoverflow.com/questions/66098678/understanding-aarch64-assembly-function-call-how-is-stack-operated

[条件分岐に使用するeq,neなどの一覧](https://www.mztn.org/dragon/arm6408cond.html#suffix)

memo:
- Zero Flag (Z)
- Carry Flag (C)
- Negative Flag (N)
- Overflow Flag (V)

### ebnf
```
program    = stmt*
stmt       = expr ";"
           | "return" expr ";"
           | "if" "(" expr ")" stmt ("else" stmt)?
           | "while "(" expr ")" stmt
           | "for" "(" expr? ";" expr? ";" expr? ")" stmt
           | "{" stmt* "}"
expr       = assign
assign     = andor ("=" assign)?
andor      = equality ("&&" equality | "||" equality)*
equality   = relational ("==" relational | "!=" relational)*
relational = add ("<" add | "<=" add | ">" add | ">=" add)*
add        = mul ("+" mul | "-" mul)*
mul        = unary ("*" unary | "/" unary)*
unary      = ("+" | "-")? primary
primary    = num
           | ident ("(" (expr ","?)* ")")?
           | "(" expr ")"
```