#ifndef     LA_SORTIE
# define    LA_SORTIE

extern "C" {

    typedef struct {
        int x;
        int y;
    } position;

    typedef enum {
        SOL = 0,
        MUR = 1,
        ARRIVEE = 2
    } type;

    void    la_sortie(type const [10][10], position const *);

}

#endif  /* LA_SORTIE */
