import http from 'k6/http';
import { sleep, check } from 'k6';

// Test configuration
export const options = {
  stages: [
    { duration: '30s', target: 10 }, // Ramp up to 10 users over 30 seconds (reduced from 20)
    { duration: '1m', target: 10 },  // Stay at 10 users for 1 minute
    { duration: '30s', target: 0 },  // Ramp down to 0 users
  ],
  thresholds: {
    http_req_duration: ['p(95)<500'], // 95% of requests should be below 500ms
    http_req_failed: ['rate<0.01'],   // Less than 1% of requests should fail
  },
};

// Function to generate a unique username
function generateUniqueUsername() {
  return `testuser_${Date.now()}_${Math.floor(Math.random() * 10000)}`;
}

export default function() {
  // Create a unique user for each virtual user
  const testUser = {
    username: generateUniqueUsername(),
    password: 'testpassword123'
  };

  let authToken = null;

  // Test 1: Create a user
  const registerRes = http.post('http://localhost:9000/api/register', JSON.stringify(testUser), {
    headers: { 'Content-Type': 'application/json' },
  });
  
  check(registerRes, {
    'user registration successful': (r) => r.status === 200 || r.status === 201,
  });
  
  // Test 2: Login
  const loginRes = http.post('http://localhost:9000/api/login', JSON.stringify(testUser), {
    headers: { 'Content-Type': 'application/json' },
  });
  
  check(loginRes, {
    'login successful': (r) => r.status === 200,
  });
  
  // Extract the token if login was successful
  if (loginRes.status === 200) {
    try {
      const loginData = JSON.parse(loginRes.body);
      authToken = loginData.token;
    } catch (e) {
      console.error('Failed to parse login response', e);
    }
  }
  
  // Test 3: Get todos (requires authentication)
  if (authToken) {
    const todosRes = http.get('http://localhost:9000/api/todos', {
      headers: {
        'Authorization': `Bearer ${authToken}`,
        'Content-Type': 'application/json',
      },
    });
    
    check(todosRes, {
      'get todos successful': (r) => r.status === 200,
    });
    
    // Test 4: Create a todo - FIXED endpoint from /api/todos to /api/todo
    const newTodo = {
      task: `Test task ${Date.now()}`,
      description: 'Created during performance testing',
      important: false,
      date: new Date().toISOString().split('T')[0], // Add current date
      time: new Date().toTimeString().split(' ')[0] // Add current time
    };
    
    const createTodoRes = http.post('http://localhost:9000/api/todo', JSON.stringify(newTodo), {
      headers: {
        'Authorization': `Bearer ${authToken}`,
        'Content-Type': 'application/json',
      },
    });
    
    check(createTodoRes, {
      'create todo successful': (r) => r.status === 200 || r.status === 201,
    });
  }
  
  sleep(1);
}