// Get API Base URL from config
const API_BASE_URL = CONFIG.API_BASE_URL;
const TASK_BASE_URL = CONFIG.TASK_BASE_URL || CONFIG.API_BASE_URL;
const AUTH_BASE_URL = CONFIG.AUTH_BASE_URL || CONFIG.API_BASE_URL;

// DOM Elements
const authContainer = document.getElementById('auth-container');
const todoContainer = document.getElementById('todo-container');
const tabBtns = document.querySelectorAll('.tab-btn');
const formContainers = document.querySelectorAll('.form-container');
const loginForm = document.getElementById('login');
const registerForm = document.getElementById('register');
const loginMessage = document.getElementById('login-message');
const registerMessage = document.getElementById('register-message');
const logoutBtn = document.getElementById('logout-btn');
const addTaskForm = document.getElementById('add-task');
const tasksContainer = document.getElementById('tasks-container');
const editTaskModal = document.getElementById('edit-task-modal');
const editTaskForm = document.getElementById('edit-task-form');
const closeModalBtn = document.querySelector('.close');

// State
let token = localStorage.getItem('token');
let tasks = [];

// Check if user is logged in
function checkAuth() {
    if (token) {
        authContainer.classList.add('hidden');
        todoContainer.classList.remove('hidden');
        fetchTasks();
    } else {
        authContainer.classList.remove('hidden');
        todoContainer.classList.add('hidden');
    }
}

// Tab switching
tabBtns.forEach(btn => {
    btn.addEventListener('click', () => {
        const tabName = btn.getAttribute('data-tab');
        
        // Update active tab button
        tabBtns.forEach(b => b.classList.remove('active'));
        btn.classList.add('active');
        
        // Show corresponding form
        formContainers.forEach(container => {
            container.classList.remove('active');
            if (container.id === `${tabName}-form`) {
                container.classList.add('active');
            }
        });
    });
});

// API Requests
async function makeRequest(url, method, body = null, requiresAuth = false) {
    try {
        const headers = {
            'Content-Type': 'application/json'
        };
        
        if (requiresAuth && token) {
            headers['Authorization'] = `Bearer ${token}`;
        }
        
        const options = {
            method,
            headers
        };
        
        if (body && (method === 'POST' || method === 'PATCH')) {
            options.body = JSON.stringify(body);
        }
        
        // Determine which base URL to use based on the endpoint
        let baseUrl;
        if (url.startsWith('/login') || url.startsWith('/register')) {
            baseUrl = AUTH_BASE_URL;
        } else if (url.startsWith('/task')) {
            baseUrl = TASK_BASE_URL;
        } else {
            baseUrl = API_BASE_URL; // Fallback to the generic API URL
        }
        
        const response = await fetch(`${baseUrl}${url}`, options);
        const data = await response.json();
        
        if (!data.success) {
            throw new Error(data.error || 'Something went wrong');
        }
        
        return data;
    } catch (error) {
        console.error('API Error:', error);
        throw error;
    }
}

// Authentication
async function login(email, password) {
    try {
        const data = await makeRequest('/login', 'POST', { email, password });
        token = data.data.token;
        localStorage.setItem('token', token);
        showMessage(loginMessage, 'Login successful!', 'success');
        checkAuth();
    } catch (error) {
        showMessage(loginMessage, error.message || 'Login failed', 'error');
    }
}

async function register(email, password) {
    try {
        const data = await makeRequest('/register', 'POST', { email, password });
        showMessage(registerMessage, 'Registration successful! You can now login.', 'success');
        
        // Switch to login tab
        tabBtns[0].click();
    } catch (error) {
        showMessage(registerMessage, error.message || 'Registration failed', 'error');
    }
}

function logout() {
    token = null;
    localStorage.removeItem('token');
    checkAuth();
}

// Task Management
async function fetchTasks() {
    try {
        const data = await makeRequest('/task', 'GET', null, true);
        tasks = data.data;
        renderTasks();
    } catch (error) {
        console.error('Error fetching tasks:', error);
    }
}

