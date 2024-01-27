import React, { useState, useRef } from 'react';
import { TextField, Button, Alert, Typography, Container, Box, Checkbox, FormControlLabel } from '@mui/material';
import { Link, useNavigate } from "react-router-dom"
import ReCAPTCHA from "react-google-recaptcha"
import SingleSignOn from '../Auth/SingleSignOn';

const { REACT_APP_API_PREFIX } = process.env || "https://localhost:8080";


const SignupForm = () => {
    const navigate = useNavigate()
    const [email, setEmail] = useState('')
    const [emailError, setEmailError] = useState(false)
    const [password, setPassword] = useState('')
    const [passwordError, setPasswordError] = useState(false)
    const [confirmation, setConfirmation] = useState("")
    const [confirmationError, setConfirmationError] = useState(false)
    const [error, setError] = useState()
    const [success, setSuccess] = useState(false)
    const [accepted, setAccepted] = useState(false)
    const captchaRef = useRef(null)
    const [loading, setLoading] = useState(false)

    const handleSubmit = (event) => {
        event.preventDefault();
        const token = captchaRef.current.getValue();
        if (!token) {
            setError("You have to solve the captcha")
            return
        }

        captchaRef.current.reset();
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
                "captcha_token": token,
            })
        })
            .then(response => {
                if (!response.ok) {
                    response.json().then(content => {
                        setError(content.message)
                    });
                } else {
                    setSuccess(true)
                    navigate("/login?source=successfulRegistration")
                }
            })
            .catch(() => setError("Network communication error. Maybe backend is down?"))
            .finally(() => {
                setLoading(false)
            });
    }

    return (
        <React.Fragment>
            <Container sx={{
                width: 356,
                bgcolor: 'white',
                borderRadius: 2,
                boxShadow: '0 2px 8px 0 rgba(0, 0, 0, 0.24)',
                py: 4,
                mt: 6
            }}>
                <Box style={{ flex: 1 }}>
                    <Typography pb={3} variant='h5'>Sign up</Typography>
                    <form onSubmit={handleSubmit} action={<Link to="/login" />}>
                        <TextField
                            type="email"
                            name="email"
                            autoComplete="username"
                            variant='outlined'
                            color='primary'
                            label="Email"
                            disabled={loading}
                            value={email}
                            onChange={e => {
                                setEmailError(email === '')
                                setEmail(e.target.value)
                                setError(null)
                            }}
                            fullWidth
                            required
                            sx={{ mb: 2, backgroundColor: "#fff", borderRadius: 1 }}
                        />
                        <TextField
                            type="password"
                            name="password"
                            autoComplete="password"
                            variant='outlined'
                            color='primary'
                            label="Password"
                            error={passwordError}
                            disabled={loading}
                            onChange={e => {
                                setPassword(e.target.value)
                                setPasswordError(password.length > 0 && password.length < 8)
                                setError(null)
                            }}
                            value={password}
                            required
                            fullWidth
                            sx={{ mb: 2, backgroundColor: "#fff", borderRadius: 1 }}
                        />
                        <TextField
                            label="Confirm Password"
                            name="password"
                            autoComplete="password"
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
                            disabled={loading}
                            value={confirmation}
                            error={confirmationError}
                            fullWidth
                            sx={{ mb: 2, backgroundColor: "#fff", borderRadius: 1 }}
                            helperText={confirmationError && "New password and confirmation must match."}
                        />
                        <FormControlLabel
                            sx={{ mb: 2 }}
                            control={<Checkbox onChange={(event) => setAccepted(event.target.checked)} />}
                            label={<Typography align='left' fontSize={14}>I have read and accept the <Link to="/terms">terms of service</Link> and <Link to="/privacy">privacy policy</Link>.</Typography>}
                            disabled={loading}
                        />

                        <Box sx={{ mb: 2, placeContent: 'center', display: 'flex' }}>
                            <ReCAPTCHA
                                ref={captchaRef}
                                sitekey="6Let2RIpAAAAANGXcsSJ9aOQEaQmwKqsaZB7IAaQ"
                            />
                        </Box>
                        {success && <Alert sx={{ mb: 4 }} onClose={() => setSuccess(null)} severity="success">Signed up successfully! Navigating to login...</Alert>}
                        {error && <Alert sx={{ mb: 4 }} onClose={() => setError(null)} severity="error">{error}</Alert>}
                        <Button sx={{ mb: 4, width: '100%' }} variant="contained" color="primary" type="submit" disabled={!accepted || loading}>Sign up</Button>
                    </form>
                    <SingleSignOn />
                    {!success && <Typography>Already have an account? <Link to="/login">Login Here</Link></Typography>}
                </Box>
            </Container>
        </React.Fragment>
    )
}

export default SignupForm;