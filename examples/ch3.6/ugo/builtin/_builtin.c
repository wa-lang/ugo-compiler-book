#include <stdio.h>
#include <stdlib.h>

int ugo_builtin_exit(int x) {
	printf("ugo_builtin_exit(%d)\n", x);
	exit(x);
	return 0;
}
