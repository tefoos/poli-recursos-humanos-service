package main

import (
	"database/sql"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"hr-system/shared"
)

type EmpleadoCrud struct {
	db *sql.DB
}

func NewEmpleadoCrud(db *sql.DB) *EmpleadoCrud {
	return &EmpleadoCrud{db: db}
}

func (c *EmpleadoCrud) validateCreateEmpleado(dto shared.CreateEmpleadoDTO) error {
	if strings.TrimSpace(dto.PrimerNombre) == "" {
		return fmt.Errorf("primer nombre es requerido")
	}

	if len(dto.PrimerNombre) > 50 {
		return fmt.Errorf("primer nombre no puede exceder 50 caracteres")
	}

	if dto.SegundoNombre != nil && len(*dto.SegundoNombre) > 50 {
		return fmt.Errorf("segundo nombre no puede exceder 50 caracteres")
	}

	if strings.TrimSpace(dto.Email) == "" {
		return fmt.Errorf("email es requerido")
	}

	if len(dto.Email) > 100 {
		return fmt.Errorf("email no puede exceder 100 caracteres")
	}

	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(dto.Email) {
		return fmt.Errorf("formato de email inválido")
	}

	if strings.TrimSpace(dto.FechaNac) == "" {
		return fmt.Errorf("fecha de nacimiento es requerida")
	}

	_, err := time.Parse("2006-01-02", dto.FechaNac)
	if err != nil {
		return fmt.Errorf("formato de fecha inválido, use YYYY-MM-DD")
	}

	if dto.Sueldo <= 0 {
		return fmt.Errorf("sueldo debe ser mayor a 0")
	}

	if dto.Comision < 0 || dto.Comision > 100 {
		return fmt.Errorf("comisión debe estar entre 0 y 100")
	}

	if dto.CargoID <= 0 {
		return fmt.Errorf("cargo ID es requerido y debe ser mayor a 0")
	}

	if dto.DptoID <= 0 {
		return fmt.Errorf("departamento ID es requerido y debe ser mayor a 0")
	}

	if dto.GerenteID != nil && *dto.GerenteID <= 0 {
		return fmt.Errorf("gerente ID debe ser mayor a 0 si se proporciona")
	}

	return nil
}

func (c *EmpleadoCrud) validateUpdateEmpleado(dto shared.UpdateEmpleadoDTO) error {
	if dto.ID <= 0 {
		return fmt.Errorf("ID del empleado es requerido y debe ser mayor a 0")
	}

	createDTO := shared.CreateEmpleadoDTO{
		PrimerNombre:  dto.PrimerNombre,
		SegundoNombre: dto.SegundoNombre,
		Email:         dto.Email,
		FechaNac:      dto.FechaNac,
		Sueldo:        dto.Sueldo,
		Comision:      dto.Comision,
		CargoID:       dto.CargoID,
		GerenteID:     dto.GerenteID,
		DptoID:        dto.DptoID,
	}

	return c.validateCreateEmpleado(createDTO)
}

