import React from "react";
import { Typography, Box, Link as MuiLink, Stack } from "@mui/material";
import { ReactComponent as GoogleLogo } from './google.svg';
import { ReactComponent as FacebookLogo } from './facebook.svg';

const { REACT_APP_API_PREFIX } = process.env || "https://localhost:8080";

const SingleSignOn = () => {

    const getGoogleUrl = () => {
        return REACT_APP_API_PREFIX + `/api/public/v1/auth/google`;
    };

    const getFacebookUrl = () => {
        return REACT_APP_API_PREFIX + `/api/public/v1/auth/facebook`;
    };

    return (
        <Box sx={{ my: 2 }}>
            <Typography
                variant='h6'
                component='p'
                sx={{
                    my: 1,
                    textAlign: 'center',
                }}
            >
                Sign in with:
            </Typography>
            <Stack direction={"row"} alignItems={'center'} justifyContent={'center'}>
                <MuiLink
                    href={getGoogleUrl()}
                    sx={{
                        backgroundColor: '#fff',
                        boxShadow: '0 1px 13px 0 rgb(0 0 0 / 15%)',
                        borderRadius: 1,
                        py: 1,
                        mr: 1,
                        columnGap: '1rem',
                        textDecoration: 'none',
                        color: '#393e45',
                        cursor: 'pointer',
                        fontWeight: 500,
                        '&:hover': {
                            backgroundColor: '#f5f6f7',
                        },
                    }}
                    display='flex'
                    justifyContent='center'
                    alignItems='center'
                >
                    <GoogleLogo style={{ height: '24px', width: '48px' }} />
                </MuiLink>
                <MuiLink
                    href={getFacebookUrl()}
                    sx={{
                        backgroundColor: '#fff',
                        boxShadow: '0 1px 13px 0 rgb(0 0 0 / 15%)',
                        borderRadius: 1,
                        py: 1,
                        columnGap: '1rem',
                        textDecoration: 'none',
                        color: '#393e45',
                        cursor: 'pointer',
                        fontWeight: 500,
                        '&:hover': {
                            backgroundColor: '#f5f6f7',
                        },
                    }}
                    display='flex'
                    justifyContent='center'
                    alignItems='center'
                >
                    <FacebookLogo style={{ height: '24px', width: '48px' }} />
                </MuiLink>
            </Stack>
        </Box>
    )
}

export default SingleSignOn;