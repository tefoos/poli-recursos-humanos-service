package main

import (
	"fmt"
	"strconv"
)

func (c *Client) SelectGerenteFromListForUpdate() (*string, error) {
	gerentes, err := c.GetGerentes()
	if err != nil {
		return nil, err
	}

	fmt.Println("\nGERENTES DISPONIBLES:")
	fmt.Println("  0. Sin gerente (opcional)")
	for _, gerente := range gerentes {
		fmt.Printf("  %d. %s\n", gerente.ID, gerente.Nombre)
	}

	input := c.ReadInput("\nSeleccione el ID del gerente (0 para ninguno): ")
	if input == "0" || input == "" {
		empty := ""
		return &empty, nil
	}

	gerenteID, err := strconv.Atoi(input)
	if err != nil {
		return nil, fmt.Errorf("ID inválido")
	}

	for _, gerente := range gerentes {
		if gerente.ID == gerenteID {
			return &gerente.Nombre, nil
		}
	}

	return nil, fmt.Errorf("ID de gerente no válido")
}

func (c *Client) GetCargoIDByName(nombre string) int {
	cargos, err := c.GetCargos()
	if err != nil {
		return 0
	}

	for _, cargo := range cargos {
		if cargo.Nombre == nombre {
			return cargo.ID
		}
	}
	return 0
}

func (c *Client) GetDptoIDByName(nombre string) int {
	dptos, err := c.GetDepartamentos()
	if err != nil {
		return 0
	}

	for _, dpto := range dptos {
		if dpto.Nombre == nombre {
			return dpto.ID
		}
	}
	return 0
}

func (c *Client) GetGerenteIDByName(nombre string) int {
	gerentes, err := c.GetGerentes()
	if err != nil {
		return 0
	}

	for _, gerente := range gerentes {
		if gerente.Nombre == nombre {
			return gerente.ID
		}
	}
	return 0
}
