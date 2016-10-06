#include <algorithm>
#include <iostream>
#include <random>

#include <cstdlib>

#include "lesMeilleures.h"

int                                     main(int argc __attribute__((unused)), char **argv __attribute__((unused))) {
    std::random_device                  ran;
    std::mt19937                        gen(ran());
    std::uniform_int_distribution<int>  dis(-42, 42);
    int                                 res[2] = {-42, 42};
    int                                 a;
    int                                 b;
    int                                 c;

    do {
        a = dis(gen);
        b = dis(gen);
        c = dis(gen);
    } while (a == b && b == c && a == c);
    les_meilleures(res, a, b, c);
    if ((res[0] + res[1]) == (a + b + c - std::min(std::min(a, b), c)))
        std::cout << "Bravo !" << std::endl;
    else
        std::cout << "Perdu, recommence !" << std::endl;
    return EXIT_SUCCESS;
}
