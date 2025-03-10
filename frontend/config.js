// Configuration settings for the TODO application
const CONFIG = {
    // Environment (set to 'production' for single host, 'development' for separate hosts)
    ENV: 'development',
    
    // API Base URL - Used as fallback and in production
    API_BASE_URL: 'https://api.example.com',
    
    // Development-specific URLs
    AUTH_BASE_URL: 'http://localhost:8081',
    TASK_BASE_URL: 'http://localhost:8080',
};

// If in production mode, use a single API_BASE_URL for all endpoints
if (CONFIG.ENV === 'production') {
    CONFIG.AUTH_BASE_URL = CONFIG.API_BASE_URL;
    CONFIG.TASK_BASE_URL = CONFIG.API_BASE_URL;
} 