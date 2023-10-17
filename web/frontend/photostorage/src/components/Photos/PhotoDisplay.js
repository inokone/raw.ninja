import React from 'react';
import { CircularProgress, Alert } from "@mui/material";
import { useLocation } from "react-router-dom"
import DetailedPhotoCard from './DetailedPhotoCard';

const { REACT_APP_API_PREFIX } = process.env;

const PhotoDisplay = () => {
  const location = useLocation()
  const [error, setError] = React.useState(null)
  const [loading, setLoading] = React.useState(true)
  const [image, setImage] = React.useState(null)
  const path = location.pathname

  React.useEffect(() => {
    const loadImage = () => {
      fetch(REACT_APP_API_PREFIX + '/api/v1' + path, {
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
              setImage(content)
            })
          }
        })
        .catch(error => {
          setError(error)
          setLoading(false)
        });
    }

    loadImage()
  }, [path])

  return (
    <>
      {error !== null ? <Alert sx={{ mb: 4 }} severity="error">{error}</Alert> : null}
      {loading ? <CircularProgress mt={10} /> : <DetailedPhotoCard image={image} closable={false} />}
    </>
  );
}
export default PhotoDisplay; 