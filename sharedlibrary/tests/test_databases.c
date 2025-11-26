#include "shared.h"

int test_Databases()
{
    /* Load FreeMemory function - must use this to free memory allocated by the DLL
       because the DLL may use a different C runtime heap than the test loader */
    typedef void (*FreeMemoryFn)(char *);
    FreeMemoryFn FreeMemory = (FreeMemoryFn)load_function("FreeMemory");
    if (!FreeMemory)
    {
        fprintf(stderr, "Failed to find FreeMemory function\n");
        return 0;
    }

    typedef int (*CreateDatabaseFn)(char *, char *);
    CreateDatabaseFn CreateDatabase = (CreateDatabaseFn)load_function("CreateDatabase");
    if (!CreateDatabase)
    {
        fprintf(stderr, "Failed to find CreateDatabase function\n");
        return 0;
    }

    char *serverName = "TestServer";
    char *configJSON = "{\"id\":\"test-db\"}";

    int result = CreateDatabase(serverName, configJSON);
    if (result == 0)
    {
        printf("CreateDatabase: SUCCESS\n");
    }
    else
    {
        printf("CreateDatabase: FAILED (result = %d)\n", result);
        return 0;
    }

    typedef char *(*GetDatabaseFn)(char *, char *);
    GetDatabaseFn GetDatabase = (GetDatabaseFn)load_function("GetDatabase");
    if (!GetDatabase)
    {
        fprintf(stderr, "Failed to find GetDatabase function\n");
        return 0;
    }

    char *database = GetDatabase(serverName, "test-db");
    if (database)
    {
        printf("GetDatabase: SUCCESS (database = %s)\n", database);
        FreeMemory(database);
    }
    else
    {
        printf("GetDatabase: FAILED\n");
        return 0;
    }

    return 1;
}
