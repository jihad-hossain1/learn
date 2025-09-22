/*
 * File: student_management.cpp
 * Description: A comprehensive student management system demonstrating
 *              multiple C++ concepts including classes, file I/O, STL containers
 * Part of: C/C++ Learning Guide - Practical Project
 */

#include <iostream>
#include <fstream>
#include <vector>
#include <string>
#include <algorithm>
#include <iomanip>
using namespace std;

class Student {
private:
    int id;
    string name;
    int age;
    double gpa;
    vector<string> courses;
    
public:
    // Constructors
    Student() : id(0), name(""), age(0), gpa(0.0) {}
    
    Student(int id, string name, int age, double gpa) 
        : id(id), name(name), age(age), gpa(gpa) {}
    
    // Getters
    int getId() const { return id; }
    string getName() const { return name; }
    int getAge() const { return age; }
    double getGpa() const { return gpa; }
    vector<string> getCourses() const { return courses; }
    
    // Setters
    void setId(int newId) { id = newId; }
    void setName(const string& newName) { name = newName; }
    void setAge(int newAge) { age = newAge; }
    void setGpa(double newGpa) { 
        if (newGpa >= 0.0 && newGpa <= 4.0) {
            gpa = newGpa;
        }
    }
    
    // Course management
    void addCourse(const string& course) {
        courses.push_back(course);
    }
    
    void removeCourse(const string& course) {
        courses.erase(remove(courses.begin(), courses.end(), course), courses.end());
    }
    
    // Display methods
    void displayBasicInfo() const {
        cout << setw(5) << id << setw(20) << name 
             << setw(5) << age << setw(8) << fixed << setprecision(2) << gpa;
    }
    
    void displayFullInfo() const {
        cout << "\n=== Student Details ===" << endl;
        cout << "ID: " << id << endl;
        cout << "Name: " << name << endl;
        cout << "Age: " << age << endl;
        cout << "GPA: " << fixed << setprecision(2) << gpa << endl;
        cout << "Courses: ";
        if (courses.empty()) {
            cout << "None";
        } else {
            for (size_t i = 0; i < courses.size(); i++) {
                cout << courses[i];
                if (i < courses.size() - 1) cout << ", ";
            }
        }
        cout << endl;
    }
    
    // Utility methods
    bool isHonorStudent() const {
        return gpa >= 3.5;
    }
    
    string getGradeLevel() const {
        if (gpa >= 3.7) return "A";
        else if (gpa >= 3.0) return "B";
        else if (gpa >= 2.0) return "C";
        else if (gpa >= 1.0) return "D";
        else return "F";
    }
    
    // File I/O support
    string toFileString() const {
        string courseList = "";
        for (size_t i = 0; i < courses.size(); i++) {
            courseList += courses[i];
            if (i < courses.size() - 1) courseList += ";";
        }
        return to_string(id) + "," + name + "," + to_string(age) + "," + 
               to_string(gpa) + "," + courseList;
    }
    
    void fromFileString(const string& line) {
        size_t pos1 = line.find(',');
        size_t pos2 = line.find(',', pos1 + 1);
        size_t pos3 = line.find(',', pos2 + 1);
        size_t pos4 = line.find(',', pos3 + 1);
        
        id = stoi(line.substr(0, pos1));
        name = line.substr(pos1 + 1, pos2 - pos1 - 1);
        age = stoi(line.substr(pos2 + 1, pos3 - pos2 - 1));
        gpa = stod(line.substr(pos3 + 1, pos4 - pos3 - 1));
        
        // Parse courses
        courses.clear();
        if (pos4 != string::npos && pos4 + 1 < line.length()) {
            string courseList = line.substr(pos4 + 1);
            size_t start = 0;
            size_t end = courseList.find(';');
            
            while (end != string::npos) {
                courses.push_back(courseList.substr(start, end - start));
                start = end + 1;
                end = courseList.find(';', start);
            }
            if (start < courseList.length()) {
                courses.push_back(courseList.substr(start));
            }
        }
    }
};

class StudentManager {
private:
    vector<Student> students;
    string filename;
    
public:
    StudentManager(const string& file = "students.txt") : filename(file) {
        loadFromFile();
    }
    
    ~StudentManager() {
        saveToFile();
    }
    
