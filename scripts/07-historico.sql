CREATE TABLE historico (
    emphist_ID SERIAL PRIMARY KEY,
    emphist_fecha_retiro DATE NOT NULL DEFAULT CURRENT_DATE,
    emphist_cargo_ID INTEGER NOT NULL,
    emphist_dpto_ID INTEGER NOT NULL,
    FOREIGN KEY (emphist_cargo_ID) REFERENCES cargos(cargo_ID) ON DELETE CASCADE,
    FOREIGN KEY (emphist_dpto_ID) REFERENCES departamentos(dpto_ID) ON DELETE CASCADE
);
