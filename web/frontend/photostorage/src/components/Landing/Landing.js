import React, { Fragment } from "react";
import HeadSection from "./HeadSection";
import FeatureSection from "./FeatureSection";
import PricingSection from "./PricingSection";
import Footer from "./Footer";

function Home() {
    return (
        <Fragment>
            <HeadSection />
            <FeatureSection />
            <PricingSection />
            <Footer />
        </Fragment>
    );
}

export default Home;