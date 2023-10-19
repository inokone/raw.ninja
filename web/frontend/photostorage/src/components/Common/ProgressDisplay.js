import { Box, CircularProgress, Typography } from '@mui/material';
import React from 'react';

const ProgressDisplay = (props) => {

    return <Box mt={10}>
               { props.text ? <Typography variant='h5' mb={3}>{props.text}</Typography> : null }
               <CircularProgress 
                    value={66}
                    size={200}
                    thickness={0.5}
                    style={{ padding: "5px" }} />
           </Box>;
};
export default ProgressDisplay;