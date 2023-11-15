import React, { useState } from 'react';
import { TextField, Button, Alert, Typography, Container, Box } from '@mui/material';
import { Link } from "react-router-dom"
const { REACT_APP_API_PREFIX } = process.env || "https://localhost:8080";


const RegisterForm = () => {
    const [email, setEmail] = useState('')
    const [emailError, setEmailError] = useState(false)
    const [password, setPassword] = useState('')
    const [passwordError, setPasswordError] = useState(false)
    const [confirmation, setConfirmation] = useState("")
    const [confirmationError, setConfirmationError] = useState(false)
    const [error, setError] = useState()
    const [success, setSuccess] = useState(false)

    function handleSubmit(event) {
        event.preventDefault();
        if (emailError || passwordError || confirmationError) {
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
                <Box style={{ flex: 1 }} sx={{ m: 4 }}>
                    <Typography pb={3} variant='h4'>Registration</Typography>
                    <form onSubmit={handleSubmit} action={<Link to="/login" />}>
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
                        <TextField
                            label="Confirm Password"
                            onChange={e => {
                                let conf = e.target.value
                                setConfirmation(conf)
                                setConfirmationError(password.length > 0 && conf !== password)
                                setError(null)
                            }}
                            required
                            variant="outlined"
                            color="primary"
                            type="password"
                            value={confirmation}
                            error={confirmationError}
                            fullWidth
                            sx={{ mb: 3, backgroundColor: "#fff", borderRadius: 1 }}
                            helperText={confirmationError && "New password and confirmation must match."}
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