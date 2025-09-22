# Part 5: CRUD Operations in Node.js

This part covers implementing Create, Read, Update, and Delete operations using Redis with Node.js, including practical patterns and real-world examples.

## Basic CRUD Operations

### String-based CRUD

```typescript
// src/crud/string-crud.ts
import { RedisClientType } from 'redis';

class StringCRUD {
  private client: RedisClientType;

  constructor(client: RedisClientType) {
    this.client = client;
  }

  // Create/Update
  async set(key: string, value: string, ttl?: number): Promise<string | null> {
    try {
      if (ttl) {
        return await this.client.setEx(key, ttl, value);
      }
      return await this.client.set(key, value);
    } catch (error) {
      throw new Error(`Failed to set key ${key}: ${error.message}`);
    }
  }

  // Read
  async get(key: string): Promise<string | null> {
    try {
      return await this.client.get(key);
    } catch (error) {
      throw new Error(`Failed to get key ${key}: ${error.message}`);
    }
  }

  // Read multiple
  async getMultiple(keys: string[]): Promise<(string | null)[]> {
    try {
      return await this.client.mGet(keys);
    } catch (error) {
      throw new Error(`Failed to get multiple keys: ${error.message}`);
    }
  }

  // Update with conditions
  async setIfExists(key: string, value: string): Promise<string | null> {
    try {
      return await this.client.set(key, value, { XX: true });
    } catch (error) {
      throw new Error(`Failed to update key ${key}: ${error.message}`);
    }
  }

  async setIfNotExists(key: string, value: string): Promise<string | null> {
    try {
      return await this.client.set(key, value, { NX: true });
    } catch (error) {
      throw new Error(`Failed to create key ${key}: ${error.message}`);
    }
  }

  // Delete
  async delete(key: string): Promise<number> {
    try {
      return await this.client.del(key);
    } catch (error) {
      throw new Error(`Failed to delete key ${key}: ${error.message}`);
    }
  }

  // Delete multiple
  async deleteMultiple(keys: string[]): Promise<number> {
    try {
      return await this.client.del(keys);
    } catch (error) {
      throw new Error(`Failed to delete multiple keys: ${error.message}`);
    }
  }

  // Check existence
  async exists(key: string): Promise<boolean> {
    try {
      const result = await this.client.exists(key);
      return result === 1;
    } catch (error) {
      throw new Error(`Failed to check existence of key ${key}: ${error.message}`);
    }
  }

  // Get with TTL info
  async getWithTTL(key: string): Promise<{ value: string | null; ttl: number }> {
    try {
      const [value, ttl] = await Promise.all([
        this.client.get(key),
        this.client.ttl(key)
      ]);
      return { value, ttl };
    } catch (error) {
      throw new Error(`Failed to get key with TTL ${key}: ${error.message}`);
    }
  }
}

export default StringCRUD;
```

### Hash-based CRUD

