# Part 6: Testing Strategies in DevOps

## Introduction

Testing is a critical component of DevOps that ensures software quality, reliability, and security throughout the development lifecycle. In DevOps, testing is not just a phase but an integrated practice that occurs continuously from development to production.

## The Testing Pyramid

The testing pyramid is a fundamental concept that guides testing strategy by showing the ideal distribution of different types of tests.

```
        /\     
       /  \    End-to-End Tests (Few)
      /____\   
     /      \  
    /        \  Integration Tests (Some)
   /__________\ 
  /            \
 /              \ Unit Tests (Many)
/________________\
```

### Unit Tests (Base of Pyramid)
- **Purpose**: Test individual components in isolation
- **Characteristics**: Fast, reliable, cheap to maintain
- **Coverage**: 70-80% of total tests
- **Execution Time**: Milliseconds to seconds

### Integration Tests (Middle)
- **Purpose**: Test interactions between components
- **Characteristics**: Moderate speed, more complex setup
- **Coverage**: 15-25% of total tests
- **Execution Time**: Seconds to minutes

### End-to-End Tests (Top)
- **Purpose**: Test complete user workflows
- **Characteristics**: Slow, brittle, expensive to maintain
- **Coverage**: 5-10% of total tests
- **Execution Time**: Minutes to hours

## Types of Testing in DevOps

### 1. Unit Testing

**Definition**: Testing individual units of code (functions, methods, classes) in isolation.

**Java Example with JUnit 5**:
```java
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.DisplayName;
import org.mockito.Mock;
import org.mockito.MockitoAnnotations;
import static org.junit.jupiter.api.Assertions.*;
import static org.mockito.Mockito.*;

public class UserServiceTest {
    
    @Mock
    private UserRepository userRepository;
    
    @Mock
    private EmailService emailService;
    
    private UserService userService;
    
    @BeforeEach
    void setUp() {
        MockitoAnnotations.openMocks(this);
        userService = new UserService(userRepository, emailService);
    }
    
    @Test
    @DisplayName("Should create user successfully")
    void shouldCreateUserSuccessfully() {
        // Given
        User newUser = new User("john@example.com", "John Doe");
        User savedUser = new User(1L, "john@example.com", "John Doe");
        
        when(userRepository.existsByEmail("john@example.com")).thenReturn(false);
        when(userRepository.save(newUser)).thenReturn(savedUser);
        
        // When
        User result = userService.createUser(newUser);
        
        // Then
        assertNotNull(result);
        assertEquals(1L, result.getId());
        assertEquals("john@example.com", result.getEmail());
        verify(emailService).sendWelcomeEmail(savedUser);
    }
    
    @Test
    @DisplayName("Should throw exception when email already exists")
    void shouldThrowExceptionWhenEmailExists() {
        // Given
        User newUser = new User("john@example.com", "John Doe");
        when(userRepository.existsByEmail("john@example.com")).thenReturn(true);
        
        // When & Then
        assertThrows(UserAlreadyExistsException.class, () -> {
            userService.createUser(newUser);
        });
        
        verify(userRepository, never()).save(any(User.class));
        verify(emailService, never()).sendWelcomeEmail(any(User.class));
    }
}
```

**JavaScript Example with Jest**:
```javascript
// userService.js
class UserService {
  constructor(userRepository, emailService) {
    this.userRepository = userRepository;
    this.emailService = emailService;
  }

  async createUser(userData) {
    const existingUser = await this.userRepository.findByEmail(userData.email);
    if (existingUser) {
      throw new Error('User already exists');
    }

    const user = await this.userRepository.create(userData);
    await this.emailService.sendWelcomeEmail(user);
    return user;
  }

  async getUserById(id) {
    if (!id) {
      throw new Error('User ID is required');
    }
    return await this.userRepository.findById(id);
  }
}

module.exports = UserService;

// userService.test.js
const UserService = require('./userService');

describe('UserService', () => {
  let userService;
  let mockUserRepository;
  let mockEmailService;

  beforeEach(() => {
    mockUserRepository = {
      findByEmail: jest.fn(),
      create: jest.fn(),
      findById: jest.fn()
    };
    
    mockEmailService = {
      sendWelcomeEmail: jest.fn()
    };
    
    userService = new UserService(mockUserRepository, mockEmailService);
  });

  describe('createUser', () => {
    it('should create user successfully', async () => {
      // Arrange
      const userData = { email: 'john@example.com', name: 'John Doe' };
      const createdUser = { id: 1, ...userData };
      
      mockUserRepository.findByEmail.mockResolvedValue(null);
      mockUserRepository.create.mockResolvedValue(createdUser);
      mockEmailService.sendWelcomeEmail.mockResolvedValue();

      // Act
      const result = await userService.createUser(userData);

      // Assert
      expect(result).toEqual(createdUser);
      expect(mockUserRepository.findByEmail).toHaveBeenCalledWith('john@example.com');
      expect(mockUserRepository.create).toHaveBeenCalledWith(userData);
      expect(mockEmailService.sendWelcomeEmail).toHaveBeenCalledWith(createdUser);
    });

    it('should throw error when user already exists', async () => {
      // Arrange
      const userData = { email: 'john@example.com', name: 'John Doe' };
      const existingUser = { id: 1, ...userData };
      
      mockUserRepository.findByEmail.mockResolvedValue(existingUser);

      // Act & Assert
      await expect(userService.createUser(userData))
        .rejects.toThrow('User already exists');
      
      expect(mockUserRepository.create).not.toHaveBeenCalled();
      expect(mockEmailService.sendWelcomeEmail).not.toHaveBeenCalled();
    });
  });

  describe('getUserById', () => {
    it('should return user when valid ID provided', async () => {
      // Arrange
      const userId = 1;
      const user = { id: userId, email: 'john@example.com', name: 'John Doe' };
      mockUserRepository.findById.mockResolvedValue(user);

      // Act
      const result = await userService.getUserById(userId);

      // Assert
      expect(result).toEqual(user);
      expect(mockUserRepository.findById).toHaveBeenCalledWith(userId);
    });

    it('should throw error when no ID provided', async () => {
      // Act & Assert
      await expect(userService.getUserById())
        .rejects.toThrow('User ID is required');
      
      expect(mockUserRepository.findById).not.toHaveBeenCalled();
    });
  });
});
```

