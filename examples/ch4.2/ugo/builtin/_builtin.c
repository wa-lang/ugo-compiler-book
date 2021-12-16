#include <stdio.h>
#include <stdlib.h>

int ugo_builtin_println(int x) {
	return printf("%d\n", x);
}
int ugo_builtin_exit(int x) {
	exit(x);
	return 0;
}
