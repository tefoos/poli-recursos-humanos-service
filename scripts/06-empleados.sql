CREATE TABLE empleados (
    empl_ID SERIAL PRIMARY KEY,
    empl_primer_nombre VARCHAR(50) NOT NULL,
    empl_segundo_nombre VARCHAR(50),
    empl_email VARCHAR(100) NOT NULL UNIQUE,
    empl_fecha_nac DATE NOT NULL,
    empl_sueldo DECIMAL(10,2) NOT NULL,
    empl_comision DECIMAL(5,2) DEFAULT 0.00,
    empl_cargo_ID INTEGER NOT NULL,
    empl_Gerente_ID INTEGER,
    empl_dpto_ID INTEGER NOT NULL,
    is_deleted BOOLEAN DEFAULT FALSE,
    FOREIGN KEY (empl_cargo_ID) REFERENCES cargos(cargo_ID) ON DELETE RESTRICT,
    FOREIGN KEY (empl_Gerente_ID) REFERENCES empleados(empl_ID) ON DELETE SET NULL,
    FOREIGN KEY (empl_dpto_ID) REFERENCES departamentos(dpto_ID) ON DELETE RESTRICT
);
