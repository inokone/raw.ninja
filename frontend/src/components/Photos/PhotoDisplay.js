import React from 'react';
import { Alert } from "@mui/material";
import { useLocation } from "react-router-dom"
import DetailedPhotoCard from './DetailedPhotoCard';
import ProgressDisplay from '../Common/ProgressDisplay';

const { REACT_APP_API_PREFIX } = process.env || "https://localhost:8080";

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
            response.json().then(content => {
              setError(content.message)
            });
          } else {
            response.json().then(content => {
              setImage(content)
              setLoading(false)
            })
          }
        })
        .catch(error => {
          setError(error.message)
          setLoading(false)
        });
    }

    if(!image && !loading && !error) {
      loadImage()
    }
  }, [path, image, loading, error])

  return (
    <React.Fragment>
      {error && <Alert sx={{ mb: 4 }} onClose={() => setError(null)} severity="error">{error}</Alert>}
      {loading && <ProgressDisplay /> }
      {image && <DetailedPhotoCard image={image} closable={false} />}
    </React.Fragment>
  );
}
export default PhotoDisplay; 