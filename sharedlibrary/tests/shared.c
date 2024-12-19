#include "shared.h"

void *handle = NULL;

void *load_function(const char *func_name)
{
    void *func = dlsym(handle, func_name);
    if (!func)
    {
        fprintf(stderr, "Failed to load function %s: %s\n", func_name, dlerror());
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
