import React, { useState, useRef } from "react";
import { TextField, Button, Alert, Box, Container, Typography } from "@mui/material";
import { Link, useNavigate } from "react-router-dom"
import ReCAPTCHA from "react-google-recaptcha"

const { REACT_APP_API_PREFIX } = process.env || "https://localhost:8080";

const Login = ({ setUser }) => {
    const navigate = useNavigate();
    const [email, setEmail] = useState("")
    const [password, setPassword] = useState("")
    const [emailError, setEmailError] = useState(false)
    const [passwordError, setPasswordError] = useState(false)
    const [error, setError] = useState()
    const [success, setSuccess] = useState(false)
    const captchaRef = useRef(null)
    const [loading, setLoading] = useState(false)


    const handleSubmit = (event) => {
        event.preventDefault()
        const token = captchaRef.current.getValue();
        if (!token) {
            setError("You have to solve the captcha")
            return
        }

        captchaRef.current.reset();

        setError(null)
        setEmailError(email === '')
        setPasswordError(password === '')

        if (email && password) {
            setLoading(true)
            fetch(REACT_APP_API_PREFIX + '/api/v1/auth/login', {
                method: "POST",
                mode: "cors",
                credentials: "include",
                headers: {
                    "Content-Type": "application/json",
                },
                body: JSON.stringify({
                    "email": email,
                    "password": password,
                    "captcha_token": token,
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
                        updateLoggedinUser("/home")
                    }
                })
                .catch(error => {
                    setError(error.message)
                })
                .finally(() => {
                    setLoading(false)
                });
        }
    }

    const updateLoggedinUser = (redirectPath) => {
        fetch(REACT_APP_API_PREFIX + '/api/v1/account/profile', {
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
                <Box style={{ flex: 1 }} sx={{ m: 4 }}>
                    <form autoComplete="off" onSubmit={handleSubmit} sx={{ backgroundColor: "#fff" }}>
                        <TextField
                            label="Email"
                            onChange={e => {
                                setEmail(e.target.value)
                                setError(null)
                            }}
                            required
                            disabled={loading}
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
                            disabled={loading}
                            variant="outlined"
                            color="primary"
                            type="password"
                            value={password}
                            error={passwordError}
                            fullWidth
                            sx={{ mb: 3, backgroundColor: "#fff", borderRadius: 1 }}
                        />
                        <Box sx={{ mb: 4, placeContent: 'center', display: 'flex' }}>
                            <ReCAPTCHA
                                ref={captchaRef}
                                sitekey="6Let2RIpAAAAANGXcsSJ9aOQEaQmwKqsaZB7IAaQ"
                            />
                        </Box>
                        {success && <Alert sx={{ mb: 4 }} severity="success">Logged in successfully!</Alert>}
                        {error && <Alert sx={{ mb: 4 }} severity="error">{error}</Alert>}
                        <Button sx={{ mb: 2 }} disabled={loading} variant="contained" color="primary" type="submit">Login</Button>
                    </form>
                    <Typography sx={{ mb: 2 }}><Link to="/password/recover">Forgot password?</Link> - <Link to="/signup">Sign up</Link></Typography>
                </Box>
            </Container>
        </React.Fragment>
    );
}

export default Login;