# Network Security for Full Stack Development

## Table of Contents
1. [Security Fundamentals](#security-fundamentals)
2. [SSL/TLS and HTTPS](#ssltls-and-https)
3. [Common Web Vulnerabilities](#common-web-vulnerabilities)
4. [Authentication Security](#authentication-security)
5. [API Security Best Practices](#api-security-best-practices)
6. [Network Security Headers](#network-security-headers)
7. [Encryption and Hashing](#encryption-and-hashing)
8. [Security Testing](#security-testing)

## Security Fundamentals

### CIA Triad

#### Confidentiality
- **Definition**: Information is accessible only to authorized users
- **Implementation**: Encryption, access controls, authentication
- **Example**: User passwords, personal data, API keys

#### Integrity
- **Definition**: Information remains accurate and unaltered
- **Implementation**: Hashing, digital signatures, checksums
- **Example**: Database records, file transfers, software updates

#### Availability
- **Definition**: Information and services are accessible when needed
- **Implementation**: Redundancy, load balancing, DDoS protection
- **Example**: Website uptime, API availability, database access

### Security Principles

1. **Least Privilege**: Grant minimum necessary permissions
2. **Defense in Depth**: Multiple layers of security
3. **Fail Securely**: Default to secure state when errors occur
4. **Security by Design**: Build security from the ground up
5. **Zero Trust**: Never trust, always verify

## SSL/TLS and HTTPS

### What is SSL/TLS?

- **SSL (Secure Sockets Layer)**: Deprecated protocol (SSL 3.0 and earlier)
- **TLS (Transport Layer Security)**: Modern encryption protocol
- **Current versions**: TLS 1.2 (minimum), TLS 1.3 (preferred)

### How HTTPS Works

#### TLS Handshake Process
```
Client                                Server
  |                                     |
  |  1. Client Hello                   |
  |------------------------------------>|
  |     - TLS version                   |
  |     - Cipher suites                 |
  |     - Random number                 |
  |                                     |
  |  2. Server Hello                   |
  |<------------------------------------|
  |     - Chosen cipher suite           |
  |     - Server certificate            |
  |     - Random number                 |
  |                                     |
  |  3. Key Exchange                    |
  |<----------------------------------->|
  |     - Verify certificate            |
  |     - Generate session keys         |
  |                                     |
  |  4. Encrypted Communication        |
  |<====================================|
```

### SSL Certificate Types

#### Domain Validation (DV)
- **Validation**: Domain ownership only
- **Use case**: Basic websites, blogs
- **Trust level**: Low
- **Example**: Let's Encrypt certificates

#### Organization Validation (OV)
- **Validation**: Domain + organization identity
- **Use case**: Business websites
- **Trust level**: Medium
- **Example**: Commercial websites

#### Extended Validation (EV)
- **Validation**: Extensive organization verification
- **Use case**: E-commerce, banking
- **Trust level**: High
- **Example**: Shows organization name in browser

### Certificate Implementation

#### Let's Encrypt (Free SSL)
```bash
# Install Certbot
sudo apt-get install certbot python3-certbot-nginx

# Obtain certificate
sudo certbot --nginx -d example.com -d www.example.com

# Auto-renewal
sudo crontab -e
0 12 * * * /usr/bin/certbot renew --quiet
```

#### Nginx SSL Configuration
```nginx
server {
    listen 443 ssl http2;
    server_name example.com www.example.com;
    
    # SSL Certificate
    ssl_certificate /etc/letsencrypt/live/example.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/example.com/privkey.pem;
    
    # SSL Configuration
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers ECDHE-RSA-AES256-GCM-SHA512:DHE-RSA-AES256-GCM-SHA512;
    ssl_prefer_server_ciphers off;
    ssl_session_cache shared:SSL:10m;
    ssl_session_timeout 10m;
    
    # Security Headers
    add_header Strict-Transport-Security "max-age=31536000; includeSubDomains" always;
    add_header X-Frame-Options DENY always;
    add_header X-Content-Type-Options nosniff always;
    add_header X-XSS-Protection "1; mode=block" always;
    
    location / {
        proxy_pass http://localhost:3000;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}

# Redirect HTTP to HTTPS
server {
    listen 80;
    server_name example.com www.example.com;
    return 301 https://$server_name$request_uri;
}
```

## Common Web Vulnerabilities

### 1. SQL Injection

#### Vulnerable Code
```javascript
// DON'T DO THIS
const query = `SELECT * FROM users WHERE email = '${email}' AND password = '${password}'`;
db.query(query, (err, results) => {
    // Handle results
});

// Malicious input: email = "admin@example.com'; DROP TABLE users; --"
```

#### Secure Code
```javascript
// Use parameterized queries
const query = 'SELECT * FROM users WHERE email = ? AND password = ?';
db.query(query, [email, hashedPassword], (err, results) => {
    // Handle results
});

// Or use ORM
const user = await User.findOne({
    where: {
        email: email,
        password: hashedPassword
    }
});
```

### 2. Cross-Site Scripting (XSS)

#### Types of XSS

**Stored XSS**
```javascript
// Vulnerable: Storing user input without sanitization
const comment = req.body.comment; // "<script>alert('XSS')</script>"
db.query('INSERT INTO comments (content) VALUES (?)', [comment]);

// Later displayed without escaping
res.send(`<div>${comment}</div>`); // Executes script
```

**Reflected XSS**
```javascript
// Vulnerable: Reflecting user input
const search = req.query.q;
res.send(`<h1>Search results for: ${search}</h1>`);
// URL: /search?q=<script>alert('XSS')</script>
```

**DOM-based XSS**
```javascript
// Vulnerable client-side code
const urlParams = new URLSearchParams(window.location.search);
const message = urlParams.get('message');
document.getElementById('content').innerHTML = message; // Dangerous
```

#### XSS Prevention
```javascript
// Server-side: Sanitize input
const DOMPurify = require('isomorphic-dompurify');
const clean = DOMPurify.sanitize(userInput);

// Template engines with auto-escaping
// Handlebars
<div>{{userInput}}</div> <!-- Auto-escaped -->
<div>{{{userInput}}}</div> <!-- Raw HTML - dangerous -->

// React (auto-escapes by default)
const UserComment = ({ comment }) => {
    return <div>{comment}</div>; // Safe
    // return <div dangerouslySetInnerHTML={{__html: comment}} />; // Dangerous
};

// Content Security Policy
app.use((req, res, next) => {
    res.setHeader(
        'Content-Security-Policy',
        "default-src 'self'; script-src 'self' 'unsafe-inline'; style-src 'self' 'unsafe-inline'"
    );
    next();
});
```

### 3. Cross-Site Request Forgery (CSRF)

#### Attack Example
```html
<!-- Malicious website -->
<form action="https://bank.com/transfer" method="POST" style="display:none;">
    <input name="to" value="attacker-account" />
    <input name="amount" value="10000" />
</form>
<script>
    document.forms[0].submit();
</script>
```

#### CSRF Prevention
```javascript
// CSRF Token
const csrf = require('csurf');
const csrfProtection = csrf({ cookie: true });

app.use(csrfProtection);

app.get('/form', (req, res) => {
    res.render('form', { csrfToken: req.csrfToken() });
});

// In template
<form method="POST" action="/transfer">
    <input type="hidden" name="_csrf" value="{{csrfToken}}">
    <input name="amount" type="number">
    <button type="submit">Transfer</button>
</form>

// SameSite Cookies
app.use(session({
    cookie: {
        sameSite: 'strict', // or 'lax'
        secure: true, // HTTPS only
        httpOnly: true
    }
}));
```

### 4. Insecure Direct Object References (IDOR)

#### Vulnerable Code
```javascript
// DON'T DO THIS
app.get('/api/users/:id/profile', (req, res) => {
    const userId = req.params.id;
    // No authorization check!
    const user = await User.findById(userId);
    res.json(user);
});
```

#### Secure Code
```javascript
// Proper authorization
app.get('/api/users/:id/profile', authenticateToken, async (req, res) => {
    const requestedUserId = req.params.id;
    const currentUserId = req.user.id;
    
    // Check if user can access this profile
    if (requestedUserId !== currentUserId && !req.user.isAdmin) {
        return res.status(403).json({ error: 'Access denied' });
    }
    
    const user = await User.findById(requestedUserId);
    res.json(user);
});
```

### 5. Security Misconfiguration

#### Common Issues
```javascript
// DON'T DO THIS

// 1. Exposing sensitive information
app.get('/debug', (req, res) => {
    res.json({
        env: process.env, // Exposes all environment variables
        config: config    // Exposes configuration
    });
});

// 2. Default credentials
const dbConfig = {
    username: 'admin',
    password: 'admin123' // Default password
};

// 3. Verbose error messages in production
app.use((err, req, res, next) => {
    res.status(500).json({
        error: err.message,
        stack: err.stack // Exposes internal structure
    });
});
```

#### Secure Configuration
```javascript
// Environment-based configuration
const config = {
    database: {
        host: process.env.DB_HOST,
        username: process.env.DB_USER,
        password: process.env.DB_PASS
    },
    jwt: {
        secret: process.env.JWT_SECRET
    }
};

// Proper error handling
app.use((err, req, res, next) => {
    console.error(err.stack); // Log for debugging
    
    if (process.env.NODE_ENV === 'production') {
        res.status(500).json({ error: 'Internal server error' });
    } else {
        res.status(500).json({
            error: err.message,
            stack: err.stack
        });
    }
});

// Remove sensitive headers
app.disable('x-powered-by');
```

## Authentication Security

### Password Security

#### Password Hashing
```javascript
const bcrypt = require('bcrypt');
const saltRounds = 12;

// Hash password
const hashPassword = async (password) => {
    return await bcrypt.hash(password, saltRounds);
};

// Verify password
const verifyPassword = async (password, hash) => {
    return await bcrypt.compare(password, hash);
};

// Registration
app.post('/register', async (req, res) => {
    const { email, password } = req.body;
    
    // Password strength validation
    if (password.length < 8) {
        return res.status(400).json({ error: 'Password too short' });
    }
    
    const hashedPassword = await hashPassword(password);
    const user = await User.create({ email, password: hashedPassword });
    
    res.status(201).json({ message: 'User created' });
});
```

#### JWT Security
```javascript
const jwt = require('jsonwebtoken');
const crypto = require('crypto');

// Generate secure secret
const jwtSecret = crypto.randomBytes(64).toString('hex');

// Create token with expiration
const createToken = (payload) => {
    return jwt.sign(payload, jwtSecret, {
        expiresIn: '15m', // Short expiration
        issuer: 'your-app',
        audience: 'your-app-users'
    });
};

// Refresh token (longer expiration, stored securely)
const createRefreshToken = (payload) => {
    return jwt.sign(payload, jwtSecret, {
        expiresIn: '7d'
    });
};

// Token verification middleware
const authenticateToken = (req, res, next) => {
    const authHeader = req.headers['authorization'];
    const token = authHeader && authHeader.split(' ')[1];
    
    if (!token) {
        return res.status(401).json({ error: 'Access token required' });
    }
    
    jwt.verify(token, jwtSecret, (err, user) => {
        if (err) {
            return res.status(403).json({ error: 'Invalid token' });
        }
        req.user = user;
        next();
    });
};
```

### Multi-Factor Authentication (MFA)

#### TOTP (Time-based One-Time Password)
```javascript
const speakeasy = require('speakeasy');
const QRCode = require('qrcode');

// Generate secret for user
app.post('/api/mfa/setup', authenticateToken, async (req, res) => {
    const secret = speakeasy.generateSecret({
        name: `YourApp (${req.user.email})`,
        issuer: 'YourApp'
    });
    
    // Store secret temporarily (confirm before saving permanently)
    await User.update(
        { tempMfaSecret: secret.base32 },
        { where: { id: req.user.id } }
    );
    
    // Generate QR code
    const qrCodeUrl = await QRCode.toDataURL(secret.otpauth_url);
    
    res.json({
        secret: secret.base32,
        qrCode: qrCodeUrl
    });
});

// Verify and enable MFA
app.post('/api/mfa/verify', authenticateToken, async (req, res) => {
    const { token } = req.body;
    const user = await User.findById(req.user.id);
    
    const verified = speakeasy.totp.verify({
        secret: user.tempMfaSecret,
        encoding: 'base32',
        token: token,
        window: 2 // Allow 2 time steps
    });
    
    if (verified) {
        await User.update(
            { 
                mfaSecret: user.tempMfaSecret,
                mfaEnabled: true,
                tempMfaSecret: null
            },
            { where: { id: req.user.id } }
        );
        
        res.json({ message: 'MFA enabled successfully' });
    } else {
        res.status(400).json({ error: 'Invalid token' });
    }
});
```

## API Security Best Practices

### Rate Limiting
```javascript
const rateLimit = require('express-rate-limit');

// General rate limiting
const generalLimiter = rateLimit({
    windowMs: 15 * 60 * 1000, // 15 minutes
    max: 100, // Limit each IP to 100 requests per windowMs
    message: 'Too many requests from this IP'
});

// Strict rate limiting for auth endpoints
const authLimiter = rateLimit({
    windowMs: 15 * 60 * 1000,
    max: 5, // Only 5 login attempts per 15 minutes
    skipSuccessfulRequests: true
});

app.use('/api/', generalLimiter);
app.use('/api/auth/', authLimiter);
```

### Input Validation
```javascript
const Joi = require('joi');

// Validation schemas
const userSchema = Joi.object({
    email: Joi.string().email().required(),
    password: Joi.string().min(8).pattern(new RegExp('^(?=.*[a-z])(?=.*[A-Z])(?=.*[0-9])(?=.*[!@#\$%\^&\*])')).required(),
    name: Joi.string().min(2).max(50).required()
});

// Validation middleware
const validateInput = (schema) => {
    return (req, res, next) => {
        const { error } = schema.validate(req.body);
        if (error) {
            return res.status(400).json({
                error: 'Validation failed',
                details: error.details.map(detail => detail.message)
            });
        }
        next();
    };
};

app.post('/api/users', validateInput(userSchema), async (req, res) => {
    // Handle validated input
});
```

### API Key Management
```javascript
// API key generation
const generateApiKey = () => {
    return crypto.randomBytes(32).toString('hex');
};

// API key middleware
const validateApiKey = async (req, res, next) => {
    const apiKey = req.headers['x-api-key'];
    
    if (!apiKey) {
        return res.status(401).json({ error: 'API key required' });
    }
    
    // Hash the API key for comparison
    const hashedKey = crypto.createHash('sha256').update(apiKey).digest('hex');
    
    const keyRecord = await ApiKey.findOne({
        where: { 
            keyHash: hashedKey,
            isActive: true,
            expiresAt: { [Op.gt]: new Date() }
        }
    });
    
    if (!keyRecord) {
        return res.status(401).json({ error: 'Invalid API key' });
    }
    
    // Update last used timestamp
    await keyRecord.update({ lastUsedAt: new Date() });
    
    req.apiKey = keyRecord;
    next();
};
```

## Network Security Headers

### Essential Security Headers
```javascript
const helmet = require('helmet');

// Use Helmet for basic security headers
app.use(helmet());

// Or configure manually
app.use((req, res, next) => {
    // Prevent clickjacking
    res.setHeader('X-Frame-Options', 'DENY');
    
    // Prevent MIME type sniffing
    res.setHeader('X-Content-Type-Options', 'nosniff');
    
    // Enable XSS protection
    res.setHeader('X-XSS-Protection', '1; mode=block');
    
    // Force HTTPS
    res.setHeader('Strict-Transport-Security', 'max-age=31536000; includeSubDomains; preload');
    
    // Content Security Policy
    res.setHeader('Content-Security-Policy', 
        "default-src 'self'; " +
        "script-src 'self' 'unsafe-inline'; " +
        "style-src 'self' 'unsafe-inline'; " +
        "img-src 'self' data: https:; " +
        "font-src 'self'; " +
        "connect-src 'self';"
    );
    
    // Referrer Policy
    res.setHeader('Referrer-Policy', 'strict-origin-when-cross-origin');
    
    // Permissions Policy
    res.setHeader('Permissions-Policy', 
        'camera=(), microphone=(), geolocation=()'
    );
    
    next();
});
```

### Content Security Policy (CSP)
```javascript
// Basic CSP
const csp = {
    directives: {
        defaultSrc: ["'self'"],
        scriptSrc: ["'self'", "'unsafe-inline'", "https://cdn.jsdelivr.net"],
        styleSrc: ["'self'", "'unsafe-inline'", "https://fonts.googleapis.com"],
        fontSrc: ["'self'", "https://fonts.gstatic.com"],
        imgSrc: ["'self'", "data:", "https:"],
        connectSrc: ["'self'", "https://api.example.com"],
        frameSrc: ["'none'"],
        objectSrc: ["'none'"],
        upgradeInsecureRequests: []
    }
};

app.use(helmet.contentSecurityPolicy(csp));
```

## Encryption and Hashing

### Symmetric Encryption
```javascript
const crypto = require('crypto');

class Encryption {
    constructor() {
        this.algorithm = 'aes-256-gcm';
        this.secretKey = crypto.scryptSync(process.env.ENCRYPTION_PASSWORD, 'salt', 32);
    }
    
    encrypt(text) {
        const iv = crypto.randomBytes(16);
        const cipher = crypto.createCipher(this.algorithm, this.secretKey, iv);
        
        let encrypted = cipher.update(text, 'utf8', 'hex');
        encrypted += cipher.final('hex');
        
        const authTag = cipher.getAuthTag();
        
        return {
            encrypted,
            iv: iv.toString('hex'),
            authTag: authTag.toString('hex')
        };
    }
    
    decrypt(encryptedData) {
        const { encrypted, iv, authTag } = encryptedData;
        
        const decipher = crypto.createDecipher(
            this.algorithm, 
            this.secretKey, 
            Buffer.from(iv, 'hex')
        );
        
        decipher.setAuthTag(Buffer.from(authTag, 'hex'));
        
        let decrypted = decipher.update(encrypted, 'hex', 'utf8');
        decrypted += decipher.final('utf8');
        
        return decrypted;
    }
}
```

### Hashing
```javascript
// SHA-256 hashing
const hashData = (data) => {
    return crypto.createHash('sha256').update(data).digest('hex');
};

// HMAC for message authentication
const createHMAC = (data, secret) => {
    return crypto.createHmac('sha256', secret).update(data).digest('hex');
};

// Verify HMAC
const verifyHMAC = (data, secret, expectedHmac) => {
    const computedHmac = createHMAC(data, secret);
    return crypto.timingSafeEqual(
        Buffer.from(computedHmac, 'hex'),
        Buffer.from(expectedHmac, 'hex')
    );
};
```

## Security Testing

### Automated Security Testing
```javascript
// Security testing with Jest
const request = require('supertest');
const app = require('../app');

describe('Security Tests', () => {
    test('should reject requests without authentication', async () => {
        const response = await request(app)
            .get('/api/protected')
            .expect(401);
    });
    
    test('should prevent SQL injection', async () => {
        const maliciousInput = "'; DROP TABLE users; --";
        
        const response = await request(app)
            .post('/api/users')
            .send({ email: maliciousInput })
            .expect(400);
    });
    
    test('should include security headers', async () => {
        const response = await request(app)
            .get('/')
            .expect(200);
            
        expect(response.headers['x-frame-options']).toBe('DENY');
        expect(response.headers['x-content-type-options']).toBe('nosniff');
    });
    
    test('should rate limit requests', async () => {
        // Make multiple requests quickly
        const promises = Array(10).fill().map(() => 
            request(app).post('/api/auth/login')
        );
        
        const responses = await Promise.all(promises);
        const rateLimited = responses.some(res => res.status === 429);
        
        expect(rateLimited).toBe(true);
    });
});
```

### Security Checklist

#### Authentication & Authorization
- [ ] Strong password requirements
- [ ] Password hashing with salt
- [ ] JWT with short expiration
- [ ] Refresh token mechanism
- [ ] Multi-factor authentication
- [ ] Role-based access control
- [ ] Session management

#### Input Validation
- [ ] Server-side validation
- [ ] SQL injection prevention
- [ ] XSS prevention
- [ ] CSRF protection
- [ ] File upload security
- [ ] Input sanitization

#### Network Security
- [ ] HTTPS everywhere
- [ ] Security headers
- [ ] CORS configuration
- [ ] Rate limiting
- [ ] API key management
- [ ] Network monitoring

#### Data Protection
- [ ] Encryption at rest
- [ ] Encryption in transit
- [ ] Secure key management
- [ ] Data backup security
- [ ] PII protection
- [ ] GDPR compliance

## Key Takeaways

1. **Security is not optional**: Build security from the ground up
2. **Defense in depth**: Use multiple layers of protection
3. **Keep it simple**: Complex security is often broken security
4. **Stay updated**: Security threats evolve constantly
5. **Test regularly**: Automated security testing is essential
6. **Monitor continuously**: Detect and respond to threats quickly

## Next Steps

In the next part, we'll cover:
- Performance optimization
- Caching strategies
- CDN implementation
- Load balancing

---

*Continue to [Part 4: Performance and Optimization](04-performance-optimization.md)*