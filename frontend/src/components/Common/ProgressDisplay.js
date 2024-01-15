import { Box, CircularProgress, Typography } from '@mui/material';
import React from 'react';
import PropTypes from "prop-types";

const ProgressDisplay = ({ text }) => {

    return <Box mt={10}>
               { text ? <Typography variant='h5' mb={3}>{text}</Typography> : null }
               <CircularProgress 
                    value={66}
                    size={200}
                    thickness={0.5}
                    style={{ padding: "5px" }} />
           </Box>;
};

ProgressDisplay.propTypes = {
    text: PropTypes.string
};

export default ProgressDisplay;