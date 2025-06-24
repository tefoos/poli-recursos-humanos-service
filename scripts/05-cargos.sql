CREATE TABLE cargos (
    cargo_ID SERIAL PRIMARY KEY,
    cargo_nombre VARCHAR(100) NOT NULL UNIQUE,
    cargo_sueldo_minimo DECIMAL(10,2) NOT NULL,
    cargo_sueldo_maximo DECIMAL(10,2) NOT NULL,
    CHECK (cargo_sueldo_maximo >= cargo_sueldo_minimo)
);
