package shared

type CreateEmpleadoDTO struct {
	PrimerNombre  string  `json:"empl_primer_nombre"`
	SegundoNombre *string `json:"empl_segundo_nombre"`
	Email         string  `json:"empl_email"`
	FechaNac      string  `json:"empl_fecha_nac"`
	Sueldo        float64 `json:"empl_sueldo"`
	Comision      float64 `json:"empl_comision"`
	CargoID       int     `json:"empl_cargo_id"`
	GerenteID     *int    `json:"empl_gerente_id"`
	DptoID        int     `json:"empl_dpto_id"`
}

type UpdateEmpleadoDTO struct {
	ID            int     `json:"empl_id"`
	PrimerNombre  string  `json:"empl_primer_nombre"`
	SegundoNombre *string `json:"empl_segundo_nombre"`
	Email         string  `json:"empl_email"`
	FechaNac      string  `json:"empl_fecha_nac"`
	Sueldo        float64 `json:"empl_sueldo"`
	Comision      float64 `json:"empl_comision"`
	CargoID       int     `json:"empl_cargo_id"`
	GerenteID     *int    `json:"empl_gerente_id"`
	DptoID        int     `json:"empl_dpto_id"`
}

type SelectEmpleadoDTO struct {
	ID int `json:"empl_id"`
}

type DeleteEmpleadoDTO struct {
	ID int `json:"empl_id"`
}

type EmpleadoResponseDTO struct {
	ID            int     `json:"empl_id"`
	PrimerNombre  string  `json:"empl_primer_nombre"`
	SegundoNombre *string `json:"empl_segundo_nombre"`
	Email         string  `json:"empl_email"`
	FechaNac      string  `json:"empl_fecha_nac"`
	Sueldo        float64 `json:"empl_sueldo"`
	Comision      float64 `json:"empl_comision"`
	CargoID       int     `json:"empl_cargo_id"`
	GerenteID     *int    `json:"empl_gerente_id"`
	DptoID        int     `json:"empl_dpto_id"`
	IsDeleted     bool    `json:"is_deleted"`
}

type EmpleadoDetailResponseDTO struct {
	ID                 int     `json:"empl_id"`
	PrimerNombre       string  `json:"empl_primer_nombre"`
	SegundoNombre      *string `json:"empl_segundo_nombre"`
	Email              string  `json:"empl_email"`
	FechaNac           string  `json:"empl_fecha_nac"`
	Sueldo             float64 `json:"empl_sueldo"`
	Comision           float64 `json:"empl_comision"`
	CargoNombre        string  `json:"cargo_nombre"`
	GerenteNombre      *string `json:"gerente_nombre"`
	DepartamentoNombre string  `json:"departamento_nombre"`
	IsDeleted          bool    `json:"is_deleted"`
}

type CargoDTO struct {
	ID     int    `json:"cargo_id"`
	Nombre string `json:"cargo_nombre"`
}

type DepartamentoDTO struct {
	ID     int    `json:"dpto_id"`
	Nombre string `json:"dpto_nombre"`
}

type GerenteDTO struct {
	ID     int    `json:"empl_id"`
	Nombre string `json:"nombre_completo"`
}

type CreateEmpleadoResponseDTO struct {
	ID                 int     `json:"empl_id"`
	PrimerNombre       string  `json:"empl_primer_nombre"`
	SegundoNombre      *string `json:"empl_segundo_nombre"`
	FechaNac           string  `json:"empl_fecha_nac"`
	CargoNombre        string  `json:"cargo_nombre"`
	DepartamentoNombre string  `json:"departamento_nombre"`
	GerenteNombre      *string `json:"gerente_nombre"`
	Sueldo             float64 `json:"empl_sueldo"`
	Comision           float64 `json:"empl_comision"`
}

type UpdateEmpleadoResponseDTO struct {
	ID                 int     `json:"empl_id"`
	PrimerNombre       string  `json:"empl_primer_nombre"`
	SegundoNombre      *string `json:"empl_segundo_nombre"`
	FechaNac           string  `json:"empl_fecha_nac"`
	CargoNombre        string  `json:"cargo_nombre"`
	DepartamentoNombre string  `json:"departamento_nombre"`
	GerenteNombre      *string `json:"gerente_nombre"`
	Sueldo             float64 `json:"empl_sueldo"`
	Comision           float64 `json:"empl_comision"`
}

type Request struct {
	Operation string      `json:"operation"`
	Data      interface{} `json:"data"`
}

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}
