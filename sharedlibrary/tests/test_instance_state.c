#include "shared.h"

int test_ServerInstanceStateMethods()
{
    typedef int (*LoadServerInstanceStateFn)(char *, char *);
    LoadServerInstanceStateFn LoadServerInstanceState = (LoadServerInstanceStateFn)load_function("LoadServerInstanceState");
    if (!LoadServerInstanceState)
    {
        fprintf(stderr, "Failed to find LoadServerInstanceState function\n");
        return 0;
    }

    char *serverName = "TestServer";
    char *stateJSON = "{\"databases\":{\"test-db\":{\"id\":\"test-db\"}}}";
    int result = LoadServerInstanceState(serverName, stateJSON);
    if (result == 0)
    {
        printf("LoadServerInstanceState: SUCCESS\n");
    }
    else
    {
        printf("LoadServerInstanceState: FAILED (result = %d)\n", result);
        return 0;
    }

    typedef char *(*GetServerInstanceStateFn)(char *);
    GetServerInstanceStateFn GetServerInstanceState = (GetServerInstanceStateFn)load_function("GetServerInstanceState");
    if (!GetServerInstanceState)
    {
        fprintf(stderr, "Failed to find GetServerInstanceState function\n");
        return 0;
    }

    char *state = GetServerInstanceState(serverName);
    if (state)
    {
        printf("GetServerInstanceState: SUCCESS (state = %s)\n", state);
    }
    else
    {
        printf("GetServerInstanceState: FAILED\n");
        return 0;
    }

    const char *expected_state = "{\"databases\":{\"test-db\":{\"id\":\"test-db\",\"_ts\":0,\"_rid\":\"\",\"_etag\":\"\",\"_self\":\"\"}},\"collections\":{\"test-db\":{}},\"documents\":{\"test-db\":{}},\"triggers\":{\"test-db\":{}},\"sprocs\":{\"test-db\":{}},\"udfs\":{\"test-db\":{}}}";
    char *compact_state = compact_json(state);
    if (!compact_state)
    {
        free(state);
        return 0;
    }

    if (strcmp(compact_state, expected_state) == 0)
    {
        printf("GetServerInstanceState: State matches expected value.\n");
    }
    else
    {
        printf("GetServerInstanceState: State does not match expected value.\n");
        printf("Expected: %s\n", expected_state);
        printf("Actual:   %s\n", compact_state);
        return 0;
    }

    free(state);
    free(compact_state);
    return 1;
}
