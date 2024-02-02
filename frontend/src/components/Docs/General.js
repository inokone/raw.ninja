import * as React from 'react';

import { Box, Typography } from '@mui/material';

const GeneralDocs = () => {
    return (
        <Box maxWidth='md' sx={{ bgcolor: 'white', borderRadius: '8px', mx: 'auto', marginTop: 2, paddingTop: 2, paddingBottom: 4 }}>
            <Typography variant='h4' sx={{ marginBottom: 6, marginTop: 2 }}>General Information</Typography>
            <Typography variant='h5' align='left' sx={{ margin: 2 }}>Introduction</Typography>
            <Typography align='left' sx={{ margin: 2, marginBottom: 6 }}>RAW.Ninja is intended to be a low-cost cloud image storage for professional photographers. The application focuses on securely storing RAW image files for a predefined time (e.g. time interval set in a customer contract for a photo shoot).</Typography>
            <Typography variant='h5' align='left' sx={{ margin: 2 }}>Account</Typography>
            <Typography align='left' sx={{ margin: 2, marginBottom: 2 }}>You can create an account on the site with an email address and password, or with identity providers. The application supports Google and Facebook as such providers for now.</Typography>
            <Typography align='left' sx={{ margin: 2, marginBottom: 2 }}>You can delete all your data by sending an email to support@raw.ninja from the email address you are registered with, requesting to delete all your data.</Typography>
            <Typography align='left' sx={{ margin: 2, marginBottom: 6 }}>For accounts created with email address and password your upload rights are revoked until your email address is confirmed.</Typography>
            <Typography variant='h5' align='left' sx={{ margin: 2 }}>Uploads</Typography>
            <Typography align='left' sx={{ margin: 2, marginBottom: 2 }}>You can upload as many photos as fits in your quota. The limit for a single upload is 20 files, and for each upload the site creates a URL you can use for your reference.</Typography>
            <Typography align='left' sx={{ margin: 2, marginBottom: 2 }}>All photos uploaded to the site are stored</Typography>
            <Typography align='left' sx={{ margin: 2, marginBottom: 2 }}><li>Encrypted</li></Typography>
            <Typography align='left' sx={{ margin: 2, marginBottom: 2 }}><li>With 99.999999999% (11 nines) durability</li></Typography>
            <Typography align='left' sx={{ margin: 2, marginBottom: 6 }}></Typography>            
            <Typography variant='h5' align='left' sx={{ margin: 2 }}>Happy uploading!</Typography>
        </Box>
    )
}

export default GeneralDocs