import React from 'react';
import { Alert } from "@mui/material";
import { useLocation } from "react-router-dom"
import DetailedPhotoCard from './DetailedPhotoCard';
import ProgressDisplay from '../Common/ProgressDisplay';

const { REACT_APP_API_PREFIX } = process.env;

const PhotoDisplay = () => {
  const location = useLocation()
  const [error, setError] = React.useState(null)
  const [loading, setLoading] = React.useState(false)
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
            throw new Error(response.status + ": " + response.statusText);
          } else {
            response.json().then(content => {
              setImage(content)
              setLoading(false)
            })
          }
        })
        .catch(error => {
          setError(error)
          setLoading(false)
        });
    }

    if(!image && !loading && !error) {
      loadImage()
    }
  }, [path, image, loading, error])

  return (
    <>
      {error && <Alert sx={{ mb: 4 }} severity="error">{error}</Alert>}
      {loading && <ProgressDisplay /> }
      {image && <DetailedPhotoCard image={image} closable={false} />}
    </>
  );
}
export default PhotoDisplay; 