```typescript
// src/crud/hash-crud.ts
import { RedisClientType } from 'redis';

interface HashData {
  [field: string]: string | number;
}

class HashCRUD {
  private client: RedisClientType;

  constructor(client: RedisClientType) {
    this.client = client;
  }

  // Create/Update entire hash
  async setHash(key: string, data: HashData): Promise<number> {
    try {
      return await this.client.hSet(key, data);
    } catch (error) {
      throw new Error(`Failed to set hash ${key}: ${error.message}`);
    }
  }

  // Create/Update single field
  async setField(key: string, field: string, value: string | number): Promise<number> {
    try {
      return await this.client.hSet(key, field, value);
    } catch (error) {
      throw new Error(`Failed to set field ${field} in hash ${key}: ${error.message}`);
    }
  }

  // Create field only if it doesn't exist
  async setFieldIfNotExists(key: string, field: string, value: string | number): Promise<number> {
    try {
      return await this.client.hSetNX(key, field, value);
    } catch (error) {
      throw new Error(`Failed to set field ${field} if not exists in hash ${key}: ${error.message}`);
    }
  }

  // Read entire hash
  async getHash(key: string): Promise<HashData> {
    try {
      return await this.client.hGetAll(key);
    } catch (error) {
      throw new Error(`Failed to get hash ${key}: ${error.message}`);
    }
  }

  // Read single field
  async getField(key: string, field: string): Promise<string | undefined> {
    try {
      return await this.client.hGet(key, field);
    } catch (error) {
      throw new Error(`Failed to get field ${field} from hash ${key}: ${error.message}`);
    }
  }

  // Read multiple fields
  async getFields(key: string, fields: string[]): Promise<(string | undefined)[]> {
    try {
      return await this.client.hmGet(key, fields);
    } catch (error) {
      throw new Error(`Failed to get fields from hash ${key}: ${error.message}`);
    }
  }

  // Get all field names
  async getFieldNames(key: string): Promise<string[]> {
    try {
      return await this.client.hKeys(key);
    } catch (error) {
      throw new Error(`Failed to get field names from hash ${key}: ${error.message}`);
    }
  }

  // Get all values
  async getValues(key: string): Promise<string[]> {
    try {
      return await this.client.hVals(key);
    } catch (error) {
      throw new Error(`Failed to get values from hash ${key}: ${error.message}`);
    }
  }

  // Update field with increment
  async incrementField(key: string, field: string, increment: number = 1): Promise<number> {
    try {
      return await this.client.hIncrBy(key, field, increment);
    } catch (error) {
      throw new Error(`Failed to increment field ${field} in hash ${key}: ${error.message}`);
    }
  }

  // Update field with float increment
  async incrementFieldFloat(key: string, field: string, increment: number): Promise<number> {
    try {
      return await this.client.hIncrByFloat(key, field, increment);
    } catch (error) {
      throw new Error(`Failed to increment float field ${field} in hash ${key}: ${error.message}`);
    }
  }

  // Delete field
  async deleteField(key: string, field: string): Promise<number> {
    try {
      return await this.client.hDel(key, field);
    } catch (error) {
      throw new Error(`Failed to delete field ${field} from hash ${key}: ${error.message}`);
    }
  }

  // Delete multiple fields
  async deleteFields(key: string, fields: string[]): Promise<number> {
    try {
      return await this.client.hDel(key, fields);
    } catch (error) {
      throw new Error(`Failed to delete fields from hash ${key}: ${error.message}`);
    }
  }

  // Delete entire hash
  async deleteHash(key: string): Promise<number> {
    try {
      return await this.client.del(key);
    } catch (error) {
      throw new Error(`Failed to delete hash ${key}: ${error.message}`);
    }
  }

  // Check if field exists
  async fieldExists(key: string, field: string): Promise<boolean> {
    try {
      const result = await this.client.hExists(key, field);
      return result === 1;
    } catch (error) {
      throw new Error(`Failed to check field existence ${field} in hash ${key}: ${error.message}`);
    }
  }

  // Get hash length
  async getLength(key: string): Promise<number> {
    try {
      return await this.client.hLen(key);
    } catch (error) {
      throw new Error(`Failed to get length of hash ${key}: ${error.message}`);
    }
  }

  // Scan hash fields
  async scanFields(key: string, pattern: string = '*', count: number = 10): Promise<{ cursor: number; fields: string[] }> {
    try {
      const result = await this.client.hScan(key, 0, {
        MATCH: pattern,
        COUNT: count
      });
      return {
        cursor: result.cursor,
        fields: result.tuples.map(tuple => tuple.field)
      };
    } catch (error) {
      throw new Error(`Failed to scan fields in hash ${key}: ${error.message}`);
    }
  }
}

export default HashCRUD;
```

### List-based CRUD

