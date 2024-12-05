const express = require('express');
const jwt = require('jsonwebtoken');
const router = express.Router();

const SECRET_KEY = 'your-secret-key';

// Mock user data
const users = [{ username: 'admin', password: 'password', role: 'admin' }];

// Login route
router.post('/', (req, res) => {
    const { username, password } = req.body;

    const user = users.find(u => u.username === username && u.password === password);

    if (!user) {
        return res.status(401).json({ error: 'Invalid credentials.' });
    }

    const token = jwt.sign({ username: user.username, role: user.role }, SECRET_KEY, {
        expiresIn: '1h',
    });

    res.json({ token });
});

module.exports = router;