async function createTask(title, body) {
    try {
        await makeRequest('/task', 'POST', { title, body }, true);
        addTaskForm.reset();
        fetchTasks();
    } catch (error) {
        console.error('Error creating task:', error);
    }
}

async function updateTask(id, title, body, completed) {
    try {
        await makeRequest(`/task/${id}`, 'PATCH', { title, body, completed }, true);
        closeModal();
        fetchTasks();
    } catch (error) {
        console.error('Error updating task:', error);
    }
}

async function deleteTask(id) {
    if (confirm('Are you sure you want to delete this task?')) {
        try {
            await makeRequest(`/task/${id}`, 'DELETE', null, true);
            fetchTasks();
        } catch (error) {
            console.error('Error deleting task:', error);
        }
    }
}

// UI Functions
function renderTasks() {
    tasksContainer.innerHTML = '';
    
    if (tasks.length === 0) {
        tasksContainer.innerHTML = '<p>No tasks found. Add a new task to get started!</p>';
        return;
    }
    
    tasks.forEach(task => {
        const taskElement = document.createElement('div');
        taskElement.className = `task-item ${task.completed ? 'completed' : ''}`;
        
        taskElement.innerHTML = `
            <div class="task-content">
                <div class="task-title">${task.title}</div>
                <div class="task-body">${task.body}</div>
            </div>
            <div class="task-actions">
                <button class="edit-btn" data-id="${task.id}">Edit</button>
                <button class="delete-btn" data-id="${task.id}">Delete</button>
            </div>
        `;
        
        tasksContainer.appendChild(taskElement);
    });
    
    // Add event listeners to buttons
    document.querySelectorAll('.edit-btn').forEach(btn => {
        btn.addEventListener('click', () => openEditModal(btn.getAttribute('data-id')));
    });
    
    document.querySelectorAll('.delete-btn').forEach(btn => {
        btn.addEventListener('click', () => deleteTask(btn.getAttribute('data-id')));
    });
}

function openEditModal(taskId) {
    const task = tasks.find(t => t.id === parseInt(taskId));
    
    if (task) {
        document.getElementById('edit-task-id').value = task.id;
        document.getElementById('edit-task-title').value = task.title;
        document.getElementById('edit-task-body').value = task.body;
        document.getElementById('edit-task-completed').checked = task.completed;
        
        editTaskModal.classList.remove('hidden');
    }
}

function closeModal() {
    editTaskModal.classList.add('hidden');
}

function showMessage(element, message, type) {
    element.textContent = message;
    element.className = `message ${type}`;
    
    // Clear message after 3 seconds
    setTimeout(() => {
        element.textContent = '';
        element.className = 'message';
    }, 3000);
}

// Event Listeners
loginForm.addEventListener('submit', (e) => {
    e.preventDefault();
    const email = document.getElementById('login-email').value;
    const password = document.getElementById('login-password').value;
    login(email, password);
});

registerForm.addEventListener('submit', (e) => {
    e.preventDefault();
    const email = document.getElementById('register-email').value;
    const password = document.getElementById('register-password').value;
    register(email, password);
});

logoutBtn.addEventListener('click', logout);

addTaskForm.addEventListener('submit', (e) => {
    e.preventDefault();
    const title = document.getElementById('task-title').value;
    const body = document.getElementById('task-body').value;
    createTask(title, body);
});

editTaskForm.addEventListener('submit', (e) => {
    e.preventDefault();
    const id = document.getElementById('edit-task-id').value;
    const title = document.getElementById('edit-task-title').value;
    const body = document.getElementById('edit-task-body').value;
    const completed = document.getElementById('edit-task-completed').checked;
    updateTask(id, title, body, completed);
});

closeModalBtn.addEventListener('click', closeModal);

// Close modal when clicking outside
window.addEventListener('click', (e) => {
    if (e.target === editTaskModal) {
        closeModal();
    }
});

// Initialize the app
document.addEventListener('DOMContentLoaded', checkAuth); 