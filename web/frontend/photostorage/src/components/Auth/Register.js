import React, { useState } from 'react';
import { TextField, Button, Stack, Alert, Typography, Container, Box } from '@mui/material';
import { Link } from "react-router-dom"
const { REACT_APP_API_PREFIX } = process.env;


const RegisterForm = () => {
    const [firstName, setFirstName] = useState('')
    const [lastName, setLastName] = useState('')
    const [email, setEmail] = useState('')
    const [password, setPassword] = useState('')
    const [error, setError] = useState()
    const [success, setSuccess] = useState(false)

    function handleSubmit(event) {
        event.preventDefault();
        fetch(REACT_APP_API_PREFIX + '/api/v1/auth/signup', {
            method: "POST",
            mode: "cors",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify({
                "email": email,
                "password": password,
                "phone": "+123456789"
            })
        })
            .then(response => {
                if (!response.ok) {
                    response.json().then(content => {
                        setError(content.message)
                    })
                } else {
                    setError(null)
                    setSuccess(true)
                }
            })
            .catch(error => console.error(error));
    }

    return (
        <React.Fragment>
            <Container maxWidth="sm">
                <Box sx={{ width: 500, m: 4 }}>
                    <Typography pb={3} variant='h4'>Registration</Typography>
                    <form onSubmit={handleSubmit} action={<Link to="/login" />}>
                        <Stack spacing={2} direction="row" sx={{ marginBottom: 4 }}>
                            <TextField
                                type="text"
                                variant='outlined'
                                color='primary'
                                label="First Name"
                                onChange={e => setFirstName(e.target.value)}
                                value={firstName}
                                fullWidth
                                required
                                sx={{ backgroundColor: "#fff" }}
                            />
                            <TextField
                                type="text"
                                variant='outlined'
                                color='primary'
                                label="Last Name"
                                onChange={e => setLastName(e.target.value)}
                                value={lastName}
                                fullWidth
                                required
                                sx={{ backgroundColor: "#fff" }}
                            />
                        </Stack>
                        <TextField
                            type="email"
                            variant='outlined'
                            color='primary'
                            label="Email"
                            error={error}
                            onChange={e => setEmail(e.target.value)}
                            value={email}
                            fullWidth
                            required
                            sx={{ mb: 4, backgroundColor: "#fff" }}
                        />
                        <TextField
                            type="password"
                            variant='outlined'
                            color='primary'
                            label="Password"
                            onChange={e => setPassword(e.target.value)}
                            value={password}
                            required
                            fullWidth
                            sx={{ mb: 4, backgroundColor: "#fff" }}
                        />
                        {success ? <Alert sx={{ mb: 4 }} severity="success">Signed up successfully! Please <Link to="/login">log in</Link>!</Alert> : null}
                        {error ? <Alert sx={{ mb: 4 }} severity="error">{error}</Alert> : null}
                        <Button sx={{ mb: 4 }} variant="contained" color="primary" type="submit">Register</Button>
                    </form>
                    {!success ? <Typography>Already have an account? <Link to="/login">Login Here</Link></Typography> : null}
                </Box>
            </Container>
        </React.Fragment>
    )
}

export default RegisterForm;