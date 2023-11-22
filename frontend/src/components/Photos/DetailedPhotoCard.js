import React from 'react';
import Image from 'mui-image'
import { Button, Tooltip, Box, Typography, IconButton, Alert } from "@mui/material";
import Grid from '@mui/material/Unstable_Grid2';
import MetadataDisplay from './MetadataDisplay';
import CloseIcon from '@mui/icons-material/Close';
import { useNavigate } from 'react-router-dom';
import { ReactComponent as PeaLogo } from '../../photopea.svg'


const { REACT_APP_API_PREFIX } = process.env || "https://localhost:8080";

const DetailedPhotoCard = ({ image, setImage, closable, onClose }) => {
  const navigate = useNavigate();
  const [error, setError] = React.useState(null)

  const handleDownloadClick = () => {
    fetch(REACT_APP_API_PREFIX + '/api/v1/photos/' + image.id + '/download', {
      method: "GET",
      mode: "cors",
      credentials: "include"
    })
      .then(response => {
        if (!response.ok) {
          return new Promise((resolve, reject) => {
            reject(response.status + ":" + response.statusText)
          })
        }
        return response.blob()
      })
      .then(blob => {
        var url = window.URL.createObjectURL(blob);
        var a = document.createElement('a');
        a.href = url;
        a.setAttribute(
          'download',
          image.descriptor.filename,
        );
        document.body.appendChild(a);
        a.click();
        a.remove();
      })
      .catch(error => {
        console.log(error)
      });
  }

  const handleDeleteClick = () => {
    fetch(REACT_APP_API_PREFIX + '/api/v1/photos/' + image.id, {
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

  const handleEditClick = () => {
    navigate('/editor/' + image.id + '?format=' + image.descriptor.format)
  }

  const handleFavoriteClick = () => {
    const updatedImage = { ...image };
    updatedImage.descriptor.favorite = !image.descriptor.favorite
    if (setImage) {
      setImage(updatedImage)
    } else {
      uploadImage(updatedImage)
    }
  }

  const uploadImage = (image) => {
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
          response.json().then(content => {
            setError(content.message)
          });
        }
      })
      .catch(error => {
        setError(error.message)
      });
  }

  return (
    <React.Fragment>
      {error && <Alert sx={{ mb: 4 }} severity="error">{error}</Alert>}
      <Box sx={{ bgcolor: 'rgba(0, 0, 0, 0.34)', borderRadius: '4px', my: 1, mr: 1.5 }}>
        <Box sx={{ bgcolor: 'rgba(0, 0, 0, 0.54)', color: 'white', padding: '4px', borderRadius: '4px'}}>
          <Grid container>
            <Grid xs={11}><Typography variant="h5" >{image.descriptor.filename}</Typography></Grid>
            <Grid xs={1} sx={{ position: 'relative' }}>
              {closable ?
                <IconButton onClick={onClose} sx={{ color: 'white', position: 'absolute', right: '0' }}>
                  <CloseIcon />
                </IconButton>
                : null}
            </Grid>
          </Grid>
        </Box>
        <Grid container spacing={2} padding={3}>
          <Grid sm={12} md={9}>
            <Image src={image.descriptor.thumbnail} height="80vh" />
          </Grid>
          <Grid sm={12} md={3}>
            <MetadataDisplay metadata={image.descriptor} />
            <Grid container mt={4} spacing={1} padding={1}>
              <Grid xs={4}>
                <Tooltip title="Download RAW file">
                  <Button variant='contained' color='primary' onClick={handleDownloadClick}>Download</Button>
                </Tooltip>
              </Grid>
              <Grid xs={4}>
                <Tooltip title="Mark file as favorite">
                  <Button variant='contained' color='secondary' onClick={handleFavoriteClick}>{!image.descriptor.favorite ? "Favorite" : "Unfavorite"}</Button>
                </Tooltip>
              </Grid>
              <Grid xs={4}>
                <Tooltip title="Delete RAW file">
                  <Button variant='contained' color='secondary' onClick={handleDeleteClick}>Delete</Button>
                </Tooltip>
              </Grid>
              <Grid xs={4} sx={{display: { xs: 'none', md: 'flex' }}}>
                <Tooltip title="Edit RAW">
                  <Button variant='contained' color='secondary' onClick={handleEditClick}>Edit with <PeaLogo /></Button>
                </Tooltip>
              </Grid>
            </Grid>
          </Grid>
        </Grid>
      </Box>
    </React.Fragment>
  );
}
export default DetailedPhotoCard; 