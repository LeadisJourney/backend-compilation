#include <iostream>

#include <cstdlib>

#include "maze.hpp"

#include "laSortie.h"

position        g_cur_pos;
type const    (*g_dat)[10][10];

extern "C" {

    void    deplacement_haut() {
        if ((*g_dat)[g_cur_pos.y - 1][g_cur_pos.x] != MUR) {    //
            std::cout << "\"up\", ";                            // Attention au (0, 0) en bas ou en haut (coordonnees cartesiennes) 
            --(g_cur_pos.y);                                    //
        }
    }
    void    deplacement_bas() {
        if ((*g_dat)[g_cur_pos.y + 1][g_cur_pos.x] != MUR) {    //
        std::cout << "\"dn\", ";                                // Attention au (0, 0) en bas ou en haut (coordonnees cartesiennes) 
            ++(g_cur_pos.y);                                    //
        }
    }
    void    deplacement_gauche() {
        if ((*g_dat)[g_cur_pos.y][g_cur_pos.x - 1] != MUR) {
            std::cout << "\"lt\", ";
            --(g_cur_pos.x);
        }
    }
    void    deplacement_droite() {
        if ((*g_dat)[g_cur_pos.y][g_cur_pos.x + 1] != MUR) {
            std::cout << "\"rt\", ";
            ++(g_cur_pos.x);
        }
    }

}

void                                    success() {
    if ((*g_dat)[g_cur_pos.y][g_cur_pos.x] == ARRIVEE)
        std:: cout << "\"success\"";
    else
        std:: cout << "\"fail\"";
}

int                                     main(int argc __attribute__((unused)), char **argv __attribute__((unused))) {
    maze<10, 10>                        maz;
    position                            lea{4, 8};
    unsigned int                        i = 0;

    g_cur_pos.x = lea.x;
    g_cur_pos.y = lea.y;
    g_dat = reinterpret_cast<type const (*)[10][10]>(&maz.data()) ;
    std::cout << "\"Graphic\": [{ mapX: 10, mapY: 10, leadisX: 4, leadisY: 8, leadisZ: 0, leadisMovement: [";
    la_sortie(reinterpret_cast<type const (* const)[10]>(maz.data()), &lea);
    success();
    std::cout << "], elemList: [";
    for (unsigned int y = 0; y < 10; ++y)
        for (unsigned int x = 0; x < 10; ++x)
            if (maz.data()[y][x] == decltype(maz)::WALL) { // Utiliser l'operateur[] quand il sera implemente
                std::cout << (i ? ", " : "") << "{name: \"Wall\", id: " << i << ", x: " << x << ", y: " << y << ", visible: true}";
                ++i;
            }
    std::cout << "] }]" << std::endl;
    return EXIT_SUCCESS;
}
