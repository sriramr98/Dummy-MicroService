# TODO Application

A simple TODO application built with HTML, CSS, and JavaScript that allows users to register, login, and manage their tasks.

## Features

- User Authentication (Login/Register)
- Create, Read, Update, and Delete tasks
- Mark tasks as completed
- Responsive design

## Setup Instructions

1. Clone this repository to your local machine.
2. Configure the API endpoints in `config.js` (see API Configuration section below).
3. Open the `index.html` file in your web browser.

## API Configuration

The application supports different configurations for development and production environments:

### Development Environment (Separate Hosts)

In development, you might have authentication and task services running on different hosts. Configure `config.js` as follows:

```javascript
// Configuration settings for the TODO application
const CONFIG = {
    // Set to 'development' for separate hosts
    ENV: 'development',
    
    // API Base URL - Used as fallback
    API_BASE_URL: 'https://api.example.com',
    
    // Development-specific URLs
    AUTH_BASE_URL: 'http://localhost:8080',
    TASK_BASE_URL: 'http://localhost:8081',
};
```

### Production Environment (Single Host)

In production, you'll typically have a single API endpoint. Configure `config.js` as follows:

```javascript
// Configuration settings for the TODO application
const CONFIG = {
    // Set to 'production' for single host
    ENV: 'production',
    
    // API Base URL - Used for all endpoints in production
    API_BASE_URL: 'https://api.example.com',
    
    // These will be ignored in production mode
    AUTH_BASE_URL: 'http://localhost:8080',
    TASK_BASE_URL: 'http://localhost:8081',
};
```

The application automatically uses the appropriate URLs based on the `ENV` setting.

### Using Environment Variables

If you need to use environment variables with this application, you have a few options:

1. **Development Server**: If you're using a development server like Node.js with Express, you can dynamically generate the `config.js` file with environment variables before serving it.

2. **Build Process**: If you integrate a build process (like webpack), you can use environment variables during the build.

3. **Server-side Configuration**: Have your web server (like Nginx or Apache) dynamically generate the config file when serving the application.

Example for a simple Node.js server that generates config.js:

```javascript
const express = require('express');
const app = express();
const port = process.env.PORT || 3000;

// Serve static files
app.use(express.static('public'));

// Generate config.js with environment variables
app.get('/config.js', (req, res) => {
  res.set('Content-Type', 'application/javascript');
  res.send(`
    // Configuration settings for the TODO application
    const CONFIG = {
      // Environment
      ENV: '${process.env.NODE_ENV || 'development'}',
      
      // API Base URL
      API_BASE_URL: '${process.env.API_BASE_URL || 'https://api.example.com'}',
      
      // Service-specific URLs
      AUTH_BASE_URL: '${process.env.AUTH_BASE_URL || 'http://localhost:8080'}',
      TASK_BASE_URL: '${process.env.TASK_BASE_URL || 'http://localhost:8081'}',
    };
    
    // If in production mode, use a single API_BASE_URL for all endpoints
    if (CONFIG.ENV === 'production') {
      CONFIG.AUTH_BASE_URL = CONFIG.API_BASE_URL;
      CONFIG.TASK_BASE_URL = CONFIG.API_BASE_URL;
    }
  `);
});

app.listen(port, () => {
  console.log(`Server running at http://localhost:${port}`);
});
```

## API Endpoints

The application uses the following API endpoints:

- **Authentication** (served from AUTH_BASE_URL)
  - `POST /login` - User login
  - `POST /register` - User registration

- **Task Management** (served from TASK_BASE_URL)
  - `GET /task` - Get all tasks
  - `GET /task/{taskId}` - Get a specific task
  - `POST /task` - Create a new task
  - `PATCH /task/{taskId}` - Update a task
  - `DELETE /task/{taskId}` - Delete a task

## Usage

1. **Registration**
   - Click on the "Register" tab
   - Enter your email and password
   - Click the "Register" button

2. **Login**
   - Enter your email and password
   - Click the "Login" button

3. **Adding a Task**
   - Fill in the title and description in the "Add New Task" form
   - Click the "Add Task" button

4. **Editing a Task**
   - Click the "Edit" button on a task
   - Update the task details in the modal
   - Check/uncheck the "Completed" checkbox to mark the task as completed/incomplete
   - Click the "Update Task" button

5. **Deleting a Task**
   - Click the "Delete" button on a task
   - Confirm the deletion when prompted

6. **Logout**
   - Click the "Logout" button in the top-right corner

## Browser Support

This application works on all modern browsers that support ES6 features, including:
- Chrome
- Firefox
- Safari
- Edge 