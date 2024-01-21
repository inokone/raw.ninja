import * as React from 'react';
import { Grid } from '@mui/material';
import AlbumCard from './AlbumCard';

const Collections = ({ collections, onClick }) => {

    return (
        <>
            {collections !== null && collections.length > 0 &&
                <Grid container>
                    {collections.map((collection) => {
                        return (<Grid item key={collection.id} xs={6} md={4} lg={2} xl={2}><AlbumCard album={collection} onClick={onClick} /></Grid>)
                    })}
                </Grid>
            }
        </>
    )
}

export default Collections;