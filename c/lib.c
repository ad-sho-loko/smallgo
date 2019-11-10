#include <stdio.h>
#include <stdlib.h>

void EXPECT(char* msg, int got, int want){
    fprintf(stderr, "[Test: %s] => %d\n", msg, got);
    if (want != got){
        fprintf(stderr, "expected %d, but %d\n", want, got);
        exit(1);
    }
}