```typescript
// src/crud/list-crud.ts
import { RedisClientType } from 'redis';

class ListCRUD {
  private client: RedisClientType;

  constructor(client: RedisClientType) {
    this.client = client;
  }

  // Create/Add to list (left/beginning)
  async pushLeft(key: string, ...values: string[]): Promise<number> {
    try {
      return await this.client.lPush(key, values);
    } catch (error) {
      throw new Error(`Failed to push left to list ${key}: ${error.message}`);
    }
  }

  // Create/Add to list (right/end)
  async pushRight(key: string, ...values: string[]): Promise<number> {
    try {
      return await this.client.rPush(key, values);
    } catch (error) {
      throw new Error(`Failed to push right to list ${key}: ${error.message}`);
    }
  }

  // Read entire list
  async getAll(key: string): Promise<string[]> {
    try {
      return await this.client.lRange(key, 0, -1);
    } catch (error) {
      throw new Error(`Failed to get all from list ${key}: ${error.message}`);
    }
  }

  // Read range from list
  async getRange(key: string, start: number, end: number): Promise<string[]> {
    try {
      return await this.client.lRange(key, start, end);
    } catch (error) {
      throw new Error(`Failed to get range from list ${key}: ${error.message}`);
    }
  }

  // Read element at index
  async getByIndex(key: string, index: number): Promise<string | null> {
    try {
      return await this.client.lIndex(key, index);
    } catch (error) {
      throw new Error(`Failed to get index ${index} from list ${key}: ${error.message}`);
    }
  }

  // Update element at index
  async setByIndex(key: string, index: number, value: string): Promise<string> {
    try {
      return await this.client.lSet(key, index, value);
    } catch (error) {
      throw new Error(`Failed to set index ${index} in list ${key}: ${error.message}`);
    }
  }

  // Insert before/after element
  async insertBefore(key: string, pivot: string, value: string): Promise<number> {
    try {
      return await this.client.lInsert(key, 'BEFORE', pivot, value);
    } catch (error) {
      throw new Error(`Failed to insert before in list ${key}: ${error.message}`);
    }
  }

  async insertAfter(key: string, pivot: string, value: string): Promise<number> {
    try {
      return await this.client.lInsert(key, 'AFTER', pivot, value);
    } catch (error) {
      throw new Error(`Failed to insert after in list ${key}: ${error.message}`);
    }
  }

  // Remove and return from left
  async popLeft(key: string): Promise<string | null> {
    try {
      return await this.client.lPop(key);
    } catch (error) {
      throw new Error(`Failed to pop left from list ${key}: ${error.message}`);
    }
  }

  // Remove and return from right
  async popRight(key: string): Promise<string | null> {
    try {
      return await this.client.rPop(key);
    } catch (error) {
      throw new Error(`Failed to pop right from list ${key}: ${error.message}`);
    }
  }

  // Blocking pop operations
  async blockingPopLeft(key: string, timeout: number): Promise<{ key: string; element: string } | null> {
    try {
      return await this.client.blPop({ key, timeout });
    } catch (error) {
      throw new Error(`Failed to blocking pop left from list ${key}: ${error.message}`);
    }
  }

  async blockingPopRight(key: string, timeout: number): Promise<{ key: string; element: string } | null> {
    try {
      return await this.client.brPop({ key, timeout });
    } catch (error) {
      throw new Error(`Failed to blocking pop right from list ${key}: ${error.message}`);
    }
  }

  // Remove elements by value
  async removeElements(key: string, count: number, value: string): Promise<number> {
    try {
      return await this.client.lRem(key, count, value);
    } catch (error) {
      throw new Error(`Failed to remove elements from list ${key}: ${error.message}`);
    }
  }

  // Trim list to range
  async trim(key: string, start: number, end: number): Promise<string> {
    try {
      return await this.client.lTrim(key, start, end);
    } catch (error) {
      throw new Error(`Failed to trim list ${key}: ${error.message}`);
    }
  }

  // Get list length
  async getLength(key: string): Promise<number> {
    try {
      return await this.client.lLen(key);
    } catch (error) {
      throw new Error(`Failed to get length of list ${key}: ${error.message}`);
    }
  }

  // Delete entire list
  async deleteList(key: string): Promise<number> {
    try {
      return await this.client.del(key);
    } catch (error) {
      throw new Error(`Failed to delete list ${key}: ${error.message}`);
    }
  }

  // Move element between lists
  async moveElement(source: string, destination: string, from: 'LEFT' | 'RIGHT', to: 'LEFT' | 'RIGHT'): Promise<string | null> {
    try {
      return await this.client.lMove(source, destination, from, to);
    } catch (error) {
      throw new Error(`Failed to move element from ${source} to ${destination}: ${error.message}`);
    }
  }
}

export default ListCRUD;
```

### Set-based CRUD

