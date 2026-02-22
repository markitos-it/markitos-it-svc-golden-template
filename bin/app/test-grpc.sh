#!/bin/bash

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
MAGENTA='\033[0;35m'
WHITE='\033[1;37m'
NC='\033[0m' # No Color

# Variables
HOST="localhost"
PORT="3000"
DELAY=2
SERVICE="goldens.GoldenService"
DB_HOST="localhost"
DB_PORT="55432"
DB_USER="markitos-it-svc-goldens"
DB_PASS="markitos-it-svc-goldens"
DB_NAME="markitos-it-svc-goldens"
BACKUP_FILE="/tmp/golden_db_backup.sql"

# Functions
print_header() {
    echo -e "\n${WHITE}🚀 gRPC Service Test Suite${NC}\n"
}

print_test() {
    echo -e "\n${BLUE}── $1${NC}"
    echo -e "${CYAN}──────────────────────────────────────────────────────────${NC}"
}

print_success() {
    echo -e "${GREEN}✓ Success${NC}: $1"
}

print_error() {
    echo -e "${RED}✗ Error${NC}: $1"
}

print_info() {
    echo -e "${YELLOW}ℹ Info${NC}: $1"
}

print_separator() {
    echo -e "${CYAN}────────────────────────────────────────────────────────${NC}"
}

wait_for_next_test() {
    echo -e "\n${MAGENTA}⏳ Esperando $DELAY segundos antes del siguiente test...${NC}"
    sleep $DELAY
}

check_server() {
    print_test "TEST 0: Verificando conexión al servidor"
    
    if grpcurl -plaintext ${HOST}:${PORT} list &>/dev/null; then
        print_success "Servidor gRPC disponible en ${HOST}:${PORT}"
        return 0
    else
        print_error "No se puede conectar al servidor gRPC en ${HOST}:${PORT}"
        echo -e "${YELLOW}Asegúrate de que el servicio está corriendo:${NC}"
        echo -e "${CYAN}  docker-compose up${NC}"
        return 1
    fi
}

test_get_all_goldens() {
    print_test "TEST 1: GetAllGoldens - Obtener todos los registros"
    
    echo -e "${CYAN}Comando:${NC}"
    echo -e "  ${WHITE}grpcurl -plaintext ${HOST}:${PORT} ${SERVICE}/GetAllGoldens${NC}\n"
    
    echo -e "${CYAN}Respuesta:${NC}"
    
    if grpcurl -plaintext ${HOST}:${PORT} ${SERVICE}/GetAllGoldens; then
        print_success "GetAllGoldens completado"
        return 0
    else
        print_error "GetAllGoldens falló"
        return 1
    fi
}

test_get_golden_by_id() {
    local GOLDEN_ID=$1
    
    print_test "TEST 2: GetGoldenById - Obtener registro por ID"
    
    echo -e "${CYAN}Comando:${NC}"
    echo -e "  ${WHITE}grpcurl -plaintext -d '{\"id\":\"${GOLDEN_ID}\"}' ${HOST}:${PORT} ${SERVICE}/GetGoldenById${NC}\n"
    
    echo -e "${CYAN}Respuesta:${NC}"
    
    if grpcurl -plaintext -d "{\"id\":\"${GOLDEN_ID}\"}" ${HOST}:${PORT} ${SERVICE}/GetGoldenById; then
        print_success "GetGoldenById completado"
        return 0
    else
        print_error "GetGoldenById falló - Verifica que el ID existe"
        return 1
    fi
}

test_list_services() {
    print_test "TEST 3: Listar servicios disponibles"
    
    echo -e "${CYAN}Comando:${NC}"
    echo -e "  ${WHITE}grpcurl -plaintext ${HOST}:${PORT} list${NC}\n"
    
    echo -e "${CYAN}Respuesta:${NC}"
    
    if grpcurl -plaintext ${HOST}:${PORT} list; then
        print_success "Listado de servicios completado"
        return 0
    else
        print_error "No se pudo listar los servicios"
        return 1
    fi
}

backup_database() {
    print_info "Realizando backup de la base de datos..."
    export PGPASSWORD="${DB_PASS}"
    
    if pg_dump -h ${DB_HOST} -p ${DB_PORT} -U ${DB_USER} -d ${DB_NAME} > ${BACKUP_FILE} 2>/dev/null; then
        print_success "Backup realizado en ${BACKUP_FILE}"
        return 0
    else
        print_error "No se pudo hacer backup. ¿Está instalado pg_dump?"
        return 1
    fi
}

restore_database() {
    print_info "Restaurando base de datos al estado anterior..."
    export PGPASSWORD="${DB_PASS}"
    
    if psql -h ${DB_HOST} -p ${DB_PORT} -U ${DB_USER} -d ${DB_NAME} < ${BACKUP_FILE} 2>/dev/null; then
        print_success "Base de datos restaurada correctamente"
        rm -f ${BACKUP_FILE}
        return 0
    else
        print_error "No se pudo restaurar la base de datos. ¿Está psql disponible?"
        return 1
    fi
}

# Main
main() {
    print_header
    
    print_info "Conectando a ${CYAN}${HOST}:${PORT}${NC}"
    print_separator
    
    # Check if server is running
    if ! check_server; then
        exit 1
    fi
    
    # Backup database before tests
    echo ""
    if ! backup_database; then
        print_error "Continuando sin backup..."
    fi
    
    wait_for_next_test
    
    # Test GetAllGoldens
    test_get_all_goldens
    
    wait_for_next_test
    
    # Extract first ID from response (if possible)
    FIRST_ID=$(grpcurl -plaintext ${HOST}:${PORT} ${SERVICE}/GetAllGoldens 2>/dev/null | grep '"id"' | head -1 | grep -o '"[^"]*"' | sed 's/"//g' | tail -1)
    
    if [ -z "$FIRST_ID" ]; then
        print_info "No se encontraron registros. Usando un ID de ejemplo."
        FIRST_ID="example-id-123"
    else
        print_info "Usando el primer ID encontrado: ${CYAN}${FIRST_ID}${NC}"
    fi
    
    # Test GetGoldenById
    test_get_golden_by_id "$FIRST_ID"
    
    wait_for_next_test
    
    # Test List Services
    test_list_services
    
    print_separator
    
    # Restore database after tests
    echo ""
    if [ -f ${BACKUP_FILE} ]; then
        restore_database
    fi
    
    echo -e "\n${GREEN}✓ Suite de pruebas completada${NC}\n"
}

# Check if grpcurl is installed
if ! command -v grpcurl &> /dev/null; then
    echo -e "${RED}Error: grpcurl no está instalado${NC}"
    echo -e "${YELLOW}Instálalo con:${NC}"
    echo -e "  ${CYAN}go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest${NC}"
    exit 1
fi

# Run main function
main
