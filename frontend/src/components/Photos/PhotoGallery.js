import { useState } from "react";
import * as React from 'react';
import { useNavigate } from 'react-router-dom';
import { Box } from "@mui/material";

import PhotoAlbum from "react-photo-album";

import Lightbox from "yet-another-react-lightbox";
import "yet-another-react-lightbox/styles.css";
import Fullscreen from "yet-another-react-lightbox/plugins/fullscreen";
import Slideshow from "yet-another-react-lightbox/plugins/slideshow";
import Thumbnails from "yet-another-react-lightbox/plugins/thumbnails";
import Captions from "yet-another-react-lightbox/plugins/captions";
import Zoom from "yet-another-react-lightbox/plugins/zoom";
import "yet-another-react-lightbox/plugins/thumbnails.css";
import "yet-another-react-lightbox/plugins/captions.css";
import PhotoCard from "./PhotoCard";
import EditIcon from '@mui/icons-material/Edit';
import DeleteIcon from '@mui/icons-material/Delete';


const { REACT_APP_API_PREFIX } = process.env || "https://localhost:8080";

export default function PhotoGallery({ photos, setPhoto, setSelected }) {
    const navigate = useNavigate();
    const [index, setIndex] = useState(-1);

    const handleFullscreenClick = (photo) => setIndex(photos.indexOf(photo));

    const handleEditClick = (photo) => {
        navigate('/editor/' + photo.id + '?format=' + photo.format)
    }

    const handleDeleteClick = (photo) => {
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
    }

    return (
        <Box sx={{ width: { lg: photos.length < 4 ? photos.length / 4 : '100%'}}}>
            <PhotoAlbum
                photos={photos}
                layout="rows"
                targetRowHeight={200}
                spacing={3}
                renderPhoto={({ photo, imageProps }) => (
                    <PhotoCard photo={photo} setPhoto={setPhoto} setSelected={setSelected} imageProps={imageProps} onClick={() => handleFullscreenClick(photo)} />
                )}
            />

            <Lightbox
                slides={photos}
                open={index >= 0}
                index={index}
                close={() => setIndex(-1)}
                plugins={[Fullscreen, Captions, Slideshow, Thumbnails, Zoom]}
                toolbar={{
                    buttons: [
                        <button key="edit" type="button" className="yarl__button" onClick={() => handleEditClick(photos[index])}>
                            <EditIcon />
                        </button>,
                        <button key="delete" type="button" className="yarl__button" onClick={() => handleDeleteClick(photos[index])}>
                            <DeleteIcon />
                        </button>,
                        "close",
                    ],
                }}
            />
        </Box>
    );
}