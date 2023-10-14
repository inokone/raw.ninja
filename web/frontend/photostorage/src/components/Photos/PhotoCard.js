import * as React from 'react';
import DownloadIcon from '@mui/icons-material/Download';
import FavoriteIcon from '@mui/icons-material/Favorite';
import { useNavigate } from "react-router-dom"
import { Card, CardMedia, Box, CardActions, CardActionArea, Typography, IconButton, Tooltip } from "@mui/material";
import { makeStyles } from '@mui/styles';

const { REACT_APP_API_PREFIX } = process.env;

const useStyles = makeStyles((theme) => ({
  favorite: {
    color: 'red',
  },
  nonfavorite: {
    color: 'gray',
  },
}));

const PhotoCard = (props) => {
  const navigate = useNavigate()
  const classes = useStyles();

  const handleFavoriteClick = (image) => {
    const updatedImage = { ...image };
    updatedImage.descriptor.favorite = !image.descriptor.favorite
    props.setImage(updatedImage)
  }

  const handleDownloadClick = (image) => {
    fetch(REACT_APP_API_PREFIX + '/api/v1/photos/' + image.id + '/download', {
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
          image.descriptor.filename,
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
    <Card sx={{ maxWidth: 250 }}>
      <Box sx={{ position: 'relative' }}>
        <CardMedia
          component="img"
          height="200px"
          image={props.image.descriptor.thumbnail}
          loading="lazy"
          alt={props.image.descriptor.filename}
          onClick={() => handleClick(props.image.id)}
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
          <Typography variant="body1">{props.image.descriptor.filename}</Typography>
          <Typography variant="body2">{new Date(props.image.descriptor.uploaded).toLocaleDateString()}</Typography>
        </Box>
      </Box>
      <CardActionArea>
        <CardActions disableSpacing>
          <Tooltip title="Prevent file from being deleted by life-cycle rules">
            <IconButton aria-label="Add to favorites" onClick={() => handleFavoriteClick(props.image)}>
              <FavoriteIcon className={props.image.descriptor.favorite ? classes.favorite : classes.nonfavorite} />
            </IconButton>
          </Tooltip>
          <Tooltip title="Download RAW file">
            <IconButton aria-label="Download RAW image" onClick={() => handleDownloadClick(props.image)}>
              <DownloadIcon />
            </IconButton>
          </Tooltip>
        </CardActions>
      </CardActionArea>
    </Card>
  );
}
export default PhotoCard;