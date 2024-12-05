const express = require('express');
const jwt = require('jsonwebtoken');
const router = express.Router();

const { jwtSecret, refreshSecret } = require('../config');
if (!jwtSecret || !refreshSecret) {
    throw new Error('Missing JWT_SECRET or REFRESH_SECRET environment variables');
}

// TODO: use db in production
const users = [{ username: 'admin', password: 'password', role: 'admin' }];

// Login route
router.post('/', (req, res) => {
    const { username, password } = req.body;

    // Check if a token exists in the headers
    const authHeader = req.headers['authorization'];
    if (authHeader && authHeader.startsWith('Bearer ')) {
        const token = authHeader.split(' ')[1];
        try {
            const decoded = jwt.verify(token, jwtSecret);

            if (decoded.username === username) {
                return res.status(200).json({ message: 'User already logged in.', token });
            }
        } catch (err) {
            // Invalid token, proceed to login
        }
    }

    const user = users.find(u => u.username === username && u.password === password);

    if (!user) {
        return res.status(401).json({ error: 'Invalid credentials.' });
    }

    const token = jwt.sign({ username: user.username, role: user.role }, jwtSecret, {
        expiresIn: '15m',
    });

    res.json({ token });
});

module.exports = router;