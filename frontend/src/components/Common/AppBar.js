import * as React from 'react';
import AppBar from '@mui/material/AppBar';
import Box from '@mui/material/Box';
import Toolbar from '@mui/material/Toolbar';
import IconButton from '@mui/material/IconButton';
import Typography from '@mui/material/Typography';
import Menu from '@mui/material/Menu';
import MenuIcon from '@mui/icons-material/Menu';
import Container from '@mui/material/Container';
import Avatar from '@mui/material/Avatar';
import Button from '@mui/material/Button';
import Tooltip from '@mui/material/Tooltip';
import MenuItem from '@mui/material/MenuItem';
import ViewModuleIcon from '@mui/icons-material/ViewModule';
import CloudUploadIcon from '@mui/icons-material/CloudUpload';
import CollectionsIcon from '@mui/icons-material/Collections';
import OpeningSearchField from './OpeningSearchField';
import { SvgIcon } from '@mui/material';
import { useNavigate } from "react-router-dom";
import stringToColor from 'string-to-color';
import withStyles from "@mui/styles/withStyles";
import PropTypes from "prop-types";
import Logo from './Logo';


const settings = {
  Profile: 'profile',
  Logout: 'logout'
};

const styles = theme => ({
  appBar: {
    boxShadow: theme.shadows[6],
    backgroundColor: theme.palette.common.white
  },
  toolbar: {
    display: "flex",
    justifyContent: "space-between"
  },
  menuButtonText: {
    fontSize: theme.typography.body1.fontSize,
    fontWeight: theme.typography.h6.fontWeight
  },
  brandText: {
    fontFamily: "Orbitron",
    fontWeight: 600
  },
  noDecoration: {
    textDecoration: "none !important"
  },
});

