import * as React from 'react';
import PhotoGrid from '../Photos/PhotoGrid';
import { Typography } from '@mui/material';

const { REACT_APP_API_PREFIX } = process.env || "https://localhost:8080";

const Favorites = ({ user }) => {
    const [imageCount, setImageCount] = React.useState(null)

    const populate = () => {
        if (!user) {
            return new Promise((resolve, reject) => {
                reject("User not set!")
            })
        }
        return fetch(REACT_APP_API_PREFIX + '/api/v1/search/favorites', {
            method: "GET",
            mode: "cors",
            credentials: "include"
        })
    }

    const handleDataLoaded = (images) => {
        setImageCount(images.length)
    }

    return (
        <React.Fragment>
            {imageCount > 0 && <Typography variant='h4'>Favorite photos</Typography>}
            <PhotoGrid populator={populate} data={[]} onDataLoaded={handleDataLoaded}/>
        </React.Fragment>
    )
}

export default Favorites;