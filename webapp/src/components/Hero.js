import React from "react";
import LazyHero from 'react-lazy-hero';

const heroImgSrc = "https://c.pxhere.com/photos/7a/bb/bridge_light_suspension_bridge_east_river_relexion-87244.jpg!d"

const Hero = (props) => {
  const lowerContent = (
    <p></p>
  )
  return (
    <div className="splash">
        <LazyHero imageSrc={heroImgSrc} opacity={.2} minHeight="90vh" transitionDuration={800}>
          <div className="shadow-30">
            <h1 className="logotype">Fabrikam</h1>
            <h2>Your trusted partner in candidate screening</h2>
            <div>
              {props.children}
            </div>
          </div>
        </LazyHero>
        {lowerContent}
    </div>
  )
}

export default Hero;
