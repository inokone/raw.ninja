import React, { Fragment } from "react";
import PropTypes from "prop-types";
import classNames from "classnames";
import { Grid, Typography, Card, Button, Box, SvgIcon } from "@mui/material";
import { Link } from "react-router-dom"
import withStyles from "@mui/styles/withStyles";
import useMediaQuery from "@mui/material/useMediaQuery";
import Logo from '../Common/Logo'

const styles = (theme) => ({
    extraLargeButton: {
        paddingTop: theme.spacing(1.5),
        paddingBottom: theme.spacing(1.5),
        [theme.breakpoints.up("xs")]: {
            paddingTop: theme.spacing(1),
            paddingBottom: theme.spacing(1),
        },
        [theme.breakpoints.up("lg")]: {
            paddingTop: theme.spacing(2),
            paddingBottom: theme.spacing(2),
        },
        fontSize: theme.typography.h6.fontSize,
    },
    card: {
        boxShadow: theme.shadows[4],
        marginLeft: theme.spacing(2),
        marginRight: theme.spacing(2),
        [theme.breakpoints.up("xs")]: {
            paddingTop: theme.spacing(3),
            paddingBottom: theme.spacing(3),
        },
        [theme.breakpoints.up("sm")]: {
            paddingTop: theme.spacing(5),
            paddingBottom: theme.spacing(5),
            paddingLeft: theme.spacing(4),
            paddingRight: theme.spacing(4),
        },
        [theme.breakpoints.up("md")]: {
            paddingTop: theme.spacing(5.5),
            paddingBottom: theme.spacing(5.5),
            paddingLeft: theme.spacing(5),
            paddingRight: theme.spacing(5),
        },
        [theme.breakpoints.up("lg")]: {
            paddingTop: theme.spacing(6),
            paddingBottom: theme.spacing(6),
            paddingLeft: theme.spacing(6),
            paddingRight: theme.spacing(6),
        },
        [theme.breakpoints.down("xl")]: {
            width: "auto",
        },
    },
    wrapper: {
        position: "relative",
        backgroundColor: theme.palette.secondary.main,
        paddingBottom: theme.spacing(2),
    },
    image: {
        maxWidth: "100%",
        verticalAlign: "middle",
        borderRadius: theme.shape.borderRadius,
        boxShadow: theme.shadows[4],
    },
    brandText: {
        fontFamily: "Orbitron",
        fontWeight: 600
    },
    container: {
        marginTop: theme.spacing(6),
        marginBottom: theme.spacing(12),
        [theme.breakpoints.down("lg")]: {
            marginBottom: theme.spacing(9),
        },
        [theme.breakpoints.down("md")]: {
            marginBottom: theme.spacing(6),
        },
        [theme.breakpoints.down("md")]: {
            marginBottom: theme.spacing(3),
        },
    },
    containerFix: {
        [theme.breakpoints.up("md")]: {
            maxWidth: "none !important",
        },
    },
    mainLogo: {
        marginRight: 5,
        color: theme.palette.primary.main,
        fontSize: theme.typography.h3.fontSize
    },
});

function HeadSection(props) {
    const { classes, theme } = props;
    const isWidthUpLg = useMediaQuery(theme.breakpoints.up("lg"));

    return (
        <Fragment>
            <div className={classNames("lg-p-top", classes.wrapper)}>
                <div className={classNames("container-fluid", classes.container)}>
                    <Box display="flex" justifyContent="center" className="row">
                        <Card
                            className={classes.card}
                            data-aos-delay="200"
                            data-aos="zoom-in"
                        >
                            <div className={classNames(classes.containerFix, "container")}>
                                <Box justifyContent="space-between" className="row">
                                    <Grid item xs={12} md={3}>
                                        <Box
                                            display="flex"
                                            flexDirection="column"
                                            justifyContent="space-between"
                                            height="100%"
                                        >
                                            <Box mb={4}>
                                                <SvgIcon className={classes.mainLogo} component={Logo} />
                                                <Typography variant="h3"
                                                    className={classes.brandText}
                                                    display="inline"
                                                    color="primary">
                                                    RAW
                                                </Typography>
                                                <Typography variant="h3"
                                                    className={classes.brandText}
                                                    display="inline"
                                                    color="secondary">
                                                    Ninja
                                                </Typography>
                                            </Box>
                                            <div>
                                                <Box mb={2}>
                                                    <Typography
                                                        variant={isWidthUpLg ? "h6" : "body1"}
                                                        color="textSecondary"
                                                    >
                                                        Low-cost cloud image
                                                        storage for professional photographers.
                                                    </Typography>
                                                </Box>
                                                <Button
                                                    variant="contained"
                                                    color="primary"
                                                    fullWidth
                                                    className={classes.extraLargeButton}
                                                    href="/signup"
                                                >
                                                    Sign up for FREE
                                                </Button>
                                                <Typography mt={2} variant={isWidthUpLg ? "h6" : "body1"} color="textSecondary">
                                                    Already have an account?  <Link to="/login">Login Here</Link>
                                                </Typography>
                                            </div>
                                        </Box>
                                    </Grid>
                                </Box>
                            </div>
                        </Card>
                    </Box>
                </div>
            </div>
        </Fragment>
    );
}

HeadSection.propTypes = {
    classes: PropTypes.object,
    theme: PropTypes.object,
};

export default withStyles(styles, { withTheme: true })(HeadSection);