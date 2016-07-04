void    les_meilleures(int resultats[2], int a, int b, int c) {
    /* A supprimer */
    if (a > b && a > c) {
        resultats[0] = a;
        if (b > c)
            resultats[1] = b;
        else
            resultats[1] = c;
    } else if (b > a && b > c) {
        resultats[0] = b;
        if (a > c)
            resultats[1] = c;
        else
            resultats[1] = a;
    } else {
        resultats[0] = c;
        if (a > b)
            resultats[1] = a;
        else
            resultats[1] = b;
    }
    /* A supprimer */
}