**Python Example with pytest**:
```python
# user_service.py
from typing import Optional
from dataclasses import dataclass

@dataclass
class User:
    id: Optional[int]
    email: str
    name: str

class UserAlreadyExistsError(Exception):
    pass

class UserService:
    def __init__(self, user_repository, email_service):
        self.user_repository = user_repository
        self.email_service = email_service
    
    def create_user(self, user_data: dict) -> User:
        existing_user = self.user_repository.find_by_email(user_data['email'])
        if existing_user:
            raise UserAlreadyExistsError("User already exists")
        
        user = self.user_repository.create(user_data)
        self.email_service.send_welcome_email(user)
        return user
    
    def get_user_by_id(self, user_id: int) -> Optional[User]:
        if not user_id:
            raise ValueError("User ID is required")
        return self.user_repository.find_by_id(user_id)

# test_user_service.py
import pytest
from unittest.mock import Mock, MagicMock
from user_service import UserService, User, UserAlreadyExistsError

class TestUserService:
    
    @pytest.fixture
    def user_repository(self):
        return Mock()
    
    @pytest.fixture
    def email_service(self):
        return Mock()
    
    @pytest.fixture
    def user_service(self, user_repository, email_service):
        return UserService(user_repository, email_service)
    
    def test_create_user_success(self, user_service, user_repository, email_service):
        # Arrange
        user_data = {'email': 'john@example.com', 'name': 'John Doe'}
        created_user = User(id=1, email='john@example.com', name='John Doe')
        
        user_repository.find_by_email.return_value = None
        user_repository.create.return_value = created_user
        
        # Act
        result = user_service.create_user(user_data)
        
        # Assert
        assert result == created_user
        user_repository.find_by_email.assert_called_once_with('john@example.com')
        user_repository.create.assert_called_once_with(user_data)
        email_service.send_welcome_email.assert_called_once_with(created_user)
    
    def test_create_user_already_exists(self, user_service, user_repository, email_service):
        # Arrange
        user_data = {'email': 'john@example.com', 'name': 'John Doe'}
        existing_user = User(id=1, email='john@example.com', name='John Doe')
        
        user_repository.find_by_email.return_value = existing_user
        
        # Act & Assert
        with pytest.raises(UserAlreadyExistsError, match="User already exists"):
            user_service.create_user(user_data)
        
        user_repository.create.assert_not_called()
        email_service.send_welcome_email.assert_not_called()
    
    def test_get_user_by_id_success(self, user_service, user_repository):
        # Arrange
        user_id = 1
        user = User(id=user_id, email='john@example.com', name='John Doe')
        user_repository.find_by_id.return_value = user
        
        # Act
        result = user_service.get_user_by_id(user_id)
        
        # Assert
        assert result == user
        user_repository.find_by_id.assert_called_once_with(user_id)
    
    @pytest.mark.parametrize("invalid_id", [None, 0, ""])
    def test_get_user_by_id_invalid_id(self, user_service, user_repository, invalid_id):
        # Act & Assert
        with pytest.raises(ValueError, match="User ID is required"):
            user_service.get_user_by_id(invalid_id)
        
        user_repository.find_by_id.assert_not_called()
```

### 2. Integration Testing

**Definition**: Testing the interaction between different components or services.

