import { useState } from "react";
import * as React from 'react';
import PropTypes from "prop-types";
import { useNavigate } from 'react-router-dom';
import { Box } from "@mui/material";
import PhotoAlbum from "react-photo-album";

import PhotoCard from "./PhotoCard";
import Lightbox from "./Lightbox";
import RatingLightbox from "./RatingLightbox";

const { REACT_APP_API_PREFIX } = process.env || "https://localhost:8080";

const PhotoGallery = ({ photos, updatePhoto, setSelected, config }) => {
    const navigate = useNavigate();
    const [index, setIndex] = useState(-1);
    
    const handleFullscreenClick = React.useCallback((photo) => {
        setIndex(photos.indexOf(photo));
    }, [photos, setIndex])

    const handleRatingChange = React.useCallback((photo, rating) => {
        photo.rating = rating
        updatePhoto(photo)
    }, [updatePhoto])

    const handleEditClick = React.useCallback((photo) => {
        navigate('/editor/' + photo.id + '?format=' + photo.format, {
            state: {
                photo_id: photo.id,
                photo_format: photo.format,
                photo_name: photo.title
            }
        })
    }, [navigate])

    const handleDeleteClick = React.useCallback((photo) => {
        fetch(REACT_APP_API_PREFIX + '/api/v1/photos/' + photo.id, {
            method: "DELETE",
            mode: "cors",
            credentials: "include"
        })
            .then(response => {
                if (!response.ok) {
                    return new Promise((resolve, reject) => {
                        reject(response.status + ":" + response.statusText)
                    })
                }
                navigate(0)
            });
    }, [navigate])

    return (
        <Box sx={{ width: { lg: photos.length < 4 ? photos.length / 4 : '100%'}}}>
            {photos &&
            <>
                <PhotoAlbum
                    photos={photos}
                    layout="rows"
                    targetRowHeight={200}
                    spacing={3}
                    renderPhoto={({ photo, imageProps }) => (
                        <PhotoCard 
                            config={config && config.card} 
                            photo={photo} 
                            updatePhoto={updatePhoto} 
                            setSelected={setSelected} 
                            imageProps={imageProps} 
                            onClick={() => handleFullscreenClick(photo)} />
                    )}
                />
                {config && config.lightbox && config.lightbox.ratingEnabled ?
                    <RatingLightbox
                        photos={photos}
                        index={index}
                        setIndex={setIndex}
                        onRatingChange={handleRatingChange} /> :
                    <Lightbox
                        photos={photos}
                        config={config && config.lightbox}
                        index={index}
                        setIndex={setIndex}
                        onDeleteClick={handleDeleteClick}
                        onEditClick={handleEditClick} />     
                }
            </>}
        </Box>
    );
}

PhotoGallery.propTypes = {
    photos: PropTypes.array,
    updatePhoto: PropTypes.func.isRequired,
    setSelected: PropTypes.func.isRequired,
    config: PropTypes.object
};

export default PhotoGallery;