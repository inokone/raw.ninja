import React from "react";
import classNames from "classnames";
import { Grid, Typography } from "@mui/material";
import { withStyles } from "@mui/styles";
import PriceCard from "./PriceCard";
import calculateSpacing from "./calculateSpacing";
import useMediaQuery from "@mui/material/useMediaQuery";
import useWidth from "./useWidth";

const styles = (theme) => ({
    containerFix: {
        [theme.breakpoints.down("lg")]: {
            paddingLeft: theme.spacing(6),
            paddingRight: theme.spacing(6),
        },
        [theme.breakpoints.down("md")]: {
            paddingLeft: theme.spacing(4),
            paddingRight: theme.spacing(4),
        },
        [theme.breakpoints.down("sm")]: {
            paddingLeft: theme.spacing(2),
            paddingRight: theme.spacing(2),
        },
        overflow: "hidden",
        paddingTop: theme.spacing(1),
        paddingBottom: theme.spacing(1),
    },
    cardWrapper: {
        [theme.breakpoints.down("sm")]: {
            marginLeft: "auto",
            marginRight: "auto",
            maxWidth: 340,
        },
    },
    cardWrapperHighlighted: {
        [theme.breakpoints.down("sm")]: {
            marginLeft: "auto",
            marginRight: "auto",
            maxWidth: 360,
        },
    },
});

function PricingSection(props) {
    const { classes, theme } = props;
    const width = useWidth();
    const isWidthUpMd = useMediaQuery(theme.breakpoints.up("md"));
    return (
        <div className="lg-p-top" style={{ backgroundColor: "#FFFFFF" }}>
            <Typography variant="h3" align="center" className="lg-mg-bottom">
                Pricing
            </Typography>
            <div className={classNames("container-fluid", classes.containerFix)}>
                <Grid
                    container
                    spacing={calculateSpacing(width, theme)}
                    className={classNames(classes.gridContainer, "lg-mg-bottom")}
                >
                    <Grid
                        item
                        xs={12}
                        sm={6}
                        lg={4}
                        className={classes.cardWrapper}
                        data-aos="zoom-in-up"
                    >
                        <PriceCard
                            title="Free"
                            pricing={
                                <span>
                                    $0.00
                                    <Typography display="inline"> / month</Typography>
                                </span>
                            }
                            features={["Up to 500 Mb storage", "Store RAW images", "Edit RAW images"]}
                        />
                    </Grid>
                    <Grid
                        item
                        className={classes.cardWrapperHighlighted}
                        xs={12}
                        sm={6}
                        lg={4}
                        data-aos="zoom-in-up"
                        data-aos-delay="200"
                    >
                        <PriceCard
                            highlighted
                            title="Premium"
                            pricing={
                                <span>
                                    Not Available
                                </span>
                            }
                            features={["Up to 1Tb storage", "Image lifecycle rules", "Image sharing options"]}
                        />
                    </Grid>
                    <Grid
                        item
                        className={classes.cardWrapper}
                        xs={12}
                        sm={6}
                        lg={4}
                        data-aos="zoom-in-up"
                        data-aos-delay={isWidthUpMd ? "400" : "0"}
                    >
                        <PriceCard
                            title="Business"
                            pricing={
                                <span>
                                    Not Available
                                </span>
                            }
                            features={["Up to 5Tb storage", "Up to 5 users", "Shared workspace"]}
                        />
                    </Grid>
                </Grid>
            </div>
        </div>
    );
}

PricingSection.propTypes = {};

export default withStyles(styles, { withTheme: true })(PricingSection);