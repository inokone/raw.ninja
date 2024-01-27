import * as React from 'react';

import { Box, Typography } from '@mui/material';

const PhotoDocs = () => {
    return (
        <Box maxWidth='md' sx={{ bgcolor: 'white', borderRadius: '8px', mx: 'auto', marginTop: 2, paddingTop: 2, paddingBottom: 4 }}>
            <Typography variant='h4' sx={{ marginBottom: 6, marginTop: 2 }}>Photos page</Typography>
            <Typography variant='h5' align='left' sx={{ margin: 2 }}>Introduction</Typography>
            <Typography align='left' sx={{ margin: 2, marginBottom: 6 }}>Photos is a central hub for managing your uploaded photos and raw image files. The uploaded photos are displayed in chronological order, allowing you to easily track and manage your collection. For each photo and raw image RAW.Ninja generates a thumbnail image. This image is used in the photos page, while the original image is only accessed for editing and downloading.</Typography>
            <Typography variant='h5' align='left' sx={{ margin: 2 }}>Upload Quota</Typography>
            <Typography align='left' sx={{ margin: 2, marginBottom: 2 }}>The space for your uploaded images is limited, and the limit is called quota. If the size of your uploaded photos reach the quota you can not upload more images. If you delete a photo, you get back the part of your quota the image used.</Typography>
            <Typography align='left' sx={{ margin: 2, marginBottom: 2 }}><b>Quota Calculation:</b> Both the original photo and its generated thumbnail contribute to your upload quota.</Typography>
            <Typography align='left' sx={{ margin: 2, marginBottom: 6 }}><b>Frozen Photo Quota:</b> If you mark a photo as frozen, it will not be available for download or editing, but it will still be visible in the photo list. Frozen photos consume only half of your upload quota.</Typography>
            <Typography variant='h5' align='left' sx={{ margin: 2 }}>Freeze and Unfreeze</Typography>
            <Typography align='left' sx={{ margin: 2, marginBottom: 2 }}><b>Freezing Photos:</b> You can mark a photo as frozen to save upload quota. However, frozen photos cannot be downloaded or edited.</Typography>
            <Typography align='left' sx={{ margin: 2, marginBottom: 6 }}><b>Unfreezing Photos:</b> To download or edit a frozen photo, it needs to be unfrozen first. Keep in mind that every photo unfrozen will not be eligible for freezing again for the next 30 days.</Typography>
            <Typography variant='h5' align='left' sx={{ margin: 2 }}>Happy uploading!</Typography>
        </Box>
    )
}

export default PhotoDocs

