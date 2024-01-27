import * as React from 'react';

import { Box, Typography, List, ListItem } from '@mui/material';

const PrivacyPolicy = () => {
    return (
        <Box maxWidth='md' sx={{ bgcolor: 'white', borderRadius: '8px', mx: 'auto', marginTop: 2, paddingTop: 2, paddingBottom: 4 }}>
            <Typography variant='h4' sx={{ marginBottom: 6, marginTop: 2 }}>Privacy Policy</Typography>
            <Typography variant='h5' align='left' sx={{ margin: 2 }}>Introduction</Typography>
            <Typography align='left' sx={{ margin: 2, marginBottom: 6 }}>This Privacy Policy describes how RAW.Ninja ("we", "us" or "our") collects, uses, and discloses your personal information when you use our web application for storing photos and RAW images.</Typography>
            <Typography variant='h5' align='left' sx={{ margin: 2 }}>Collection of Personal Information</Typography>
            <Typography align='left' sx={{ margin: 2, marginBottom: 2 }}>We may collect the following personal information from you:</Typography>
            <List sx={{ listStyleType: 'disc', pl: 6 }}>
                <ListItem sx={{ display: 'list-item' }}><Typography align='left'><b>Registration information:</b> When you register for an account, we collect your username, password, email address, and optionally your first name and last name.</Typography></ListItem>
                <ListItem sx={{ display: 'list-item' }}><Typography align='left'><b>Content you upload:</b> When you upload photos or RAW images to our application, we collect the content of those images and any description you provide to the images.</Typography></ListItem>
                <ListItem sx={{ display: 'list-item' }}><Typography align='left'><b>Device information:</b> We may collect information about your device, such as your IP address, browser type, and operating system.</Typography></ListItem>
                <ListItem sx={{ display: 'list-item' }}><Typography align='left'><b>Usage information:</b> We may collect information about how you use our application, such as the pages you visit and the features you use.</Typography></ListItem>
            </List>
            <Typography variant='h5' align='left' sx={{ margin: 2 }}>Use of Personal Information</Typography>
            <Typography align='left' sx={{ margin: 2, marginBottom: 2 }}>We use your personal information for the following purposes:</Typography>
            <List sx={{ listStyleType: 'disc', pl: 6 }}>
                <ListItem sx={{ display: 'list-item' }}><Typography align='left'><b>To provide our application to you:</b> We use your registration information to create your account and to provide you with access to our application.</Typography></ListItem>
                <ListItem sx={{ display: 'list-item' }}><Typography align='left'><b>To store your photos and RAW images:</b> We store your content on our servers so that you can access it from anywhere.</Typography></ListItem>
                <ListItem sx={{ display: 'list-item' }}><Typography align='left'><b>To improve our application:</b>  We use your usage information to help us improve our application and to make sure it is working properly.</Typography></ListItem>
                <ListItem sx={{ display: 'list-item' }}><Typography align='left'><b>To send you marketing communications:</b> We may send you email or other communications about our products and services, but you can opt out of receiving these communications at any time.</Typography></ListItem>
            </List>
            <Typography variant='h5' align='left' sx={{ margin: 2 }}>Sharing of Personal Information</Typography>
            <Typography align='left' sx={{ margin: 2, marginBottom: 2 }}>We do not share your personal information with third parties except in the following limited circumstances:</Typography>
            <List sx={{ listStyleType: 'disc', pl: 6 }}>
                <ListItem sx={{ display: 'list-item' }}><Typography align='left'><b>With your consent:</b> We will share your personal information with third parties if you give us your consent to do so.</Typography></ListItem>
                <ListItem sx={{ display: 'list-item' }}><Typography align='left'><b>Service providers:</b> Third parties that provide services on our behalf or help us operate the Service or our business (such as hosting, information technology, customer support, email delivery, consumer research, marketing, and website analytics).</Typography></ListItem>
                <ListItem sx={{ display: 'list-item' }}><Typography align='left'><b>To comply with the law:</b> We may share your personal information with third parties if we are required to do so by law.</Typography></ListItem>
                <ListItem sx={{ display: 'list-item' }}><Typography align='left'><b>For fraud prevention:</b> We may share your personal information with third parties to prevent fraud.</Typography></ListItem>
            </List>
            <Typography variant='h5' align='left' sx={{ margin: 2 }}>Data Security</Typography>
            <Typography align='left' sx={{ margin: 2, marginBottom: 6 }}>We take precautions to protect your personal information from unauthorized access, use, disclosure, alteration, or destruction. We use industry - standard security measures, such as firewalls and encryption, to protect your information.</Typography>
            <Typography variant='h5' align='left' sx={{ margin: 2 }}>Your Choices</Typography>
            <Typography align='left' sx={{ margin: 2, marginBottom: 2 }}>You have the following choices regarding your personal information:</Typography>
            <List sx={{ listStyleType: 'disc', pl: 6 }}>
                <ListItem sx={{ display: 'list-item' }}><Typography align='left'><b>Access:</b> You have the right to access your personal information that we collect and retain.</Typography></ListItem>
                <ListItem sx={{ display: 'list-item' }}><Typography align='left'><b>Correction:</b> You have the right to correct any inaccurate or incomplete personal information that we collect and retain.</Typography></ListItem>
                <ListItem sx={{ display: 'list-item' }}><Typography align='left'><b>Deletion:</b> You have the right to request that we delete your personal information from our records.</Typography></ListItem>
                <ListItem sx={{ display: 'list-item' }}><Typography align='left'><b>Objection:</b> You have the right to object to the processing of your personal information, or to request that we restrict the processing of your personal information.</Typography></ListItem>
                <ListItem sx={{ display: 'list-item' }}><Typography align='left'><b>Withdrawal of consent:</b> You have the right to withdraw your consent to the processing of your personal information.</Typography></ListItem>
            </List>
            <Typography align='left' sx={{ margin: 2, marginBottom: 2 }}>To exercise any of these choices, please contact us at support@raw.ninja.</Typography>
            <Typography variant='h5' align='left' sx={{ margin: 2 }}>Changes to this Privacy Policy</Typography>
            <Typography align='left' sx={{ margin: 2, marginBottom: 6 }}> We may update this Privacy Policy from time to time. We will notify you of any material changes by posting the new Privacy Policy on our website. You should review this Privacy Policy periodically to ensure that you are aware of the latest updates.</Typography>
            <Typography variant='h5' align='left' sx={{ margin: 2 }}>Contact Us</Typography>
            <Typography align='left' sx={{ margin: 2, marginBottom: 2 }}>If you have any questions about this Privacy Policy, please contact us at support@raw.ninja.</Typography>
        </Box>

    )
}

export default PrivacyPolicy