const ResponsiveAppBar = ({ theme, classes, user, setQuery }) => {
  const navigate = useNavigate()
  const [anchorElNav, setAnchorElNav] = React.useState(null);
  const [anchorElUser, setAnchorElUser] = React.useState(null);
  const adminRoleID = 0

  const isAdmin = (user) => {
    return user.role.id === adminRoleID
  }

  const profileMenu = () => {
    let res = {}
    if (isAdmin(user)) {
      res.Admin = 'admin'
    }
    return { ...res, ...settings }
  }

  const handleOpenNavMenu = (event) => {
    setAnchorElNav(event.currentTarget);
  };

  const handleCloseNavMenu = () => {
    setAnchorElNav(null);
  };

  const handleOpenUserMenu = (event) => {
    setAnchorElUser(event.currentTarget);
  };

  const handleCloseUserMenu = () => {
    setAnchorElUser(null);
  };

  const isAuthenticated = () => {
    return user !== null
  }

  const getInitials = () => {
    if (!isAuthenticated())
      return ""
    if (user.first_name)
      return user.first_name + " " + user.last_name
    return user.email
  }

  const handleMenuClick = (page) => {
    setAnchorElNav(null);
    navigate('/' + page);
  };

  const handleUserClick = (page) => {
    setAnchorElUser(null);
    navigate('/' + page);
  };

  return (
    <AppBar position="fixed" sx={{ zIndex: (theme) => theme.zIndex.drawer + 1 }}>
      <Container maxWidth="xxl" sx={{ bgcolor: theme.palette.common.black }}>
        <Toolbar disableGutters variant="dense">
          <SvgIcon sx={{ mr: 0.5, color: theme.palette.primary.main }} component={Logo} />
          <Typography
            variant="h4"
            className={classes.brandText}
            display="inline"
            color="primary"
            onClick={() => navigate(user ? '/home' : '/')}
          >
            RAW
          </Typography>
          <Typography
            variant="h4"
            className={classes.brandText}
            display="inline"
            color="secondary"
            onClick={() => navigate(user ? '/home' : '/')}
          >
            Ninja
          </Typography>
          <Typography
            ml={3}
            variant="h4"
            className={classes.brandText}
            display="inline"
            color="orange"
            onClick={() => navigate(user ? '/home' : '/')}
          >
            BETA
          </Typography>
          <Box sx={{ flexGrow: 1 }} />
          {isAuthenticated() &&
            <React.Fragment>
              <Box sx={{ display: { xs: 'none', md: 'flex' } }}><OpeningSearchField setQuery={setQuery} /></Box>
              <Box sx={{ display: { xs: 'none', md: 'flex' } }}>
                <Button
                  key="upload"
                  onClick={() => handleMenuClick("upload")}
                  sx={{ m: 0, color: theme.palette.secondary.main }}>
                  <Typography>UPLOAD</Typography>
                </Button>
                <Button
                  key="photos"
                  onClick={() => handleMenuClick("photos")}
                  sx={{ m: 0, color: theme.palette.secondary.main }}>
                  <Typography>PHOTOS</Typography>
                </Button>
                <Button
                  key="albums"
                  onClick={() => handleMenuClick("albums")}
                  sx={{ m: 0, color: theme.palette.secondary.main }}>
                  <Typography>ALBUMS</Typography>
                </Button>
              </Box>
              <Box sx={{ display: { xs: 'flex', md: 'none' } }}>
                <IconButton
                  size="large"
                  aria-label="account of current user"
                  aria-controls="menu-appbar"
                  aria-haspopup="true"
                  onClick={handleOpenNavMenu}
                >
                  <MenuIcon sx={{ color: theme.palette.secondary.main }} />
                </IconButton>
                <Menu
                  id="menu-appbar"
                  anchorEl={anchorElNav}
                  anchorOrigin={{
                    vertical: 'bottom',
                    horizontal: 'left',
                  }}
                  keepMounted
                  transformOrigin={{
                    vertical: 'top',
                    horizontal: 'left',
                  }}
                  open={Boolean(anchorElNav)}
                  onClose={handleCloseNavMenu}
                  sx={{
                    display: { xs: 'block', md: 'none' }
                  }}
                >
                  <MenuItem key="upload" onClick={() => handleMenuClick('upload')}>
                    <CloudUploadIcon sx={{ mr: 1 }} />
                    <Typography textAlign="center" sx={{ fontFamily: ['"Montserrat"', 'Open Sans'].join(',') }}>Upload</Typography>
                  </MenuItem>
                  <MenuItem key="photos" onClick={() => handleMenuClick('photos')}>
                    <ViewModuleIcon sx={{ mr: 1 }} />
                    <Typography textAlign="center" sx={{ fontFamily: ['"Montserrat"', 'Open Sans'].join(',') }}>Photos</Typography>
                  </MenuItem>
                  <MenuItem key="albums" onClick={() => handleMenuClick('albums')}>
                    <CollectionsIcon sx={{ mr: 1 }} />
                    <Typography textAlign="center" sx={{ fontFamily: ['"Montserrat"', 'Open Sans'].join(',') }}>Albums</Typography>
                  </MenuItem>
                </Menu>
              </Box>
              <Box sx={{ md: 0, ml: 1 }}>
                <Tooltip title={"Open profile"}>
                  <IconButton onClick={handleOpenUserMenu} sx={{ p: 0 }}>
                    <Avatar alt={getInitials()} src="/static/images/avatar/3.jpg" sx={{ bgcolor: stringToColor(getInitials()) }} />
                  </IconButton>
                </Tooltip>
                <Menu
                  sx={{ mt: '45px' }}
                  id="menu-appbar"
                  anchorEl={anchorElUser}
                  anchorOrigin={{
                    vertical: 'top',
                    horizontal: 'right',
                  }}
                  keepMounted
                  transformOrigin={{
                    vertical: 'top',
                    horizontal: 'right',
                  }}
                  open={Boolean(anchorElUser)}
                  onClose={handleCloseUserMenu}
                >
                  {Object.entries(profileMenu()).map(([setting, path], i) => (
                    <MenuItem key={setting} onClick={() => handleUserClick(path)} sx={{ fontFamily: ['"Montserrat"', 'Open Sans'].join(',') }}>
                      <Typography textAlign="center" sx={{ fontFamily: ['"Montserrat"', 'Open Sans'].join(',') }}>{setting}</Typography>
                    </MenuItem>
                  ))}
                </Menu>
              </Box>
            </React.Fragment>}
        </Toolbar>
      </Container>
    </AppBar>
  );
}

ResponsiveAppBar.propTypes = {
  theme: PropTypes.object.isRequired,
  classes: PropTypes.object.isRequired,
};

export default withStyles(styles, { withTheme: true })(ResponsiveAppBar);