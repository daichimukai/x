/*
 * Written to demonstrate the difference between:
 *
 * 0. without redirect,
 *
 *    $ ./a.out | ts
 *    stderr
 *    Mar 11 21:55:03 stdout
 *    $
 *
 * 1. redirect fd 1 to /dev/null followed by fd 2 to fd 1,
 *
 *    $ ./a.out >/dev/null 2>&1 | ts
 *    $
 *
 * 2. and redirect fd 2 to fd 1 followed by fd 1 to /dev/null.
 *
 *    $ ./a.out 2>&1 >/dev/null | ts
 *    Mar 11 21:55:38 stderr
 *    $
 */

#include <stdio.h>

int main() {
	fprintf(stdout, "stdout\n");
	fprintf(stderr, "stderr\n");

	return 0;
}
