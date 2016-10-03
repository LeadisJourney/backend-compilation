#ifndef     MAZE_HPP
# define    MAZE_HPP

#include "laSortie.h"

template <unsigned int W, unsigned int H>
class maze {
    static_assert(W > 3 && H > 3, "Maze's size too small!");
public:
    enum type : int {
        GROUND = ::SOL,
        WALL = ::MUR,
        END = ::ARRIVEE
    };
private:
    type    _dat[H][W];

public:
    maze();
    bool            valid(position const &);
    void            generate();
    type const    (&data() const)[W][H] ;
    auto            operator[](unsigned int);
};

# include "maze.tpp"

#endif /* MAZE_HPP */