```typescript
// src/crud/set-crud.ts
import { RedisClientType } from 'redis';

class SetCRUD {
  private client: RedisClientType;

  constructor(client: RedisClientType) {
    this.client = client;
  }

  // Create/Add members to set
  async addMembers(key: string, ...members: string[]): Promise<number> {
    try {
      return await this.client.sAdd(key, members);
    } catch (error) {
      throw new Error(`Failed to add members to set ${key}: ${error.message}`);
    }
  }

  // Read all members
  async getAllMembers(key: string): Promise<string[]> {
    try {
      return await this.client.sMembers(key);
    } catch (error) {
      throw new Error(`Failed to get all members from set ${key}: ${error.message}`);
    }
  }

  // Check if member exists
  async isMember(key: string, member: string): Promise<boolean> {
    try {
      const result = await this.client.sIsMember(key, member);
      return result === 1;
    } catch (error) {
      throw new Error(`Failed to check membership in set ${key}: ${error.message}`);
    }
  }

  // Check multiple members
  async areMembers(key: string, members: string[]): Promise<boolean[]> {
    try {
      const results = await this.client.smIsMember(key, members);
      return results.map(result => result === 1);
    } catch (error) {
      throw new Error(`Failed to check multiple memberships in set ${key}: ${error.message}`);
    }
  }

  // Get random member(s)
  async getRandomMember(key: string): Promise<string | null> {
    try {
      return await this.client.sRandMember(key);
    } catch (error) {
      throw new Error(`Failed to get random member from set ${key}: ${error.message}`);
    }
  }

  async getRandomMembers(key: string, count: number): Promise<string[]> {
    try {
      return await this.client.sRandMemberCount(key, count);
    } catch (error) {
      throw new Error(`Failed to get random members from set ${key}: ${error.message}`);
    }
  }

  // Remove and return random member
  async popRandomMember(key: string): Promise<string | null> {
    try {
      return await this.client.sPop(key);
    } catch (error) {
      throw new Error(`Failed to pop random member from set ${key}: ${error.message}`);
    }
  }

  async popRandomMembers(key: string, count: number): Promise<string[]> {
    try {
      return await this.client.sPopCount(key, count);
    } catch (error) {
      throw new Error(`Failed to pop random members from set ${key}: ${error.message}`);
    }
  }

  // Remove specific member(s)
  async removeMembers(key: string, ...members: string[]): Promise<number> {
    try {
      return await this.client.sRem(key, members);
    } catch (error) {
      throw new Error(`Failed to remove members from set ${key}: ${error.message}`);
    }
  }

  // Get set size
  async getSize(key: string): Promise<number> {
    try {
      return await this.client.sCard(key);
    } catch (error) {
      throw new Error(`Failed to get size of set ${key}: ${error.message}`);
    }
  }

  // Set operations
  async union(keys: string[]): Promise<string[]> {
    try {
      return await this.client.sUnion(keys);
    } catch (error) {
      throw new Error(`Failed to get union of sets: ${error.message}`);
    }
  }

  async intersection(keys: string[]): Promise<string[]> {
    try {
      return await this.client.sInter(keys);
    } catch (error) {
      throw new Error(`Failed to get intersection of sets: ${error.message}`);
    }
  }

  async difference(keys: string[]): Promise<string[]> {
    try {
      return await this.client.sDiff(keys);
    } catch (error) {
      throw new Error(`Failed to get difference of sets: ${error.message}`);
    }
  }

  // Store set operations results
  async storeUnion(destination: string, keys: string[]): Promise<number> {
    try {
      return await this.client.sUnionStore(destination, keys);
    } catch (error) {
      throw new Error(`Failed to store union result: ${error.message}`);
    }
  }

  async storeIntersection(destination: string, keys: string[]): Promise<number> {
    try {
      return await this.client.sInterStore(destination, keys);
    } catch (error) {
      throw new Error(`Failed to store intersection result: ${error.message}`);
    }
  }

  async storeDifference(destination: string, keys: string[]): Promise<number> {
    try {
      return await this.client.sDiffStore(destination, keys);
    } catch (error) {
      throw new Error(`Failed to store difference result: ${error.message}`);
    }
  }

  // Move member between sets
  async moveMember(source: string, destination: string, member: string): Promise<number> {
    try {
      return await this.client.sMove(source, destination, member);
    } catch (error) {
      throw new Error(`Failed to move member from ${source} to ${destination}: ${error.message}`);
    }
  }

  // Delete entire set
  async deleteSet(key: string): Promise<number> {
    try {
      return await this.client.del(key);
    } catch (error) {
      throw new Error(`Failed to delete set ${key}: ${error.message}`);
    }
  }

  // Scan set members
  async scanMembers(key: string, pattern: string = '*', count: number = 10): Promise<{ cursor: number; members: string[] }> {
    try {
      const result = await this.client.sScan(key, 0, {
        MATCH: pattern,
        COUNT: count
      });
      return {
        cursor: result.cursor,
        members: result.members
      };
    } catch (error) {
      throw new Error(`Failed to scan members in set ${key}: ${error.message}`);
    }
  }
}

export default SetCRUD;
```

