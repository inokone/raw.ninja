import React from 'react';

import { styled, alpha } from '@mui/material/styles';
import SearchIcon from '@mui/icons-material/Search';
import InputBase from '@mui/material/InputBase';
import { useNavigate } from 'react-router-dom';

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
  color: theme.palette.secondary.main
}));

const StyledInputBase = styled(InputBase)(({ theme }) => ({
  '& .MuiInputBase-input': {
    padding: theme.spacing(1, 1, 1, 0),
    paddingLeft: `calc(1em + ${theme.spacing(4)})`,
    transition: theme.transitions.create('width'),
    width: '100%',
    color: theme.palette.secondary.main,
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

const handleQueryChange = (event, navigate) => {
  navigate("/search", {
    state: {
      query: event.target.value
    }
  });
}

const OpeningSearchField = ({ setQuery }) => {
  const navigate = useNavigate()

  return (<Search>
    <SearchIconWrapper>
      <SearchIcon onClick={input => input && input.focus()} />
    </SearchIconWrapper>
    <StyledInputBase
      placeholder="Search…"
      name="photosearch"
      autoComplete="photosearch"
      inputProps={{ 'aria-label': 'search' }}
      onChange={e => handleQueryChange(e, navigate)}
    />
  </Search>);
}

export default OpeningSearchField