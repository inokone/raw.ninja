import * as React from 'react';
import { Card, CardMedia, Box, Typography, Tooltip, Grid} from "@mui/material";
import noimage from './no-image.svg';
import CollectionsIcon from '@mui/icons-material/Collections';


const AlbumCard = ({ album, onClick }) => {

    const handleClick = (id) => {
        onClick(id)
    }

    const dateOf = (data) => {
        return new Date(data).toLocaleDateString()
    }

    return (
        <Card style={{ flex: 1 }} sx={{position: 'relative', cursor: "pointer", margin: 1, bgcolor: "lightgrey" }}>
            <Box>
                <Typography 
                    variant='h7'
                    sx={{
                        position: 'absolute',
                        top: 0,
                        left: 0,
                        width: '100%',
                        background:
                            'linear-gradient(to bottom, rgba(0,0,0,0.7) 0%, ' +
                            'rgba(0,0,0,0.3) 70%, rgba(0,0,0,0) 100%)',
                        color: 'white',
                        padding: 1
                    }} 
                >{album.name}</Typography>
                <Box sx={{
                    position: 'absolute',
                    bottom: -5,
                    right: 5,
                    color: 'white',
                    fontWeight: 'fontWeightMedium', 
                    fontSize: '18px',
                    display: 'inline'
                }}>
                    <Grid container>
                    <Typography marginBottom={1} marginRight={'3px'}>{album.photo_count}</Typography>
                    <CollectionsIcon fontSize='small'/>
                    </Grid>
                </Box>
                <Tooltip title={dateOf(album.created_at)}>
                    <CardMedia
                        component="img"
                        height="200px"
                        image={album.thumbnail ? album.thumbnail.url : noimage}
                        loading="lazy"
                        onClick={() => handleClick(album.id)}
                    />
                </Tooltip>
            </Box>
        </Card>
    );
}

export default AlbumCard;