/*
 * File: 02_data_types.c
 * Description: Demonstrates basic data types and variables in C
 * Part of: C/C++ Learning Guide - Part 2
 */

#include <stdio.h>

int main() {
    // Integer types
    int age = 25;
    short year = 2024;
    long population = 8000000000L;
    
    // Floating point types
    float height = 5.9f;
    double weight = 70.5;
    
    // Character type
    char grade = 'A';
    
    // Display all variables
    printf("=== Basic Data Types Demo ===\n");
    printf("Age (int): %d\n", age);
    printf("Year (short): %hd\n", year);
    printf("Population (long): %ld\n", population);
    printf("Height (float): %.1f feet\n", height);
    printf("Weight (double): %.1f kg\n", weight);
    printf("Grade (char): %c\n", grade);
    
    // Size of data types
    printf("\n=== Size of Data Types ===\n");
    printf("Size of int: %zu bytes\n", sizeof(int));
    printf("Size of float: %zu bytes\n", sizeof(float));
    printf("Size of double: %zu bytes\n", sizeof(double));
    printf("Size of char: %zu bytes\n", sizeof(char));
    
    // User input example
    int user_number;
    printf("\nEnter a number: ");
    scanf("%d", &user_number);
    printf("You entered: %d\n", user_number);
    printf("Double of your number: %d\n", user_number * 2);
    
    return 0;
}

/*
 * Key Points:
 * - Different data types have different sizes and ranges
 * - Use appropriate format specifiers with printf/scanf
 * - Always use & operator with scanf for basic data types
 */