func (c *EmpleadoCrud) Insert(dto shared.CreateEmpleadoDTO) (*shared.CreateEmpleadoResponseDTO, error) {
	if err := c.validateCreateEmpleado(dto); err != nil {
		return nil, fmt.Errorf("validación fallida: %v", err)
	}

	query := `
		INSERT INTO empleados (empl_primer_nombre, empl_segundo_nombre, empl_email,
		empl_fecha_nac, empl_sueldo, empl_comision, empl_cargo_id, empl_gerente_id, empl_dpto_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING empl_id`

	var newID int
	err := c.db.QueryRow(query, dto.PrimerNombre, dto.SegundoNombre, dto.Email,
		dto.FechaNac, dto.Sueldo, dto.Comision, dto.CargoID, dto.GerenteID, dto.DptoID).Scan(&newID)

	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			return nil, fmt.Errorf("email ya existe en el sistema")
		}
		if strings.Contains(err.Error(), "foreign key constraint") {
			return nil, fmt.Errorf("ID de cargo, gerente o departamento no válido")
		}
		return nil, fmt.Errorf("error insertando empleado: %v", err)
	}

	detailQuery := `
		SELECT e.empl_id, e.empl_primer_nombre, e.empl_segundo_nombre, e.empl_fecha_nac,
		       c.cargo_nombre,
		       d.dpto_nombre,
		       CASE WHEN g.empl_id IS NOT NULL
		            THEN CONCAT(g.empl_primer_nombre, ' ', COALESCE(g.empl_segundo_nombre, ''))
		            ELSE NULL
		       END as gerente_nombre,
		       e.empl_sueldo, e.empl_comision
		FROM empleados e
		INNER JOIN cargos c ON e.empl_cargo_id = c.cargo_id
		INNER JOIN departamentos d ON e.empl_dpto_id = d.dpto_id
		LEFT JOIN empleados g ON e.empl_gerente_id = g.empl_id
		WHERE e.empl_id = $1`

	var response shared.CreateEmpleadoResponseDTO
	err = c.db.QueryRow(detailQuery, newID).Scan(
		&response.ID, &response.PrimerNombre, &response.SegundoNombre, &response.FechaNac,
		&response.CargoNombre, &response.DepartamentoNombre, &response.GerenteNombre,
		&response.Sueldo, &response.Comision)

	if err != nil {
		return nil, fmt.Errorf("error obteniendo detalles del empleado creado: %v", err)
	}

	return &response, nil
}

func (c *EmpleadoCrud) Update(dto shared.UpdateEmpleadoDTO) (*shared.UpdateEmpleadoResponseDTO, error) {
	if err := c.validateUpdateEmpleado(dto); err != nil {
		return nil, fmt.Errorf("validación fallida: %v", err)
	}

	query := `
		UPDATE empleados
		SET empl_primer_nombre=$1, empl_segundo_nombre=$2, empl_email=$3,
		    empl_fecha_nac=$4, empl_sueldo=$5, empl_comision=$6,
		    empl_cargo_id=$7, empl_gerente_id=$8, empl_dpto_id=$9
		WHERE empl_id=$10 AND is_deleted=false`

	result, err := c.db.Exec(query, dto.PrimerNombre, dto.SegundoNombre, dto.Email,
		dto.FechaNac, dto.Sueldo, dto.Comision, dto.CargoID, dto.GerenteID, dto.DptoID, dto.ID)

	if err != nil {
		if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
			return nil, fmt.Errorf("email ya existe en el sistema")
		}
		if strings.Contains(err.Error(), "foreign key constraint") {
			return nil, fmt.Errorf("ID de cargo, gerente o departamento no válido")
		}
		return nil, fmt.Errorf("error actualizando empleado: %v", err)
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return nil, fmt.Errorf("empleado no encontrado o ya está eliminado")
	}

	detailQuery := `
		SELECT e.empl_id, e.empl_primer_nombre, e.empl_segundo_nombre, e.empl_fecha_nac,
		       c.cargo_nombre,
		       d.dpto_nombre,
		       CASE WHEN g.empl_id IS NOT NULL
		            THEN CONCAT(g.empl_primer_nombre, ' ', COALESCE(g.empl_segundo_nombre, ''))
		            ELSE NULL
		       END as gerente_nombre,
		       e.empl_sueldo, e.empl_comision
		FROM empleados e
		INNER JOIN cargos c ON e.empl_cargo_id = c.cargo_id
		INNER JOIN departamentos d ON e.empl_dpto_id = d.dpto_id
		LEFT JOIN empleados g ON e.empl_gerente_id = g.empl_id
		WHERE e.empl_id = $1`

	var response shared.UpdateEmpleadoResponseDTO
	err = c.db.QueryRow(detailQuery, dto.ID).Scan(
		&response.ID, &response.PrimerNombre, &response.SegundoNombre, &response.FechaNac,
		&response.CargoNombre, &response.DepartamentoNombre, &response.GerenteNombre,
		&response.Sueldo, &response.Comision)

	if err != nil {
		return nil, fmt.Errorf("error obteniendo detalles del empleado actualizado: %v", err)
	}

	return &response, nil
}

