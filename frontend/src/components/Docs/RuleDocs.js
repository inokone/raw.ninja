import * as React from 'react';

import { Box, Typography, List, ListItem } from '@mui/material';
import AddIcon from '@mui/icons-material/Add';
import { Link } from 'react-router-dom';

const RuleDocs = () => {
    return (
        <Box maxWidth='md' sx={{ bgcolor: 'white', borderRadius: '8px', marginLeft: 'auto', marginRight: 'auto', marginTop: 2, paddingTop: 2, paddingBottom: 4 }}>
            <Typography variant='h4' sx={{ marginBottom: 6, marginTop: 2 }}>Optimize your storage with Lifecycle Rules</Typography>
            <Typography variant='h5' align='left' sx={{ margin: 2 }}>Introduction</Typography>
            <Typography align='left' sx={{ margin: 2, marginBottom: 6 }}>RAW.Ninja empowers you to manage your photos effectively and efficiently. With the introduction of lifecycle rules, you can automate the management of your photo storage, optimizing your quota and ensuring you only pay for the storage you actually need.</Typography>
            <Typography variant='h5' align='left' sx={{ margin: 2 }}>What are Lifecycle Rules?</Typography>
            <Typography align='left' sx={{ margin: 2, marginBottom: 6 }}>Lifecycle rules are predefined instructions that govern the automatic management of your photos. These rules specify the conditions under which specific actions should be taken, such as moving photos to a cheaper storage tier or deleting them altogether.</Typography>
            <Typography variant='h5' align='left' sx={{ margin: 2 }}>Benefits of Lifecycle Rules</Typography>
            <List sx={{ listStyleType: 'disc', pl: 6 }}>
                <ListItem sx={{ display: 'list-item' }}><Typography align='left'><b>Optimize Storage Quota:</b> By automating the removal of older or less frequently accessed photos, you can significantly reduce your storage consumption and save money.</Typography></ListItem>
                <ListItem sx={{ display: 'list-item' }}><Typography align='left'><b>Reduce Manual Effort:</b> Lifecycle rules eliminate the need for manual photo management, saving you time and effort.</Typography></ListItem>
                <ListItem sx={{ display: 'list-item' }}><Typography align='left'><b>Improve Data Organization:</b> By moving photos to different storage tiers based on their usage, you can maintain a more organized and efficient storage system.</Typography></ListItem>
            </List>
            <Typography variant='h5' align='left' sx={{ margin: 2 }}>Creating Lifecycle Rules</Typography>
            <Typography align='left' sx={{ margin: 2 }}>To create a lifecycle rule, follow these steps:</Typography>
            <List sx={{ listStyleType: 'numeric', pl: 6 }}>
                <ListItem sx={{ display: 'list-item' }}><Typography align='left'>Navigate to the <Link to="/rulesets">Rules</Link> page using your profile menu.</Typography></ListItem>
                <ListItem sx={{ display: 'list-item' }}><div style={{
                    display: 'flex',
                    alignItems: 'center',
                    flexWrap: 'wrap',
                }}><Typography align='left'>Click on the </Typography><AddIcon fontSize='small' sx={{marginLeft: '3px', marginRight: '3px'}}/><Typography> button on the bottom right corner of the page.</Typography></div></ListItem>
                <ListItem sx={{ display: 'list-item' }}><Typography align='left'><b>Define the Timing:</b> Determine the trigger for the rule, the age of a photo.</Typography></ListItem>
                <ListItem sx={{ display: 'list-item' }}><Typography align='left'><b>Specify the Action:</b> Decide whether the rule should move the photo to a cheaper storage tier or delete it altogether.</Typography></ListItem>
                <ListItem sx={{ display: 'list-item' }}><Typography align='left'><b>Combine Rules:</b> Combine up to three lifecycle rules into a rule set to apply multiple actions to different types of photos.</Typography></ListItem>
            </List>
            <Typography variant='h5' align='left' sx={{ margin: 2 }}>Associating Rule Sets with Collections and Albums</Typography>
            <Typography align='left' sx={{ margin: 2 }}>You can associate rule sets with specific collections or albums to apply the rules only to those items. This allows for targeted optimization of your photo storage based on your specific needs e.g. storing photos for a contract.</Typography>
            <Typography variant='h5' align='left' sx={{ margin: 2, marginTop: 6 }}>Applying Rule Sets to All Uploaded Photos</Typography>
            <Typography align='left' sx={{ margin: 2 }}>For comprehensive storage optimization, you can apply rule sets to all your uploaded photos. This ensures that all photos are managed according to your defined rules, streamlining your storage management.</Typography>
        </Box>
    )
}

export default RuleDocs