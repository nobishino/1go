#!/bin/bash
assert() {
  expected="$1"
  input="$2"

  ./main "$input" > tmp.s
  cc -o tmp tmp.s
  ./tmp
  actual="$?"

  if [ "$actual" = "$expected" ]; then
    echo "$input => $actual"
  else
    echo "$input => $expected expected, but got $actual"
    exit 1
  fi
}
go test ./...

assert 0 0
assert 42 42
assert 9 4+5
assert 13 4+21-12
assert 41 " 12 + 34 - 5 "
assert 4 " 1*2 + 5 /2"
assert 47 '5+6*7'
assert 15 '5*(9-6)'
assert 4 '(3+5)/2'
assert 10 '-10+20'
assert 1 '1==1'
assert 0 '1==2'
assert 0 '1!=1'

echo OK