## User Management Example

### User Model with Redis

```typescript
// src/models/user.ts
import { RedisClientType } from 'redis';
import HashCRUD from '../crud/hash-crud';
import SetCRUD from '../crud/set-crud';

interface User {
  id: string;
  username: string;
  email: string;
  firstName: string;
  lastName: string;
  createdAt: string;
  lastLogin?: string;
  isActive: boolean;
}

interface UserStats {
  loginCount: number;
  postsCount: number;
  followersCount: number;
  followingCount: number;
}

class UserModel {
  private hashCrud: HashCRUD;
  private setCrud: SetCRUD;
  private client: RedisClientType;

  constructor(client: RedisClientType) {
    this.client = client;
    this.hashCrud = new HashCRUD(client);
    this.setCrud = new SetCRUD(client);
  }

  // Create user
  async createUser(user: Omit<User, 'createdAt'>): Promise<User> {
    const userData: User = {
      ...user,
      createdAt: new Date().toISOString()
    };

    try {
      // Store user data
      await this.hashCrud.setHash(`user:${user.id}`, {
        username: userData.username,
        email: userData.email,
        firstName: userData.firstName,
        lastName: userData.lastName,
        createdAt: userData.createdAt,
        isActive: userData.isActive.toString()
      });

      // Add to username index
      await this.client.set(`username:${userData.username}`, user.id);

      // Add to email index
      await this.client.set(`email:${userData.email}`, user.id);

      // Add to active users set if active
      if (userData.isActive) {
        await this.setCrud.addMembers('users:active', user.id);
      }

      // Initialize user stats
      await this.hashCrud.setHash(`user:${user.id}:stats`, {
        loginCount: 0,
        postsCount: 0,
        followersCount: 0,
        followingCount: 0
      });

      return userData;
    } catch (error) {
      throw new Error(`Failed to create user: ${error.message}`);
    }
  }

  // Read user by ID
  async getUserById(id: string): Promise<User | null> {
    try {
      const userData = await this.hashCrud.getHash(`user:${id}`);
      
      if (Object.keys(userData).length === 0) {
        return null;
      }

      return {
        id,
        username: userData.username,
        email: userData.email,
        firstName: userData.firstName,
        lastName: userData.lastName,
        createdAt: userData.createdAt,
        lastLogin: userData.lastLogin,
        isActive: userData.isActive === 'true'
      };
    } catch (error) {
      throw new Error(`Failed to get user by ID: ${error.message}`);
    }
  }

  // Read user by username
  async getUserByUsername(username: string): Promise<User | null> {
    try {
      const userId = await this.client.get(`username:${username}`);
      if (!userId) {
        return null;
      }
      return await this.getUserById(userId);
    } catch (error) {
      throw new Error(`Failed to get user by username: ${error.message}`);
    }
  }

  // Read user by email
  async getUserByEmail(email: string): Promise<User | null> {
    try {
      const userId = await this.client.get(`email:${email}`);
      if (!userId) {
        return null;
      }
      return await this.getUserById(userId);
    } catch (error) {
      throw new Error(`Failed to get user by email: ${error.message}`);
    }
  }

  // Update user
  async updateUser(id: string, updates: Partial<Omit<User, 'id' | 'createdAt'>>): Promise<User | null> {
    try {
      const existingUser = await this.getUserById(id);
      if (!existingUser) {
        return null;
      }

      // Handle username change
      if (updates.username && updates.username !== existingUser.username) {
        await this.client.del(`username:${existingUser.username}`);
        await this.client.set(`username:${updates.username}`, id);
      }

      // Handle email change
      if (updates.email && updates.email !== existingUser.email) {
        await this.client.del(`email:${existingUser.email}`);
        await this.client.set(`email:${updates.email}`, id);
      }

      // Handle active status change
      if (updates.isActive !== undefined && updates.isActive !== existingUser.isActive) {
        if (updates.isActive) {
          await this.setCrud.addMembers('users:active', id);
        } else {
          await this.setCrud.removeMembers('users:active', id);
        }
      }

      // Update user data
      const updateData: any = {};
      Object.keys(updates).forEach(key => {
        if (updates[key] !== undefined) {
          updateData[key] = typeof updates[key] === 'boolean' 
            ? updates[key].toString() 
            : updates[key];
        }
      });

      await this.hashCrud.setHash(`user:${id}`, updateData);

      return await this.getUserById(id);
    } catch (error) {
      throw new Error(`Failed to update user: ${error.message}`);
    }
  }

  // Delete user
  async deleteUser(id: string): Promise<boolean> {
    try {
      const user = await this.getUserById(id);
      if (!user) {
        return false;
      }

      // Remove from indexes
      await this.client.del(`username:${user.username}`);
      await this.client.del(`email:${user.email}`);

      // Remove from active users set
      await this.setCrud.removeMembers('users:active', id);

      // Delete user data and stats
      await Promise.all([
        this.hashCrud.deleteHash(`user:${id}`),
        this.hashCrud.deleteHash(`user:${id}:stats`)
      ]);

      return true;
    } catch (error) {
      throw new Error(`Failed to delete user: ${error.message}`);
    }
  }

  // Get user stats
  async getUserStats(id: string): Promise<UserStats | null> {
    try {
      const stats = await this.hashCrud.getHash(`user:${id}:stats`);
      
      if (Object.keys(stats).length === 0) {
        return null;
      }

      return {
        loginCount: parseInt(stats.loginCount) || 0,
        postsCount: parseInt(stats.postsCount) || 0,
        followersCount: parseInt(stats.followersCount) || 0,
        followingCount: parseInt(stats.followingCount) || 0
      };
    } catch (error) {
      throw new Error(`Failed to get user stats: ${error.message}`);
    }
  }

  // Update login info
  async updateLastLogin(id: string): Promise<void> {
    try {
      const now = new Date().toISOString();
      await this.hashCrud.setField(`user:${id}`, 'lastLogin', now);
      await this.hashCrud.incrementField(`user:${id}:stats`, 'loginCount');
    } catch (error) {
      throw new Error(`Failed to update last login: ${error.message}`);
    }
  }

  // Get active users
  async getActiveUsers(): Promise<string[]> {
    try {
      return await this.setCrud.getAllMembers('users:active');
    } catch (error) {
      throw new Error(`Failed to get active users: ${error.message}`);
    }
  }

  // Search users by pattern
  async searchUsers(pattern: string): Promise<User[]> {
    try {
      const keys = await this.client.keys(`user:*`);
      const userIds = keys
        .filter(key => key.startsWith('user:') && !key.includes(':stats'))
        .map(key => key.split(':')[1]);

      const users: User[] = [];
      for (const id of userIds) {
        const user = await this.getUserById(id);
        if (user && (
          user.username.toLowerCase().includes(pattern.toLowerCase()) ||
          user.email.toLowerCase().includes(pattern.toLowerCase()) ||
          user.firstName.toLowerCase().includes(pattern.toLowerCase()) ||
          user.lastName.toLowerCase().includes(pattern.toLowerCase())
        )) {
          users.push(user);
        }
      }

      return users;
    } catch (error) {
      throw new Error(`Failed to search users: ${error.message}`);
    }
  }
}

export default UserModel;
```

