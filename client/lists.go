package main

import (
	"encoding/json"
	"fmt"
	"hr-system/shared"
	"strconv"
)

type CargoConDatosDTO struct {
	ID     int    `json:"cargo_id"`
	Nombre string `json:"cargo_nombre"`
}

type DepartamentoConDatosDTO struct {
	ID        int    `json:"dpto_id"`
	Nombre    string `json:"dpto_nombre"`
	Direccion string `json:"direccion"`
	Ciudad    string `json:"ciudad"`
}

func (c *Client) GetCargos() ([]shared.CargoDTO, error) {
	req := shared.Request{
		Operation: "LIST_CARGOS",
		Data:      nil,
	}
	response, err := c.SendRequest(req)
	if err != nil {
		return nil, err
	}
	if !response.Success {
		return nil, fmt.Errorf(response.Message)
	}
	dataBytes, _ := json.Marshal(response.Data)
	var cargos []shared.CargoDTO
	err = json.Unmarshal(dataBytes, &cargos)
	if err != nil {
		return nil, fmt.Errorf("error procesando lista de cargos: %v", err)
	}
	return cargos, nil
}

func (c *Client) GetDepartamentos() ([]shared.DepartamentoDTO, error) {
	req := shared.Request{
		Operation: "LIST_DEPARTAMENTOS",
		Data:      nil,
	}
	response, err := c.SendRequest(req)
	if err != nil {
		return nil, err
	}
	if !response.Success {
		return nil, fmt.Errorf(response.Message)
	}
	dataBytes, _ := json.Marshal(response.Data)
	var departamentos []shared.DepartamentoDTO
	err = json.Unmarshal(dataBytes, &departamentos)
	if err != nil {
		return nil, fmt.Errorf("error procesando lista de departamentos: %v", err)
	}
	return departamentos, nil
}

func (c *Client) GetDepartamentosConDatos() ([]DepartamentoConDatosDTO, error) {
	req := shared.Request{
		Operation: "LIST_DEPARTAMENTOS_CON_DATOS",
		Data:      nil,
	}
	response, err := c.SendRequest(req)
	if err != nil {
		return nil, err
	}
	if !response.Success {
		return nil, fmt.Errorf(response.Message)
	}
	dataBytes, _ := json.Marshal(response.Data)
	var departamentos []DepartamentoConDatosDTO
	err = json.Unmarshal(dataBytes, &departamentos)
	if err != nil {
		return nil, fmt.Errorf("error procesando lista de departamentos con datos: %v", err)
	}
	return departamentos, nil
}

func (c *Client) GetGerentes() ([]shared.GerenteDTO, error) {
	req := shared.Request{
		Operation: "LIST_GERENTES",
		Data:      nil,
	}
	response, err := c.SendRequest(req)
	if err != nil {
		return nil, err
	}
	if !response.Success {
		return nil, fmt.Errorf(response.Message)
	}
	dataBytes, _ := json.Marshal(response.Data)
	var gerentes []shared.GerenteDTO
	err = json.Unmarshal(dataBytes, &gerentes)
	if err != nil {
		return nil, fmt.Errorf("error procesando lista de gerentes: %v", err)
	}
	return gerentes, nil
}

func (c *Client) SelectCargoFromList() (int, error) {
	cargosConDatos, err := c.GetCargos()
	if err != nil {
		return 0, fmt.Errorf("error obteniendo cargos: %v", err)
	}
	fmt.Println("\nCARGOS DISPONIBLES:")
	for _, cargo := range cargosConDatos {
		fmt.Printf("  %d. %s \n", cargo.ID, cargo.Nombre)
	}
	for {
		cargoID, err := c.ReadIntInput("\nSeleccione el ID del cargo: ")
		if err != nil {
			fmt.Printf("ID inválido\n")
			continue
		}
		for _, cargo := range cargosConDatos {
			if cargo.ID == cargoID {
				return cargoID, nil
			}
		}
		fmt.Printf("ID de cargo no válido. Seleccione uno de la lista.\n")
	}
}

func (c *Client) SelectDepartamentoFromList() (int, error) {
	dptosConDatos, err := c.GetDepartamentosConDatos()
	if err != nil {
		return 0, fmt.Errorf("error obteniendo departamentos: %v", err)
	}
	fmt.Println("\nDEPARTAMENTOS DISPONIBLES:")
	for _, dpto := range dptosConDatos {
		fmt.Printf("  %d. %s - %s - %s\n", dpto.ID, dpto.Nombre, dpto.Direccion, dpto.Ciudad)
	}
	for {
		dptoID, err := c.ReadIntInput("\nSeleccione el ID del departamento: ")
		if err != nil {
			fmt.Printf("ID inválido\n")
			continue
		}
		for _, dpto := range dptosConDatos {
			if dpto.ID == dptoID {
				return dptoID, nil
			}
		}
		fmt.Printf("ID de departamento no válido. Seleccione uno de la lista.\n")
	}
}

func (c *Client) SelectGerenteFromList() (*int, error) {
	gerentes, err := c.GetGerentes()
	if err != nil {
		return nil, fmt.Errorf("error obteniendo gerentes: %v", err)
	}
	fmt.Println("\nGERENTES DISPONIBLES:")
	fmt.Println("  0. Sin gerente (opcional)")
	for _, gerente := range gerentes {
		fmt.Printf("  %d. %s\n", gerente.ID, gerente.Nombre)
	}
	for {
		input := c.ReadInput("\nSeleccione el ID del gerente (0 para ninguno): ")
		if input == "0" || input == "" {
			return nil, nil
		}
		gerenteID, err := strconv.Atoi(input)
		if err != nil {
			fmt.Printf("ID inválido\n")
			continue
		}
		for _, gerente := range gerentes {
			if gerente.ID == gerenteID {
				return &gerenteID, nil
			}
		}
		fmt.Printf("ID de gerente no válido. Seleccione uno de la lista o 0.\n")
	}
}
