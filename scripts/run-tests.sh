#!/bin/bash
# Revision: 2025-12-03 | Author: Claude | Version: 1.0.0
# Stumpf.Works NAS - Comprehensive Test Runner

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Script directory
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"

# Test configuration
COVERAGE_THRESHOLD=80
COVERAGE_DIR="$PROJECT_ROOT/coverage"
BACKEND_DIR="$PROJECT_ROOT/backend"
FRONTEND_DIR="$PROJECT_ROOT/frontend"

# Create coverage directory
mkdir -p "$COVERAGE_DIR"

# Print header
print_header() {
    echo -e "${BLUE}╔════════════════════════════════════════════════════════════════╗${NC}"
    echo -e "${BLUE}║  Stumpf.Works NAS - Test Suite                                ║${NC}"
    echo -e "${BLUE}╚════════════════════════════════════════════════════════════════╝${NC}"
    echo ""
}

# Print section
print_section() {
    echo -e "\n${BLUE}━━━ $1 ━━━${NC}\n"
}

# Print success
print_success() {
    echo -e "${GREEN}✓ $1${NC}"
}

# Print error
print_error() {
    echo -e "${RED}✗ $1${NC}"
}

# Print warning
print_warning() {
    echo -e "${YELLOW}⚠ $1${NC}"
}

# Parse command line arguments
BACKEND_ONLY=false
FRONTEND_ONLY=false
COVERAGE_ONLY=false
VERBOSE=false
RACE=false

while [[ $# -gt 0 ]]; do
    case $1 in
        --backend)
            BACKEND_ONLY=true
            shift
            ;;
        --frontend)
            FRONTEND_ONLY=true
            shift
            ;;
        --coverage)
            COVERAGE_ONLY=true
            shift
            ;;
        --verbose|-v)
            VERBOSE=true
            shift
            ;;
        --race)
            RACE=true
            shift
            ;;
        --help|-h)
            echo "Usage: $0 [OPTIONS]"
            echo ""
            echo "Options:"
            echo "  --backend       Run only backend tests"
            echo "  --frontend      Run only frontend tests"
            echo "  --coverage      Generate coverage report only"
            echo "  --verbose, -v   Verbose output"
            echo "  --race          Enable race detector (Go tests)"
            echo "  --help, -h      Show this help message"
            echo ""
            exit 0
            ;;
        *)
            echo "Unknown option: $1"
            echo "Use --help for usage information"
            exit 1
            ;;
    esac
done

# Start
print_header

# Run backend tests
run_backend_tests() {
    print_section "Backend Tests (Go)"

    cd "$BACKEND_DIR"

    # Check if tests exist
    TEST_FILES=$(find . -name "*_test.go" | wc -l)
    echo "Found $TEST_FILES test files"

    if [ "$TEST_FILES" -eq 0 ]; then
        print_warning "No test files found in backend"
        return 0
    fi

    # Build test flags
    TEST_FLAGS="-v"
    if [ "$RACE" = true ]; then
        TEST_FLAGS="$TEST_FLAGS -race"
    fi

    # Run tests
    echo "Running tests with coverage..."
    if go test $TEST_FLAGS -coverprofile="$COVERAGE_DIR/backend-coverage.out" ./...; then
        print_success "Backend tests passed"

        # Generate coverage report
        COVERAGE=$(go tool cover -func="$COVERAGE_DIR/backend-coverage.out" | grep total | awk '{print $3}' | sed 's/%//')
        echo "Coverage: $COVERAGE%"

        # Check coverage threshold
        if (( $(echo "$COVERAGE < $COVERAGE_THRESHOLD" | bc -l) )); then
            print_warning "Coverage ($COVERAGE%) is below threshold ($COVERAGE_THRESHOLD%)"
        else
            print_success "Coverage meets threshold ($COVERAGE_THRESHOLD%)"
        fi

        # Generate HTML coverage report
        go tool cover -html="$COVERAGE_DIR/backend-coverage.out" -o "$COVERAGE_DIR/backend-coverage.html"
        print_success "HTML coverage report: $COVERAGE_DIR/backend-coverage.html"

        return 0
    else
        print_error "Backend tests failed"
        return 1
    fi
}

# Run frontend tests
run_frontend_tests() {
    print_section "Frontend Tests (Vitest)"

    cd "$FRONTEND_DIR"

    # Check if node_modules exists
    if [ ! -d "node_modules" ]; then
        print_warning "node_modules not found, installing dependencies..."
        npm install
    fi

    # Run tests
    echo "Running tests with coverage..."
    if npm run test:coverage -- --reporter=verbose --coverage.reporter=json --coverage.reporter=html --coverage.reporter=text; then
        print_success "Frontend tests passed"

        # Move coverage to central location
        if [ -d "coverage" ]; then
            cp -r coverage "$COVERAGE_DIR/frontend"
            print_success "Coverage report: $COVERAGE_DIR/frontend/index.html"
        fi

        return 0
    else
        print_error "Frontend tests failed"
        return 1
    fi
}

# Generate coverage summary
generate_coverage_summary() {
    print_section "Coverage Summary"

    echo "Backend Coverage Report:"
    if [ -f "$COVERAGE_DIR/backend-coverage.out" ]; then
        cd "$BACKEND_DIR"
        go tool cover -func="$COVERAGE_DIR/backend-coverage.out" | tail -10
    else
        print_warning "No backend coverage data found"
    fi

    echo ""
    echo "Frontend Coverage Report:"
    if [ -f "$COVERAGE_DIR/frontend/coverage-summary.json" ]; then
        cat "$COVERAGE_DIR/frontend/coverage-summary.json" | jq '.total'
    else
        print_warning "No frontend coverage data found"
    fi

    echo ""
    echo "Coverage reports available at:"
    echo "  Backend:  file://$COVERAGE_DIR/backend-coverage.html"
    echo "  Frontend: file://$COVERAGE_DIR/frontend/index.html"
}

# Main execution
BACKEND_RESULT=0
FRONTEND_RESULT=0

if [ "$COVERAGE_ONLY" = true ]; then
    generate_coverage_summary
    exit 0
fi

if [ "$BACKEND_ONLY" = false ] && [ "$FRONTEND_ONLY" = false ]; then
    # Run both
    run_backend_tests || BACKEND_RESULT=$?
    run_frontend_tests || FRONTEND_RESULT=$?
elif [ "$BACKEND_ONLY" = true ]; then
    run_backend_tests || BACKEND_RESULT=$?
elif [ "$FRONTEND_ONLY" = true ]; then
    run_frontend_tests || FRONTEND_RESULT=$?
fi

# Generate summary
generate_coverage_summary

# Final result
echo ""
print_section "Test Results"

if [ $BACKEND_RESULT -eq 0 ] && [ $FRONTEND_RESULT -eq 0 ]; then
    print_success "All tests passed!"
    exit 0
else
    if [ $BACKEND_RESULT -ne 0 ]; then
        print_error "Backend tests failed"
    fi
    if [ $FRONTEND_RESULT -ne 0 ]; then
        print_error "Frontend tests failed"
    fi
    exit 1
fi