## Usage Examples

### Basic CRUD Example

```typescript
// src/examples/crud-example.ts
import { createClient } from 'redis';
import UserModel from '../models/user';
import StringCRUD from '../crud/string-crud';
import HashCRUD from '../crud/hash-crud';

async function crudExample() {
  const client = createClient({ socket: { host: 'localhost', port: 6379 } });
  await client.connect();

  try {
    // String CRUD operations
    const stringCrud = new StringCRUD(client);
    
    await stringCrud.set('app:name', 'My Redis App');
    const appName = await stringCrud.get('app:name');
    console.log('App name:', appName);

    // Hash CRUD operations
    const hashCrud = new HashCRUD(client);
    
    await hashCrud.setHash('config:app', {
      version: '1.0.0',
      environment: 'development',
      debug: 'true'
    });
    
    const config = await hashCrud.getHash('config:app');
    console.log('App config:', config);

    // User model operations
    const userModel = new UserModel(client);
    
    const newUser = await userModel.createUser({
      id: '1000',
      username: 'johndoe',
      email: 'john@example.com',
      firstName: 'John',
      lastName: 'Doe',
      isActive: true
    });
    
    console.log('Created user:', newUser);
    
    const user = await userModel.getUserByUsername('johndoe');
    console.log('Retrieved user:', user);
    
    await userModel.updateLastLogin('1000');
    
    const stats = await userModel.getUserStats('1000');
    console.log('User stats:', stats);
    
  } catch (error) {
    console.error('Error:', error);
  } finally {
    await client.quit();
  }
}

crudExample();
```

