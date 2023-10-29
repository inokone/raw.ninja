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
import OpeningSearchField from './OpeningSearchField';
import { SvgIcon } from '@mui/material';
import { useNavigate } from "react-router-dom";
import stringToColor from 'string-to-color';

const settings = {
  Profile: 'profile',
  Account: 'account',
  Logout: 'logout'
};

const ResponsiveAppBar = ({ user, setQuery }) => {
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

  const SvgComponent = (props) => (
    <svg
      xmlns="http://www.w3.org/2000/svg"
      width={32}
      height={32}
      viewBox="0 0 24 24"
      {...props}
    >
      <path d="M11.29.01c-.589.038-1.232.124-1.767.239l-.146.03-.109.142C8.295 1.68 7.666 3.87 7.41 6.895c-.092 1.07-.118 1.805-.116 3.133.002.614.01 1.169.017 1.23l.013.115L8.4 9.51l1.076-1.867.162-.103c1.825-1.182 4.02-2.338 5.752-3.03 1.874-.75 3.421-1.107 4.686-1.082.182.003.332.003.33-.002A12.102 12.102 0 0 0 16.093.709a12.02 12.02 0 0 0-3.03-.672A19.302 19.302 0 0 0 11.29.01z" />
      <path d="M8.439.53A12.034 12.034 0 0 0 3.49 3.535 12.009 12.009 0 0 0 .642 8.107l-.103.297.06.15c.396 1.001 1.332 2.165 2.681 3.328 1.469 1.264 3.382 2.54 5.692 3.791.097.054.144.071.137.05A501.5 501.5 0 0 0 7.04 12.14a.745.745 0 0 1-.078-.16l-.025-.685c-.017-.44-.02-1.016-.013-1.754.011-1.082.023-1.431.079-2.241.212-3.015.816-5.386 1.718-6.75A.865.865 0 0 0 8.795.43c0-.014-.013-.01-.356.099zM19.46 3.811a14.3 14.3 0 0 1-.262.028c-1.884.186-4.623 1.248-7.596 2.942a46.2 46.2 0 0 0-1.368.818c-.035.024.194.028 2.115.028h2.155l.349.182a39.4 39.4 0 0 1 4.13 2.49c2.262 1.583 3.802 3.112 4.526 4.49.06.117.114.21.116.207.015-.013.105-.403.158-.672.073-.377.142-.86.184-1.285.045-.47.045-1.606 0-2.082-.152-1.582-.539-2.947-1.216-4.306a11.632 11.632 0 0 0-1.358-2.133c-.227-.283-.526-.628-.563-.647-.064-.037-.527-.071-.91-.07-.221.003-.428.008-.46.01z" />
      <path d="M14.898 8.269c.006.016.49.859 1.075 1.872.855 1.482 1.066 1.859 1.071 1.92.043.552.058 2.52.025 3.38-.104 2.67-.434 4.776-1.004 6.378-.223.625-.516 1.234-.782 1.628a.478.478 0 0 0-.068.117 12.032 12.032 0 0 0 7.234-5.647c.34-.597.696-1.387.932-2.073l.087-.252-.06-.152a6.916 6.916 0 0 0-.45-.868c-.562-.9-1.558-1.936-2.86-2.977-1.25-.998-2.832-2.04-4.558-2.996-.621-.344-.655-.36-.642-.33zM.318 9.235a12.169 12.169 0 0 0-.094 5.09c.388 1.993 1.31 3.9 2.632 5.445l.255.299c.086.099.435.136 1.105.121.846-.02 1.653-.168 2.767-.514 1.861-.574 4.217-1.701 6.564-3.14l.257-.158-2.16-.009-2.164-.01-.3-.157C4.8 13.902 1.618 11.34.5 9.208a2.82 2.82 0 0 0-.118-.212c-.004 0-.032.107-.064.239z" />
      <path d="m15.605 14.49-1.078 1.866-.26.167a39.447 39.447 0 0 1-3.978 2.221C7.709 19.982 5.563 20.6 3.94 20.57a3.044 3.044 0 0 0-.34.002c.003.017.474.443.68.618 2.788 2.343 6.438 3.283 10.044 2.588l.304-.058.11-.144c.796-1.03 1.364-2.685 1.683-4.908.137-.964.22-1.889.274-3.105.02-.502.025-2.641.002-2.825l-.013-.116-1.079 1.868z" />
    </svg>
  );

  return (
    <AppBar position="static">
      <Container maxWidth="xxl">
        <Toolbar disableGutters variant="dense">
          <SvgIcon sx={{ display: { xs: 'none', md: 'flex' }, mr: 1 }} component={SvgComponent} />
          <Typography
            variant="h5"
            noWrap
            onClick={() => navigate('/')}
            sx={{
              mr: 2,
              cursor: "pointer",
              display: { xs: 'none', md: 'flex' },
              fontWeight: 700,
              letterSpacing: '.3rem',
              color: 'inherit',
              textDecoration: 'none',
              fontFamily: ['"Montserrat"', 'Open Sans'].join(',')
            }}
          >
            PhotoStore
          </Typography>
          <Box sx={{ display: { xs: 'flex', md: 'none' } }}>
            <IconButton
              size="large"
              aria-label="account of current user"
              aria-controls="menu-appbar"
              aria-haspopup="true"
              onClick={handleOpenNavMenu}
              color="inherit"
            >
              <MenuIcon />
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
            </Menu>
          </Box>
          <SvgIcon sx={{ display: { xs: 'flex', md: 'none' }, mr: 1 }} component={SvgComponent} />
          <Typography
            variant="h5"
            noWrap
            onClick={() => navigate('/')}
            sx={{
              mr: 2,
              cursor: "pointer",
              display: { xs: 'flex', md: 'none' },
              fontWeight: 700,
              letterSpacing: '.3rem',
              color: 'inherit',
              textDecoration: 'none',
              fontFamily: ['"Montserrat"', 'Open Sans'].join(',')
            }}
          >
            PhotoStore
          </Typography>
          <Box sx={{ flexGrow: 1 }} />
          {isAuthenticated() ?
            <OpeningSearchField setQuery={setQuery} /> : null}
          <Box sx={{ display: { xs: 'none', md: 'flex' } }}>
            <Button
              key="upload"
              onClick={() => handleMenuClick("upload")}
              sx={{ color: 'white', m: 0 }}>
              <CloudUploadIcon />
            </Button>
            <Button
              key="photos"
              onClick={() => handleMenuClick("photos")}
              sx={{ color: 'white', m: 0 }}>
              <ViewModuleIcon />
            </Button>
          </Box>
          {isAuthenticated() ?
            <Box sx={{ md: 0, ml: 1 }}>
              <Tooltip title={"Open profile"}>
                <IconButton onClick={handleOpenUserMenu} sx={{ p: 0 }}>
                  <Avatar alt={getInitials()} src="/static/images/avatar/3.jpg" sx={{ bgcolor: stringToColor(user.email) }} />
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
            : null}
        </Toolbar>
      </Container>
    </AppBar>
  );
}
export default ResponsiveAppBar;