**Spring Boot Integration Test Example**:
```java
@SpringBootTest(webEnvironment = SpringBootTest.WebEnvironment.RANDOM_PORT)
@TestPropertySource(locations = "classpath:application-test.properties")
@Testcontainers
public class UserControllerIntegrationTest {
    
    @Container
    static PostgreSQLContainer<?> postgres = new PostgreSQLContainer<>("postgres:13")
            .withDatabaseName("testdb")
            .withUsername("test")
            .withPassword("test");
    
    @Autowired
    private TestRestTemplate restTemplate;
    
    @Autowired
    private UserRepository userRepository;
    
    @DynamicPropertySource
    static void configureProperties(DynamicPropertyRegistry registry) {
        registry.add("spring.datasource.url", postgres::getJdbcUrl);
        registry.add("spring.datasource.username", postgres::getUsername);
        registry.add("spring.datasource.password", postgres::getPassword);
    }
    
    @BeforeEach
    void setUp() {
        userRepository.deleteAll();
    }
    
    @Test
    void shouldCreateUserSuccessfully() {
        // Given
        CreateUserRequest request = new CreateUserRequest("john@example.com", "John Doe");
        
        // When
        ResponseEntity<UserResponse> response = restTemplate.postForEntity(
            "/api/users", request, UserResponse.class);
        
        // Then
        assertEquals(HttpStatus.CREATED, response.getStatusCode());
        assertNotNull(response.getBody());
        assertEquals("john@example.com", response.getBody().getEmail());
        assertEquals("John Doe", response.getBody().getName());
        
        // Verify in database
        Optional<User> savedUser = userRepository.findByEmail("john@example.com");
        assertTrue(savedUser.isPresent());
        assertEquals("John Doe", savedUser.get().getName());
    }
    
    @Test
    void shouldReturnConflictWhenUserAlreadyExists() {
        // Given
        User existingUser = new User("john@example.com", "John Doe");
        userRepository.save(existingUser);
        
        CreateUserRequest request = new CreateUserRequest("john@example.com", "Jane Doe");
        
        // When
        ResponseEntity<ErrorResponse> response = restTemplate.postForEntity(
            "/api/users", request, ErrorResponse.class);
        
        // Then
        assertEquals(HttpStatus.CONFLICT, response.getStatusCode());
        assertNotNull(response.getBody());
        assertEquals("User already exists", response.getBody().getMessage());
    }
}
```

**Node.js Integration Test with Supertest**:
```javascript
// integration.test.js
const request = require('supertest');
const app = require('../app');
const { setupTestDB, cleanupTestDB } = require('./helpers/db');

describe('User API Integration Tests', () => {
  beforeAll(async () => {
    await setupTestDB();
  });

  afterAll(async () => {
    await cleanupTestDB();
  });

  beforeEach(async () => {
    // Clean database before each test
    await request(app).delete('/api/test/cleanup');
  });

  describe('POST /api/users', () => {
    it('should create user successfully', async () => {
      const userData = {
        email: 'john@example.com',
        name: 'John Doe',
        password: 'securePassword123'
      };

      const response = await request(app)
        .post('/api/users')
        .send(userData)
        .expect(201);

      expect(response.body).toMatchObject({
        id: expect.any(Number),
        email: 'john@example.com',
        name: 'John Doe'
      });
      expect(response.body.password).toBeUndefined();
    });

    it('should return 400 for invalid email', async () => {
      const userData = {
        email: 'invalid-email',
        name: 'John Doe',
        password: 'securePassword123'
      };

      const response = await request(app)
        .post('/api/users')
        .send(userData)
        .expect(400);

      expect(response.body.error).toContain('Invalid email format');
    });
  });

  describe('GET /api/users/:id', () => {
    it('should return user by id', async () => {
      // Create user first
      const createResponse = await request(app)
        .post('/api/users')
        .send({
          email: 'john@example.com',
          name: 'John Doe',
          password: 'securePassword123'
        });

      const userId = createResponse.body.id;

      // Get user by id
      const response = await request(app)
        .get(`/api/users/${userId}`)
        .expect(200);

      expect(response.body).toMatchObject({
        id: userId,
        email: 'john@example.com',
        name: 'John Doe'
      });
    });

    it('should return 404 for non-existent user', async () => {
      const response = await request(app)
        .get('/api/users/99999')
        .expect(404);

      expect(response.body.error).toContain('User not found');
    });
  });
});
```

### 3. End-to-End (E2E) Testing

**Definition**: Testing complete user workflows from start to finish.

**Cypress Example**:
```javascript
// cypress/integration/user-registration.spec.js
describe('User Registration Flow', () => {
  beforeEach(() => {
    // Reset database state
    cy.task('db:seed');
    cy.visit('/register');
  });

  it('should register new user successfully', () => {
    // Fill registration form
    cy.get('[data-testid="email-input"]')
      .type('john@example.com');
    
    cy.get('[data-testid="name-input"]')
      .type('John Doe');
    
    cy.get('[data-testid="password-input"]')
      .type('SecurePassword123!');
    
    cy.get('[data-testid="confirm-password-input"]')
      .type('SecurePassword123!');
    
    // Submit form
    cy.get('[data-testid="register-button"]')
      .click();
    
    // Verify success
    cy.get('[data-testid="success-message"]')
      .should('contain', 'Registration successful');
    
    // Verify redirect to dashboard
    cy.url().should('include', '/dashboard');
    
    // Verify user is logged in
    cy.get('[data-testid="user-menu"]')
      .should('contain', 'John Doe');
  });

  it('should show validation errors for invalid input', () => {
    // Submit empty form
    cy.get('[data-testid="register-button"]')
      .click();
    
    // Check validation errors
    cy.get('[data-testid="email-error"]')
      .should('contain', 'Email is required');
    
    cy.get('[data-testid="name-error"]')
      .should('contain', 'Name is required');
    
    cy.get('[data-testid="password-error"]')
      .should('contain', 'Password is required');
  });

  it('should handle duplicate email registration', () => {
    // Try to register with existing email
    cy.get('[data-testid="email-input"]')
      .type('existing@example.com');
    
    cy.get('[data-testid="name-input"]')
      .type('John Doe');
    
    cy.get('[data-testid="password-input"]')
      .type('SecurePassword123!');
    
    cy.get('[data-testid="confirm-password-input"]')
      .type('SecurePassword123!');
    
    cy.get('[data-testid="register-button"]')
      .click();
    
    // Verify error message
    cy.get('[data-testid="error-message"]')
      .should('contain', 'Email already exists');
    
    // Verify form is still visible
    cy.get('[data-testid="register-form"]')
      .should('be.visible');
  });
});

// cypress/support/commands.js
Cypress.Commands.add('login', (email, password) => {
  cy.request({
    method: 'POST',
    url: '/api/auth/login',
    body: { email, password }
  }).then((response) => {
    window.localStorage.setItem('authToken', response.body.token);
  });
});

Cypress.Commands.add('createUser', (userData) => {
  cy.request({
    method: 'POST',
    url: '/api/users',
    body: userData
  });
});
```