    // Student management operations
    void addStudent() {
        int id, age;
        string name;
        double gpa;
        
        cout << "\n=== Add New Student ===" << endl;
        cout << "Enter ID: ";
        cin >> id;
        
        // Check if ID already exists
        if (findStudentById(id) != nullptr) {
            cout << "Error: Student with ID " << id << " already exists!" << endl;
            return;
        }
        
        cin.ignore(); // Clear input buffer
        cout << "Enter Name: ";
        getline(cin, name);
        cout << "Enter Age: ";
        cin >> age;
        cout << "Enter GPA (0.0-4.0): ";
        cin >> gpa;
        
        Student newStudent(id, name, age, gpa);
        students.push_back(newStudent);
        cout << "Student added successfully!" << endl;
    }
    
    void removeStudent() {
        int id;
        cout << "\nEnter student ID to remove: ";
        cin >> id;
        
        auto it = find_if(students.begin(), students.end(),
                         [id](const Student& s) { return s.getId() == id; });
        
        if (it != students.end()) {
            cout << "Removing student: " << it->getName() << endl;
            students.erase(it);
            cout << "Student removed successfully!" << endl;
        } else {
            cout << "Student with ID " << id << " not found!" << endl;
        }
    }
    
    void displayAllStudents() const {
        if (students.empty()) {
            cout << "\nNo students in the system." << endl;
            return;
        }
        
        cout << "\n=== All Students ===" << endl;
        cout << setw(5) << "ID" << setw(20) << "Name" 
             << setw(5) << "Age" << setw(8) << "GPA" << setw(8) << "Grade" << endl;
        cout << string(46, '-') << endl;
        
        for (const auto& student : students) {
            student.displayBasicInfo();
            cout << setw(8) << student.getGradeLevel() << endl;
        }
    }
    
    void searchStudent() const {
        int choice;
        cout << "\n=== Search Student ===" << endl;
        cout << "1. Search by ID" << endl;
        cout << "2. Search by Name" << endl;
        cout << "Choose option: ";
        cin >> choice;
        
        if (choice == 1) {
            int id;
            cout << "Enter student ID: ";
            cin >> id;
            
            const Student* student = findStudentById(id);
            if (student) {
                student->displayFullInfo();
            } else {
                cout << "Student with ID " << id << " not found!" << endl;
            }
        } else if (choice == 2) {
            string name;
            cin.ignore();
            cout << "Enter student name: ";
            getline(cin, name);
            
            bool found = false;
            for (const auto& student : students) {
                if (student.getName().find(name) != string::npos) {
                    student.displayFullInfo();
                    found = true;
                }
            }
            if (!found) {
                cout << "No students found with name containing '" << name << "'" << endl;
            }
        }
    }
    
    void manageCourses() {
        int id;
        cout << "\nEnter student ID: ";
        cin >> id;
        
        Student* student = findStudentById(id);
        if (!student) {
            cout << "Student not found!" << endl;
            return;
        }
        
        int choice;
        cout << "\n=== Course Management for " << student->getName() << " ===" << endl;
        cout << "1. Add Course" << endl;
        cout << "2. Remove Course" << endl;
        cout << "3. View Courses" << endl;
        cout << "Choose option: ";
        cin >> choice;
        
        if (choice == 1) {
            string course;
            cin.ignore();
            cout << "Enter course name: ";
            getline(cin, course);
            student->addCourse(course);
            cout << "Course added successfully!" << endl;
        } else if (choice == 2) {
            string course;
            cin.ignore();
            cout << "Enter course name to remove: ";
            getline(cin, course);
            student->removeCourse(course);
            cout << "Course removed successfully!" << endl;
        } else if (choice == 3) {
            student->displayFullInfo();
        }
    }
    
