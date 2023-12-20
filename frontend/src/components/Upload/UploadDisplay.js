import * as React from 'react';
import PhotoGrid from '../Photos/PhotoGrid';
import { useLocation } from "react-router-dom";
import { Typography } from '@mui/material';


const { REACT_APP_API_PREFIX } = process.env || "https://localhost:8080";

const UploadDisplay = ({user}) => {
    const location = useLocation()
    const path = location.pathname

    const populate = () => {
        if (!user) {
            return new Promise((resolve, reject) => {
                reject("User not set!")
            })
        }
        return fetch(REACT_APP_API_PREFIX + '/api/v1' + path, {
                method: "GET",
                mode: "cors",
                credentials: "include"
            }) 
    }

    return (
        <>
            <Typography variant='h4'>Upload</Typography>
            <PhotoGrid populator={populate} data={[]}></PhotoGrid>
        </>
    )
}

export default UploadDisplay;