const express = require('express');
const jwt = require('jsonwebtoken');
const router = express.Router();

const { jwtSecret, refreshSecret } = require('../config');
if (!jwtSecret || !refreshSecret) {
    throw new Error('Missing JWT_SECRET or REFRESH_SECRET environment variables');
}

// Middleware to authenticate the token
function authenticateToken(req, res, next) {
    const authHeader = req.headers['authorization'];
    if (!authHeader || !authHeader.startsWith('Bearer ')) {
        return res.status(401).json({ error: 'Authorization header missing or malformed.' });
    }

    const token = authHeader.split(' ')[1];
    jwt.verify(token, jwtSecret, (err, user) => {
        if (err) return res.status(403).json({ error: 'Token invalid or expired.' });

        if (global.tokenBlacklist.has(token)) {
            return res.status(403).json({ error: 'Token has been revoked.' });
        }

        req.user = user;
        next();
    });
}

// Logout route
router.post('/', authenticateToken, (req, res) => {
    const authHeader = req.headers['authorization'];
    const token = authHeader.split(' ')[1];

    global.tokenBlacklist.add(token);
    res.status(200).json({ message: 'Logged out successfully.' });
});

module.exports = router;