    void generateReports() const {
        if (students.empty()) {
            cout << "\nNo students to generate reports." << endl;
            return;
        }
        
        cout << "\n=== Student Reports ===" << endl;
        
        // Honor students
        cout << "\nHonor Students (GPA >= 3.5):" << endl;
        bool hasHonorStudents = false;
        for (const auto& student : students) {
            if (student.isHonorStudent()) {
                cout << "- " << student.getName() << " (GPA: " 
                     << fixed << setprecision(2) << student.getGpa() << ")" << endl;
                hasHonorStudents = true;
            }
        }
        if (!hasHonorStudents) {
            cout << "No honor students found." << endl;
        }
        
        // Statistics
        double totalGpa = 0;
        int totalAge = 0;
        for (const auto& student : students) {
            totalGpa += student.getGpa();
            totalAge += student.getAge();
        }
        
        cout << "\nClass Statistics:" << endl;
        cout << "Total Students: " << students.size() << endl;
        cout << "Average GPA: " << fixed << setprecision(2) 
             << totalGpa / students.size() << endl;
        cout << "Average Age: " << fixed << setprecision(1) 
             << static_cast<double>(totalAge) / students.size() << endl;
    }
    
    // File operations
    void saveToFile() const {
        ofstream file(filename);
        if (file.is_open()) {
            for (const auto& student : students) {
                file << student.toFileString() << endl;
            }
            file.close();
            cout << "Data saved to " << filename << endl;
        } else {
            cout << "Error: Could not save to file!" << endl;
        }
    }
    
    void loadFromFile() {
        ifstream file(filename);
        if (file.is_open()) {
            string line;
            students.clear();
            
            while (getline(file, line)) {
                if (!line.empty()) {
                    Student student;
                    student.fromFileString(line);
                    students.push_back(student);
                }
            }
            file.close();
            cout << "Loaded " << students.size() << " students from " << filename << endl;
        }
    }
    
private:
    Student* findStudentById(int id) {
        auto it = find_if(students.begin(), students.end(),
                         [id](const Student& s) { return s.getId() == id; });
        return (it != students.end()) ? &(*it) : nullptr;
    }
    
    const Student* findStudentById(int id) const {
        auto it = find_if(students.begin(), students.end(),
                         [id](const Student& s) { return s.getId() == id; });
        return (it != students.end()) ? &(*it) : nullptr;
    }
};

void displayMenu() {
    cout << "\n=== Student Management System ===" << endl;
    cout << "1. Add Student" << endl;
    cout << "2. Remove Student" << endl;
    cout << "3. Display All Students" << endl;
    cout << "4. Search Student" << endl;
    cout << "5. Manage Courses" << endl;
    cout << "6. Generate Reports" << endl;
    cout << "7. Save Data" << endl;
    cout << "8. Load Data" << endl;
    cout << "9. Exit" << endl;
    cout << "Choose an option (1-9): ";
}

int main() {
    cout << "Welcome to Student Management System!" << endl;
    
    StudentManager manager;
    int choice;
    
    do {
        displayMenu();
        cin >> choice;
        
        switch (choice) {
            case 1:
                manager.addStudent();
                break;
            case 2:
                manager.removeStudent();
                break;
            case 3:
                manager.displayAllStudents();
                break;
            case 4:
                manager.searchStudent();
                break;
            case 5:
                manager.manageCourses();
                break;
            case 6:
                manager.generateReports();
                break;
            case 7:
                manager.saveToFile();
                break;
            case 8:
                manager.loadFromFile();
                break;
            case 9:
                cout << "\nThank you for using Student Management System!" << endl;
                break;
            default:
                cout << "Invalid option! Please try again." << endl;
        }
        
        if (choice != 9) {
            cout << "\nPress Enter to continue...";
            cin.ignore();
            cin.get();
        }
        
    } while (choice != 9);
    
    return 0;
}

/*
 * This project demonstrates:
 * 
 * 1. Object-Oriented Programming:
 *    - Classes with private/public members
 *    - Constructors and destructors
 *    - Encapsulation and data hiding
 *    - Method overloading
 * 
 * 2. STL Usage:
 *    - vector for dynamic arrays
 *    - string for text manipulation
 *    - algorithms (find_if, remove)
 *    - iterators
 * 
 * 3. File I/O:
 *    - Reading from and writing to files
 *    - Data persistence
 *    - Error handling
 * 
 * 4. Advanced C++ Features:
 *    - Lambda expressions
 *    - Range-based for loops
 *    - Auto keyword
 *    - Smart pointers concepts
 * 
 * 5. Programming Best Practices:
 *    - Code organization
 *    - Error handling
 *    - User input validation
 *    - Modular design
 */