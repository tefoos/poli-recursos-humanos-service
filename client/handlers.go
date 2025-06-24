package main

import (
	"encoding/json"
	"fmt"
	"strings"

	"hr-system/shared"
)

func (c *Client) HandleInsert() {
	fmt.Println("\n--- CREAR EMPLEADO ---")

	primerNombre := c.ReadValidatedName("Primer nombre: ", true)
	segundoNombre := c.ReadValidatedName("Segundo nombre (opcional): ", false)
	email := c.ReadEmailInput("Email: ")
	fechaNac := c.ReadDateInput("Fecha de nacimiento (YYYY-MM-DD): ")
	sueldo := c.ReadPositiveFloatInput("Sueldo: ")
	comision := c.ReadRangeFloatInput("Comisión (%): ", 0, 100)

	cargoID, err := c.SelectCargoFromList()
	if err != nil {
		fmt.Printf("Error obteniendo cargos: %v\n", err)
		return
	}

	gerenteID, err := c.SelectGerenteFromList()
	if err != nil {
		fmt.Printf("Error obteniendo gerentes: %v\n", err)
		return
	}

	dptoID, err := c.SelectDepartamentoFromList()
	if err != nil {
		fmt.Printf("Error obteniendo departamentos: %v\n", err)
		return
	}

	dto := shared.CreateEmpleadoDTO{
		PrimerNombre:  *primerNombre,
		SegundoNombre: segundoNombre,
		Email:         email,
		FechaNac:      fechaNac,
		Sueldo:        sueldo,
		Comision:      comision,
		CargoID:       cargoID,
		GerenteID:     gerenteID,
		DptoID:        dptoID,
	}

	req := shared.Request{
		Operation: "INSERT",
		Data:      dto,
	}

	response, err := c.SendRequest(req)
	if err != nil {
		fmt.Printf("Error enviando petición: %v\n", err)
		return
	}

	c.PrintResponse(response)
}

func (c *Client) HandleUpdate() {
	fmt.Println("\n--- ACTUALIZAR EMPLEADO ---")

	empleadoID, err := c.ReadIntInput("ID del empleado a actualizar: ")
	if err != nil {
		fmt.Printf("Error en ID: %v\n", err)
		return
	}

	fmt.Println("Obteniendo datos actuales del empleado...")
	currentData, err := c.GetEmpleadoActual(empleadoID)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("\nDatos actuales del empleado:\n")
	fmt.Printf("1. Primer nombre: %s\n", currentData.PrimerNombre)
	if currentData.SegundoNombre != nil {
		fmt.Printf("2. Segundo nombre: %s\n", *currentData.SegundoNombre)
	} else {
		fmt.Printf("2. Segundo nombre: (vacío)\n")
	}
	fmt.Printf("3. Email: %s\n", currentData.Email)
	fmt.Printf("4. Fecha nacimiento: %s\n", currentData.FechaNac)
	fmt.Printf("5. Cargo: %s\n", currentData.CargoNombre)
	fmt.Printf("6. Departamento: %s\n", currentData.DepartamentoNombre)
	if currentData.GerenteNombre != nil {
		fmt.Printf("7. Gerente: %s\n", *currentData.GerenteNombre)
	} else {
		fmt.Printf("7. Gerente: (sin gerente)\n")
	}
	fmt.Printf("8. Sueldo: %.2f\n", currentData.Sueldo)
	fmt.Printf("9. Comisión: %.2f\n", currentData.Comision)

	campos := c.ReadInput("\n¿Qué campos desea actualizar? (ej: 1,3,5 o 'todo'): ")

	if strings.ToLower(strings.TrimSpace(campos)) == "todo" {
		c.UpdateAllFields(empleadoID)
		return
	}

	c.UpdateSelectedFields(empleadoID, campos, currentData)
}

func (c *Client) GetEmpleadoActual(empleadoID int) (*shared.EmpleadoDetailResponseDTO, error) {
	dto := shared.SelectEmpleadoDTO{
		ID: empleadoID,
	}

	req := shared.Request{
		Operation: "SELECT",
		Data:      dto,
	}

	response, err := c.SendRequest(req)
	if err != nil {
		return nil, err
	}

	if !response.Success {
		return nil, fmt.Errorf(response.Message)
	}

	dataBytes, _ := json.Marshal(response.Data)
	var empleado shared.EmpleadoDetailResponseDTO
	err = json.Unmarshal(dataBytes, &empleado)
	if err != nil {
		return nil, fmt.Errorf("error procesando datos del empleado: %v", err)
	}

	if empleado.IsDeleted {
		return nil, fmt.Errorf("no se puede actualizar un empleado eliminado")
	}

	if len(empleado.FechaNac) > 10 {
		empleado.FechaNac = empleado.FechaNac[:10]
	}

	return &empleado, nil
}

