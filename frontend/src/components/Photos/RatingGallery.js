import * as React from 'react';
import PhotoGrid from './PhotoGrid';
import { Typography } from '@mui/material';

const { REACT_APP_API_PREFIX } = process.env || "https://localhost:8080";

const RatingGallery = ({user}) => {

  const populate = () => {
    if (!user) {
      return new Promise((_, reject) => {
        reject("User not set!")
      })
    }
    return fetch(REACT_APP_API_PREFIX + '/api/v1/photos/', {
      method: "GET",
      mode: "cors",
      credentials: "include"
    })
  }

  const gridConfig = {
    card: {
      ratingEnabled: true,
      editingEnabled: false,
      selectionEnabled: false,
      favoriteEnabled: false
    },
    lightbox: {
      editingEnabled: false,
      ratingEnabled: true,
      deletingEnabled: false,
    }
  }

  return (
    <>
      <Typography variant='h4' sx={{ marginBottom: 2, marginTop: 2 }}>Ratings</Typography>
      <PhotoGrid populator={populate} config={gridConfig}/>
    </>
  )
}

export default RatingGallery;
