#include "shared.h"

void test_ServerInstanceStateMethods()
{
    typedef int (*LoadServerInstanceStateFn)(char *, char *);
    LoadServerInstanceStateFn LoadServerInstanceState = (LoadServerInstanceStateFn)load_function("LoadServerInstanceState");
    if (!LoadServerInstanceState)
    {
        fprintf(stderr, "Failed to find LoadServerInstanceState function\n");
        return;
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
    }

    typedef char *(*GetServerInstanceStateFn)(char *);
    GetServerInstanceStateFn GetServerInstanceState = (GetServerInstanceStateFn)load_function("GetServerInstanceState");
    if (!GetServerInstanceState)
    {
        fprintf(stderr, "Failed to find GetServerInstanceState function\n");
        return;
    }

    char *state = GetServerInstanceState(serverName);
    if (state)
    {
        printf("GetServerInstanceState: SUCCESS (state = %s)\n", state);
    }
    else
    {
        printf("GetServerInstanceState: FAILED\n");
    }

    const char *expected_state = "{\"databases\":{\"test-db\":{\"id\":\"test-db\",\"_ts\":0,\"_rid\":\"\",\"_etag\":\"\",\"_self\":\"\"}},\"collections\":{\"test-db\":{}},\"documents\":{\"test-db\":{}}}";
    char *compact_state = compact_json(state);
    if (!compact_state)
    {
        free(state);
        return;
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
    }

    free(state);
    free(compact_state);
}
