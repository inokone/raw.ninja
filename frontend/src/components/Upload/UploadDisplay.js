import * as React from 'react';
import PhotoGrid from '../Photos/PhotoGrid';
import { useLocation } from "react-router-dom";
import { Typography } from '@mui/material';


const { REACT_APP_API_PREFIX } = process.env || "https://localhost:8080";

const UploadDisplay = ({user}) => {
    const [title, setTitle] = React.useState(null)
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

    const handleDataLoaded = (data) => {
        setTitle(data.name)
    }

    return (
        <>
            {title && <Typography sx={{ marginBottom: 4, marginTop: 2 }} variant='h4'>Upload {title}</Typography>}
            <PhotoGrid populator={populate} data={[]} onDataLoaded={handleDataLoaded}></PhotoGrid>
        </>
    )
}

export default UploadDisplay;