**Playwright Example**:
```javascript
// tests/user-registration.spec.js
const { test, expect } = require('@playwright/test');

test.describe('User Registration', () => {
  test.beforeEach(async ({ page }) => {
    // Setup test data
    await page.request.post('/api/test/reset-db');
    await page.goto('/register');
  });

  test('should register new user successfully', async ({ page }) => {
    // Fill form
    await page.fill('[data-testid="email"]', 'john@example.com');
    await page.fill('[data-testid="name"]', 'John Doe');
    await page.fill('[data-testid="password"]', 'SecurePassword123!');
    await page.fill('[data-testid="confirmPassword"]', 'SecurePassword123!');
    
    // Submit
    await page.click('[data-testid="submit"]');
    
    // Verify success
    await expect(page.locator('[data-testid="success-message"]'))
      .toContainText('Registration successful');
    
    // Verify navigation
    await expect(page).toHaveURL(/.*\/dashboard/);
    
    // Verify user menu
    await expect(page.locator('[data-testid="user-menu"]'))
      .toContainText('John Doe');
  });

  test('should validate form fields', async ({ page }) => {
    // Submit empty form
    await page.click('[data-testid="submit"]');
    
    // Check validation messages
    await expect(page.locator('[data-testid="email-error"]'))
      .toContainText('Email is required');
    
    await expect(page.locator('[data-testid="name-error"]'))
      .toContainText('Name is required');
    
    await expect(page.locator('[data-testid="password-error"]'))
      .toContainText('Password is required');
  });

  test('should handle server errors gracefully', async ({ page }) => {
    // Mock server error
    await page.route('/api/users', route => {
      route.fulfill({
        status: 500,
        contentType: 'application/json',
        body: JSON.stringify({ error: 'Internal server error' })
      });
    });
    
    // Fill and submit form
    await page.fill('[data-testid="email"]', 'john@example.com');
    await page.fill('[data-testid="name"]', 'John Doe');
    await page.fill('[data-testid="password"]', 'SecurePassword123!');
    await page.fill('[data-testid="confirmPassword"]', 'SecurePassword123!');
    await page.click('[data-testid="submit"]');
    
    // Verify error handling
    await expect(page.locator('[data-testid="error-message"]'))
      .toContainText('Something went wrong. Please try again.');
  });
});
```

### 4. Performance Testing

**Definition**: Testing system performance under various load conditions.

**JMeter Test Plan Example**:
```xml
<?xml version="1.0" encoding="UTF-8"?>
<jmeterTestPlan version="1.2">
  <hashTree>
    <TestPlan guiclass="TestPlanGui" testclass="TestPlan" testname="API Load Test">
      <stringProp name="TestPlan.comments">Load test for user API</stringProp>
      <boolProp name="TestPlan.functional_mode">false</boolProp>
      <boolProp name="TestPlan.serialize_threadgroups">false</boolProp>
      <elementProp name="TestPlan.arguments" elementType="Arguments" guiclass="ArgumentsPanel">
        <collectionProp name="Arguments.arguments"/>
      </elementProp>
      <stringProp name="TestPlan.user_define_classpath"></stringProp>
    </TestPlan>
    <hashTree>
      <ThreadGroup guiclass="ThreadGroupGui" testclass="ThreadGroup" testname="User API Load">
        <stringProp name="ThreadGroup.on_sample_error">continue</stringProp>
        <elementProp name="ThreadGroup.main_controller" elementType="LoopController">
          <boolProp name="LoopController.continue_forever">false</boolProp>
          <stringProp name="LoopController.loops">10</stringProp>
        </elementProp>
        <stringProp name="ThreadGroup.num_threads">100</stringProp>
        <stringProp name="ThreadGroup.ramp_time">60</stringProp>
        <longProp name="ThreadGroup.start_time">1</longProp>
        <longProp name="ThreadGroup.end_time">1</longProp>
        <boolProp name="ThreadGroup.scheduler">false</boolProp>
        <stringProp name="ThreadGroup.duration"></stringProp>
        <stringProp name="ThreadGroup.delay"></stringProp>
      </ThreadGroup>
      <hashTree>
        <HTTPSamplerProxy guiclass="HttpTestSampleGui" testclass="HTTPSamplerProxy" testname="Create User">
          <elementProp name="HTTPsampler.Arguments" elementType="Arguments">
            <collectionProp name="Arguments.arguments">
              <elementProp name="" elementType="HTTPArgument">
                <boolProp name="HTTPArgument.always_encode">false</boolProp>
                <stringProp name="Argument.value">{
  "email": "user${__threadNum}@example.com",
  "name": "User ${__threadNum}",
  "password": "password123"
}</stringProp>
                <stringProp name="Argument.metadata">=</stringProp>
              </elementProp>
            </collectionProp>
          </elementProp>
          <stringProp name="HTTPSampler.domain">localhost</stringProp>
          <stringProp name="HTTPSampler.port">8080</stringProp>
          <stringProp name="HTTPSampler.protocol">http</stringProp>
          <stringProp name="HTTPSampler.contentEncoding"></stringProp>
          <stringProp name="HTTPSampler.path">/api/users</stringProp>
          <stringProp name="HTTPSampler.method">POST</stringProp>
          <boolProp name="HTTPSampler.follow_redirects">true</boolProp>
          <boolProp name="HTTPSampler.auto_redirects">false</boolProp>
          <boolProp name="HTTPSampler.use_keepalive">true</boolProp>
          <boolProp name="HTTPSampler.DO_MULTIPART_POST">false</boolProp>
          <stringProp name="HTTPSampler.embedded_url_re"></stringProp>
          <stringProp name="HTTPSampler.connect_timeout"></stringProp>
          <stringProp name="HTTPSampler.response_timeout"></stringProp>
        </HTTPSamplerProxy>
      </hashTree>
    </hashTree>
  </hashTree>
</jmeterTestPlan>
```

