import * as React from 'react';
import { CircularProgress, Alert, Grid } from "@mui/material";
import PhotoCard from './PhotoCard';


const { REACT_APP_API_PREFIX } = process.env;

const PhotoList = () => {
  const [error, setError] = React.useState(null)
  const [loading, setLoading] = React.useState(true)
  const [images, setImages] = React.useState(null)

  const loadImages = () => {
    fetch(REACT_APP_API_PREFIX + '/api/v1/photos/', {
      method: "GET",
      mode: "cors",
      credentials: "include"
    })
      .then(response => {
        if (!response.ok) {
          if (response.status !== 200) {
            setError(response.status + ": " + response.statusText);
          } else {
            response.json().then(content => setError(content.message))
          }
          setLoading(false)
        } else {
          response.json().then(content => {
            setLoading(false)
            setImages(content)
          })
        }
      })
      .catch(error => {
        setError(error)
        setLoading(false)
      });
  }

  React.useEffect(() => {
    loadImages()
  }, [])

  const setImage = (image) => {
    fetch(REACT_APP_API_PREFIX + '/api/v1/photos/' + image.id, {
      method: "PUT",
      mode: "cors",
      credentials: "include",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(image)
    })
      .then(response => {
        if (!response.ok) {
          if (response.status !== 200) {
            setError(response.status + ": " + response.statusText);
          } else {
            response.json().then(content => setError(content.message))
          }
        } else {
          let newImages = images.slice()
          for (let i = 0; i < images.length; i++) {
            if (images[i].id === image.id) {
              newImages[i] = image
              setImages(newImages)
              return
            }
          }
        }
      })
      .catch(error => {
        setError(error)
      });
  }

  return (
    <>
      {error !== null ? <Alert sx={{ mb: 4 }} severity="error">{error}</Alert> : null}
      {loading ? <CircularProgress /> :
        <Grid container spacing={1} sx={{ flexGrow: 1, pl: 2, pt: 3 }} >
          {images.map((image) => {
            return (
              <Grid item key={image.id} xs={6} sm={4} md={3} lg={2}>
                <PhotoCard image={image} setImage={setImage} />
              </Grid>
            );
          })}
        </Grid>}
    </>
  );
}

export default PhotoList;