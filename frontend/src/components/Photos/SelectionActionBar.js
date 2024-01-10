import * as React from 'react';
import { useTheme } from '@mui/material/styles';
import Drawer from '@mui/material/Drawer';
import List from '@mui/material/List';
import Tooltip from '@mui/material/Tooltip';
import CollectionsIcon from '@mui/icons-material/Collections';
import DeleteIcon from '@mui/icons-material/Delete';
import ClearIcon from '@mui/icons-material/Clear';
import ListItem from '@mui/material/ListItem';
import ListItemButton from '@mui/material/ListItemButton';
import ListItemIcon from '@mui/material/ListItemIcon';

const drawerWidth = 60;

const SelectionActionBar = ({ open, handleCreate, handleDelete, handleClear }) => {
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
                <ListItem key="collection" disablePadding>
                    <Tooltip title="Create collection from selection">
                        <ListItemButton onClick={() => handleCreate()}>
                            <ListItemIcon>
                                <CollectionsIcon sx={{ color: theme.palette.background.paper }} />
                            </ListItemIcon>
                        </ListItemButton>
                    </Tooltip>
                </ListItem>
                <ListItem key="delete" disablePadding>
                    <Tooltip title="Delete selected photos">
                        <ListItemButton onClick={() => handleDelete()}>
                            <ListItemIcon>
                                <DeleteIcon sx={{ color: theme.palette.background.paper }} />
                            </ListItemIcon>
                        </ListItemButton>
                    </Tooltip>
                </ListItem>
                <ListItem key="clear" disablePadding>
                    <Tooltip title="Clear selection">
                        <ListItemButton onClick={() => handleClear()}>
                            <ListItemIcon>
                                <ClearIcon sx={{ color: theme.palette.background.paper }} />
                            </ListItemIcon>
                        </ListItemButton>
                    </Tooltip>
                </ListItem>
            </List>
        </Drawer>

    );
}

export default SelectionActionBar
