#include "shared.h"

int test_CreateServerInstance();
int test_StopServerInstance();
int test_ServerInstanceStateMethods();
int test_Databases();

int main(int argc, char *argv[])
{
    if (argc < 2)
    {
        fprintf(stderr, "Usage: %s <path_to_shared_library>\n", argv[0]);
        return EXIT_FAILURE;
    }

    const char *libPath = argv[1];
    handle = load_library(libPath);
    if (!handle)
    {
        fprintf(stderr, "Failed to load shared library: %s\n", get_load_error());
        return EXIT_FAILURE;
    }

    printf("Running tests for library: %s\n", libPath);
    int results[] = {
        test_CreateServerInstance(),
        test_Databases(),
        test_ServerInstanceStateMethods(),
        test_StopServerInstance(),
    };

    int numTests = sizeof(results) / sizeof(results[0]);
    int numPassed = 0;
    for (int i = 0; i < numTests; i++)
    {
        if (results[i])
        {
            numPassed++;
        }
    }

    printf("Tests passed: %d/%d\n", numPassed, numTests);

    close_library(handle);
    return EXIT_SUCCESS;
}
