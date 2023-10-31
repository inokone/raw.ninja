import React, { useState } from 'react';
import { TextField, Button, Stack, Alert, Typography, Container, Box } from '@mui/material';
import { Link } from "react-router-dom"
const { REACT_APP_API_PREFIX } = process.env;


const RegisterForm = () => {
    const [firstName, setFirstName] = useState('')
    const [firstNameError, setFirstNameError] = useState(false)
    const [lastName, setLastName] = useState('')
    const [lastNameError, setLastNameError] = useState(false)
    const [email, setEmail] = useState('')
    const [emailError, setEmailError] = useState(false)
    const [password, setPassword] = useState('')
    const [passwordError, setPasswordError] = useState(false)
    const [error, setError] = useState()
    const [success, setSuccess] = useState(false)

    function handleSubmit(event) {
        event.preventDefault();
        if (emailError || firstNameError || lastNameError || passwordError) {
            return
        }
        setError(null)
        fetch(REACT_APP_API_PREFIX + '/api/v1/account/signup', {
            method: "POST",
            mode: "cors",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify({
                "email": email,
                "password": password,
                "firstname": firstName,
                "lastname": lastName,
            })
        })
            .then(response => {
                if (!response.ok) {
                    response.json().then(content => {
                        setError(content.message)
                    });
                } else {
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
                                onChange={e => {
                                    setFirstNameError(firstName === '')
                                    setFirstName(e.target.value)
                                    setError(null)
                                }}
                                error={firstNameError}
                                value={firstName}
                                fullWidth
                                required
                                sx={{ backgroundColor: "#fff", borderRadius: 1 }}
                            />
                            <TextField
                                type="text"
                                variant='outlined'
                                color='primary'
                                label="Last Name"
                                onChange={e => {
                                    setLastNameError(lastName === '')
                                    setLastName(e.target.value)
                                    setError(null)
                                }}
                                error={lastNameError}
                                value={lastName}
                                fullWidth
                                required
                                sx={{ backgroundColor: "#fff", borderRadius: 1 }}
                            />
                        </Stack>
                        <TextField
                            type="email"
                            variant='outlined'
                            color='primary'
                            label="Email"
                            value={email}
                            onChange={e => {
                                setEmailError(email === '')
                                setEmail(e.target.value)
                                setError(null)
                            }}
                            fullWidth
                            required
                            sx={{ mb: 4, backgroundColor: "#fff", borderRadius: 1 }}
                        />
                        <TextField
                            type="password"
                            variant='outlined'
                            color='primary'
                            label="Password"
                            error={passwordError}
                            onChange={e => {
                                setPassword(e.target.value)
                                setPasswordError(password.length > 0 && password.length < 8)
                                setError(null)
                            }}
                            value={password}
                            required
                            fullWidth
                            sx={{ mb: 4, backgroundColor: "#fff", borderRadius: 1 }}
                        />
                        {success && <Alert sx={{ mb: 4 }} severity="success">Signed up successfully! Please <Link to="/login">log in</Link>!</Alert>}
                        {error && <Alert sx={{ mb: 4 }} severity="error">{error}</Alert>}
                        <Button sx={{ mb: 4 }} variant="contained" color="primary" type="submit">Register</Button>
                    </form>
                    {!success && <Typography>Already have an account? <Link to="/login">Login Here</Link></Typography>}
                </Box>
            </Container>
        </React.Fragment>
    )
}

export default RegisterForm;