**K6 Performance Test Example**:
```javascript
// load-test.js
import http from 'k6/http';
import { check, sleep } from 'k6';
import { Rate } from 'k6/metrics';

// Custom metrics
const errorRate = new Rate('errors');

// Test configuration
export const options = {
  stages: [
    { duration: '2m', target: 10 },   // Ramp up to 10 users
    { duration: '5m', target: 10 },   // Stay at 10 users
    { duration: '2m', target: 50 },   // Ramp up to 50 users
    { duration: '5m', target: 50 },   // Stay at 50 users
    { duration: '2m', target: 100 },  // Ramp up to 100 users
    { duration: '5m', target: 100 },  // Stay at 100 users
    { duration: '2m', target: 0 },    // Ramp down to 0 users
  ],
  thresholds: {
    http_req_duration: ['p(95)<500'], // 95% of requests must complete below 500ms
    http_req_failed: ['rate<0.1'],    // Error rate must be below 10%
    errors: ['rate<0.1'],             // Custom error rate must be below 10%
  },
};

const BASE_URL = 'http://localhost:8080';

export function setup() {
  // Setup test data
  console.log('Setting up test data...');
  return { baseUrl: BASE_URL };
}

export default function (data) {
  const userId = Math.floor(Math.random() * 10000);
  
  // Test user creation
  const createUserPayload = {
    email: `user${userId}@example.com`,
    name: `User ${userId}`,
    password: 'password123'
  };
  
  const createResponse = http.post(
    `${data.baseUrl}/api/users`,
    JSON.stringify(createUserPayload),
    {
      headers: {
        'Content-Type': 'application/json',
      },
    }
  );
  
  const createSuccess = check(createResponse, {
    'user creation status is 201': (r) => r.status === 201,
    'user creation response time < 500ms': (r) => r.timings.duration < 500,
    'user creation response has id': (r) => JSON.parse(r.body).id !== undefined,
  });
  
  errorRate.add(!createSuccess);
  
  if (createSuccess) {
    const createdUser = JSON.parse(createResponse.body);
    
    // Test user retrieval
    const getResponse = http.get(`${data.baseUrl}/api/users/${createdUser.id}`);
    
    const getSuccess = check(getResponse, {
      'user retrieval status is 200': (r) => r.status === 200,
      'user retrieval response time < 200ms': (r) => r.timings.duration < 200,
      'user retrieval returns correct user': (r) => {
        const user = JSON.parse(r.body);
        return user.id === createdUser.id && user.email === createdUser.email;
      },
    });
    
    errorRate.add(!getSuccess);
  }
  
  sleep(1); // Wait 1 second between iterations
}

export function teardown(data) {
  // Cleanup test data
  console.log('Cleaning up test data...');
}
```

### 5. Security Testing

**Definition**: Testing for security vulnerabilities and compliance.

