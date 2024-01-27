import React from "react"; 
import PropTypes from "prop-types";
import { Container, Box, Typography } from "@mui/material";
import { Link } from "react-router-dom";
import LoginForm from "./LoginForm";
import SingleSignOn from "./SingleSignOn";

const Login = ({ setUser }) => {

    return (
        <Box>
            <Container sx={{ 
                width: 356,
                bgcolor: 'white', 
                borderRadius: 2, 
                boxShadow: '0 2px 8px 0 rgba(0, 0, 0, 0.24)',
                py: 4, 
                mt: 6 }}>
                <LoginForm setUser={setUser} />
                <SingleSignOn />
                <Typography><Link to="/password/recover">Forgot password?</Link> - <Link to="/signup">Sign up</Link></Typography>
            </Container>
        </Box>
    );
}

Login.propTypes = {
    setUser: PropTypes.func.isRequired
};

export default Login;