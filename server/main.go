package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"hr-system/shared"
	"log"
	"net"
	"os"

	_ "github.com/lib/pq"
)

type Server struct {
	db   *sql.DB
	crud *EmpleadoCrud
	port string
}

func NewServer(port string) *Server {
	return &Server{
		port: port,
	}
}

func (s *Server) connectDB() error {
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return fmt.Errorf("error conectando a la base de datos: %v", err)
	}
	err = db.Ping()
	if err != nil {
		return fmt.Errorf("error haciendo ping a la base de datos: %v", err)
	}
	s.db = db
	s.crud = NewEmpleadoCrud(db)
	log.Println("✓ Conexión a PostgreSQL establecida exitosamente")
	return nil
}

func (s *Server) Start() error {
	if err := s.connectDB(); err != nil {
		return err
	}
	defer s.db.Close()
	listener, err := net.Listen("tcp", ":"+s.port)
	if err != nil {
		return fmt.Errorf("error iniciando servidor: %v", err)
	}
	defer listener.Close()
	log.Printf("✓ Servidor iniciado en puerto %s", s.port)
	log.Println("✓ Esperando conexiones de clientes...")
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Error aceptando conexión: %v", err)
			continue
		}
		go s.handleClient(conn)
	}
}

func (s *Server) handleClient(conn net.Conn) {
	defer conn.Close()
	clientAddr := conn.RemoteAddr().String()
	log.Printf("✓ Cliente conectado desde: %s", clientAddr)
	decoder := json.NewDecoder(conn)
	encoder := json.NewEncoder(conn)
	for {
		var req shared.Request
		if err := decoder.Decode(&req); err != nil {
			log.Printf("Error decodificando request: %v", err)
			break
		}
		log.Printf("Operación recibida: %s", req.Operation)
		response := s.processRequest(req)
		if err := encoder.Encode(response); err != nil {
			log.Printf("Error enviando response: %v", err)
			break
		}
	}
	log.Printf("✓ Cliente %s desconectado", clientAddr)
}

func (s *Server) processRequest(req shared.Request) shared.Response {
	switch req.Operation {
	case "INSERT":
		return s.handleInsert(req.Data)
	case "UPDATE":
		return s.handleUpdate(req.Data)
	case "SELECT":
		return s.handleSelect(req.Data)
	case "DELETE":
		return s.handleDelete(req.Data)
	case "LIST_CARGOS":
		return s.handleListCargos()
	case "LIST_DEPARTAMENTOS_CON_DATOS":
		return s.handleListDepartamentosConDatos()
	case "LIST_GERENTES":
		return s.handleListGerentes()

	default:
		return shared.Response{
			Success: false,
			Message: "Operación no válida. Operaciones disponibles: INSERT, UPDATE, SELECT, DELETE, LIST_CARGOS, LIST_CARGOS_CON_DATOS, LIST_DEPARTAMENTOS, LIST_DEPARTAMENTOS_CON_DATOS, LIST_GERENTES",
		}
	}
}

func (s *Server) handleInsert(data interface{}) shared.Response {
	var dto shared.CreateEmpleadoDTO
	jsonData, err := json.Marshal(data)
	if err != nil {
		return shared.Response{
			Success: false,
			Message: fmt.Sprintf("Error procesando datos: %v", err),
		}
	}
	if err := json.Unmarshal(jsonData, &dto); err != nil {
		return shared.Response{
			Success: false,
			Message: fmt.Sprintf("Error en formato de datos: %v", err),
		}
	}
	result, err := s.crud.Insert(dto)
	if err != nil {
		return shared.Response{
			Success: false,
			Message: err.Error(),
		}
	}
	return shared.Response{
		Success: true,
		Message: "Empleado creado exitosamente",
		Data:    result,
	}
}

func (s *Server) handleUpdate(data interface{}) shared.Response {
	var dto shared.UpdateEmpleadoDTO
	jsonData, err := json.Marshal(data)
	if err != nil {
		return shared.Response{
			Success: false,
			Message: fmt.Sprintf("Error procesando datos: %v", err),
		}
	}
	if err := json.Unmarshal(jsonData, &dto); err != nil {
		return shared.Response{
			Success: false,
			Message: fmt.Sprintf("Error en formato de datos: %v", err),
		}
	}
	result, err := s.crud.Update(dto)
	if err != nil {
		return shared.Response{
			Success: false,
			Message: err.Error(),
		}
	}
	return shared.Response{
		Success: true,
		Message: "Empleado actualizado exitosamente",
		Data:    result,
	}
}

func (s *Server) handleSelect(data interface{}) shared.Response {
	dataMap, ok := data.(map[string]interface{})
	if !ok {
		return shared.Response{
			Success: false,
			Message: "Formato de datos inválido para SELECT",
		}
	}
	idFloat, ok := dataMap["empl_id"].(float64)
	if !ok {
		return shared.Response{
			Success: false,
			Message: "ID del empleado es requerido para SELECT",
		}
	}
	id := int(idFloat)
	result, err := s.crud.Select(id)
	if err != nil {
		return shared.Response{
			Success: false,
			Message: err.Error(),
		}
	}
	return shared.Response{
		Success: true,
		Message: "Empleado encontrado",
		Data:    result,
	}
}

func (s *Server) handleDelete(data interface{}) shared.Response {
	dataMap, ok := data.(map[string]interface{})
	if !ok {
		return shared.Response{
			Success: false,
			Message: "Formato de datos inválido para DELETE",
		}
	}
	idFloat, ok := dataMap["empl_id"].(float64)
	if !ok {
		return shared.Response{
			Success: false,
			Message: "ID del empleado es requerido para DELETE",
		}
	}
	id := int(idFloat)
	err := s.crud.Delete(id)
	if err != nil {
		return shared.Response{
			Success: false,
			Message: err.Error(),
		}
	}
	return shared.Response{
		Success: true,
		Message: "Empleado eliminado exitosamente y guardado en histórico",
	}
}

func (s *Server) handleListCargos() shared.Response {
	cargos, err := s.crud.ListCargos()
	if err != nil {
		return shared.Response{
			Success: false,
			Message: err.Error(),
		}
	}
	return shared.Response{
		Success: true,
		Message: "Lista de cargos obtenida",
		Data:    cargos,
	}
}

func (s *Server) handleListDepartamentosConDatos() shared.Response {
	departamentos, err := s.crud.ListDepartamentosConDatos()
	if err != nil {
		return shared.Response{
			Success: false,
			Message: err.Error(),
		}
	}
	return shared.Response{
		Success: true,
		Message: "Departamentos con datos obtenidos exitosamente",
		Data:    departamentos,
	}
}

func (s *Server) handleListGerentes() shared.Response {
	gerentes, err := s.crud.ListGerentes()
	if err != nil {
		return shared.Response{
			Success: false,
			Message: err.Error(),
		}
	}
	return shared.Response{
		Success: true,
		Message: "Lista de gerentes obtenida",
		Data:    gerentes,
	}
}

func main() {
	port := os.Getenv("SERVER_PORT")
	server := NewServer(port)
	log.Println("=== SERVIDOR DE RECURSOS HUMANOS ===")
	log.Println("Iniciando servidor...")
	if err := server.Start(); err != nil {
		log.Fatalf("Error iniciando servidor: %v", err)
	}
}
