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
    int                                 gr1;
    int                                 gr2;

    do {
        a = dis(gen);
        b = dis(gen);
        c = dis(gen);
    } while (a == b && b == c && a == c);
    gr1 = std::max(a, gr2 = std::max(b, c));
    gr2 = std::min(std::max(std::min(b, c), a), gr2);
    
    
    les_meilleures(res, a, b, c);
    if ((res[0] + res[1]) == (a + b + c - std::min(std::min(a, b), c)))
        std::cout << "Bravo !" << std::endl;
    else
        std::cout << "Perdu, recommence !" << std::endl;
    std::cout << "\"Graphic\": [" << a << ", " << b << ", " << c << ", " << res[0] << ", " << res[1]  << ", " << gr1 << ", " << gr2 << "]" << std::endl;
    return EXIT_SUCCESS;
}
