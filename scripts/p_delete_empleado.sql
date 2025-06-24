CREATE OR REPLACE FUNCTION p_delete_empleado(p_empl_id INTEGER)
RETURNS TABLE(
    success BOOLEAN,
    message TEXT,
    empl_id INTEGER,
    empl_primer_nombre VARCHAR(50),
    empl_segundo_nombre VARCHAR(50),
    empl_cargo_id INTEGER,
    empl_dpto_id INTEGER
)
LANGUAGE plpgsql
AS $$
DECLARE
    v_empleado_exists BOOLEAN;
    v_empleado_record RECORD;
BEGIN
    SELECT EXISTS(
        SELECT 1 FROM empleados
        WHERE empleados.empl_id = p_empl_id AND empleados.is_deleted = false
    ) INTO v_empleado_exists;

    IF NOT v_empleado_exists THEN
        RETURN QUERY SELECT
            false::BOOLEAN,
            'Empleado no encontrado o ya está eliminado'::TEXT,
            NULL::INTEGER,
            NULL::VARCHAR(50),
            NULL::VARCHAR(50),
            NULL::INTEGER,
            NULL::INTEGER;
        RETURN;
    END IF;

    SELECT empleados.empl_id, empleados.empl_primer_nombre, empleados.empl_segundo_nombre,
           empleados.empl_cargo_id, empleados.empl_dpto_id
    INTO v_empleado_record
    FROM empleados
    WHERE empleados.empl_id = p_empl_id AND empleados.is_deleted = false;

    INSERT INTO historico (emphist_cargo_id, emphist_dpto_id)
    VALUES (v_empleado_record.empl_cargo_id, v_empleado_record.empl_dpto_id);

    UPDATE empleados
    SET is_deleted = true
    WHERE empleados.empl_id = p_empl_id;

    RETURN QUERY SELECT
        true::BOOLEAN,
        'Empleado eliminado exitosamente y guardado en histórico'::TEXT,
        v_empleado_record.empl_id,
        v_empleado_record.empl_primer_nombre,
        v_empleado_record.empl_segundo_nombre,
        v_empleado_record.empl_cargo_id,
        v_empleado_record.empl_dpto_id;

EXCEPTION
    WHEN OTHERS THEN
        RETURN QUERY SELECT
            false::BOOLEAN,
            ('Error eliminando empleado: ' || SQLERRM)::TEXT,
            NULL::INTEGER,
            NULL::VARCHAR(50),
            NULL::VARCHAR(50),
            NULL::INTEGER,
            NULL::INTEGER;
END;
$$;
