import * as React from 'react';
import { CircularProgress, Alert, Card, CardMedia, Box, CardActions, CardActionArea, Typography, IconButton, Grid } from "@mui/material";
import DownloadIcon from '@mui/icons-material/Download';
import FavoriteIcon from '@mui/icons-material/Favorite';
import { useNavigate } from "react-router-dom"


const { REACT_APP_API_PREFIX } = process.env;

const PhotoCard = (props) => {
  const navigate = useNavigate()

  const handleDownloadClick = (id, filename) => {
    fetch(REACT_APP_API_PREFIX + '/api/v1/photos/' + id + '/download', {
        method: "GET",
        mode: "cors",
        credentials: "include"
    })
    .then(response => response.blob())
    .then(blob => {
        var url = window.URL.createObjectURL(blob);
        var a = document.createElement('a');
        a.href = url;
        a.setAttribute(
          'download',
          filename,
        );
        document.body.appendChild(a);
        a.click();    
        a.remove();     
    });
  }

  const handleClick = (id) => {
    navigate("/photos/" + id)
  }

  return (
    <Card sx={{ maxWidth: 200 }}>
      <Box sx={{ position: 'relative' }}>
        <CardMedia
          component="img"
          height="200px"
          image={props.source}
          loading="lazy"
          alt={props.filename}
          onClick={() => handleClick(props.id)} 
        />
        <Box
          sx={{
            position: 'absolute',
            bottom: 0,
            left: 0,
            width: '100%',
            bgcolor: 'rgba(0, 0, 0, 0.54)',
            color: 'white',
            padding: '10px',
          }}
        >
          <Typography variant="body1">{props.filename}</Typography>
          <Typography variant="body2">{new Date(props.date).toLocaleDateString()}</Typography>
        </Box>
      </Box>
      <CardActionArea>
      <CardActions disableSpacing>
        <IconButton aria-label="Add to favorites">
          <FavoriteIcon />
        </IconButton>
        <IconButton aria-label="Download RAW image" onClick={() => handleDownloadClick(props.id, props.filename)}>
          <DownloadIcon />
        </IconButton>
      </CardActions>
    </CardActionArea>
    </Card>
  );
}

export default function PhotoList() {
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

  if (!error && images === null) {
    loadImages()
  }

  return (
    <>
      {error !== null ? <Alert sx={{mb: 4}} severity="error">{error}</Alert>:null}
      {loading ? <CircularProgress /> : 
        <Grid container spacing={1} sx={{ flexGrow: 1 }}>
          {images.map((image) => {
            return (
              <Grid item key={image.id} xs={6} sm={4} md={3} lg={2}>
                <PhotoCard id={image.id} source={image.descriptor.thumbnail} filename={image.descriptor.filename} date={image.descriptor.uploaded}/>
              </Grid>
            );
          })}
        </Grid>}
    </>
  );
}