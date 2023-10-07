import React, {useState} from "react";
import { TextField, Button, Alert } from "@mui/material";
import { Link, useNavigate } from "react-router-dom"
const { REACT_APP_API_PREFIX } = process.env;
 
const Login = () => {
    const navigate = useNavigate();
    const [email, setEmail] = useState("")
    const [password, setPassword] = useState("")
    const [emailError, setEmailError] = useState(false)
    const [passwordError, setPasswordError] = useState(false)
    const [error,setError]=useState()
    const [success,setSuccess]=useState(false)

    const setLoggedinUser = (response) => {
        console.log(response)
        navigate("/")
    }

    const getLoggedinUser = (response) => {
        fetch(REACT_APP_API_PREFIX + '/api/v1/auth/profile', {
            method: "GET",
            mode: "cors",
            credentials: "include",
            headers: {
                "Content-Type": "application/json"
            }
        })
        .then(response => {
            if (!response.ok) {

                response.json().then(content => setError(content["message"]))
            } else {
                response.json().then(content => setLoggedinUser(content))
            }
        })
        .catch(error => setError(error));
    }
 
    const handleSubmit = (event) => {
        event.preventDefault()
 
        setEmailError(email === '')
        setPasswordError(password === '')
 
        if (email && password) {
            fetch(REACT_APP_API_PREFIX + '/api/v1/auth/login', {
                method: "POST",
                mode: "cors",
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
                    response.json().then(content => setError(content["message"]))
                } else {
                    response.json().then(content => {
                        setError(null)
                        setSuccess(true)
                        getLoggedinUser(content)
                    })
                }
            })
            .catch(error => setError(error));
        }
    }
     
    return ( 
        <React.Fragment>
        <form autoComplete="off" onSubmit={handleSubmit} sx={{ backgroundColor: "#fff" }}>
            <TextField 
                label="Email"
                onChange={e => setEmail(e.target.value)}
                required
                variant="outlined"
                color="primary"
                type="email"
                sx={{mb: 3, backgroundColor: "#fff"}}
                fullWidth
                value={email}
                error={emailError}
            />
            <></>
            <TextField 
                label="Password"
                onChange={e => setPassword(e.target.value)}
                required
                variant="outlined"
                color="primary"
                type="password"
                value={password}
                error={passwordError}
                fullWidth
                sx={{mb: 3, backgroundColor: "#fff"}}
            />
            {success?<Alert sx={{mb: 4}} severity="success">Logged in successfully!</Alert>:null}
            {error?<Alert sx={{mb: 4}} severity="error">{error}</Alert>:null}
            <Button sx={{mb: 4}} variant="contained" color="primary" type="submit">Login</Button>
        </form>
        <small>Need an account? <Link to="/register">Register here</Link></small>
        </React.Fragment>
     );
}
 
export default Login;