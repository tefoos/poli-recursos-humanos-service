CREATE TABLE localizaciones (
    localiz_ID SERIAL PRIMARY KEY,
    localiz_ciudad_ID INTEGER NOT NULL,
    localiz_direccion VARCHAR(255) NOT NULL,
    FOREIGN KEY (localiz_ciudad_ID) REFERENCES ciudades(ciud_ID) ON DELETE CASCADE
);
