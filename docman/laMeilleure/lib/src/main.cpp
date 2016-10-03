#include <algorithm>
#include <iostream>
#include <random>

#include <cstdlib>

#include "laMeilleure.h"

int                                     main(int argc __attribute__((unused)), char **argv __attribute__((unused))) {
    std::random_device                  ran;
    std::mt19937                        gen(ran());
    std::uniform_int_distribution<int>  dis(1, 42);
    int                                 a;
    int                                 b;
    int                                 c;
    int                                 res;

    do {
        a = dis(gen);
        b = dis(gen);
        c = dis(gen);
    } while (a == b && b == c && a == c);
    res = la_meilleure(a, b, c);
    if (res == std::max(std::max(a, b), c))
        std::cout << "Bravo !" << std::endl;
    else
        std::cout << "Perdu, recommence !" << std::endl;
    return EXIT_SUCCESS;
}
