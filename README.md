[![test](https://github.com/nobishino/1go/actions/workflows/test.yml/badge.svg?branch=main)](https://github.com/nobishino/1go/actions/workflows/test.yml)

# 1go

Learn you go compiler for great good!

## コマンドメモ

```sh
docker run --rm -it -v `pwd`:/go/1go golang:1.16.2
```

## 現在の文法

```ebnf
program    = stmt*
stmt       = expr ";" | "return" expr ";"
expr       = assign
assign     = equality ("=" assign)?
equality   = relational ("==" relational | "!=" relational)*
relational = add ("<" add | "<=" add | ">" add | ">=" add)*
add        = mul ("+" mul | "-" mul)*
mul        = unary ("*" unary | "/" unary)*
unary      = ("+" | "-")? primary
primary    = num | ident | "(" expr ")"
```