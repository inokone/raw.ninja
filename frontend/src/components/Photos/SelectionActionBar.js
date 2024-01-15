import * as React from 'react';
import PropTypes from "prop-types";
import { useTheme } from '@mui/material/styles';
import Drawer from '@mui/material/Drawer';
import List from '@mui/material/List';
import Tooltip from '@mui/material/Tooltip';
import ListItem from '@mui/material/ListItem';
import ListItemButton from '@mui/material/ListItemButton';
import ListItemIcon from '@mui/material/ListItemIcon';

const drawerWidth = 60;

const SelectionActionBar = ({ open, items, actions }) => {
    const theme = useTheme();

    return (
        <Drawer
            sx={{
                paddingTop: '48px',
                width: drawerWidth,
                flexShrink: 0,
                '& .MuiDrawer-paper': {
                    width: drawerWidth,
                    boxSizing: 'border-box',
                },
            }}
            PaperProps={{
                sx: {
                    backgroundColor: theme.palette.secondary.main,
                    color: theme.palette.secondary,
                }
            }}
            variant="persistent"
            anchor="left"
            open={open}
        >
            <List sx={{ marginTop: '48px', width: drawerWidth }}>
                {actions.map((action) => (
                <ListItem disablePadding key={action.tooltip}>
                    <Tooltip title={action.tooltip}>
                        <ListItemButton onClick={() => action.action(items)}>
                            <ListItemIcon>
                                {action.icon}
                            </ListItemIcon>
                        </ListItemButton>
                    </Tooltip>
                </ListItem>))}
            </List>
        </Drawer>
    );
}

SelectionActionBar.propTypes = {
    open: PropTypes.bool.isRequired,
    items: PropTypes.array.isRequired,
    actions: PropTypes.array.isRequired
};

export default SelectionActionBar