### Express API Example

```typescript
// src/examples/express-api.ts
import express from 'express';
import { createClient } from 'redis';
import UserModel from '../models/user';

const app = express();
app.use(express.json());

// Initialize Redis
const client = createClient({ socket: { host: 'localhost', port: 6379 } });
client.connect();

const userModel = new UserModel(client);

// Create user
app.post('/users', async (req, res) => {
  try {
    const user = await userModel.createUser(req.body);
    res.status(201).json(user);
  } catch (error) {
    res.status(400).json({ error: error.message });
  }
});

// Get user by ID
app.get('/users/:id', async (req, res) => {
  try {
    const user = await userModel.getUserById(req.params.id);
    if (!user) {
      return res.status(404).json({ error: 'User not found' });
    }
    res.json(user);
  } catch (error) {
    res.status(500).json({ error: error.message });
  }
});

// Update user
app.put('/users/:id', async (req, res) => {
  try {
    const user = await userModel.updateUser(req.params.id, req.body);
    if (!user) {
      return res.status(404).json({ error: 'User not found' });
    }
    res.json(user);
  } catch (error) {
    res.status(400).json({ error: error.message });
  }
});

// Delete user
app.delete('/users/:id', async (req, res) => {
  try {
    const deleted = await userModel.deleteUser(req.params.id);
    if (!deleted) {
      return res.status(404).json({ error: 'User not found' });
    }
    res.status(204).send();
  } catch (error) {
    res.status(500).json({ error: error.message });
  }
});

// Get user stats
app.get('/users/:id/stats', async (req, res) => {
  try {
    const stats = await userModel.getUserStats(req.params.id);
    if (!stats) {
      return res.status(404).json({ error: 'User not found' });
    }
    res.json(stats);
  } catch (error) {
    res.status(500).json({ error: error.message });
  }
});

// Search users
app.get('/users/search/:pattern', async (req, res) => {
  try {
    const users = await userModel.searchUsers(req.params.pattern);
    res.json(users);
  } catch (error) {
    res.status(500).json({ error: error.message });
  }
});

const PORT = process.env.PORT || 3000;
app.listen(PORT, () => {
  console.log(`Server running on port ${PORT}`);
});
```

## Next Steps

Now that you understand CRUD operations with Redis in Node.js, you're ready to explore:

- [Part 6: Connection Management and Error Handling](./part-06-connection-management.md)

## Quick Reference

### String Operations
```typescript
await client.set('key', 'value');
await client.get('key');
await client.del('key');
await client.setEx('key', 3600, 'value'); // with TTL
```

### Hash Operations
```typescript
await client.hSet('hash', 'field', 'value');
await client.hGet('hash', 'field');
await client.hGetAll('hash');
await client.hDel('hash', 'field');
```

### List Operations
```typescript
await client.lPush('list', 'value');
await client.rPush('list', 'value');
await client.lRange('list', 0, -1);
await client.lPop('list');
```

### Set Operations
```typescript
await client.sAdd('set', 'member');
await client.sMembers('set');
await client.sIsMember('set', 'member');
await client.sRem('set', 'member');
```

---

**Previous**: [← Part 4: Node.js Redis Client Setup](./part-04-nodejs-setup.md)  
**Next**: [Part 6: Connection Management and Error Handling →](./part-06-connection-management.md)