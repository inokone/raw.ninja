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

const MetadataDisplay = ({ metadata }) => {
  const classes = useStyles();

  const formatShutterSpeed = (shutterSpeed) => {
    let validDividers = [2, 4, 8, 15, 30, 60, 125, 250, 500, 1000, 2000, 4000, 8000]
    if (shutterSpeed < 1) {
      let fraction = 1 / shutterSpeed
      let lastDivider = 2
      for (let i = 0; i < validDividers.length; i++) {
        let divider = validDividers[i]
        if (fraction < divider) {
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

  const formatBytes = (bytes, decimals = 2) => {
    let negative = (bytes < 0)
    if (negative) {
      bytes = -bytes
    }
    if (!+bytes) return '0 Bytes'

    const k = 1024
    const dm = decimals < 0 ? 0 : decimals
    const sizes = ['Bytes', 'KiB', 'MiB', 'GiB', 'TiB', 'PiB', 'EiB', 'ZiB', 'YiB']

    const i = Math.floor(Math.log(bytes) / Math.log(k))

    let displayText = `${parseFloat((bytes / Math.pow(k, i)).toFixed(dm))} ${sizes[i]}`
    return negative ? "-" + displayText : displayText
  }

  return (
    <React.Fragment>
      <Typography variant="h5">General</Typography>
      <Table size="small">
        <TableBody>
          <TableRow>
            <TableCell className={classes.tableDataCell}>Format</TableCell>
            <TableCell className={classes.tableDataCell}>{metadata.format}</TableCell>
          </TableRow>
          <TableRow>
            <TableCell className={classes.tableDataCell}>Uploaded</TableCell>
            <TableCell className={classes.tableDataCell}>{new Date(metadata.uploaded).toLocaleString()}</TableCell>
          </TableRow>
          <TableRow>
            <TableCell className={classes.tableDataCell}>Taken</TableCell>
            <TableCell className={classes.tableDataCell}>{new Date(metadata.metadata.timestamp).toLocaleString()}</TableCell>
          </TableRow>
          <TableRow>
            <TableCell className={classes.tableDataCell}>Size</TableCell>
            <TableCell className={classes.tableDataCell}>{formatBytes(metadata.metadata.data_size)}</TableCell>
          </TableRow>
        </TableBody>
      </Table>
      <Typography variant="h5" mt={3}>Image</Typography>
      <Table size="small">
        <TableBody>
          <TableRow>
            <TableCell className={classes.tableDataCell}>Width</TableCell>
            <TableCell className={classes.tableDataCell}>{metadata.metadata.width + " px"}</TableCell>
          </TableRow>
          <TableRow>
            <TableCell className={classes.tableDataCell}>Height</TableCell>
            <TableCell className={classes.tableDataCell}>{metadata.metadata.height + " px"}</TableCell>
          </TableRow>
          <TableRow>
            <TableCell className={classes.tableDataCell}>ISO</TableCell>
            <TableCell className={classes.tableDataCell}>{metadata.metadata.ISO}</TableCell>
          </TableRow>
          <TableRow>
            <TableCell className={classes.tableDataCell}>Aperture</TableCell>
            <TableCell className={classes.tableDataCell}>ƒ/{Math.round((metadata.metadata.aperture + Number.EPSILON) * 100) / 100}</TableCell>
          </TableRow>
          <TableRow>
            <TableCell className={classes.tableDataCell}>Shutter Speed</TableCell>
            <TableCell className={classes.tableDataCell}>{formatShutterSpeed(metadata.metadata.shutter)}</TableCell>
          </TableRow>
        </TableBody>
      </Table>
      <Typography variant="h5" mt={3}>Camera</Typography>
      <Table size="small">
        <TableBody>
          <TableRow>
            <TableCell className={classes.tableDataCell}>Manufacturer</TableCell>
            <TableCell className={classes.tableDataCell}>{metadata.metadata.camera_make}</TableCell>
          </TableRow>
          <TableRow>
            <TableCell className={classes.tableDataCell}>Model</TableCell>
            <TableCell className={classes.tableDataCell}>{metadata.metadata.camera_model}</TableCell>
          </TableRow>
          <TableRow>
            <TableCell className={classes.tableDataCell}>Software Version</TableCell>
            <TableCell className={classes.tableDataCell}>{metadata.metadata.camera_sw}</TableCell>
          </TableRow>
        </TableBody>
      </Table>
      <Typography variant="h5" mt={3}>Lens</Typography>
      <Table size="small">
        <TableBody>
          <TableRow>
            <TableCell className={classes.tableDataCell}>Manufacturer</TableCell>
            <TableCell className={classes.tableDataCell}>{metadata.metadata.lens_make}</TableCell>
          </TableRow>
          <TableRow>
            <TableCell className={classes.tableDataCell}>Model</TableCell>
            <TableCell className={classes.tableDataCell}>{metadata.metadata.lens_model}</TableCell>
          </TableRow>
        </TableBody>
      </Table>
    </React.Fragment>
  );
}
export default MetadataDisplay; 