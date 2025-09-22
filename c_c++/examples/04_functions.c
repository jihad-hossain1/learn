/*
 * File: 04_functions.c
 * Description: Demonstrates functions, parameters, and recursion
 * Part of: C/C++ Learning Guide - Part 4
 */

#include <stdio.h>
#include <math.h>

// Function declarations (prototypes)
int add(int a, int b);
float calculateArea(float radius);
void greetUser(char name[]);
void swapValues(int *a, int *b);
int factorial(int n);
int fibonacci(int n);
bool isPrime(int n);
void printMenu();

int main() {
    printf("=== Functions Demo ===\n\n");
    
    // 1. Simple function call
    printf("1. Basic Addition\n");
    int result = add(15, 25);
    printf("15 + 25 = %d\n\n", result);
    
    // 2. Function with float parameters
    printf("2. Circle Area Calculator\n");
    float radius = 5.0;
    float area = calculateArea(radius);
    printf("Area of circle with radius %.1f = %.2f\n\n", radius, area);
    
    // 3. Function with string parameter
    printf("3. Greeting Function\n");
    char name[] = "Alice";
    greetUser(name);
    printf("\n");
    
    // 4. Pass by reference (using pointers)
    printf("4. Swap Function (Pass by Reference)\n");
    int x = 10, y = 20;
    printf("Before swap: x = %d, y = %d\n", x, y);
    swapValues(&x, &y);
    printf("After swap: x = %d, y = %d\n\n", x, y);
    
    // 5. Recursive functions
    printf("5. Recursive Functions\n");
    int num = 5;
    printf("Factorial of %d = %d\n", num, factorial(num));
    
    printf("Fibonacci sequence (first 10 numbers): ");
    for (int i = 0; i < 10; i++) {
        printf("%d ", fibonacci(i));
    }
    printf("\n\n");
    
    // 6. Prime number checker
    printf("6. Prime Number Checker\n");
    for (int i = 2; i <= 20; i++) {
        if (isPrime(i)) {
            printf("%d ", i);
        }
    }
    printf("\n\n");
    
    // 7. Menu system using functions
    printf("7. Menu System\n");
    printMenu();
    
    return 0;
}

// Function definitions

// Simple addition function
int add(int a, int b) {
    return a + b;
}

// Calculate area of circle
float calculateArea(float radius) {
    const float PI = 3.14159;
    return PI * radius * radius;
}

// Greet user by name
void greetUser(char name[]) {
    printf("Hello, %s! Welcome to C programming.\n", name);
}

// Swap two integers using pointers (pass by reference)
void swapValues(int *a, int *b) {
    int temp = *a;
    *a = *b;
    *b = temp;
}

// Recursive factorial function
int factorial(int n) {
    if (n <= 1) {
        return 1;  // Base case
    }
    return n * factorial(n - 1);  // Recursive case
}

// Recursive Fibonacci function
int fibonacci(int n) {
    if (n <= 1) {
        return n;  // Base cases: fib(0) = 0, fib(1) = 1
    }
    return fibonacci(n - 1) + fibonacci(n - 2);  // Recursive case
}

// Check if a number is prime
bool isPrime(int n) {
    if (n < 2) {
        return false;
    }
    
    for (int i = 2; i <= sqrt(n); i++) {
        if (n % i == 0) {
            return false;
        }
    }
    return true;
}

// Print a simple menu
void printMenu() {
    printf("=== Calculator Menu ===\n");
    printf("1. Addition\n");
    printf("2. Subtraction\n");
    printf("3. Multiplication\n");
    printf("4. Division\n");
    printf("5. Exit\n");
    printf("Choose an option: ");
}

/*
 * Key Points:
 * - Function prototype declares the function before main()
 * - Function definition provides the actual implementation
 * - Pass by value: Function receives a copy of the argument
 * - Pass by reference: Function receives the address of the argument
 * - Recursive functions call themselves with a base case to stop
 * - void functions don't return a value
 * - Functions help organize code and promote reusability
 */