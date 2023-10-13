import DownloadIcon from '@mui/icons-material/Download';
import FavoriteIcon from '@mui/icons-material/Favorite';
import { useNavigate } from "react-router-dom"
import {Card, CardMedia, Box, CardActions, CardActionArea, Typography, IconButton,  Tooltip } from "@mui/material";

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
    <Card sx={{ maxWidth: 250 }}>
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
        <Tooltip title="Prevent file from being deleted by life-cycle rules">
          <IconButton aria-label="Add to favorites">
            <FavoriteIcon />
          </IconButton>
        </Tooltip>
        <Tooltip title="Download RAW file">
          <IconButton aria-label="Download RAW image" onClick={() => handleDownloadClick(props.id, props.filename)}>
            <DownloadIcon />
          </IconButton>
        </Tooltip>
      </CardActions>
    </CardActionArea>
    </Card>
  );
}
export default PhotoCard;