**OWASP ZAP Automation Example**:
```python
#!/usr/bin/env python3
# security_test.py
import time
import requests
from zapv2 import ZAPv2

class SecurityTester:
    def __init__(self, target_url, zap_proxy_url='http://127.0.0.1:8080'):
        self.target_url = target_url
        self.zap = ZAPv2(proxies={'http': zap_proxy_url, 'https': zap_proxy_url})
        
    def run_security_tests(self):
        print(f"Starting security tests for {self.target_url}")
        
        # Start ZAP session
        self.zap.core.new_session()
        
        # Spider the application
        print("Starting spider scan...")
        spider_id = self.zap.spider.scan(self.target_url)
        
        # Wait for spider to complete
        while int(self.zap.spider.status(spider_id)) < 100:
            print(f"Spider progress: {self.zap.spider.status(spider_id)}%")
            time.sleep(2)
        
        print("Spider scan completed")
        
        # Run active scan
        print("Starting active scan...")
        active_scan_id = self.zap.ascan.scan(self.target_url)
        
        # Wait for active scan to complete
        while int(self.zap.ascan.status(active_scan_id)) < 100:
            print(f"Active scan progress: {self.zap.ascan.status(active_scan_id)}%")
            time.sleep(5)
        
        print("Active scan completed")
        
        # Generate report
        self.generate_report()
        
    def generate_report(self):
        alerts = self.zap.core.alerts()
        
        high_risk = [alert for alert in alerts if alert['risk'] == 'High']
        medium_risk = [alert for alert in alerts if alert['risk'] == 'Medium']
        low_risk = [alert for alert in alerts if alert['risk'] == 'Low']
        
        print(f"\nSecurity Scan Results:")
        print(f"High Risk Issues: {len(high_risk)}")
        print(f"Medium Risk Issues: {len(medium_risk)}")
        print(f"Low Risk Issues: {len(low_risk)}")
        
        if high_risk:
            print("\nHigh Risk Issues:")
            for alert in high_risk:
                print(f"- {alert['alert']} at {alert['url']}")
        
        # Save HTML report
        html_report = self.zap.core.htmlreport()
        with open('security_report.html', 'w') as f:
            f.write(html_report)
        
        print("\nDetailed report saved to security_report.html")
        
        # Fail if high risk issues found
        if high_risk:
            raise Exception(f"Security scan failed: {len(high_risk)} high risk issues found")

if __name__ == "__main__":
    import sys
    
    if len(sys.argv) != 2:
        print("Usage: python security_test.py <target_url>")
        sys.exit(1)
    
    target_url = sys.argv[1]
    tester = SecurityTester(target_url)
    
    try:
        tester.run_security_tests()
        print("Security tests passed!")
    except Exception as e:
        print(f"Security tests failed: {e}")
        sys.exit(1)
```

## Test Automation in CI/CD

### GitHub Actions Example

```yaml
# .github/workflows/test.yml
name: Test Pipeline

on:
  push:
    branches: [ main, develop ]
  pull_request:
    branches: [ main ]

jobs:
  unit-tests:
    runs-on: ubuntu-latest
    
    services:
      postgres:
        image: postgres:13
        env:
          POSTGRES_PASSWORD: postgres
          POSTGRES_DB: testdb
        options: >-
          --health-cmd pg_isready
          --health-interval 10s
          --health-timeout 5s
          --health-retries 5
        ports:
          - 5432:5432
    
    steps:
    - uses: actions/checkout@v3
    
    - name: Set up JDK 11
      uses: actions/setup-java@v3
      with:
        java-version: '11'
        distribution: 'temurin'
    
    - name: Cache Maven dependencies
      uses: actions/cache@v3
      with:
        path: ~/.m2
        key: ${{ runner.os }}-m2-${{ hashFiles('**/pom.xml') }}
        restore-keys: ${{ runner.os }}-m2
    
    - name: Run unit tests
      run: mvn test
    
    - name: Generate test report
      uses: dorny/test-reporter@v1
      if: success() || failure()
      with:
        name: Maven Tests
        path: target/surefire-reports/*.xml
        reporter: java-junit
    
    - name: Upload coverage to Codecov
      uses: codecov/codecov-action@v3
      with:
        file: target/site/jacoco/jacoco.xml

  integration-tests:
    runs-on: ubuntu-latest
    needs: unit-tests
    
    services:
      postgres:
        image: postgres:13
        env:
          POSTGRES_PASSWORD: postgres
          POSTGRES_DB: testdb
        options: >-
          --health-cmd pg_isready
          --health-interval 10s
          --health-timeout 5s
          --health-retries 5
        ports:
          - 5432:5432
      
      redis:
        image: redis:alpine
        options: >-
          --health-cmd "redis-cli ping"
          --health-interval 10s
          --health-timeout 5s
          --health-retries 5
        ports:
          - 6379:6379
    
    steps:
    - uses: actions/checkout@v3
    
    - name: Set up JDK 11
      uses: actions/setup-java@v3
      with:
        java-version: '11'
        distribution: 'temurin'
    
    - name: Cache Maven dependencies
      uses: actions/cache@v3
      with:
        path: ~/.m2
        key: ${{ runner.os }}-m2-${{ hashFiles('**/pom.xml') }}
        restore-keys: ${{ runner.os }}-m2
    
    - name: Run integration tests
      run: mvn verify -Pintegration-tests
      env:
        DATABASE_URL: jdbc:postgresql://localhost:5432/testdb
        DATABASE_USERNAME: postgres
        DATABASE_PASSWORD: postgres
        REDIS_URL: redis://localhost:6379

  e2e-tests:
    runs-on: ubuntu-latest
    needs: integration-tests
    
    steps:
    - uses: actions/checkout@v3
    
    - name: Set up Node.js
      uses: actions/setup-node@v3
      with:
        node-version: '16'
        cache: 'npm'
    
    - name: Install dependencies
      run: npm ci
    
    - name: Build application
      run: npm run build
    
    - name: Start application
      run: |
        npm start &
        npx wait-on http://localhost:3000
    
    - name: Run E2E tests
      run: npx cypress run
      env:
        CYPRESS_baseUrl: http://localhost:3000
    
    - name: Upload E2E test videos
      uses: actions/upload-artifact@v3
      if: failure()
      with:
        name: cypress-videos
        path: cypress/videos

  security-tests:
    runs-on: ubuntu-latest
    needs: integration-tests
    
    steps:
    - uses: actions/checkout@v3
    
    - name: Run Trivy vulnerability scanner
      uses: aquasecurity/trivy-action@master
      with:
        scan-type: 'fs'
        scan-ref: '.'
        format: 'sarif'
        output: 'trivy-results.sarif'
    
    - name: Upload Trivy scan results to GitHub Security tab
      uses: github/codeql-action/upload-sarif@v2
      with:
        sarif_file: 'trivy-results.sarif'
    
    - name: Run OWASP Dependency Check
      run: |
        mvn org.owasp:dependency-check-maven:check
    
    - name: Upload dependency check results
      uses: actions/upload-artifact@v3
      if: always()
      with:
        name: dependency-check-report
        path: target/dependency-check-report.html

  performance-tests:
    runs-on: ubuntu-latest
    needs: integration-tests
    if: github.ref == 'refs/heads/main'
    
    steps:
    - uses: actions/checkout@v3
    
    - name: Set up Node.js
      uses: actions/setup-node@v3
      with:
        node-version: '16'
    
    - name: Install k6
      run: |
        sudo apt-key adv --keyserver hkp://keyserver.ubuntu.com:80 --recv-keys C5AD17C747E3415A3642D57D77C6C491D6AC1D69
        echo "deb https://dl.k6.io/deb stable main" | sudo tee /etc/apt/sources.list.d/k6.list
        sudo apt-get update
        sudo apt-get install k6
    
    - name: Start application
      run: |
        npm start &
        npx wait-on http://localhost:3000
    
    - name: Run performance tests
      run: k6 run tests/performance/load-test.js
    
    - name: Upload performance test results
      uses: actions/upload-artifact@v3
      if: always()
      with:
        name: performance-test-results
        path: performance-results.json
```

