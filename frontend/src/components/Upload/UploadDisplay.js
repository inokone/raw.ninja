import * as React from 'react';
import PropTypes from "prop-types";
import PhotoGrid from '../Photos/PhotoGrid';
import { useLocation, useNavigate } from "react-router-dom";
import { Typography, Fab, Box } from '@mui/material';
import StarIcon from "@mui/icons-material/Star";


const { REACT_APP_API_PREFIX } = process.env || "https://localhost:8080";

const UploadDisplay = ({ user }) => {
    const [title, setTitle] = React.useState(null)
    const location = useLocation()
    const navigate = useNavigate()
    const path = location.pathname

    const populate = () => {
        if (!user) {
            return new Promise((_, reject) => {
                reject("User not set!")
            })
        }
        return fetch(REACT_APP_API_PREFIX + '/api/v1' + path, {
            method: "GET",
            mode: "cors",
            credentials: "include"
        })
    }

    const handleDataLoaded = (data) => {
        setTitle(data.name)
    }

    const handleRateClick = () => {
        navigate(path + "/ratings")
    }

    return (
        <>
            {title && <Typography sx={{ marginBottom: 4, marginTop: 2 }} variant='h4'>Upload {title}</Typography>}
            <PhotoGrid populator={populate} onDataLoaded={handleDataLoaded}></PhotoGrid>
            <Box sx={{
                '& > :not(style)': { m: 1 },
                position: "fixed",
                bottom: 16,
                right: 16
            }}>
                <Fab onClick={handleRateClick} color="primary" aria-label="rate">
                    <StarIcon />
                </Fab>
            </Box>
        </>
    )
}

UploadDisplay.propTypes = {
    user: PropTypes.object,
};

export default UploadDisplay;