#include<stdio.h>


int main(){
    
    // do while
    int number;
    do{
        printf("Enter a positive number: ");
        scanf("%d",&number);
    } while(number <= 0);
    
    printf("\n");
    return 0;
}