int la_meilleure(int a, int b, int c) {
    return (a > b) ? (a > c ? c : a) : (b > c ? b : c); // à supprimer
    // A decommenter : return /* a, b, ou c*/;
}
