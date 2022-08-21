#!/bin/bash
assert() {
	expected="$1"
	input="$2"

	./bin/gocc "$input" > ./bin/tmp.s
	cc -o ./bin/tmp ./bin/tmp.s
	./bin/tmp
	actual="$?"

	if [ "$actual" = "$expected" ]; then
		echo "[OK] $input => $actual"
	else
		echo "[FAIL] $input => $expected expected, but got $actual"
		exit 1
	fi
}

assert_lib() {
  lib="$1"
  expected="$2"
  input="$3"

 	./bin/gocc "$input" > ./bin/tmp.s
 	cc -w -o ./bin/lib.s -S "$lib"
 	cc -o ./bin/tmp ./bin/tmp.s ./bin/lib.s
 	./bin/tmp
 	actual="$?"

 	if [ "$actual" = "$expected" ]; then
 		echo "[OK] $input => $actual"
 	else
 		echo "[FAIL] $input => $expected expected, but got $actual"
 		exit 1
 	fi
}



#assert 0 "return 0;"
#assert 42 "return 42;"
#
#assert 5 "return 1+2+2;"
#assert 5 "return 10-7+2;"
#assert 3 "return  1 + 3 + 5 - 4 -2;"
#
#assert 47 "return 5+6*7;"
#assert 15 "return 5*(9-6);"
#assert 4 "return (3+5)/2;"
#assert 10 "return -10+20;"
#assert 6 "return (-10) * (+10) / -50 + 4;"
#assert 2 "return +2;"
#assert 12 "return +8-(-4);"
#
#assert 0 'return 0==1;'
#assert 1 'return 42==42;'
#assert 1 'return 0!=1;'
#assert 0 'return 42!=42;'
#
#assert 1 'return 0<1;'
#assert 0 'return 1<1;'
#assert 0 'return 2<1;'
#assert 1 'return 0<=1;'
#assert 1 'return 1<=1;'
#assert 0 'return 2<=1;'
#
#assert 1 'return 1>0;'
#assert 0 'return 1>1;'
#assert 0 'return 1>2;'
#assert 1 'return 1>=0;'
#assert 1 'return 1>=1;'
#assert 0 'return 1>=2;'
#
#assert 8 "a=8;return a;"
#assert 3 "a=1;b=a+2; return b;"
#assert 6 "a=3;b=3;return a+b;"
#assert 9 "a=6;a=a+3; return a;"
assert 6 "a=1; b=2; c=3; d = a+b+c; return d;"
#assert 6 "a=1; b=2; c=3; return a+b+c;"
##assert 26 "{a=1;b=1;c=1;d=1;e=1;f=1;g=1;h=1;i=1;j=1;k=1;l=1;m=1;n=1;o=1;p=1;q=1;r=1;s=1;t=1;u=1;v=1;w=1;x=1;y=1;z=1;return a+b+c+d+e+f+g+h+i+j+k+l+m+n+o+p+q+r+s+t+u+v+w+x+y+z;}"
##assert 2 "a=2;b=9;c=3;d=0;e=2;f=4;g=8;h=9;i=7;j=0;k=8;l=8;m=2;n=5;o=1;p=6;q=2;r=1;s=9;t=8;u=1;v=6;w=2;x=7;y=9;z=2;return a/b+c-d/e*f+g+h+i+j*k-l+m-n-o/p/q/r/s-t+u*v*w-x-y-z;"
#assert 10 "five=5;result=five*2; return result;"
#
#assert 10 "{return 10; return 100;}"
#
#assert 20 "if ( 8 > 2 ) return 20; return 10;"
#assert 10 "if ( 2 > 8 ) return 20; return 10;"
#assert 20 "if ( 8 > 2 ) return 20; else return 10;"
#assert 10 "if ( 8 < 2 ) return 20; else return 10;"
#assert 10 "if ( 8==8 ) return 10;"
#
#assert 2 "cond = 2; if ( cond == 1 ) return 1; else if ( cond == 2 ) return 2; else return 3;"
#assert 10 "i=0; while ( i<10 ) i=i+1; return i;"
#assert 2 "x=2; while(x ==1) x=x+1; return x;"
#assert 10 "total = 0; for (i=0;i<5;i=i+1) total = total + i; return total;"
#
#assert 10 "{ return 10; }"
#assert 20 "result = 0; if ( 1 > 0 ) { result = 10; result = result * 2; } else { result = 30; } return result;"
#assert 100 "count = 0; result = 0; while( count < 10 ) { result = result + 10; count = count + 1; } return result;"
#assert 50 "result = 0; for(i=0; i<5; i=i+1) { result = result + 10; } return result;"
#assert 5 "
#result = 0;
#for (;;) {
#  if (result > 4) {
#    return result;
#  }
#  result = result + 1;
#}
#return result;
#"
#
#assert 2 "if ( 1==1 && 2 != 3) {return 2;} else {return 1;}"
#assert 3 "if (( (1==1) && (2==2) )|| (1==0)) {return 3;} else {return 100;}"
#assert 1 "return (1==1||0==1);"
#assert 0 "return (1==1&&0==1);"
#assert_lib "./lib/foo.c" 2 "return foo();"
#assert_lib "./lib/foo.c" 1 "foo(); return 1;"
#assert_lib "./lib/foo.c" 3 "two = foo(); return two + 1;" # ccでビルドしたライブラリの戻り値はw0に保存されるため現状取り出せない
#assert_lib "./lib/foo.c" 2 "for (i=0; i<10; i=i+1) { return foo(); }"
#assert_lib "./lib/foo.c" 49 "_ = foo(); return 49;"
#assert_lib "./lib/foo.c" 2 "foo(); foo(); return 2;"
#assert_lib "./lib/foo.c" 30 "
#for (i=0; i<10; i=i+1) {
#  _ = foo();
#}
#return 30;"
#assert_lib "./lib/foo.c" 50 "i=0; while(i<10) {i=i+1;} foo(); return 50;"
#assert_lib "./lib/foo.c" 4 "foo();foo();foo();foo();foo();foo();foo();foo();foo();foo();foo();foo();foo();foo();foo();foo();foo();foo();foo();foo();foo();foo();foo();foo();foo();foo();foo();foo();foo();foo();foo();foo();foo();foo();foo();foo();foo();foo(); return 4;"
#assert_lib "./lib/foo.c" 5 "return 5;"
#assert_lib "./lib/foo.c" 5 "for (i=0; i<10; i=i+1) { foo(); } return 5;"
#
#echo OK
