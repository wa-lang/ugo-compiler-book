extern int exit(int);

int ugo_builtin_exit(int x) {
	exit(x);
	return 0;
}
