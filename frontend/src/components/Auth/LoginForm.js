import React, { useState, useRef } from "react";
import PropTypes from "prop-types";
import { TextField, Button, Alert, Box } from "@mui/material";
import { useNavigate, useSearchParams } from "react-router-dom"
import ReCAPTCHA from "react-google-recaptcha"

const { REACT_APP_API_PREFIX } = process.env || "https://localhost:8080";

const LoginForm = ({ setUser }) => {
    const navigate = useNavigate();
    const [email, setEmail] = useState("")
    const [password, setPassword] = useState("")
    const [emailError, setEmailError] = useState(false)
    const [passwordError, setPasswordError] = useState(false)
    const [error, setError] = useState()
    const [success, setSuccess] = useState()
    const captchaRef = useRef(null)
    const [loading, setLoading] = useState(false)
    const [queryParameters] = useSearchParams()

    const handleSubmit = (event) => {
        event.preventDefault()
        const token = captchaRef.current.getValue();
        if (!token) {
            setError("You have to solve the captcha")
            return
        }
        captchaRef.current.reset();
        setError(null)
        setSuccess(null)
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
                        setSuccess("Logged in successfully!")
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

    React.useEffect(() => {
        let source = queryParameters.get("source")
        if (source === "successfulRegistration") {
            setSuccess("Successfully signed up, please log in!")
        }
    }, [queryParameters, success])

    return (
                <Box style={{ flex: 1}}>
                    <form autoComplete="off" onSubmit={handleSubmit}>
                        <TextField
                            label="Email"
                            name="email"
                            onChange={e => {
                                setEmail(e.target.value)
                                setError(null)
                                setSuccess(null)
                            }}
                            required
                            disabled={loading}
                            autoComplete="username"
                            variant="outlined"
                            color="primary"
                            type="email"
                            sx={{ mb: 2, borderRadius: 1 }}
                            fullWidth
                            value={email}
                            error={emailError}
                        />
                        <TextField
                            label="Password"
                            name="password"
                            onChange={e => {
                                setPassword(e.target.value)
                                setError(null)
                                setSuccess(null)
                            }}
                            required
                            disabled={loading}
                            autoComplete="password"
                            variant="outlined"
                            color="primary"
                            type="password"
                            value={password}
                            error={passwordError}
                            fullWidth
                            sx={{ mb: 2, borderRadius: 1 }}
                        />
                        <Box sx={{ mb: 2, placeContent: 'center', display: 'flex' }}>
                            <ReCAPTCHA
                                ref={captchaRef}
                                sitekey="6Let2RIpAAAAANGXcsSJ9aOQEaQmwKqsaZB7IAaQ"
                            />
                        </Box>
                        {success && <Alert sx={{ mb: 4 }} onClose={() => setSuccess(null)} severity="success">{success}</Alert>}
                        {error && <Alert sx={{ mb: 4 }} onClose={() => setError(null)} severity="error">{error}</Alert>}
                        <Button sx={{ mb: 4, width: '100%', fontSize: '16px' }} disabled={loading} variant="contained" color="primary" type="submit">Login</Button>
                    </form>
                </Box>
    );
}

LoginForm.propTypes = {
    setUser: PropTypes.func.isRequired
};

export default LoginForm;