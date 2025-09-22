# Makefile for C/C++ Learning Examples
# This file helps automate the compilation process

# Compiler settings
CC = gcc
CXX = g++
CFLAGS = -Wall -Wextra -std=c99
CXXFLAGS = -Wall -Wextra -std=c++17

# Directories
EXAMPLE_DIR = examples
PROJECT_DIR = projects
BIN_DIR = bin

# Create bin directory if it doesn't exist
$(shell mkdir -p $(BIN_DIR))

# Default target
all: c_examples cpp_examples projects

# C Examples
c_examples: hello_world data_types control_structures functions arrays_strings

hello_world: $(EXAMPLE_DIR)/01_hello_world.c
	$(CC) $(CFLAGS) $< -o $(BIN_DIR)/hello_world

data_types: $(EXAMPLE_DIR)/02_data_types.c
	$(CC) $(CFLAGS) $< -o $(BIN_DIR)/data_types

control_structures: $(EXAMPLE_DIR)/03_control_structures.c
	$(CC) $(CFLAGS) $< -o $(BIN_DIR)/control_structures

functions: $(EXAMPLE_DIR)/04_functions.c
	$(CC) $(CFLAGS) $< -lm -o $(BIN_DIR)/functions

arrays_strings: $(EXAMPLE_DIR)/05_arrays_strings.c
	$(CC) $(CFLAGS) $< -o $(BIN_DIR)/arrays_strings

# C++ Examples
cpp_examples: cpp_basics

cpp_basics: $(EXAMPLE_DIR)/06_cpp_basics.cpp
	$(CXX) $(CXXFLAGS) $< -o $(BIN_DIR)/cpp_basics

# Projects
projects: student_management

student_management: $(PROJECT_DIR)/student_management.cpp
	$(CXX) $(CXXFLAGS) $< -o $(BIN_DIR)/student_management

# Run targets
run-hello: hello_world
	./$(BIN_DIR)/hello_world

run-data-types: data_types
	./$(BIN_DIR)/data_types

run-control: control_structures
	./$(BIN_DIR)/control_structures

run-functions: functions
	./$(BIN_DIR)/functions

run-arrays: arrays_strings
	./$(BIN_DIR)/arrays_strings

run-cpp: cpp_basics
	./$(BIN_DIR)/cpp_basics

run-student: student_management
	./$(BIN_DIR)/student_management

# Clean compiled files
clean:
	rm -rf $(BIN_DIR)/*

# Help target
help:
	@echo "Available targets:"
	@echo "  all              - Compile all examples and projects"
	@echo "  c_examples       - Compile all C examples"
	@echo "  cpp_examples     - Compile all C++ examples"
	@echo "  projects         - Compile all projects"
	@echo ""
	@echo "Individual targets:"
	@echo "  hello_world      - Compile hello world example"
	@echo "  data_types       - Compile data types example"
	@echo "  control_structures - Compile control structures example"
	@echo "  functions        - Compile functions example"
	@echo "  arrays_strings   - Compile arrays and strings example"
	@echo "  cpp_basics       - Compile C++ basics example"
	@echo "  student_management - Compile student management project"
	@echo ""
	@echo "Run targets:"
	@echo "  run-hello        - Run hello world example"
	@echo "  run-data-types   - Run data types example"
	@echo "  run-control      - Run control structures example"
	@echo "  run-functions    - Run functions example"
	@echo "  run-arrays       - Run arrays and strings example"
	@echo "  run-cpp          - Run C++ basics example"
	@echo "  run-student      - Run student management project"
	@echo ""
	@echo "  clean            - Remove compiled files"
	@echo "  help             - Show this help message"

# Declare phony targets
.PHONY: all c_examples cpp_examples projects clean help
.PHONY: hello_world data_types control_structures functions arrays_strings
.PHONY: cpp_basics student_management
.PHONY: run-hello run-data-types run-control run-functions run-arrays run-cpp run-student