import * as React from 'react';
import { CircularProgress, Alert, Card, CardMedia, CardContent, CardActions, CardActionArea, Typography, IconButton, Grid } from "@mui/material";
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
      <CardActionArea>
        <CardMedia
          component="img"
          height="200"
          width="200"
          image={props.source}
          loading="lazy"
          alt={props.filename}
          onClick={() => handleClick(props.id)} 
        />
        <CardContent>
          <Typography variant="h6" component="div">
            {props.filename}
          </Typography>
        </CardContent>
        <CardActions disableSpacing>
          <IconButton aria-label="Add to favorites">
            <FavoriteIcon />
          </IconButton>
          <IconButton aria-label="Download RAW image">
            <DownloadIcon onClick={() => handleDownloadClick(props.id, props.filename)}/>
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
      <Grid container spacing={2}>
        {images.map((image) => {
          return (
            <Grid item xs={6} md={4}>
              <PhotoCard id={image.id} source={image.descriptor.thumbnail} filename={image.descriptor.filename}/>
            </Grid>
          );
        })}
      </Grid>}
    </>
  );


}