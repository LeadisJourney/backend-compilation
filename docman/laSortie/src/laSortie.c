typedef struct {
    int x;
    int y;
} position;

typedef enum {
    SOL = 0,
    MUR = 1,
    ARRIVEE = 2
} type;

void    deplacement_haut();
void    deplacement_bas();
void    deplacement_gauche();
void    deplacement_droite();

/* Exemple :
  0 : vide
  1 : mur
  2 : arrivee
  
  
  1 1 1 1 1 1 1 1 1 1
  1 0 0 0 2 0 0 0 0 1
  1 0 0 0 0 0 0 0 0 1
  1 0 1 1 1 1 1 0 0 1
  1 0 0 0 1 0 0 0 0 1
  1 0 0 0 1 0 0 0 0 1
  1 1 1 0 1 0 1 1 1 1
  1 0 0 0 1 0 0 0 0 1
  1 0 0 0 0 0 0 0 0 1
  1 1 1 1 1 1 1 1 1 1
 */

/* A supprimer */
#include <stdlib.h>
#include <time.h>
/* A supprimer */
void    la_sortie(type const labyrinthe[10][10], position const *leadis) {

    /* A supprimer */
    srand(time(NULL));
    (void)labyrinthe;
    (void)leadis;
    if (rand() % 2)
        deplacement_gauche();
    else
        deplacement_droite();
    deplacement_haut();
    deplacement_haut();
    deplacement_haut();
    deplacement_gauche();
    deplacement_gauche();
    deplacement_haut();
    deplacement_haut();
    deplacement_haut();
    deplacement_droite();
    deplacement_droite();
    deplacement_droite();
    deplacement_haut();
    /* A supprimer */
}
