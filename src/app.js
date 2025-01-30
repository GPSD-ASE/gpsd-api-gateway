const express = require('express');
const fs = require('fs');
const https = require('https');
const app = express();
const PORT = process.env.GPSD_API_GATEWAY_PORT;

const options = {
    key: fs.readFileSync('/etc/ssl/certs/api.gpsd.gateway.com.key'),
    cert: fs.readFileSync('/etc/ssl/certs/api.gpsd.gateway.com.crt')
};

global.tokenBlacklist = new Set();
global.refreshTokens = new Map();

app.use(express.json());
app.use((req, res, next) => {
    res.setHeader('Content-Type', 'application/json');
    next();
});

app.get('/api/', (req, res) => {
    res.status(200).json({ message: 'API Gateway is running over HTTPS.' });
});

// Routes
app.use('/api/user/login', require('./routes/login'));
app.use('/api/user/register', require('./routes/register'));
app.use('/api/user/refresh', require('./routes/refresh'));
app.use('/api/user/logout', require('./routes/logout'));

// Start HTTPS Server
https.createServer(options, app).listen(PORT, '0.0.0.0', () => {
    console.log(`Server is running on https://0.0.0.0:${PORT}`);
});