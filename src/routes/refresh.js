const express = require('express');
const jwt = require('jsonwebtoken');
const router = express.Router();

const { jwtSecret, refreshSecret } = require('../config');
if (!jwtSecret || !refreshSecret) {
    throw new Error('Missing JWT_SECRET or REFRESH_SECRET environment variables');
}

// Middleware to verify the refresh token
function verifyRefreshToken(req, res, next) {
    const authHeader = req.headers['authorization'];
    if (!authHeader || !authHeader.startsWith('Bearer ')) {
        return res.status(401).json({ error: 'Authorization header missing or malformed.' });
    }

const currentToken = authHeader.split(' ')[1]; // Extract the token part

    if (!global.refreshTokens.has(currentToken)) {
        console.log("This is the currentToken: ", currentToken)
        return res.status(403).json({ error: 'Invalid current token.' });
    }

    jwt.verify(currentToken, refreshSecret, (err, user) => {
        if (err) return res.status(403).json({ error: 'Refresh token invalid or expired.' });

        req.user = user;
        next();
    });
}

// Refresh route
router.post('/', verifyRefreshToken, (req, res) => {
    const user = req.user;

    // Generate a new access token
    const accessToken = jwt.sign(
        { username: user.username, role: user.role },
        jwtSecret,
        { expiresIn: '15m' }
    );

    res.status(200).json({ accessToken });
});

module.exports = router;