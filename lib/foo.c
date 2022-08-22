#include <stdio.h>
int foo() {
    printf("foo\n");
    return 2;
}

void noop() {
    printf("noop\n");
}

int add1(int a) {
    return a + 1;
}

int addAB(int a, int b) {
    return a + b;
}

int addABC(int a, int b, int c) {
    return a + b + c;
}

void wprintf(char *c) {
    printf("%s", c);
}