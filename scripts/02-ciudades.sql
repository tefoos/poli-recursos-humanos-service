CREATE TABLE ciudades (
    ciud_ID SERIAL PRIMARY KEY,
    ciud_pais_ID INTEGER NOT NULL,
    ciud_nombre VARCHAR(100) NOT NULL,
    FOREIGN KEY (ciud_pais_ID) REFERENCES paises(pais_ID) ON DELETE CASCADE
);
