import React, { useState } from "react";
import { TextField, Button, Alert, Box, Container, Typography } from "@mui/material";
import { Link, useNavigate } from "react-router-dom"

const { REACT_APP_API_PREFIX } = process.env;

const Login = ({ setUser }) => {
    const navigate = useNavigate();
    const [email, setEmail] = useState("")
    const [password, setPassword] = useState("")
    const [emailError, setEmailError] = useState(false)
    const [passwordError, setPasswordError] = useState(false)
    const [error, setError] = useState()
    const [success, setSuccess] = useState(false)

    const handleSubmit = (event) => {
        event.preventDefault()
        setError(null)


        setEmailError(email === '')
        setPasswordError(password === '')

        if (email && password) {
            fetch(REACT_APP_API_PREFIX + '/api/v1/auth/login', {
                method: "POST",
                mode: "cors",
                credentials: "include",
                headers: {
                    "Content-Type": "application/json",
                },
                body: JSON.stringify({
                    "email": email,
                    "password": password
                })
            })
                .then(response => {
                    if (!response.ok) {
                        response.json().then(content => {
                            setError(content.message)
                        });
                    } else {
                        setError(null)
                        setSuccess(true)
                        updateLoggedinUser("/")
                    }
                })
                .catch(error => {
                    setError(error.message)
                });
        }
    }

    const updateLoggedinUser = (redirectPath) => {
        fetch(REACT_APP_API_PREFIX + '/api/v1/auth/profile', {
            method: "GET",
            mode: "cors",
            credentials: "include",
            headers: {
                "Content-Type": "application/json",
            }
        })
            .then(response => {
                if (!response.ok) {
                    response.json().then(content => {
                        setError(content.message)
                    });
                } else {
                    response.json().then(content => {
                        setUser(content)
                        navigate(redirectPath)
                    })
                }
            })
            .catch(error => setError(error));
    }

    return (
        <React.Fragment>
            <Container maxWidth="sm">
                <Box sx={{ width: 500, m: 4 }}>
                    <form autoComplete="off" onSubmit={handleSubmit} sx={{ backgroundColor: "#fff" }}>
                        <TextField
                            label="Email"
                            onChange={e => {
                                setEmail(e.target.value)
                                setError(null)
                            }}
                            required
                            variant="outlined"
                            color="primary"
                            type="email"
                            sx={{ mb: 3, backgroundColor: "#fff", borderRadius: 1 }}
                            fullWidth
                            value={email}
                            error={emailError}
                        />
                        <TextField
                            label="Password"
                            onChange={e => {
                                setPassword(e.target.value)
                                setError(null)
                            }}
                            required
                            variant="outlined"
                            color="primary"
                            type="password"
                            value={password}
                            error={passwordError}
                            fullWidth
                            sx={{ mb: 3, backgroundColor: "#fff", borderRadius: 1 }}
                        />
                        {success ? <Alert sx={{ mb: 4 }} severity="success">Logged in successfully!</Alert> : null}
                        {error ? <Alert sx={{ mb: 4 }} severity="error">{error}</Alert> : null}
                        <Button sx={{ mb: 4 }} variant="contained" color="primary" type="submit">Login</Button>
                    </form>
                    <Typography>Need an account? <Link to="/register">Register here</Link></Typography>
                </Box>
            </Container>
        </React.Fragment>
    );
}

export default Login;