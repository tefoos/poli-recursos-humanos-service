CREATE INDEX idx_ciudades_pais ON ciudades(ciud_pais_ID);
CREATE INDEX idx_localizaciones_ciudad ON localizaciones(localiz_ciudad_ID);
CREATE INDEX idx_departamentos_localiz ON departamentos(dpto_localiz_ID);
CREATE INDEX idx_empleados_cargo ON empleados(empl_cargo_ID);
CREATE INDEX idx_empleados_gerente ON empleados(empl_Gerente_ID);
CREATE INDEX idx_empleados_dpto ON empleados(empl_dpto_ID);
CREATE INDEX idx_historico_cargo ON historico(emphist_cargo_ID);
CREATE INDEX idx_historico_dpto ON historico(emphist_dpto_ID);

CREATE INDEX idx_empleados_email ON empleados(empl_email);
CREATE INDEX idx_empleados_nombre ON empleados(empl_primer_nombre);
CREATE INDEX idx_departamentos_nombre ON departamentos(dpto_nombre);
CREATE INDEX idx_cargos_nombre ON cargos(cargo_nombre);
CREATE INDEX idx_paises_nombre ON paises(pais_nombre);
CREATE INDEX idx_ciudades_nombre ON ciudades(ciud_nombre);

CREATE INDEX idx_empleados_dpto_cargo ON empleados(empl_dpto_ID, empl_cargo_ID);
CREATE INDEX idx_historico_fecha_cargo ON historico(emphist_fecha_retiro, emphist_cargo_ID);
