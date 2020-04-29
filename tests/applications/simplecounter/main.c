#include <stdio.h>
#include <unistd.h>

/**
 * This program will continously count from 0 to 2^16 inclusively
 * and output the current value every 1 seconds to stdout
 *  Parameters: none
 *  Return: 
 *    - status code, where 0 is success (however, unreachable
 *      during normal execution of program)
 */ 
int main() {
    printf("Welcome from simplecounter\n");
    fflush(stdout);

    unsigned short int var = 0; // Should be equivalent to uint16

    while (1) {
        printf("%i\n", var);
        fflush(stdout);

        sleep(1);
        var++;
    }

    return 0;
}