#ifndef SHARED_H
#define SHARED_H

#include <stdio.h>
#include <stdlib.h>
#include <dlfcn.h>
#include <string.h>
#include <ctype.h>

extern void *handle;

void *load_function(const char *func_name);
char *compact_json(const char *json);

#endif
