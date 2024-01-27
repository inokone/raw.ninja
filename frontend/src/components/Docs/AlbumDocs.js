import * as React from 'react';
import CollectionsIcon from '@mui/icons-material/Collections'
import { Link } from 'react-router-dom';

import { Box, Typography, List, ListItem } from '@mui/material';

const AlbumDocs = () => {
    return (
        <Box maxWidth='md' sx={{ bgcolor: 'white', borderRadius: '8px', mx: 'auto', marginTop: 2, paddingTop: 2, paddingBottom: 4 }}>
            <Typography variant='h4' sx={{ marginBottom: 6, marginTop: 2 }}>Organize your photos with Albums</Typography>
            <Typography variant='h5' align='left' sx={{ margin: 2 }}>Introduction</Typography>
            <Typography align='left' sx={{ margin: 2, marginBottom: 6 }}>RAW.Ninja empowers you to organize and manage your ever-growing collection of photos with ease. Introducing albums and lifecycle rules, powerful tools for creating customized photo organization and ensuring efficient storage utilization.</Typography>
            <Typography variant='h5' align='left' sx={{ margin: 2 }}>Albums: Categorize and Navigate Your Photos</Typography>
            <Typography align='left' sx={{ margin: 2, marginBottom: 6 }}>Create albums to group your photos based on events, places, or any other meaningful criteria. Easily add or remove photos from albums, allowing for flexible organization and effortless navigation.</Typography>
            <Typography variant='h5' align='left' sx={{ margin: 2 }}>Lifecycle Rules: Optimize Storage with Automation</Typography>
            <Typography align='left' sx={{ margin: 2, marginBottom: 6 }}>Set up lifecycle rules to automate the management of your photos based on specific conditions. For instance, you can automatically move older or less frequently accessed photos to a lower-cost storage tier, or even delete them to free up space.</Typography>
            <Typography variant='h5' align='left' sx={{ margin: 2 }}>Creating Albums</Typography>
            <List sx={{ listStyleType: 'numeric', pl: 6 }}>
                <ListItem sx={{ display: 'list-item' }}><Typography align='left'>Navigate to the <Link to="/photos">Photos</Link> page.</Typography></ListItem>
                <ListItem sx={{ display: 'list-item' }}><Typography align='left'>Select the photos you want to include in the album.</Typography></ListItem>
                <ListItem sx={{ display: 'list-item' }}><Typography align='left'>Click the <CollectionsIcon fontSize='small'/> button on the side bar.</Typography></ListItem>
                <ListItem sx={{ display: 'list-item' }}><Typography align='left'>Give your album a name.</Typography></ListItem>
                <ListItem sx={{ display: 'list-item' }}><Typography align='left'>Optionally, add tags to categorize your album.</Typography></ListItem>
            </List>
            <Typography variant='h5' align='left' sx={{ margin: 2 }}>Benefits of Albums and Lifecycle Rules</Typography>
            <List sx={{ listStyleType: 'disc', pl: 6 }}>
                <ListItem sx={{ display: 'list-item' }}><Typography align='left'><b>Simplified Photo Organization:</b> Easily categorize and navigate your photos using intuitive albums.</Typography></ListItem>
                <ListItem sx={{ display: 'list-item' }}><Typography align='left'></Typography><b>Efficient Storage Management:</b> Automate storage optimization with lifecycle rules, saving you time and effort.</ListItem>
                <ListItem sx={{ display: 'list-item' }}><Typography align='left'></Typography><b>Enhanced Photo Accessibility:</b> Quickly find the photos you need by browsing your organized albums.</ListItem>
                <ListItem sx={{ display: 'list-item' }}><Typography align='left'></Typography><b>Tailored Storage Management:</b> Adapt your storage strategy to your specific needs using flexible lifecycle rules.</ListItem>
            </List>
        </Box>
    )
}

export default AlbumDocs