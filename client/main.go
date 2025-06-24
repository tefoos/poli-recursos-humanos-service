package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"

	"hr-system/shared"
)

type Client struct {
	conn net.Conn
}

func NewClient() *Client {
	return &Client{}
}

func (c *Client) Connect(host, port string) error {
	address := fmt.Sprintf("%s:%s", host, port)
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return fmt.Errorf("error conectando al servidor: %v", err)
	}

	c.conn = conn
	log.Printf("Conectado al servidor %s", address)
	return nil
}

func (c *Client) Disconnect() {
	if c.conn != nil {
		c.conn.Close()
		log.Println("Desconectado del servidor")
	}
}

func (c *Client) SendRequest(req shared.Request) (*shared.Response, error) {
	encoder := json.NewEncoder(c.conn)
	decoder := json.NewDecoder(c.conn)

	if err := encoder.Encode(req); err != nil {
		return nil, fmt.Errorf("error enviando request: %v", err)
	}

	var response shared.Response
	if err := decoder.Decode(&response); err != nil {
		return nil, fmt.Errorf("error recibiendo response: %v", err)
	}

	return &response, nil
}

func (c *Client) ShowMenu() {
	fmt.Println("\n=== SISTEMA DE RECURSOS HUMANOS ===")
	fmt.Println("1. Crear empleado (INSERT)")
	fmt.Println("2. Actualizar empleado (UPDATE)")
	fmt.Println("3. Consultar empleado (SELECT)")
	fmt.Println("4. Eliminar empleado (DELETE)")
	fmt.Println("5. Salir")
	fmt.Print("Seleccione una opción: ")
}

func (c *Client) ReadInput(prompt string) string {
	fmt.Print(prompt)
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}

func (c *Client) ReadIntInput(prompt string) (int, error) {
	input := c.ReadInput(prompt)
	return strconv.Atoi(input)
}

func (c *Client) PrintResponse(response *shared.Response) {
	fmt.Println("\n--- RESPUESTA DEL SERVIDOR ---")
	if response.Success {
		fmt.Printf("Éxito: %s\n", response.Message)
		if response.Data != nil {
			dataJSON, _ := json.MarshalIndent(response.Data, "", "  ")
			fmt.Printf("Datos: %s\n", string(dataJSON))
		}
	} else {
		fmt.Printf("Error: %s\n", response.Message)
	}
}

func (c *Client) Run() {
	defer c.Disconnect()

	for {
		c.ShowMenu()
		option := c.ReadInput("")

		switch option {
		case "1":
			c.HandleInsert()
		case "2":
			c.HandleUpdate()
		case "3":
			c.HandleSelect()
		case "4":
			c.HandleDelete()
		case "5":
			fmt.Println("¡Hasta luego!")
			return
		default:
			fmt.Println("Opción no válida. Intente de nuevo.")
		}
	}
}

func main() {
	fmt.Println("=== CLIENTE DE RECURSOS HUMANOS ===")

	client := NewClient()

	host := "localhost"
	port := "8888"

	if err := client.Connect(host, port); err != nil {
		log.Fatalf("Error conectando: %v", err)
	}

	client.Run()
}
