#ifndef SHARED_H
#define SHARED_H

#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <ctype.h>

// Platform-specific includes for dynamic library loading
#ifdef _WIN32
    #include <windows.h>
    typedef HMODULE lib_handle_t;
#else
    #include <dlfcn.h>
    typedef void* lib_handle_t;
#endif

extern lib_handle_t handle;

void *load_function(const char *func_name);
char *compact_json(const char *json);
char *get_load_error(void);
lib_handle_t load_library(const char *path);
void close_library(lib_handle_t handle);

#endif
