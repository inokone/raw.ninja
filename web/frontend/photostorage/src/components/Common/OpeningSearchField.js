import React from 'react';

import { styled, alpha } from '@mui/material/styles';
import SearchIcon from '@mui/icons-material/Search';
import InputBase from '@mui/material/InputBase';
import { useNavigate, useLocation } from 'react-router-dom';

const Search = styled('div')(({ theme }) => ({
  position: 'relative',
  borderRadius: theme.shape.borderRadius,
  marginLeft: 0,
  width: '100%',
  [theme.breakpoints.up('sm')]: {
      width: 'auto',
  },
}));

const SearchIconWrapper = styled('div')(({ theme }) => ({
  padding: theme.spacing(0, 2),
  height: '100%',
  position: 'absolute',
  pointerEvents: 'none',
  display: 'flex',
  alignItems: 'center',
  justifyContent: 'center',
}));

const StyledInputBase = styled(InputBase)(({ theme }) => ({
color: 'inherit',
'& .MuiInputBase-input': {
    padding: theme.spacing(1, 1, 1, 0),
    paddingLeft: `calc(1em + ${theme.spacing(4)})`,
    transition: theme.transitions.create('width'),
    width: '100%',
    [theme.breakpoints.up('sm')]: {
    width: '0ch',
    cursor: 'pointer',
    '&:focus': {
      width: '20ch',
      backgroundColor: alpha(theme.palette.common.white, 0.15),
    },
    },
},
}));  

const handleQueryChange = (event, setQuery, location, navigate) => {
  setQuery(event.target.value)
  if (location.pathname !== '/search') {
        navigate("/search");
  }
}

const OpeningSearchField = (props) => {
  const navigate = useNavigate();
  const location = useLocation()

  return (<Search>
            <SearchIconWrapper>
              <SearchIcon onClick={input => input && input.focus()} />
            </SearchIconWrapper>
            <StyledInputBase
              placeholder="Search…"
              sx={{fontFamily: ['"Montserrat"', 'Open Sans'].join(',')}}
              inputProps={{ 'aria-label': 'search' }}
              onChange={e => handleQueryChange(e, props.setQuery, location, navigate)}
            />
          </Search>);
}

export default OpeningSearchField