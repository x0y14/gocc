[![CI](https://github.com/x0y14/gocc/actions/workflows/ci.yml/badge.svg)](https://github.com/x0y14/gocc/actions/workflows/ci.yml)

[macos向け9ccのc実装(未完成)](https://github.com/x0y14/9cc)のgo再実装
参考: [低レイヤを知りたい人のためのCコンパイラ作成入門](https://www.sigbus.info/compilerbook)

### ebnf
```
expr    = mul ("+" mul | "-" mul)*  
mul     = unary ("*" unary | "/" un*ary)*  
unary   = ("+" | "-")? primary*  
primary = num | "(" expr ")"  
```