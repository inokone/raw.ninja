import * as React from 'react';
import { useLocation } from "react-router-dom";
import { Stack, Chip, Typography } from "@mui/material";
import PhotoGrid from '../Photos/PhotoGrid';

const { REACT_APP_API_PREFIX } = process.env || "https://localhost:8080";

const AlbumDisplay = ({ user }) => {
    const location = useLocation()
    const path = location.pathname
    const [title, setTitle] = React.useState(null)
    const [data, setData] = React.useState(null)

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
        setData(data)
    }

    return (
        <>
            {title && <Typography sx={{ marginBottom: 2, marginTop: 2 }} variant='h4'>{title}</Typography>}
            {data && data.tags && <Stack sx={{ marginBottom: 4 }} justifyContent={'center'} direction="row" spacing={1}>
                {data.tags.map(tag => {
                    return (<Chip label={tag} />)
                })} 
            </Stack>}
            <PhotoGrid populator={populate} data={[]} onDataLoaded={handleDataLoaded} />
        </>
    )
}

export default AlbumDisplay;