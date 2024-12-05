const express = require('express');
const router = express.Router();

// Mock user database
const users = [];

// Registration route
router.post('/', (req, res) => {
    const { username, password } = req.body;

    if (!username || !password) {
        return res.status(400).json({ error: 'Username and password are required.' });
    }

    const userExists = users.find(u => u.username === username);

    if (userExists) {
        return res.status(409).json({ error: 'User already exists.' });
    }

    users.push({ username, password, role: 'user' });

    res.status(201).json({ message: 'User registered successfully.' });
});

module.exports = router;