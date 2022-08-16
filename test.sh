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

assert 0 "0;"
assert 42 "42;"

assert 5 "1+2+2;"
assert 5 "10-7+2;"
assert 3 " 1 + 3 + 5 - 4 -2;"

assert 47 "5+6*7;"
assert 15 "5*(9-6);"
assert 4 "(3+5)/2;"
assert 10 "-10+20;"
assert 6 "(-10) * (+10) / -50 + 4;"
assert 2 "+2;"
assert 12 "+8-(-4);"

assert 0 '0==1;'
assert 1 '42==42;'
assert 1 '0!=1;'
assert 0 '42!=42;'

assert 1 '0<1;'
assert 0 '1<1;'
assert 0 '2<1;'
assert 1 '0<=1;'
assert 1 '1<=1;'
assert 0 '2<=1;'

assert 1 '1>0;'
assert 0 '1>1;'
assert 0 '1>2;'
assert 1 '1>=0;'
assert 1 '1>=1;'
assert 0 '1>=2;'

assert 8 "a=8;"
assert 3 "a=1;b=a+2;"
assert 6 "a=3;b=3;a+b;"
assert 9 "a=6;a=a+3;"
assert 26 "a=1;b=1;c=1;d=1;e=1;f=1;g=1;h=1;i=1;j=1;k=1;l=1;m=1;n=1;o=1;p=1;q=1;r=1;s=1;t=1;u=1;v=1;w=1;x=1;y=1;z=1;+a+b+c+d+e+f+g+h+i+j+k+l+m+n+o+p+q+r+s+t+u+v+w+x+y+z;"
assert 2 "a=2;b=9;c=3;d=0;e=2;f=4;g=8;h=9;i=7;j=0;k=8;l=8;m=2;n=5;o=1;p=6;q=2;r=1;s=9;t=8;u=1;v=6;w=2;x=7;y=9;z=2;a/b+c-d/e*f+g+h+i+j*k-l+m-n-o/p/q/r/s-t+u*v*w-x-y-z;"
assert 10 "five=5;result=five*2;"
echo OK