func (c *Client) UpdateSelectedFields(empleadoID int, campos string, current *shared.EmpleadoDetailResponseDTO) {
	fieldsToUpdate := strings.Split(campos, ",")

	primerNombre := current.PrimerNombre
	var segundoNombre *string = current.SegundoNombre
	email := current.Email
	fechaNac := current.FechaNac
	sueldo := current.Sueldo
	comision := current.Comision

	cargoID := 0
	gerenteID := current.GerenteNombre
	dptoID := 0

	for _, field := range fieldsToUpdate {
		switch strings.TrimSpace(field) {
		case "1":
			primerNombre = *c.ReadValidatedName("Nuevo primer nombre: ", true)
		case "2":
			segundoNombre = c.ReadValidatedName("Nuevo segundo nombre (opcional): ", false)
		case "3":
			email = c.ReadEmailInput("Nuevo email: ")
		case "4":
			fechaNac = c.ReadDateInput("Nueva fecha de nacimiento (YYYY-MM-DD): ")
		case "5":
			var err error
			cargoID, err = c.SelectCargoFromList()
			if err != nil {
				fmt.Printf("Error obteniendo cargos: %v\n", err)
				return
			}
		case "6":
			var err error
			dptoID, err = c.SelectDepartamentoFromList()
			if err != nil {
				fmt.Printf("Error obteniendo departamentos: %v\n", err)
				return
			}
		case "7":
			var err error
			gerenteID, err = c.SelectGerenteFromListForUpdate()
			if err != nil {
				fmt.Printf("Error obteniendo gerentes: %v\n", err)
				return
			}
		case "8":
			sueldo = c.ReadPositiveFloatInput("Nuevo sueldo: ")
		case "9":
			comision = c.ReadRangeFloatInput("Nueva comisión (%): ", 0, 100)
		default:
			fmt.Printf("Campo '%s' no válido. Use números del 1-9.\n", field)
			return
		}
	}

	if cargoID == 0 {
		cargoID = c.GetCargoIDByName(current.CargoNombre)
	}
	if dptoID == 0 {
		dptoID = c.GetDptoIDByName(current.DepartamentoNombre)
	}

	var finalGerenteID *int
	if gerenteID != nil && *gerenteID != "" {
		id := c.GetGerenteIDByName(*gerenteID)
		finalGerenteID = &id
	}

	dto := shared.UpdateEmpleadoDTO{
		ID:            empleadoID,
		PrimerNombre:  primerNombre,
		SegundoNombre: segundoNombre,
		Email:         email,
		FechaNac:      fechaNac,
		Sueldo:        sueldo,
		Comision:      comision,
		CargoID:       cargoID,
		GerenteID:     finalGerenteID,
		DptoID:        dptoID,
	}

	req := shared.Request{
		Operation: "UPDATE",
		Data:      dto,
	}

	response, err := c.SendRequest(req)
	if err != nil {
		fmt.Printf("Error enviando petición: %v\n", err)
		return
	}

	c.PrintResponse(response)
}

func (c *Client) UpdateAllFields(empleadoID int) {
	primerNombre := c.ReadValidatedName("Primer nombre: ", true)
	segundoNombre := c.ReadValidatedName("Segundo nombre (opcional): ", false)
	email := c.ReadEmailInput("Email: ")
	fechaNac := c.ReadDateInput("Fecha de nacimiento (YYYY-MM-DD): ")
	sueldo := c.ReadPositiveFloatInput("Sueldo: ")
	comision := c.ReadRangeFloatInput("Comisión (%): ", 0, 100)

	cargoID, err := c.SelectCargoFromList()
	if err != nil {
		fmt.Printf("Error obteniendo cargos: %v\n", err)
		return
	}

	gerenteID, err := c.SelectGerenteFromList()
	if err != nil {
		fmt.Printf("Error obteniendo gerentes: %v\n", err)
		return
	}

	dptoID, err := c.SelectDepartamentoFromList()
	if err != nil {
		fmt.Printf("Error obteniendo departamentos: %v\n", err)
		return
	}

	dto := shared.UpdateEmpleadoDTO{
		ID:            empleadoID,
		PrimerNombre:  *primerNombre,
		SegundoNombre: segundoNombre,
		Email:         email,
		FechaNac:      fechaNac,
		Sueldo:        sueldo,
		Comision:      comision,
		CargoID:       cargoID,
		GerenteID:     gerenteID,
		DptoID:        dptoID,
	}

	req := shared.Request{
		Operation: "UPDATE",
		Data:      dto,
	}

	response, err := c.SendRequest(req)
	if err != nil {
		fmt.Printf("Error enviando petición: %v\n", err)
		return
	}

	c.PrintResponse(response)
}

func (c *Client) HandleSelect() {
	fmt.Println("\n--- CONSULTAR EMPLEADO ---")

	empleadoID, err := c.ReadIntInput("ID del empleado: ")
	if err != nil {
		fmt.Printf("Error en ID: %v\n", err)
		return
	}

	dto := shared.SelectEmpleadoDTO{
		ID: empleadoID,
	}

	req := shared.Request{
		Operation: "SELECT",
		Data:      dto,
	}

	response, err := c.SendRequest(req)
	if err != nil {
		fmt.Printf("Error enviando petición: %v\n", err)
		return
	}

	c.PrintResponse(response)
}

func (c *Client) HandleDelete() {
	fmt.Println("\n--- ELIMINAR EMPLEADO ---")

	empleadoID, err := c.ReadIntInput("ID del empleado a eliminar: ")
	if err != nil {
		fmt.Printf("Error en ID: %v\n", err)
		return
	}

	confirmacion := c.ReadInput("¿Está seguro? (s/N): ")
	if strings.ToLower(confirmacion) != "s" {
		fmt.Println("Operación cancelada")
		return
	}

	dto := shared.DeleteEmpleadoDTO{
		ID: empleadoID,
	}

	req := shared.Request{
		Operation: "DELETE",
		Data:      dto,
	}

	response, err := c.SendRequest(req)
	if err != nil {
		fmt.Printf("Error enviando petición: %v\n", err)
		return
	}

	c.PrintResponse(response)
}
