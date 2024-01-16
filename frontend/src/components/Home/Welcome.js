import * as React from 'react';

import { Box, Typography } from '@mui/material';

const Welcome = () => {
    return (
        <Box maxWidth='md' sx={{ bgcolor: 'white', borderRadius: '8px', marginLeft: 'auto', marginRight: 'auto', marginTop: 2, paddingTop: 2, paddingBottom: 4 }}>
            <Typography variant='h2' sx={{ marginBottom: 8, marginTop: 2 }}>{"Welcome to RAW.Ninja!"}</Typography>
            <Typography align='left' sx={{ margin: 2, marginBottom: 4 }}>We're excited to have you on board as a new member of our photo-centric community.</Typography>
            <Typography align='left' sx={{ margin: 2, marginBottom: 4 }}>Your journey with us begins now, as you embark on a seamless experience for storing and managing your raw photos. Dive into our user-friendly platform to securely upload and organize your precious data in their purest form. Explore the possibilities of preserving your visual stories, and if you ever have questions or need assistance, our team is here to help.</Typography>
            <Typography align='left' sx={{ margin: 2 }}>Thank you for entrusting us with your photos, and we look forward to being a part of your photography journey on RAW.Ninja!</Typography>
        </Box>
    )
}

export default Welcome