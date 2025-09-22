/*
 * File: 06_cpp_basics.cpp
 * Description: Introduction to C++ - Classes, Objects, and OOP basics
 * Part of: C/C++ Learning Guide - Part 9 & 10
 */

#include <iostream>
#include <string>
#include <vector>
using namespace std;

// Simple class definition
class Student {
private:
    int id;
    string name;
    double gpa;
    
public:
    // Constructor
    Student(int studentId, string studentName, double studentGpa) {
        id = studentId;
        name = studentName;
        gpa = studentGpa;
        cout << "Student " << name << " created!" << endl;
    }
    
    // Default constructor
    Student() {
        id = 0;
        name = "Unknown";
        gpa = 0.0;
        cout << "Default student created!" << endl;
    }
    
    // Destructor
    ~Student() {
        cout << "Student " << name << " destroyed!" << endl;
    }
    
    // Getter methods
    int getId() const { return id; }
    string getName() const { return name; }
    double getGpa() const { return gpa; }
    
    // Setter methods
    void setId(int newId) {
        if (newId > 0) {
            id = newId;
        }
    }
    
    void setName(string newName) {
        if (!newName.empty()) {
            name = newName;
        }
    }
    
    void setGpa(double newGpa) {
        if (newGpa >= 0.0 && newGpa <= 4.0) {
            gpa = newGpa;
        }
    }
    
    // Method to display student information
    void displayInfo() const {
        cout << "ID: " << id << ", Name: " << name << ", GPA: " << gpa << endl;
    }
    
    // Method to check if student is on honor roll
    bool isHonorRoll() const {
        return gpa >= 3.5;
    }
};

// Function overloading demonstration
int add(int a, int b) {
    return a + b;
}

double add(double a, double b) {
    return a + b;
}

int add(int a, int b, int c) {
    return a + b + c;
}

// Template function
template <typename T>
T maximum(T a, T b) {
    return (a > b) ? a : b;
}

int main() {
    cout << "=== C++ Basics Demo ===" << endl << endl;
    
    // 1. Basic C++ input/output
    cout << "1. C++ Input/Output" << endl;
    string userName;
    int userAge;
    
    cout << "Enter your name: ";
    getline(cin, userName);  // Better than cin >> for strings with spaces
    
    cout << "Enter your age: ";
    cin >> userAge;
    
    cout << "Hello, " << userName << "! You are " << userAge << " years old." << endl << endl;
    
    // 2. Creating and using objects
    cout << "2. Creating Student Objects" << endl;
    Student student1(101, "Alice Johnson", 3.8);
    Student student2(102, "Bob Smith", 3.2);
    Student student3;  // Default constructor
    
    cout << endl;
    
    // 3. Using object methods
    cout << "3. Student Information" << endl;
    student1.displayInfo();
    student2.displayInfo();
    student3.displayInfo();
    cout << endl;
    
    // 4. Modifying object data
    cout << "4. Modifying Student Data" << endl;
    student3.setId(103);
    student3.setName("Charlie Brown");
    student3.setGpa(3.9);
    
    cout << "Updated student3: ";
    student3.displayInfo();
    cout << endl;
    
    // 5. Using object methods for logic
    cout << "5. Honor Roll Check" << endl;
    vector<Student> students = {student1, student2, student3};
    
    for (const auto& student : students) {
        cout << student.getName();
        if (student.isHonorRoll()) {
            cout << " is on the honor roll!" << endl;
        } else {
            cout << " is not on the honor roll." << endl;
        }
    }
    cout << endl;
    
    // 6. Function overloading
    cout << "6. Function Overloading" << endl;
    cout << "add(5, 3) = " << add(5, 3) << endl;
    cout << "add(5.5, 3.2) = " << add(5.5, 3.2) << endl;
    cout << "add(1, 2, 3) = " << add(1, 2, 3) << endl;
    cout << endl;
    
    // 7. Template functions
    cout << "7. Template Functions" << endl;
    cout << "maximum(10, 20) = " << maximum(10, 20) << endl;
    cout << "maximum(3.14, 2.71) = " << maximum(3.14, 2.71) << endl;
    cout << "maximum('a', 'z') = " << maximum('a', 'z') << endl;
    cout << endl;
    
    // 8. References vs Pointers
    cout << "8. References vs Pointers" << endl;
    int number = 42;
    int& ref = number;  // Reference
    int* ptr = &number; // Pointer
    
    cout << "Original number: " << number << endl;
    cout << "Reference: " << ref << endl;
    cout << "Pointer value: " << *ptr << endl;
    
    ref = 100;  // Modifies original
    cout << "After modifying through reference: " << number << endl;
    
    *ptr = 200; // Modifies original
    cout << "After modifying through pointer: " << number << endl;
    cout << endl;
    
    // 9. STL Vector demonstration
    cout << "9. STL Vector" << endl;
    vector<int> numbers = {1, 2, 3, 4, 5};
    
    cout << "Original vector: ";
    for (int num : numbers) {
        cout << num << " ";
    }
    cout << endl;
    
    numbers.push_back(6);
    numbers.push_back(7);
    
    cout << "After adding elements: ";
    for (size_t i = 0; i < numbers.size(); i++) {
        cout << numbers[i] << " ";
    }
    cout << endl;
    
    cout << "Vector size: " << numbers.size() << endl;
    cout << "First element: " << numbers.front() << endl;
    cout << "Last element: " << numbers.back() << endl;
    cout << endl;
    
    // 10. String class
    cout << "10. C++ String Class" << endl;
    string greeting = "Hello";
    string name = "World";
    string message = greeting + ", " + name + "!";
    
    cout << "Concatenated string: " << message << endl;
    cout << "String length: " << message.length() << endl;
    cout << "Substring (0, 5): " << message.substr(0, 5) << endl;
    
    if (message.find("World") != string::npos) {
        cout << "Found 'World' in the message!" << endl;
    }
    
    cout << endl << "=== End of Demo ===" << endl;
    
    return 0;
}

/*
 * Key C++ Features Demonstrated:
 * 
 * 1. Classes and Objects:
 *    - Encapsulation (private/public)
 *    - Constructors and destructors
 *    - Member functions (methods)
 *    - Getter and setter methods
 * 
 * 2. C++ Improvements over C:
 *    - iostream for input/output
 *    - string class instead of char arrays
 *    - References in addition to pointers
 *    - Function overloading
 *    - Templates for generic programming
 * 
 * 3. STL (Standard Template Library):
 *    - vector for dynamic arrays
 *    - Range-based for loops
 *    - Built-in string operations
 * 
 * 4. Object-Oriented Programming:
 *    - Data hiding and encapsulation
 *    - Methods that operate on object data
 *    - Constructor/destructor lifecycle
 */