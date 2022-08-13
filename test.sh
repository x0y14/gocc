#!/bin/bash
assert() {
	expected="$1"
	input="$2"

	./bin/gocc "$input" > ./bin/tmp.s
	cc -o ./bin/tmp ./bin/tmp.s
	./bin/tmp
	actual="$?"

	if [ "$actual" = "$expected" ]; then
		echo "$input => $actual"
	else
		echo "$input => $expected expected, but got $actual"
		exit 1
	fi
}

assert 0 0
assert 42 42
assert 5 "1+2+2"
assert 5 "10-7+2"
assert 3 " 1 + 3 + 5 - 4 -2"
assert 47 "5+6*7"
assert 15 "5*(9-6)"
assert 4 "(3+5)/2"
assert 10 "-10+20"
assert 6 "(-10) * (+10) / -50 + 4"
assert 2 "+2"
assert 12 "+8-(-4)"
echo OK
