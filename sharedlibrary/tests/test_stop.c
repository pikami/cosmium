#include "shared.h"

int test_StopServerInstance()
{
    typedef int (*StopServerInstanceFn)(char *);
    StopServerInstanceFn StopServerInstance = (StopServerInstanceFn)load_function("StopServerInstance");

    if (!StopServerInstance)
    {
        fprintf(stderr, "Failed to find StopServerInstance function\n");
        return 0;
    }

    char *serverName = "TestServer";
    int result = StopServerInstance(serverName);
    if (result == 0)
    {
        printf("StopServerInstance: SUCCESS\n");
    }
    else
    {
        printf("StopServerInstance: FAILED (result = %d)\n", result);
        return 0;
    }

    return 1;
}
