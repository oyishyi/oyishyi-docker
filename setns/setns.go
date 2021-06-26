package setns

/*
#include <errno.h>
#include <sched.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <fcntl.h>

// __attribut__((constructor)) means this function will be called right after the package is imported
// in other words, this function will run before the program run
__attribute__((constructor)) void enter_namespace(void) {
	char *mydocker_pid;
	// get pid from env
	mydocker_pid = getenv("mydocker_pid");
	if (mydocker_pid) {
		//fprintf(stdout, "got env: mydocker_pid=%s\n", mydocker_pid);
	} else {
		// using env to control whether run this bunch of codes
		// if env not exist, than this function will not run
		// so that command other than "docker exec" will not run this bunch of cgo code
		fprintf(stdout, "missing env: mydocker_pid\n");
		return;
	}
	char *mydocker_cmd;
	mydocker_cmd = getenv("mydocker_cmd");
	if (mydocker_cmd) {
		fprintf(stdout, "got env: mydocker_cmd=%s\n", mydocker_cmd);
	} else {
		fprintf(stdout, "missing env: mydocker_cmd\n");
		return;
	}


	// five namespaces that need to enter
	char *namespaces[] = {"ipc", "uts", "net", "pid", "mnt"};

	char nspath[1024];

	int i; // old c compiler style
	for (i = 0; i < 5; i++) {
	    sprintf(nspath, "/proc/%s/ns/%s", mydocker_pid, namespaces[i]);
	    int fd = open(nspath, O_RDONLY);
	    // call setns to enter namespace
	    if (setns(fd, 0) == -1) {
	        fprintf(stderr, "setns %s fails: %s\n", namespaces[i], strerror(errno));
	    } else {
	        fprintf(stdout, "setns %s succeed\n", namespaces[i]);
	    }
	    close(fd);
	}
	// after enter the namespaces, run the cmd
	int res = system(mydocker_cmd);
	exit(0);
	return;
}

*/
import "C"
