template<>
maze<10, 10>::maze() :
    _dat{{WALL  , WALL  , WALL  , WALL  , WALL  , WALL  , WALL  , WALL  , WALL  , WALL  },
         {WALL  , GROUND, GROUND, GROUND, END   , GROUND, GROUND, GROUND, GROUND, WALL  },
         {WALL  , GROUND, GROUND, GROUND, GROUND, GROUND, GROUND, GROUND, GROUND, WALL  },
         {WALL  , GROUND, WALL  , WALL  , WALL  , WALL  , WALL  , WALL  , GROUND, WALL  },
         {WALL  , GROUND, GROUND, GROUND, WALL  , GROUND, GROUND, GROUND, GROUND, WALL  },
         {WALL  , GROUND, GROUND, GROUND, WALL  , GROUND, GROUND, GROUND, GROUND, WALL  },
         {WALL  , WALL  , WALL  , GROUND, WALL  , GROUND, WALL  , WALL  , WALL  , WALL  },
         {WALL  , GROUND, GROUND, GROUND, WALL  , GROUND, GROUND, GROUND, GROUND, WALL  },
         {WALL  , GROUND, GROUND, GROUND, GROUND, GROUND, GROUND, GROUND, GROUND, WALL  },
         {WALL  , WALL  , WALL  , WALL  , WALL  , WALL  , WALL  , WALL  , WALL  , WALL  } }
{
}

template<>
bool                                maze<10, 10>::valid(position const &sta __attribute__((unused))) {
    return true;
}

template <unsigned int W, unsigned int H>
void                                maze<W, H>::generate() {
    /* unimplemented */
}

template <unsigned int W, unsigned int H>
typename maze<W, H>::type const   (&maze<W, H>::data() const )[W][H] {
    return this->_dat;
}

template <unsigned int W, unsigned int H>
auto                                maze<W, H>::operator[](unsigned int i) /* -> proxy_type */ {
    /* unimplemented */
}
