#include "shared.h"

void test_CreateServerInstance();
void test_StopServerInstance();
void test_ServerInstanceStateMethods();

int main(int argc, char *argv[])
{
    if (argc < 2)
    {
        fprintf(stderr, "Usage: %s <path_to_shared_library>\n", argv[0]);
        return EXIT_FAILURE;
    }

    const char *libPath = argv[1];
    handle = dlopen(libPath, RTLD_LAZY);
    if (!handle)
    {
        fprintf(stderr, "Failed to load shared library: %s\n", dlerror());
        return EXIT_FAILURE;
    }

    printf("Running tests for library: %s\n", libPath);
    test_CreateServerInstance();
    test_ServerInstanceStateMethods();
    test_StopServerInstance();

    dlclose(handle);
    return EXIT_SUCCESS;
}
