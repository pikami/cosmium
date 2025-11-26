#include "shared.h"

int test_CreateServerInstance();
int test_StopServerInstance();
int test_ServerInstanceStateMethods();
int test_Databases();

int main(int argc, char *argv[])
{
    /* Disable stdout buffering for CI environments without a real terminal */
    setvbuf(stdout, NULL, _IONBF, 0);

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

    /* give the loaded library a short time to initialize */
#ifdef _WIN32
    Sleep(1000);
#else
    sleep(1);
#endif

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

    /* Exit explicitly before unloading the library.
       Go runtime cleanup during FreeLibrary can set a non-zero exit code on Windows. */
    int exitCode = (numPassed == numTests) ? EXIT_SUCCESS : EXIT_FAILURE;

#ifdef _WIN32
    ExitProcess(exitCode);
#else
    close_library(handle);
#endif

    return exitCode;
}
