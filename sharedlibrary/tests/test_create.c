#include "shared.h"

int test_CreateServerInstance()
{
    typedef int (*CreateServerInstanceFn)(char *, char *);
    CreateServerInstanceFn CreateServerInstance = (CreateServerInstanceFn)load_function("CreateServerInstance");

    if (!CreateServerInstance)
    {
        fprintf(stderr, "Failed to find CreateServerInstance function\n");
        return 0;
    }

    char *serverName = "TestServer";
    char *configJSON = "{\"host\":\"localhost\",\"port\":8080}";

    int result = CreateServerInstance(serverName, configJSON);
    if (result == 0)
    {
        printf("CreateServerInstance: SUCCESS\n");
    }
    else
    {
        printf("CreateServerInstance: FAILED (result = %d)\n", result);
        return 0;
    }

    return 1;
}