## Test Data Management

### Test Data Strategies

**1. Test Data Builders**:
```java
// UserTestDataBuilder.java
public class UserTestDataBuilder {
    private String email = "test@example.com";
    private String name = "Test User";
    private String password = "password123";
    private boolean active = true;
    private LocalDateTime createdAt = LocalDateTime.now();
    
    public static UserTestDataBuilder aUser() {
        return new UserTestDataBuilder();
    }
    
    public UserTestDataBuilder withEmail(String email) {
        this.email = email;
        return this;
    }
    
    public UserTestDataBuilder withName(String name) {
        this.name = name;
        return this;
    }
    
    public UserTestDataBuilder withPassword(String password) {
        this.password = password;
        return this;
    }
    
    public UserTestDataBuilder inactive() {
        this.active = false;
        return this;
    }
    
    public UserTestDataBuilder createdAt(LocalDateTime createdAt) {
        this.createdAt = createdAt;
        return this;
    }
    
    public User build() {
        return new User(email, name, password, active, createdAt);
    }
    
    public User buildAndSave(UserRepository repository) {
        return repository.save(build());
    }
}

// Usage in tests
@Test
void shouldFindActiveUsers() {
    // Given
    User activeUser = aUser()
        .withEmail("active@example.com")
        .buildAndSave(userRepository);
    
    User inactiveUser = aUser()
        .withEmail("inactive@example.com")
        .inactive()
        .buildAndSave(userRepository);
    
    // When
    List<User> activeUsers = userService.findActiveUsers();
    
    // Then
    assertThat(activeUsers)
        .hasSize(1)
        .contains(activeUser)
        .doesNotContain(inactiveUser);
}
```

**2. Test Fixtures**:
```javascript
// fixtures/users.js
const users = {
  validUser: {
    email: 'john@example.com',
    name: 'John Doe',
    password: 'SecurePassword123!'
  },
  
  adminUser: {
    email: 'admin@example.com',
    name: 'Admin User',
    password: 'AdminPassword123!',
    role: 'admin'
  },
  
  invalidUser: {
    email: 'invalid-email',
    name: '',
    password: '123'
  },
  
  createRandomUser: () => ({
    email: `user${Math.random().toString(36).substr(2, 9)}@example.com`,
    name: `User ${Math.random().toString(36).substr(2, 9)}`,
    password: 'Password123!'
  })
};

module.exports = users;

// Usage in tests
const { validUser, createRandomUser } = require('../fixtures/users');

describe('User Service', () => {
  it('should create user with valid data', async () => {
    const result = await userService.createUser(validUser);
    expect(result.email).toBe(validUser.email);
  });
  
  it('should handle multiple users', async () => {
    const user1 = await userService.createUser(createRandomUser());
    const user2 = await userService.createUser(createRandomUser());
    
    expect(user1.id).not.toBe(user2.id);
  });
});
```

### Database Testing Strategies

