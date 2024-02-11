import * as React from 'react';
import PhotoGrid from '../Photos/PhotoGrid';
import { Typography } from '@mui/material';
import { useLocation } from "react-router-dom";

const { REACT_APP_API_PREFIX } = process.env || "https://localhost:8080";

const AlbumRating = ({user}) => {
  const [title, setTitle] = React.useState(null)
  const location = useLocation()
  const path = location.pathname.replace(/\/ratings$/, '');

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

  const gridConfig = {
    card: {
      ratingEnabled: true,
      editingEnabled: false,
      selectionEnabled: false,
      favoriteEnabled: false
    },
    lightbox: {
      favoriteEnabled: false,
      editingEnabled: false,
      ratingEnabled: true,
    }
  }

  return (
    <>
      <Typography variant='h4' sx={{ marginBottom: 2, marginTop: 2 }}>{title} ratings</Typography>
      <PhotoGrid populator={populate} config={gridConfig} onDataLoaded={handleDataLoaded}/>
    </>
  )
}

export default AlbumRating;
