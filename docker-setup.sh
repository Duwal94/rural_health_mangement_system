#!/bin/bash

# Rural Health Management System - Docker Setup Script

set -e

echo "🏥 Rural Health Management System - Docker Setup"
echo "================================================"

# Function to check if Docker is installed
check_docker() {
    if ! command -v docker &> /dev/null; then
        echo "❌ Docker is not installed. Please install Docker first."
        echo "   Visit: https://docs.docker.com/get-docker/"
        exit 1
    fi
    
    if ! command -v docker-compose &> /dev/null; then
        echo "❌ Docker Compose is not installed. Please install Docker Compose first."
        echo "   Visit: https://docs.docker.com/compose/install/"
        exit 1
    fi
    
    echo "✅ Docker and Docker Compose are installed"
}

# Function to build the application
build_app() {
    echo "🔨 Building the application..."
    docker-compose build --no-cache
    echo "✅ Application built successfully"
}

# Function to start the services
start_services() {
    echo "🚀 Starting services..."
    docker-compose up -d
    echo "✅ Services started successfully"
}

# Function to check service health
check_health() {
    echo "🔍 Checking service health..."
    
    # Wait for database to be ready
    echo "   Waiting for database..."
    for i in {1..30}; do
        if docker-compose exec -T postgres pg_isready -U postgres >/dev/null 2>&1; then
            echo "   ✅ Database is ready"
            break
        fi
        echo "   ⏳ Waiting for database... ($i/30)"
        sleep 2
    done
    
    # Wait for API to be ready
    echo "   Waiting for API..."
    for i in {1..30}; do
        if curl -f http://localhost:3000/health >/dev/null 2>&1; then
            echo "   ✅ API is ready"
            break
        fi
        echo "   ⏳ Waiting for API... ($i/30)"
        sleep 2
    done
}

# Function to seed the database
seed_database() {
    echo "🌱 Seeding database with sample data..."
    docker-compose exec api go run cmd/seed/main.go
    echo "✅ Database seeded successfully"
}

# Function to show service status
show_status() {
    echo "📊 Service Status:"
    docker-compose ps
    echo ""
    echo "🌐 Application URLs:"
    echo "   API: http://localhost:3000"
    echo "   Health Check: http://localhost:3000/health"
    echo "   Database: localhost:5432"
}

# Function to show logs
show_logs() {
    echo "📝 Recent logs:"
    docker-compose logs --tail=20
}

# Function to stop services
stop_services() {
    echo "🛑 Stopping services..."
    docker-compose down
    echo "✅ Services stopped"
}

# Function to clean up
cleanup() {
    echo "🧹 Cleaning up..."
    docker-compose down -v --remove-orphans
    docker system prune -f
    echo "✅ Cleanup completed"
}

# Main script logic
case "${1:-start}" in
    "start")
        check_docker
        build_app
        start_services
        check_health
        show_status
        echo ""
        echo "🎉 Rural Health Management System is now running!"
        echo "   Try: curl http://localhost:3000/health"
        ;;
    "seed")
        seed_database
        ;;
    "status")
        show_status
        ;;
    "logs")
        show_logs
        ;;
    "stop")
        stop_services
        ;;
    "restart")
        stop_services
        sleep 2
        start_services
        check_health
        show_status
        ;;
    "clean")
        cleanup
        ;;
    "help"|"-h"|"--help")
        echo "Usage: $0 [command]"
        echo ""
        echo "Commands:"
        echo "  start    - Build and start all services (default)"
        echo "  seed     - Seed database with sample data"
        echo "  status   - Show service status"
        echo "  logs     - Show recent logs"
        echo "  stop     - Stop all services"
        echo "  restart  - Restart all services"
        echo "  clean    - Stop services and clean up"
        echo "  help     - Show this help message"
        ;;
    *)
        echo "❌ Unknown command: $1"
        echo "Run '$0 help' for usage information"
        exit 1
        ;;
esac
