import * as React from 'react';
import { useTheme } from '@mui/material/styles';
import Drawer from '@mui/material/Drawer';
import List from '@mui/material/List';
import Tooltip from '@mui/material/Tooltip';
import ListItem from '@mui/material/ListItem';
import ListItemButton from '@mui/material/ListItemButton';
import ListItemIcon from '@mui/material/ListItemIcon';

const drawerWidth = 60;

const SelectionActionBar = ({ open, actions }) => {
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
                <ListItem key="collection" disablePadding>
                    <Tooltip title={action.tooltip}>
                        <ListItemButton onClick={action.action}>
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

export default SelectionActionBar
