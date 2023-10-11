import React from 'react';
import Typography from '@mui/material/Typography';
import Table from '@mui/material/Table';
import TableBody from '@mui/material/TableBody';
import TableCell from '@mui/material/TableCell';
import TableRow from '@mui/material/TableRow';
import { makeStyles } from '@mui/styles';


const useStyles = makeStyles({
  tableDataCell: {
    color: 'white',
  },
});

const MetadataDisplay = (props) => {
  const classes = useStyles();

  const formatShutterSpeed = (shutterSpeed) => {
    let validDividers = [2,4,8,15,30,60,125,250,500,1000,2000,4000,8000]
    if (shutterSpeed < 1) {
        let fraction = 1 / shutterSpeed
        let lastDivider = 2
        for (let i = 0; i < validDividers.length; i++) {
          let divider = validDividers[i]
          if(fraction < divider){
            console.log(divider)
            if (fraction - lastDivider > divider - fraction) {
              return "1/" + divider
            } else {
              return "1/" + lastDivider
            }
          }
        }
    } else {
        return shutterSpeed.toFixed(1) + "s";
    }
  }

  return(
     <>
        <Typography variant="h5">General</Typography>
        <Table size="small">
        <TableBody>
            <TableRow>
            <TableCell className={classes.tableDataCell}>Name</TableCell>
            <TableCell className={classes.tableDataCell}>{props.metadata.filename}</TableCell>
            </TableRow>
            <TableRow>
            <TableCell className={classes.tableDataCell}>Format</TableCell>
            <TableCell className={classes.tableDataCell}>{props.metadata.format}</TableCell>
            </TableRow>
            <TableRow>
            <TableCell className={classes.tableDataCell}>Uploaded</TableCell>
            <TableCell className={classes.tableDataCell}>{new Date(props.metadata.uploaded).toLocaleString()}</TableCell>
            </TableRow>
            <TableRow>
            <TableCell className={classes.tableDataCell}>Taken</TableCell>
            <TableCell className={classes.tableDataCell}>{new Date(props.metadata.metadata.timestamp).toLocaleString()}</TableCell>
            </TableRow>
        </TableBody>
        </Table>
        <Typography variant="h5" mt={3}>Image</Typography>
        <Table size="small">
        <TableBody>
            <TableRow>
            <TableCell className={classes.tableDataCell}>Width</TableCell>
            <TableCell className={classes.tableDataCell}>{props.metadata.metadata.width + " px"}</TableCell>
            </TableRow>
            <TableRow>
            <TableCell className={classes.tableDataCell}>Height</TableCell>
            <TableCell className={classes.tableDataCell}>{props.metadata.metadata.height + " px"}</TableCell>
            </TableRow>
            <TableRow>
            <TableCell className={classes.tableDataCell}>ISO</TableCell>
            <TableCell className={classes.tableDataCell}>{props.metadata.metadata.ISO}</TableCell>
            </TableRow>
            <TableRow>
            <TableCell className={classes.tableDataCell}>Aperture</TableCell>
            <TableCell className={classes.tableDataCell}>{Math.round((props.metadata.metadata.aperture + Number.EPSILON) * 100) / 100}</TableCell>
            </TableRow>
            <TableRow>
            <TableCell className={classes.tableDataCell}>Shutter Speed</TableCell>
            <TableCell className={classes.tableDataCell}>{formatShutterSpeed(props.metadata.metadata.shutter)}</TableCell>
            </TableRow>
        </TableBody>
        </Table>
        <Typography variant="h5" mt={3}>Camera</Typography>
        <Table size="small">
        <TableBody>
            <TableRow>
            <TableCell className={classes.tableDataCell}>Manufacturer</TableCell>
            <TableCell className={classes.tableDataCell}>{props.metadata.camera_make}</TableCell>
            </TableRow>
            <TableRow>
            <TableCell className={classes.tableDataCell}>Model</TableCell>
            <TableCell className={classes.tableDataCell}>{props.metadata.camera_model}</TableCell>
            </TableRow>
            <TableRow>
            <TableCell className={classes.tableDataCell}>Software Version</TableCell>
            <TableCell className={classes.tableDataCell}>{props.metadata.camera_sw}</TableCell>
            </TableRow>
        </TableBody>
        </Table>
        <Typography variant="h5" mt={3}>Lens</Typography>
        <Table size="small">
        <TableBody>
            <TableRow>
            <TableCell className={classes.tableDataCell}>Manufacturer</TableCell>
            <TableCell className={classes.tableDataCell}>{props.metadata.lens_make}</TableCell>
            </TableRow>
            <TableRow>
            <TableCell className={classes.tableDataCell}>Model</TableCell>
            <TableCell className={classes.tableDataCell}>{props.metadata.lens_model}</TableCell>
            </TableRow>
        </TableBody>
        </Table>
        </>
    );
}
export default MetadataDisplay; 