/*
 * File: 05_arrays_strings.c
 * Description: Demonstrates arrays and string operations
 * Part of: C/C++ Learning Guide - Part 5
 */

#include <stdio.h>
#include <string.h>
#include <ctype.h>

int main() {
    printf("=== Arrays and Strings Demo ===\n\n");
    
    // 1. Basic array operations
    printf("1. Basic Array Operations\n");
    int numbers[5] = {10, 20, 30, 40, 50};
    
    printf("Array elements: ");
    for (int i = 0; i < 5; i++) {
        printf("%d ", numbers[i]);
    }
    printf("\n");
    
    // Calculate sum and average
    int sum = 0;
    for (int i = 0; i < 5; i++) {
        sum += numbers[i];
    }
    printf("Sum: %d\n", sum);
    printf("Average: %.2f\n\n", sum / 5.0);
    
    // 2. Finding maximum and minimum
    printf("2. Finding Max and Min\n");
    int scores[] = {85, 92, 78, 96, 87, 73, 89};
    int size = sizeof(scores) / sizeof(scores[0]);
    
    int max = scores[0], min = scores[0];
    int maxIndex = 0, minIndex = 0;
    
    for (int i = 1; i < size; i++) {
        if (scores[i] > max) {
            max = scores[i];
            maxIndex = i;
        }
        if (scores[i] < min) {
            min = scores[i];
            minIndex = i;
        }
    }
    
    printf("Scores: ");
    for (int i = 0; i < size; i++) {
        printf("%d ", scores[i]);
    }
    printf("\n");
    printf("Maximum: %d (at index %d)\n", max, maxIndex);
    printf("Minimum: %d (at index %d)\n\n", min, minIndex);
    
    // 3. 2D Array (Matrix)
    printf("3. 2D Array (Matrix)\n");
    int matrix[3][3] = {
        {1, 2, 3},
        {4, 5, 6},
        {7, 8, 9}
    };
    
    printf("Matrix:\n");
    for (int i = 0; i < 3; i++) {
        for (int j = 0; j < 3; j++) {
            printf("%d ", matrix[i][j]);
        }
        printf("\n");
    }
    
    // Calculate diagonal sum
    int diagonalSum = 0;
    for (int i = 0; i < 3; i++) {
        diagonalSum += matrix[i][i];
    }
    printf("Main diagonal sum: %d\n\n", diagonalSum);
    
    // 4. String operations
    printf("4. String Operations\n");
    char greeting[] = "Hello, World!";
    char name[50];
    char fullMessage[100];
    
    printf("Original string: %s\n", greeting);
    printf("String length: %zu\n", strlen(greeting));
    
    // String input
    printf("Enter your name: ");
    fgets(name, sizeof(name), stdin);
    
    // Remove newline from fgets
    name[strcspn(name, "\n")] = 0;
    
    // String concatenation
    strcpy(fullMessage, "Hello, ");
    strcat(fullMessage, name);
    strcat(fullMessage, "! Welcome to C programming.");
    
    printf("Personalized message: %s\n\n", fullMessage);
    
    // 5. String manipulation
    printf("5. String Manipulation\n");
    char text[] = "Programming is Fun!";
    printf("Original: %s\n", text);
    
    // Convert to uppercase
    char upperText[50];
    strcpy(upperText, text);
    for (int i = 0; upperText[i]; i++) {
        upperText[i] = toupper(upperText[i]);
    }
    printf("Uppercase: %s\n", upperText);
    
    // Convert to lowercase
    char lowerText[50];
    strcpy(lowerText, text);
    for (int i = 0; lowerText[i]; i++) {
        lowerText[i] = tolower(lowerText[i]);
    }
    printf("Lowercase: %s\n", lowerText);
    
    // Count vowels
    int vowelCount = 0;
    for (int i = 0; text[i]; i++) {
        char c = tolower(text[i]);
        if (c == 'a' || c == 'e' || c == 'i' || c == 'o' || c == 'u') {
            vowelCount++;
        }
    }
    printf("Number of vowels: %d\n\n", vowelCount);
    
    // 6. String comparison
    printf("6. String Comparison\n");
    char str1[] = "apple";
    char str2[] = "banana";
    char str3[] = "apple";
    
    printf("Comparing '%s' and '%s': ", str1, str2);
    if (strcmp(str1, str2) == 0) {
        printf("Equal\n");
    } else if (strcmp(str1, str2) < 0) {
        printf("'%s' comes before '%s'\n", str1, str2);
    } else {
        printf("'%s' comes after '%s'\n", str1, str2);
    }
    
    printf("Comparing '%s' and '%s': ", str1, str3);
    if (strcmp(str1, str3) == 0) {
        printf("Equal\n");
    } else {
        printf("Not equal\n");
    }
    
    // 7. Character array vs string literal
    printf("\n7. Character Arrays\n");
    char modifiable[] = "Hello";  // Can be modified
    char *readonly = "World";     // Should not be modified
    
    printf("Before modification: %s\n", modifiable);
    modifiable[0] = 'h';  // Change 'H' to 'h'
    printf("After modification: %s\n", modifiable);
    
    printf("Read-only string: %s\n", readonly);
    // readonly[0] = 'w';  // This would cause undefined behavior!
    
    return 0;
}

/*
 * Key Points:
 * - Arrays are zero-indexed (first element is at index 0)
 * - Array size must be known at compile time for static arrays
 * - Strings in C are arrays of characters ending with '\0'
 * - Use strlen() for string length, strcpy() for copying
 * - Use strcat() for concatenation, strcmp() for comparison
 * - fgets() is safer than gets() for string input
 * - String literals are read-only, character arrays are modifiable
 * - Always be careful with array bounds to avoid buffer overflow
 */