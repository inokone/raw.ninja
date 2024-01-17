import * as React from 'react';
import PropTypes from "prop-types";
import AlbumDocs from './AlbumDocs';
import PhotoDocs from './PhotoDocs';
import RuleDocs from './RuleDocs';
import { Box, Tabs, Tab} from '@mui/material';

function CustomTabPanel(props) {
    const { children, value, index, ...other } = props;

    return (
        <div
            role="tabpanel"
            hidden={value !== index}
            id={`simple-tabpanel-${index}`}
            aria-labelledby={`simple-tab-${index}`}
            {...other}
        >
            {value === index && (
                <Box sx={{ p: 3 }}>
                    {children}
                </Box>
            )}
        </div>
    );
}

CustomTabPanel.propTypes = {
    children: PropTypes.node,
    index: PropTypes.number.isRequired,
    value: PropTypes.number.isRequired,
};

function a11yProps(index) {
    return {
        id: `simple-tab-${index}`,
        'aria-controls': `simple-tabpanel-${index}`,
    };
}

const Docs = () => {
    const [value, setValue] = React.useState(0);

    const handleChange = (event, newValue) => {
        setValue(newValue);
    };

    return (
        <Box sx={{ borderBottom: 1, borderColor: 'divider', maxWidth: 'md', marginLeft: 'auto', marginRight: 'auto', marginTop: 2, paddingTop: 2, paddingBottom: 4 }}>
            <Box sx={{ borderBottom: 1, borderColor: 'divider' }}>
                <Tabs value={value} onChange={handleChange} aria-label="basic tabs example">
                    <Tab label="Photos" {...a11yProps(0)} />
                    <Tab label="Albums" {...a11yProps(1)} />
                    <Tab label="Lifecycle Rules" {...a11yProps(2)} />
                </Tabs>
            </Box>
            <CustomTabPanel value={value} index={0}>
                <PhotoDocs />
            </CustomTabPanel>
            <CustomTabPanel value={value} index={1}>
                <AlbumDocs />
            </CustomTabPanel>
            <CustomTabPanel value={value} index={2}>
                <RuleDocs/>
            </CustomTabPanel>
        </Box>
    )
}

export default Docs