func (c *EmpleadoCrud) Select(id int) (*shared.EmpleadoDetailResponseDTO, error) {
	if id <= 0 {
		return nil, fmt.Errorf("ID debe ser mayor a 0")
	}

	query := `
		SELECT e.empl_id, e.empl_primer_nombre, e.empl_segundo_nombre, e.empl_email,
		       e.empl_fecha_nac, e.empl_sueldo, e.empl_comision,
		       c.cargo_nombre,
		       CASE WHEN g.empl_id IS NOT NULL
		            THEN CONCAT(g.empl_primer_nombre, ' ', COALESCE(g.empl_segundo_nombre, ''))
		            ELSE NULL
		       END as gerente_nombre,
		       d.dpto_nombre,
		       e.is_deleted
		FROM empleados e
		INNER JOIN cargos c ON e.empl_cargo_id = c.cargo_id
		INNER JOIN departamentos d ON e.empl_dpto_id = d.dpto_id
		LEFT JOIN empleados g ON e.empl_gerente_id = g.empl_id
		WHERE e.empl_id=$1`

	var emp shared.EmpleadoDetailResponseDTO
	err := c.db.QueryRow(query, id).Scan(
		&emp.ID, &emp.PrimerNombre, &emp.SegundoNombre, &emp.Email,
		&emp.FechaNac, &emp.Sueldo, &emp.Comision,
		&emp.CargoNombre, &emp.GerenteNombre, &emp.DepartamentoNombre,
		&emp.IsDeleted)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("empleado no encontrado")
		}
		return nil, fmt.Errorf("error consultando empleado: %v", err)
	}

	return &emp, nil
}

func (c *EmpleadoCrud) Delete(id int) error {
	if id <= 0 {
		return errors.New("ID debe ser mayor a 0")
	}

	query := `SELECT success, message FROM p_delete_empleado($1)`

	var success bool
	var message string

	err := c.db.QueryRow(query, id).Scan(&success, &message)
	if err != nil {
		return fmt.Errorf("error ejecutando procedimiento almacenado: %v", err)
	}

	if !success {
		return errors.New(message)
	}

	return nil
}

func (c *EmpleadoCrud) ListCargos() ([]shared.CargoDTO, error) {
	query := `SELECT cargo_id, cargo_nombre FROM cargos ORDER BY cargo_id`

	rows, err := c.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error consultando cargos: %v", err)
	}
	defer rows.Close()

	var cargos []shared.CargoDTO
	for rows.Next() {
		var cargo shared.CargoDTO
		err := rows.Scan(&cargo.ID, &cargo.Nombre)
		if err != nil {
			return nil, fmt.Errorf("error escaneando cargo: %v", err)
		}
		cargos = append(cargos, cargo)
	}

	return cargos, nil
}

func (c *EmpleadoCrud) ListDepartamentos() ([]shared.DepartamentoDTO, error) {
	query := `SELECT dpto_id, dpto_nombre FROM departamentos ORDER BY dpto_id`

	rows, err := c.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error consultando departamentos: %v", err)
	}
	defer rows.Close()

	var departamentos []shared.DepartamentoDTO
	for rows.Next() {
		var dpto shared.DepartamentoDTO
		err := rows.Scan(&dpto.ID, &dpto.Nombre)
		if err != nil {
			return nil, fmt.Errorf("error escaneando departamento: %v", err)
		}
		departamentos = append(departamentos, dpto)
	}

	return departamentos, nil
}

func (c *EmpleadoCrud) ListGerentes() ([]shared.GerenteDTO, error) {
	query := `
		SELECT empl_id, CONCAT(empl_primer_nombre, ' ', COALESCE(empl_segundo_nombre, '')) as nombre_completo
		FROM empleados
		WHERE is_deleted=false
		ORDER BY empl_id`

	rows, err := c.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("error consultando gerentes: %v", err)
	}
	defer rows.Close()

	var gerentes []shared.GerenteDTO
	for rows.Next() {
		var gerente shared.GerenteDTO
		err := rows.Scan(&gerente.ID, &gerente.Nombre)
		if err != nil {
			return nil, fmt.Errorf("error escaneando gerente: %v", err)
		}
		gerentes = append(gerentes, gerente)
	}

	return gerentes, nil
}
