CREATE TABLE departamentos (
    dpto_ID SERIAL PRIMARY KEY,
    dpto_localiz_ID INTEGER NOT NULL,
    dpto_nombre VARCHAR(100) NOT NULL,
    FOREIGN KEY (dpto_localiz_ID) REFERENCES localizaciones(localiz_ID) ON DELETE CASCADE
);