**1. Test Containers (Java)**:
```java
@Testcontainers
class DatabaseIntegrationTest {
    
    @Container
    static PostgreSQLContainer<?> postgres = new PostgreSQLContainer<>("postgres:13")
            .withDatabaseName("testdb")
            .withUsername("test")
            .withPassword("test");
    
    @DynamicPropertySource
    static void configureProperties(DynamicPropertyRegistry registry) {
        registry.add("spring.datasource.url", postgres::getJdbcUrl);
        registry.add("spring.datasource.username", postgres::getUsername);
        registry.add("spring.datasource.password", postgres::getPassword);
    }
    
    @Test
    void shouldPersistUser() {
        // Test implementation
    }
}
```

**2. In-Memory Databases**:
```yaml
# application-test.yml
spring:
  datasource:
    url: jdbc:h2:mem:testdb
    driver-class-name: org.h2.Driver
    username: sa
    password: 
  jpa:
    hibernate:
      ddl-auto: create-drop
    show-sql: true
```

## Test Reporting and Metrics

### Test Coverage

**JaCoCo Configuration (Maven)**:
```xml
<plugin>
    <groupId>org.jacoco</groupId>
    <artifactId>jacoco-maven-plugin</artifactId>
    <version>0.8.7</version>
    <executions>
        <execution>
            <goals>
                <goal>prepare-agent</goal>
            </goals>
        </execution>
        <execution>
            <id>report</id>
            <phase>test</phase>
            <goals>
                <goal>report</goal>
            </goals>
        </execution>
        <execution>
            <id>check</id>
            <goals>
                <goal>check</goal>
            </goals>
            <configuration>
                <rules>
                    <rule>
                        <element>PACKAGE</element>
                        <limits>
                            <limit>
                                <counter>LINE</counter>
                                <value>COVEREDRATIO</value>
                                <minimum>0.80</minimum>
                            </limit>
                        </limits>
                    </rule>
                </rules>
            </configuration>
        </execution>
    </executions>
</plugin>
```

### Test Metrics Dashboard

**Grafana Dashboard Configuration**:
```json
{
  "dashboard": {
    "title": "Test Metrics Dashboard",
    "panels": [
      {
        "title": "Test Success Rate",
        "type": "stat",
        "targets": [
          {
            "expr": "(sum(test_runs_total{status=\"passed\"}) / sum(test_runs_total)) * 100",
            "legendFormat": "Success Rate"
          }
        ]
      },
      {
        "title": "Test Execution Time",
        "type": "graph",
        "targets": [
          {
            "expr": "avg(test_duration_seconds) by (test_suite)",
            "legendFormat": "{{test_suite}}"
          }
        ]
      },
      {
        "title": "Code Coverage",
        "type": "graph",
        "targets": [
          {
            "expr": "code_coverage_percentage",
            "legendFormat": "Coverage %"
          }
        ]
      }
    ]
  }
}
```

## Best Practices

### Test Design Principles

1. **FIRST Principles**
   - **Fast**: Tests should run quickly
   - **Independent**: Tests should not depend on each other
   - **Repeatable**: Tests should produce same results every time
   - **Self-Validating**: Tests should have clear pass/fail result
   - **Timely**: Tests should be written at the right time

2. **AAA Pattern**
   - **Arrange**: Set up test data and conditions
   - **Act**: Execute the code under test
   - **Assert**: Verify the expected outcome

3. **Test Naming Conventions**
   ```java
   // Good test names
   shouldReturnUserWhenValidIdProvided()
   shouldThrowExceptionWhenUserNotFound()
   shouldCalculateDiscountForPremiumCustomers()
   
   // Bad test names
   testUser()
   test1()
   userTest()
   ```

### Test Environment Management

1. **Environment Isolation**
   - Separate test environments
   - Clean state for each test
   - No shared mutable state

2. **Test Data Management**
   - Use test data builders
   - Clean up after tests
   - Avoid hardcoded test data

3. **Configuration Management**
   - Environment-specific configurations
   - Externalized test properties
   - Secure handling of test credentials

## Practical Exercises

### Exercise 1: Build a Test Suite
Create a comprehensive test suite for a REST API that includes:
1. Unit tests for business logic
2. Integration tests for API endpoints
3. End-to-end tests for user workflows
4. Performance tests for load scenarios

### Exercise 2: Test Automation Pipeline
Set up a CI/CD pipeline that:
1. Runs tests on every commit
2. Generates test reports
3. Measures code coverage
4. Fails build on test failures

### Exercise 3: Test Data Strategy
Implement a test data management strategy that:
1. Uses test data builders
2. Provides clean test environments
3. Handles test data cleanup
4. Supports parallel test execution

## Next Steps

After mastering testing strategies:

1. **Read Part 7**: [Operating Systems](./07-operating-systems.md)
2. **Practice**: Implement comprehensive test suites
3. **Automate**: Set up test automation pipelines
4. **Monitor**: Track test metrics and improve quality

## Key Takeaways

- Testing is integral to DevOps, not a separate phase
- Follow the testing pyramid for optimal test distribution
- Automate tests at all levels for fast feedback
- Use appropriate testing strategies for different scenarios
- Maintain clean, fast, and reliable tests
- Integrate testing into CI/CD pipelines
- Monitor test metrics to improve quality
- Invest in test infrastructure and tooling

---

**Continue to**: [Part 7: Operating Systems](./07-operating-systems.md)