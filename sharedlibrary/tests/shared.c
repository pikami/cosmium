#include "shared.h"

lib_handle_t handle = NULL;

char *get_load_error(void)
{
#ifdef _WIN32
    DWORD error = GetLastError();
    static char buf[256];
    FormatMessageA(FORMAT_MESSAGE_FROM_SYSTEM | FORMAT_MESSAGE_IGNORE_INSERTS,
                   NULL, error, MAKELANGID(LANG_NEUTRAL, SUBLANG_DEFAULT),
                   buf, sizeof(buf), NULL);
    return buf;
#else
    return dlerror();
#endif
}

lib_handle_t load_library(const char *path)
{
#ifdef _WIN32
    return LoadLibraryA(path);
#else
    return dlopen(path, RTLD_LAZY);
#endif
}

void close_library(lib_handle_t handle)
{
#ifdef _WIN32
    FreeLibrary(handle);
#else
    dlclose(handle);
#endif
}

void *load_function(const char *func_name)
{
#ifdef _WIN32
    void *func = (void *)GetProcAddress(handle, func_name);
#else
    void *func = dlsym(handle, func_name);
#endif
    
    if (!func)
    {
        fprintf(stderr, "Failed to load function %s: %s\n", func_name, get_load_error());
    }
    return func;
}

char *compact_json(const char *json)
{
    size_t len = strlen(json);
    char *compact = (char *)malloc(len + 1);
    if (!compact)
    {
        fprintf(stderr, "Failed to allocate memory for compacted JSON\n");
        return NULL;
    }

    char *dest = compact;
    for (const char *src = json; *src != '\0'; ++src)
    {
        if (!isspace((unsigned char)*src)) // Skip spaces, newlines, tabs, etc.
        {
            *dest++ = *src;
        }
    }
    *dest = '\0'; // Null-terminate the string

    return compact;
}
