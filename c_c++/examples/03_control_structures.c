/*
 * File: 03_control_structures.c
 * Description: Demonstrates control structures (if-else, loops, switch)
 * Part of: C/C++ Learning Guide - Part 3
 */

#include <stdio.h>

int main() {
    printf("=== Control Structures Demo ===\n\n");
    
    // 1. if-else statement
    printf("1. Grade Calculator\n");
    int score = 85;
    printf("Score: %d\n", score);
    
    if (score >= 90) {
        printf("Grade: A (Excellent!)\n");
    } else if (score >= 80) {
        printf("Grade: B (Good!)\n");
    } else if (score >= 70) {
        printf("Grade: C (Average)\n");
    } else if (score >= 60) {
        printf("Grade: D (Below Average)\n");
    } else {
        printf("Grade: F (Fail)\n");
    }
    
    // 2. for loop
    printf("\n2. Multiplication Table (5)\n");
    for (int i = 1; i <= 10; i++) {
        printf("5 x %d = %d\n", i, 5 * i);
    }
    
    // 3. while loop
    printf("\n3. Countdown using while loop\n");
    int countdown = 5;
    while (countdown > 0) {
        printf("%d... ", countdown);
        countdown--;
    }
    printf("Blast off!\n");
    
    // 4. do-while loop
    printf("\n4. Number guessing (simplified)\n");
    int guess, target = 7;
    do {
        printf("Guess a number (1-10): ");
        scanf("%d", &guess);
        if (guess != target) {
            printf("Wrong! Try again.\n");
        }
    } while (guess != target);
    printf("Correct! The number was %d\n", target);
    
    // 5. switch statement
    printf("\n5. Day of the week\n");
    int day = 3;
    printf("Day %d is: ", day);
    
    switch (day) {
        case 1:
            printf("Monday\n");
            break;
        case 2:
            printf("Tuesday\n");
            break;
        case 3:
            printf("Wednesday\n");
            break;
        case 4:
            printf("Thursday\n");
            break;
        case 5:
            printf("Friday\n");
            break;
        case 6:
            printf("Saturday\n");
            break;
        case 7:
            printf("Sunday\n");
            break;
        default:
            printf("Invalid day\n");
    }
    
    // 6. Nested loops - Pattern printing
    printf("\n6. Star Pattern\n");
    for (int i = 1; i <= 5; i++) {
        for (int j = 1; j <= i; j++) {
            printf("* ");
        }
        printf("\n");
    }
    
    // 7. Loop control - break and continue
    printf("\n7. Even numbers from 1 to 10 (using continue)\n");
    for (int i = 1; i <= 10; i++) {
        if (i % 2 != 0) {
            continue;  // Skip odd numbers
        }
        printf("%d ", i);
    }
    printf("\n");
    
    printf("\n8. Find first number divisible by 7 (using break)\n");
    for (int i = 1; i <= 50; i++) {
        if (i % 7 == 0) {
            printf("First number divisible by 7: %d\n", i);
            break;  // Exit loop when found
        }
    }
    
    return 0;
}

/*
 * Key Points:
 * - if-else: Use for conditional execution
 * - for loop: Use when you know the number of iterations
 * - while loop: Use when condition is checked before execution
 * - do-while: Use when you want to execute at least once
 * - switch: Use for multiple discrete values
 * - break: Exits the loop immediately
 * - continue: Skips current iteration and continues with next
 */